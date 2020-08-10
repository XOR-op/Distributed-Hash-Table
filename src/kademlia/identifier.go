package kademlia

import (
	"crypto/sha1"
	"math/big"
)

type Identifier struct {
	Val *big.Int
}

func (self Identifier)Xor(rhs Identifier) (reply Identifier) {
	reply.Val =new(big.Int)
	reply.Val.Xor(self.Val,rhs.Val)
	return
}

func (self Identifier)BitLen()int{
	return self.Val.BitLen()
}

func (self *Identifier)As(rhs *Identifier)  {
	self.Val=new(big.Int)
	self.Val.Set(rhs.Val)
}

func NewIdentifier(s string)(reply *Identifier)  {
	reply=new(Identifier)
	reply.Val=new(big.Int)
	tmp := sha1.Sum([]byte(s))
	reply.Val.SetBytes(tmp[:])
	return
}

func (self Identifier)LessThan(rhs *Identifier)bool  {
	return self.Val.Cmp(rhs.Val)<0
}

func NewMidIdentifier(lessBit uint) (reply *Identifier) {
	reply=new(Identifier)
	if lessBit==0{
		reply.Val=new(big.Int).SetUint64(0)
		return
	}
	reply.Val=new(big.Int)
	one:=new(big.Int).SetUint64(1)
	less:=new(big.Int).Lsh(one,lessBit-1)
	more:=new(big.Int).Lsh(one,lessBit)
	reply.Val.Add(less,more)
	return
}
