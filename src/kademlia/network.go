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

func Ping(addr *Contact) bool {
	client, err := net.DialTimeout("tcp", addr.address, dialTimeout)
	for retryCnt := 0; err != nil && retryCnt < 3; retryCnt++ {
		if terr, ok := err.(TemporaryError); ok && terr.Temporary() {
			time.Sleep(time.Duration(retryCnt+1) * retryTimeout)
			client, err = net.DialTimeout("tcp", addr.address, dialTimeout)
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

func Dial(addr *Contact) (reply *RPCClient, err error) {
	client, err := rpc.Dial("tcp", addr.address)
	for retryCnt := 0; err != nil && retryCnt < retryTimes; retryCnt++ {
		// avoid "resource temporarily unavailable" error
		if terr, ok := err.(TemporaryError); ok && terr.Temporary() {
			time.Sleep(time.Duration(retryCnt+1) * retryTimeout)
			client, err = rpc.Dial("tcp", addr.address)
		} else {
			break
		}
	}
	if err==nil {
		return &RPCClient{client, *addr.Duplicate()}, nil
	}else {
		return nil,err
	}
}
