package kademlia

type RoutingTable struct {
	selfIDRef *Identifier
	elements  [Width]*KBucket
}

func (self *RoutingTable) KClosest(targetID Identifier) (reply []*Contact, amount int) {
	/*
	defer func() {
		// sort
		sort.Slice(reply, func(i, j int) bool {
			if reply[i]==nil||reply[j]==nil{
				return reply[j]==nil
			}
			disI:=targetID.Xor(reply[i].ID)
			disJ:=targetID.Xor(reply[j].ID)
			return disI.LessThan(&disJ)
		})
	}()*/
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
	for i,_:=range reply.elements{
		reply.elements[i]=NewKBucket(K)
	}
	return
}