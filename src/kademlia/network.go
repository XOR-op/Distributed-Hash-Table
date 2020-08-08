package kademlia

import (
	"net"
	"net/rpc"
	"time"
)

const (
	dialTimeout  = 300 * time.Millisecond
	retryTimeout = 50 * time.Millisecond
	retryTimes=3
)

type TemporaryError interface {
	Temporary() bool
}

func OldPing(addr *Contact) bool {
	client, err := net.DialTimeout("tcp", addr.Address, dialTimeout)
	for retryCnt := 0; err != nil && retryCnt < 3; retryCnt++ {
		if terr, ok := err.(TemporaryError); ok && terr.Temporary() {
			time.Sleep(time.Duration(retryCnt+1) * retryTimeout)
			client, err = net.DialTimeout("tcp", addr.Address, dialTimeout)
		} else {
			return false
		}
	}
	if err == nil {
		_ = client.Close()
		return true
	}
	return false
}

//func (self *Node)Ping(addr *Contact) bool {
//	if Ping(addr) {
//		self.table.UpdateContact(addr, true)
//		return true
//	}
//	return false
//}

func (self *Node)Dial(addr *Contact) (reply *RPCClient, err error) {
	// todo should have used UDP
	client, err := rpc.Dial("tcp", addr.Address)
	for retryCnt := 0; err != nil && retryCnt < retryTimes; retryCnt++ {
		// avoid "resource temporarily unavailable" error
		if terr, ok := err.(TemporaryError); ok && terr.Temporary() {
			time.Sleep(time.Duration(retryCnt+1) * retryTimeout)
			client, err = rpc.Dial("tcp", addr.Address)
		} else {
			break
		}
	}
	if err==nil {
		return &RPCClient{client, addr.Duplicate(),&self.addr,&self.table}, nil
	}else {
		self.log.Warning("Dial",addr.Port,"failed:",err)
		return nil,err
	}
}
