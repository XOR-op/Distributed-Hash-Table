package kademlia

import (
	"github.com/sasha-s/go-deadlock"
	"time"
)

type KBucket struct {
	lookupTable    map[string]*bucketNode
	indicator      bucketNode
	Size           int
	maxSize        int
	lock           deadlock.RWMutex
	lastUpdateTime time.Time
}

func (self *KBucket)RefreshTime()  {
	self.lastUpdateTime=time.Now()
}

func (self *KBucket) fill(target *[]*Contact, curIndex *int, until int) bool {
	self.lock.RLock()
	defer self.lock.RUnlock()
	for iter := self.Head(); *curIndex < until && iter != self.VirtualNode(); *curIndex++ {
		*target = append(*target, iter.element.Duplicate())
		iter = iter.next
	}
	return *curIndex == until
}

func (self *KBucket) Add(addr *Contact) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.RefreshTime()
	val, ok := self.lookupTable[addr.Address]
	if ok {
		val.Detach()
		val.attachAfter(self.VirtualNode())
	} else {
		if self.Size < self.maxSize {
			n := newBucketNode(addr)
			n.attachAfter(self.VirtualNode())
			self.lookupTable[addr.Address] = n
		} else {
			DefaultLogger.Debug("full")
			n := self.Tail().Detach()
			if n.element.TestConn() {
				n.attachAfter(self.VirtualNode())
			} else {
				delete(self.lookupTable, n.element.Address)
				n = newBucketNode(addr)
				n.attachAfter(self.VirtualNode())
				self.lookupTable[addr.Address] = n
			}
		}
	}
}

func (self *KBucket)Drop(addr *Contact){
	self.lock.Lock()
	defer self.lock.Unlock()
	val, ok := self.lookupTable[addr.Address]
	if ok{
		val.Detach()
		delete(self.lookupTable,val.element.Address)
	}

}

func NewKBucket(maxSize int) (reply *KBucket) {
	reply = new(KBucket)
	reply.indicator = bucketNode{nil, nil, nil}
	reply.indicator.prev = &reply.indicator
	reply.indicator.next = &reply.indicator
	reply.Size = 0
	reply.lookupTable = make(map[string]*bucketNode)
	reply.maxSize = maxSize
	reply.RefreshTime()
	return
}

func (self *KBucket) Head() *bucketNode {
	return self.indicator.next
}

func (self *KBucket) Tail() *bucketNode {
	return self.indicator.prev
}

func (self *KBucket) VirtualNode() *bucketNode {
	return &self.indicator
}

type bucketNode struct {
	next    *bucketNode
	prev    *bucketNode
	element *Contact
}

func newBucketNode(val *Contact) (reply *bucketNode) {
	reply = new(bucketNode)
	reply.next = nil
	reply.prev = nil
	reply.element = val.Duplicate()
	return
}

func (self *bucketNode) attachAfter(n *bucketNode) {
	self.next = n.next
	self.prev = n
	n.next.prev = self
	n.next = self
}

func (self *bucketNode) Detach() *bucketNode {
	self.next.prev = self.prev
	self.prev.next = self.next
	self.next = nil
	self.prev = nil
	return self
}
