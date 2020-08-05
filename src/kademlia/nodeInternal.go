package kademlia

import "time"

func (self *Node) RPCFindNode(request FindNodeRequest, reply *FindNodeResponse) error {
	var stat bool
	self.table.UpdateContact(request.Auth)
	reply.KNodes, reply.Amount = self.table.KClosest(request.Target)
	reply.Stat = Triple(stat, FindEnough, FindNotEnough).(FindStat)
	return nil
}

func (self *Node) RPCFindValue(request FindValueRequest, reply *FindValueResponse) error {
	self.table.UpdateContact(request.Auth)
	if self.database.Has(request.Key) {
		reply.Stat = FindValue
		reply.KNodes = nil
		reply.Value = new(string)
		*reply.Value,_ = self.database.Get(request.Key)
	} else {
		reply.KNodes = new([K]*Contact)
		*reply.KNodes, reply.Amount = self.table.KClosest(request.ID)
		reply.Stat = Triple(reply.Amount==K, FindEnough, FindNotEnough).(FindStat)
		reply.Value = nil
	}
	return nil
}

func (self *Node) RPCPing(request PingRequest) error {
	self.table.UpdateContact(request.Auth)
	return nil
}

func (self *Node) RPCStore(request StoreRequest)error  {
	self.table.UpdateContact(request.Auth)
	self.database.Store(request.Key,request.Value,time.Now().Add(defaultExpireDuration))
	return nil
}
