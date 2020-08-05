package kademlia

import (
	log "github.com/sirupsen/logrus"
	"sync/atomic"
	"time"
)

const (
	K                     = 20
	alpha                 = 3
	Width                 = 160
	defaultExpireDuration = 30 * time.Second
)

type Node struct {
	log      Logger
	table    RoutingTable
	addr     Contact
	database Storage
}

func (self *Node) FindKClosest(key string) []*Contact {
	// initialization
	query := FindValueRequest{
		Key:  key,
		ID:   *NewIdentifier(key),
		Auth: &self.addr,
	}
	seen := make(map[string]struct{})
	list, n := self.table.KClosest(query.ID)
	pendingList := list[:]
	doneList := make([]*Contact,0)
	// begin iteration
	channel := make(chan *FindNodeResponse, K)
	hasConn:=new(int32)
	*hasConn=0
	index:=0
	for ; index < n && *hasConn < alpha; index++ {
		// initial call
		seen[pendingList[index].Address] = struct{}{}
		if OldPing(pendingList[index]) {
			doneList = append(doneList, pendingList[index])
			atomic.AddInt32(hasConn,1)
			go func(addr *Contact) {
				defer atomic.AddInt32(hasConn,-1)
				client, err := self.Dial(addr)
				if err != nil {
					log.Warning("Dial ", addr.Port, " Failed")
					return
				}
				client.FindNodeAsync(&addr.ID, channel)
				client.Close()
			}(pendingList[index])
		}
	}
	for *hasConn>0||index < len(pendingList) {
		// hasConn will always be leq than alpha
		if *hasConn>0 {
			response := <-channel
			atomic.AddInt32(hasConn, -1)
			if response.err != nil {
				continue
			}
			doneList = append(doneList, response.Auth)
			for _,v:=range response.KNodes[:response.Amount]{
				pendingList=append(pendingList,v)
			}
		}
		for *hasConn<alpha{
			seen[pendingList[index].Address] = struct{}{}
			atomic.AddInt32(hasConn,1)
			go func(addr *Contact) {
				defer atomic.AddInt32(hasConn,-1)
				client, err := self.Dial(addr)
				if err != nil {
					log.Warning("Dial ", addr.Port, " Failed")
					return
				}
				client.FindNodeAsync(&addr.ID, channel)
				client.Close()
			}(pendingList[index])
			index++
		}
	}
	SortContactSlice(doneList,&query.ID)
	if len(doneList)>K{
		return doneList[:K]
	}
	return doneList
}

func (self *Node) Get(key string) (value string, ok bool) {
}

func (self *Node) Store(key, value string) (ok bool) {

}



