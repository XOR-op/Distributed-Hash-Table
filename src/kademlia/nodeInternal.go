package kademlia

import (
	"math/rand"
	"net"
	"time"
)

func (self *Node) rpcFindNode(request FindNodeRequest, reply *FindNodeResponse) error {
	var stat bool
	self.table.UpdateContact(request.Auth)
	reply.KNodes, reply.Amount = self.table.KClosest(request.Target)
	reply.Stat = Triple(stat, FindEnough, FindNotEnough).(FindStat)
	return nil
}

func (self *Node) rpcFindValue(request FindValueRequest, reply *FindValueResponse) error {
	self.table.UpdateContact(request.Auth)
	if self.database.Has(request.Key) {
		reply.Stat = FindValue
		reply.KNodes = nil
		reply.Value = new(string)
		*reply.Value, _ = self.database.Get(request.Key)
	} else {
		reply.KNodes = new([K]*Contact)
		*reply.KNodes, reply.Amount = self.table.KClosest(request.ID)
		reply.Stat = Triple(reply.Amount == K, FindEnough, FindNotEnough).(FindStat)
		reply.Value = nil
	}
	return nil
}

func (self *Node) rpcPing(request PingRequest) error {
	self.table.UpdateContact(request.Auth)
	return nil
}

func (self *Node) rpcStore(request StoreRequest) error {
	self.table.UpdateContact(request.Auth)
	self.database.Store(request.Key, request.Value, time.Now().Add(defaultExpireDuration*time.Duration(85+rand.Intn(15))/time.Duration(100)))
	return nil
}
func (self *Node) findValueAsync(key string, addr *Contact, channel chan *FindValueResponse) {
	client, err := self.Dial(addr)
	if err != nil {
		self.log.Warning("Dial ", addr.Port, " Failed")
		return
	}
	client.FindValueAsync(key, &addr.ID, channel)
	client.Close()
}

func (self *Node) storeAsync(key, value string, addr *Contact) {
	client, err := self.Dial(addr)
	if err != nil {
		self.log.Warning("Dial ", addr.Port, " Failed")
		return
	}
	client.Store(key, value)
	client.Close()
}

func (self *Node) rpcServe(l net.Listener) {
	defer l.Close()
	defer self.log.Debug("RPC server Quited.")
	for self.running {
		conn, err := l.Accept()
		if err != nil {
			self.log.Warning("rpc.Serve: accept:", err.Error())
			return
		}
		go self.server.ServeConn(conn)
	}
}
