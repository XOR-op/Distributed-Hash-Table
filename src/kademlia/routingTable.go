package kademlia

type RoutingTable struct {
	selfIDRef *Identifier
	elements  [Width]*KBucket
	curUpdate int
}

func (self *RoutingTable) KClosest(targetID Identifier) (reply []*Contact, amount int) {
	reply=make([]*Contact,0)
	for i, _ := range reply {
		reply[i] = nil
	}
	amount = 0
	bucketID := targetID.Xor(*self.selfIDRef).BitLen() - 1
	for i := bucketID; i >= 0; i-- {
		if self.elements[i].fill(&reply, &amount, K) {
			return
		}
	}
	for i := bucketID + 1; i < Width; i++ {
		if self.elements[i].fill(&reply, &amount, K) {
			return
		}
	}
	return
}

func (self *RoutingTable) UpdateContact(addr *Contact) {
	bucketID := addr.ID.Xor(*self.selfIDRef).BitLen() - 1
	if bucketID==-1{
		return
	}
	self.elements[bucketID].Add(addr)
}

func NewRoutingTable(id *Identifier)(reply *RoutingTable)  {
	reply=new(RoutingTable)
	reply.selfIDRef=id
	reply.curUpdate=0
	for i,_:=range reply.elements{
		reply.elements[i]=NewKBucket(K)
	}
	return
}

func (self *RoutingTable)GoOn()  {
	for i:=range self.elements{
		self.elements[i].RefreshTime()
	}
}

