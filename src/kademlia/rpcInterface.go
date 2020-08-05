package kademlia


type RPCInterface struct {
	N *Node
}

func (self *RPCInterface)FindNode(arg FindNodeRequest,rep *FindNodeResponse)(err error)  {
	return self.N.RPCFindNode(arg,rep)
}

func (self *RPCInterface)Ping(arg PingRequest,_ *int)(err error)  {
	return self.N.RPCPing(arg)
}
