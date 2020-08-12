package kademlia

import (
	"sort"
	"strconv"
)

type Contact struct {
	Address string
	Port    int
	ID      Identifier
}

func NewContact(port int)(reply *Contact)  {
	reply=new(Contact)
	reply.Address="localhost:"+strconv.Itoa(port)
	reply.Port=port
	reply.ID=*NewIdentifier(reply.Address)
	return
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

func SortContactSlice(reply *[]*Contact,targetID *Identifier)  {
	sort.Slice(*reply, func(i, j int) bool {
		if (*reply)[i]==nil||(*reply)[j]==nil{
			return (*reply)[j]==nil
		}
		disI:=targetID.Xor((*reply)[i].ID)
		disJ:=targetID.Xor((*reply)[j].ID)
		return disI.LessThan(&disJ)
	})
}

func LessDistance(target *Identifier,i,j *Contact) bool {
	disI:= target.Xor(i.ID)
	disJ:= target.Xor(j.ID)
	return disI.LessThan(&disJ)
}