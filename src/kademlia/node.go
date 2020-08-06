package kademlia

import (
	"net"
	"net/rpc"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	K                     = 20
	alpha                 = 3
	Width                 = 160
	defaultExpireDuration = 30 * time.Second
	sleepDuration         = 30 * time.Millisecond
)

type Node struct {
	log      Logger
	table    RoutingTable
	addr     Contact
	database Storage
	running  bool
	server   *rpc.Server
}

func (self *Node) FindKClosest(key string) []*Contact {
	// initialization
	queryID := NewIdentifier(key)
	seen := make(map[string]struct{})
	list, n := self.table.KClosest(*queryID)
	pendingList := list[:]
	doneList := make([]*Contact, 0)
	// begin iteration
	channel := make(chan *FindNodeResponse, K)
	hasConn := new(int32)
	*hasConn = 0
	index := 0
	for ; index < n && *hasConn < alpha; index++ {
		// initial call
		seen[pendingList[index].Address] = struct{}{}
		if OldPing(pendingList[index]) {
			doneList = append(doneList, pendingList[index])
			atomic.AddInt32(hasConn, 1)
			go func(addr *Contact) {
				defer atomic.AddInt32(hasConn, -1)
				client, err := self.Dial(addr)
				if err != nil {
					self.log.Warning("Dial ", addr.Port, " Failed")
					return
				}
				client.FindNodeAsync(&addr.ID, channel)
				client.Close()
			}(pendingList[index])
		}
	}
	for *hasConn > 0 || index < len(pendingList) {
		// hasConn will always be leq than alpha
		if *hasConn > 0 {
			response := <-channel
			atomic.AddInt32(hasConn, -1)
			if response.err != nil {
				continue
			}
			doneList = append(doneList, response.Auth)
			for _, v := range response.KNodes[:response.Amount] {
				pendingList = append(pendingList, v)
			}
		}
		for *hasConn < alpha {
			if _, found := seen[pendingList[index].Address]; !found {
				seen[pendingList[index].Address] = struct{}{}
				atomic.AddInt32(hasConn, 1)
				go func(addr *Contact) {
					defer atomic.AddInt32(hasConn, -1)
					client, err := self.Dial(addr)
					if err != nil {
						self.log.Warning("Dial ", addr.Port, " Failed")
						return
					}
					client.FindNodeAsync(&addr.ID, channel)
					client.Close()
				}(pendingList[index])
			}
			index++
		}
	}
	SortContactSlice(doneList, queryID)
	if len(doneList) > K {
		return doneList[:K]
	}
	return doneList
}

func (self *Node) Get(key string) (value string, ok bool) {
	queryID := NewIdentifier(key)
	seen := make(map[string]struct{})
	list, n := self.table.KClosest(*queryID)
	shortList := list[:]
	SortContactSlice(shortList, queryID)
	closestNode := shortList[0]
	channel := make(chan *FindValueResponse, K)
	hasConn := 0
	for index := 0; index < n && hasConn < alpha; index++ {
		// initial call
		seen[shortList[index].Address] = struct{}{}
		if OldPing(shortList[index]) {
			hasConn++
			go self.findValueAsync(key, shortList[index], channel)
		}
	}
	closestNodeBackup := closestNode
	for hasConn > 0 {
		closestNodeBackup = closestNode
		response := <-channel
		hasConn--
		if response.err == nil {
			if response.Stat == FindValue {
				// store and return
				for i, _ := range shortList {
					if _, found := seen[shortList[i].Address]; found && shortList[i].Address != response.Auth.Address {
						go self.storeAsync(key, *response.Value, shortList[i])
					}
				}
				return *response.Value, true
			}
			// merge
			shortList = append(shortList, (*response.KNodes)[:alpha]...)
			SortContactSlice(shortList, queryID)
			shortList = shortList[:K]
			if LessDistance(queryID, shortList[0], closestNode) {
				closestNode = shortList[0]
			}
		} else {
			for i, _ := range shortList {
				if shortList[i].Address == response.Auth.Address {
					// remove failed contact
					length := len(shortList)
					shortList[length-1], shortList[i] = shortList[i], shortList[length-1]
					shortList = shortList[:length]
					// fix closest node
					if response.Auth.Address == closestNode.Address {
						closestNode = shortList[0]
					} else if LessDistance(queryID, shortList[0], closestNode) {
						closestNode = shortList[0]
					}
					break
				}
			}
		}
		if hasConn == 0 && closestNode == closestNodeBackup {
			// call all left
			// return false or recall alpha nodes
			for i := 0; i < len(shortList); i++ {
				if _, found := seen[shortList[i].Address]; !found {
					hasConn++
					go self.findValueAsync(key, shortList[i], channel)
				}
			}

		}
		for i := 0; i < len(shortList) && hasConn < alpha; i++ {
			if _, found := seen[shortList[i].Address]; !found {
				hasConn++
				go self.findValueAsync(key, shortList[i], channel)
			}
		}
	}
	return "NOT FOUND", false
}

func (self *Node) Store(key, value string) (ok bool) {
	result := self.FindKClosest(key)
	cnt := new(int32)
	*cnt = 0
	for _, addr := range result {
		if *cnt < alpha {
			go func(addr *Contact) {
				atomic.AddInt32(cnt, 1)
				defer atomic.AddInt32(cnt, -1)
				client, err := self.Dial(addr)
				if err == nil {
					_ = client.Store(key, value)
					client.Close()
				}
			}(addr)
		} else {
			time.Sleep(sleepDuration)
		}
	}
	return true
}

func (self *Node)Join(port int)  {
	self.table.UpdateContact(NewContact(port))
	self.FindKClosest(self.addr.Address)
}

func NewNode(port int) (reply *Node) {
	reply = new(Node)
	reply.log = *NewLogger(strconv.Itoa(port))
	reply.addr = *NewContact(port)
	reply.table = *NewRoutingTable(&reply.addr.ID)
	reply.database = *NewStorage()
	// todo udp
	reply.server = rpc.NewServer()
	_ = reply.server.Register(RPCInterface{reply})
	l, err := net.Listen("tcp", reply.addr.Address)
	if err != nil {
		reply.log.Fatal("listen failed ", err)
	}
	go reply.rpcServe(l)
	return
}
