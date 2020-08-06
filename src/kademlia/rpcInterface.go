package kademlia


type RPCInterface struct {
	N *Node
}

func (self *RPCInterface)FindNode(arg FindNodeRequest,rep *FindNodeResponse)(err error)  {
	return self.N.rpcFindNode(arg,rep)
}

func (self *RPCInterface)Ping(arg PingRequest,_ *int)(err error)  {
	return self.N.rpcPing(arg)
}
