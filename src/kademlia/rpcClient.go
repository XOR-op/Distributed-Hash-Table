package kademlia

import "net/rpc"

type RPCClient struct {
	conn    *rpc.Client
	address Contact
}

func (this *RPCClient) Close() {
	_ = this.conn.Close()
}
