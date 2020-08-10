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
	sleepDuration         = 30 * time.Millisecond
	tExpire               = 75 * time.Second
	tRefresh              = 12 * time.Second
	tDuplicate            = 10 * time.Second
	tRepublish            = 55 * time.Second
)

type Node struct {
	log      Logger
	table    RoutingTable
	addr     Contact
	database Storage
	running bool
	server  *rpc.Server
}

func (self *Node) FindKClosest(key string) []*Contact {
	self.log.Trace("start find key", key)
	return self.FindKClosestSHA1(NewIdentifier(key))
}
func (self *Node) FindKClosestSHA1(queryID *Identifier) []*Contact {
	self.log.Trace("start find id", queryID)
	// initialization
	seen := make(map[string]struct{})
	list, n := self.table.KClosest(*queryID)
	self.log.Debug("local find",n,"closest")
	pendingList := list[:n]
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
			self.log.Trace("Succeed ping",pendingList[index].Port)
			doneList = append(doneList, pendingList[index])
			atomic.AddInt32(hasConn, 1)
			go func(addr *Contact) {
				//defer atomic.AddInt32(hasConn, -1)
				client, err := self.Dial(addr)
				if err != nil {
					return
				}
				client.FindNodeAsync(&addr.ID, channel)
				client.Close()
			}(pendingList[index])
		}
	}
	for index < len(pendingList)||*hasConn>0 {
		// hasConn will always be leq than alpha
		if *hasConn > 0 {
			response := <-channel
			atomic.AddInt32(hasConn, -1)
			if response.Err != nil {
				self.log.Warning(response.Err)
				continue
			}
			doneList = append(doneList, response.Auth)
			for _, v := range response.KNodes[:response.Amount] {
				pendingList = append(pendingList, v)
			}
		}
		for *hasConn < alpha && index < len(pendingList) {
			if _, found := seen[pendingList[index].Address]; !found {
				seen[pendingList[index].Address] = struct{}{}
				atomic.AddInt32(hasConn, 1)
				go func(addr *Contact) {
					//defer atomic.AddInt32(hasConn, -1)
					client, err := self.Dial(addr)
					if err != nil {
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
	defer self.log.Debug("find ", len(doneList), "elements from id", queryID)
	return doneList
}

func (self *Node) Get(key string) (value string, ok bool) {
	self.log.Trace("start getting key", key)
	if self.database.Has(key) {
		self.log.Trace("instantly get",key)
		return self.database.Get(key)
	}
	queryID := NewIdentifier(key)
	seen := make(map[string]struct{})
	self.log.Trace("KClosest begin")
	list, n := self.table.KClosest(*queryID)
	self.log.Trace("KClosest end")
	shortList := list[:n]
	SortContactSlice(shortList, queryID)
	closestNode := shortList[0]
	channel := make(chan *FindValueResponse, K*2)
	hasConn := 0
	self.log.Trace("get() first shot")
	for index := 0; index < n && hasConn < alpha; index++ {
		// initial call
		seen[shortList[index].Address] = struct{}{}
		if OldPing(shortList[index]) {
			self.log.Trace("ping succeed",shortList[index])
			atomic.AddInt32(&GlobalRequest, 1)
			DefaultLogger.Debug("GlobalRequest", GlobalRequest, shortList[index].Port)
			hasConn++
			go self.findValueAsync(key, shortList[index], channel)
		}
	}

	closestNodeBackup := closestNode
	for hasConn > 0 {
		closestNodeBackup = closestNode
		self.log.Trace("wait")
		response := <-channel
		self.log.Trace("received")
		hasConn--
		if response.Err == nil {
			if response.Stat == FindValue {
				// store and return
				self.log.Trace("Early found")
				for i := range shortList {
					if _, found := seen[shortList[i].Address]; found && shortList[i].Address != response.Auth.Address {
						go self.storeAsync(key, response.Value, shortList[i], false)
					}
				}
				return response.Value, true
			}
			// merge
			self.log.Trace("Merge")
			shortList = append(shortList, response.KNodes[:min(alpha, len(response.KNodes))]...)
			SortContactSlice(shortList, queryID)
			shortList = shortList[:min(K, len(shortList))]
			if LessDistance(queryID, shortList[0], closestNode) {
				closestNode = shortList[0]
			}
		} else {
			self.log.Warning("response err:", response.Err)
			for i := range shortList {
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
			self.log.Trace("has to find all left nodes")
			for i := 0; i < len(shortList); i++ {
				if _, found := seen[shortList[i].Address]; !found {
					atomic.AddInt32(&GlobalRequest, 1)
					DefaultLogger.Debug("GlobalRequest", GlobalRequest, shortList[i].Port)
					hasConn++
					seen[shortList[i].Address] = struct{}{}
					go self.findValueAsync(key, shortList[i], channel)
				}
			}

		}
		self.log.Trace("new conn")
		for i := 0; i < len(shortList) && hasConn < alpha; i++ {
			if _, found := seen[shortList[i].Address]; !found {
				atomic.AddInt32(&GlobalRequest, 1)
				DefaultLogger.Debug("GlobalRequest", GlobalRequest, shortList[i].Port)
				hasConn++
				seen[shortList[i].Address] = struct{}{}
				go self.findValueAsync(key, shortList[i], channel)
			}
		}
	}
	self.log.Debug("Not found")
	return "NOT FOUND", false
}
func (self *Node) Store(key, value string) (ok bool) {
	self.database.OwnStore(key,value)
	return self.subStore(key, value, true)
}

func (self *Node) subStore(key, value string, original bool) (ok bool) {
	result := self.FindKClosest(key)
	self.log.Debug("can store", len(result), "nodes")
	cnt := new(int32)
	*cnt = 0
	for _, addr := range result {
		if *cnt < alpha {
			go func(addr *Contact) {
				atomic.AddInt32(cnt, 1)
				defer atomic.AddInt32(cnt, -1)
				client, err := self.Dial(addr)
				if err == nil {
					_ = client.Store(key, value, original)
					client.Close()
				}
			}(addr)
		} else {
			time.Sleep(sleepDuration)
		}
	}
	return true
}

func (self *Node) Join(port int) {
	self.table.GoOn()
	self.table.UpdateContact(NewContact(port))
	self.FindKClosest(self.addr.Address)
}

func (self *Node) Quit() {
	self.running = false
}

func NewNode(port int) (reply *Node) {
	reply = new(Node)
	reply.log = *NewLogger(strconv.Itoa(port))
	reply.addr = *NewContact(port)
	reply.table = *NewRoutingTable(&reply.addr.ID)
	reply.database = *NewStorage()
	reply.running = true
	// todo udp
	reply.server = rpc.NewServer()
	err := reply.server.Register(&RPCInterface{reply})
	if err != nil {
		reply.log.Fatal("register failed ", err)
	}
	l, err := net.Listen("tcp", reply.addr.Address)
	if err != nil {
		reply.log.Fatal("listen failed ", err)
	}
	go reply.rpcServe(l)
	go reply.RunDaemon()
	return
}


func (self *Node) Status() bool {
	return self.running
}
