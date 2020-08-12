package kademlia

import (
	"net/rpc"
)

//var (
//	GlobalRequest  int32 = 0
//	GlobalResponse int32 = 0
//)

type RPCClient struct {
	conn       *rpc.Client
	remoteAddr *Contact
	selfAddr   *Contact
	tablePtr   *RoutingTable
}

func (self *RPCClient) Close() {
	_ = self.conn.Close()
}

func (self *RPCClient) update(err error) {
	if err == nil {
		self.tablePtr.UpdateContact(self.remoteAddr)
	}
}

func (self *RPCClient) FindNode(id *Identifier) (reply *FindNodeResponse, err error) {
	reply = new(FindNodeResponse)
	err = self.conn.Call("RPCInterface.FindNode", FindNodeRequest{*id, self.selfAddr}, reply)
	self.update(err)
	return
}

func (self *RPCClient) FindNodeAsync(id *Identifier, channel chan *FindNodeResponse) {
	reply := new(FindNodeResponse)
	err := self.conn.Call("RPCInterface.FindNode", FindNodeRequest{*id, self.selfAddr}, reply)
	reply.Err = err
	reply.Auth = self.remoteAddr
	channel <- reply
}

func (self *RPCClient) FindValueAsync(key string, id *Identifier, channel chan *FindValueResponse) {
	reply := new(FindValueResponse)
	err := self.conn.Call("RPCInterface.FindValue", FindValueRequest{key, *id, self.selfAddr}, reply)
	reply.Err = err
	reply.Auth = self.remoteAddr
	//atomic.AddInt32(&GlobalResponse, 1)
	//DefaultLogger.Debug("GlobalResponse", GlobalResponse, self.remoteAddr.Port)
	channel <- reply
}

func (self *RPCClient) Ping() (err error) {
	arg := PingRequest{self.selfAddr}
	err = self.conn.Call("RPCInterface.Ping", arg, nil)
	self.update(err)
	return
}

func (self *RPCClient) Store(key, value string, original bool) (err error) {
	arg := StoreRequest{key, value, self.selfAddr, original}
	err = self.conn.Call("RPCInterface.Store", arg, nil)
	self.update(err)
	return
}
