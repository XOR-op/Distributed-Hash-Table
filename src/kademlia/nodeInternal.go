package kademlia

import (
	"net"
	"time"
)

const (
	daemonInterval = 5 * time.Second
)

func (self *Node) rpcFindNode(request FindNodeRequest, reply *FindNodeResponse) error {
	self.log.Trace(" invoked of RPCFindNode from", request.Auth.Port)
	self.table.UpdateContact(request.Auth)
	reply.KNodes, reply.Amount = self.table.KClosest(request.Target)
	reply.Stat = Triple(reply.Amount == K, FindEnough, FindNotEnough).(FindStat)
	return nil
}

func (self *Node) rpcFindValue(request FindValueRequest, reply *FindValueResponse) error {
	self.log.Trace(" invoked of RPCFindValue from", request.Auth.Port)
	self.table.UpdateContact(request.Auth)
	if self.database.Has(request.Key) {
		reply.Stat = FindValue
		reply.Value, _ = self.database.Get(request.Key)
	} else {
		reply.KNodes, reply.Amount = self.table.KClosest(request.ID)
		reply.Stat = Triple(reply.Amount == K, FindEnough, FindNotEnough).(FindStat)
	}
	return nil
}

func (self *Node) rpcPing(request PingRequest) error {
	self.log.Trace(" invoked of RPCPing from", request.Auth.Port)
	self.table.UpdateContact(request.Auth)
	return nil
}

func (self *Node) rpcStore(request StoreRequest) error {
	self.log.Trace(" invoked of RPCStore from", request.Auth.Port)
	self.table.UpdateContact(request.Auth)
	self.database.Store(request.Key, request.Value, request.original)
	return nil
}
func (self *Node) findValueAsync(key string, addr *Contact, channel chan *FindValueResponse) {
	client, err := self.Dial(addr)
	if err != nil {
		self.log.Warning("Dial ", addr.Port, " Failed:", err)
		rt := new(FindValueResponse)
		rt.Err = err
		rt.Auth = addr
		channel <- rt
		return
	}
	client.FindValueAsync(key, &addr.ID, channel)
	client.Close()
}

func (self *Node) storeAsync(key, value string, addr *Contact, original bool) {
	client, err := self.Dial(addr)
	if err != nil {
		self.log.Warning("Dial ", addr.Port, " Failed:", err)
		return
	}
	_ = client.Store(key, value, original)
	client.Close()
}

func (self *Node) rpcServe(l net.Listener) {
	defer l.Close()
	defer self.log.Debug("RPC server Quited.")
	for self.running {
		conn, err := l.Accept()
		if !self.running {
			conn.Close()
			return
		}
		if err != nil {
			self.log.Warning("rpc.Serve: accept:", err.Error())
			return
		}
		go self.server.ServeConn(conn)
	}
}

func (self *Node) Duplicate() {
	self.log.Trace("duplicate")
	data := self.database.NeedDuplicate()
	for k, v := range data {
		self.subStore(k, v, false)
	}
}

func (self *Node) Republish() {
	self.log.Trace("republish")
	data := self.database.NeedRepublish()
	for k, v := range data {
		self.subStore(k, v, true)
	}
}
func (self *Node) Expire() {
	self.log.Trace("expire")
	self.database.Expire()
}

func (self *Node) Refresh() {
	self.log.Trace("refresh")
	if !self.table.elements[self.table.curUpdate].lastUpdateTime.Add(tRefresh).After(time.Now()) {
		self.FindKClosestSHA1(NewMidIdentifier(uint(self.table.curUpdate)))
	}
	self.table.curUpdate=(self.table.curUpdate+1)%Width
}

func (self *Node) RunDaemon() {
	for self.running {
		time.Sleep(daemonInterval)
		if !self.running {
			return
		}
		self.Duplicate()
		if !self.running {
			return
		}
		self.Expire()
		if !self.running {
			return
		}
		self.Republish()
		if !self.running {
			return
		}
		self.Refresh()
	}
}
