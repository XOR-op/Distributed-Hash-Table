package kademlia

import "math/big"

type Identifier struct {
	val *big.Int
}

func (this *Identifier)Xor(rhs Identifier) (reply Identifier) {
	reply.val=new(big.Int)
	reply.val.Xor(this.val,rhs.val)
	return
}

func (this *Identifier)BitLen()int{
	return this.val.BitLen()
}

func (this *Identifier)As(rhs *Identifier)  {
	this.val.Set(rhs.val)
}

