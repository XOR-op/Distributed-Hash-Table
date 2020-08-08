package kademlia


type RPCInterface struct {
	N *Node
}

func (self *RPCInterface)FindNode(arg FindNodeRequest,rep *FindNodeResponse)(err error)  {
	return self.N.rpcFindNode(arg,rep)
}

func (self *RPCInterface)FindValue(arg FindValueRequest,rep *FindValueResponse)(err error)  {
	return self.N.rpcFindValue(arg,rep)
}

func (self *RPCInterface)Ping(arg PingRequest,_ *int)(err error)  {
	return self.N.rpcPing(arg)
}
func (self *RPCInterface)Store(arg StoreRequest,_ *int)(err error)  {
	return self.N.rpcStore(arg)
}
