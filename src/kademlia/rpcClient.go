package kademlia

import "net/rpc"

type RPCClient struct {
	conn     *rpc.Client
	address  *Contact
	selfAddr *Contact
	tablePtr *RoutingTable
}

func (self *RPCClient) Close() {
	_ = self.conn.Close()
}

func (self *RPCClient) update(err error) {
	if err == nil {
		self.tablePtr.UpdateContact(self.address)
	}
}

func (self *RPCClient) FindNode(id *Identifier) (reply *FindNodeResponse, err error) {
	reply = new(FindNodeResponse)
	err = self.conn.Call("RPCInterface.FindNode", FindNodeRequest{*id,self.selfAddr}, reply)
	self.update(err)
	return
}

func (self *RPCClient) FindNodeAsync(id *Identifier,channel chan *FindNodeResponse) {
	reply := new(FindNodeResponse)
	err:= self.conn.Call("RPCInterface.FindNode", FindNodeRequest{*id,self.selfAddr}, reply)
	reply.err=err
	channel<-reply
}

func (self *RPCClient) Ping() (err error) {
	arg := PingRequest{self.selfAddr}
	err = self.conn.Call("RPCInterface.Ping", arg,nil)
	self.update(err)
	return
}
