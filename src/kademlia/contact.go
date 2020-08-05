package kademlia

import "sort"

type Contact struct {
	Address string
	Port    int
	ID      Identifier
}

func (self *Contact)Duplicate()(reply *Contact)  {
	reply=new(Contact)
	reply.Address = self.Address
	reply.Port = self.Port
	reply.ID.As(&self.ID)
	return
}

func (self *Contact)TestConn() bool {
	return OldPing(self)
}

func SortContactSlice(reply []*Contact,targetID *Identifier)  {
	sort.Slice(reply, func(i, j int) bool {
		if reply[i]==nil||reply[j]==nil{
			return reply[j]==nil
		}
		disI:=targetID.Xor(reply[i].ID)
		disJ:=targetID.Xor(reply[j].ID)
		return disI.LessThan(&disJ)
	})
}