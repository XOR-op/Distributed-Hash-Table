package chord

import (
	"crypto/sha1"
	"math/big"
	"strconv"
	"strings"
)



type Identifier struct {
	ValPtr *big.Int
}


// (a,b)
func (this *Identifier) In(low, high *Identifier) bool {
	if val:=low.ValPtr.Cmp(high.ValPtr);val > 0 {
		// low > high
		return low.ValPtr.Cmp(this.ValPtr) < 0 || this.ValPtr.Cmp(high.ValPtr) < 0
	}else if val==0{
		return low.ValPtr.Cmp(this.ValPtr)!=0 // a circle
	}
	return low.ValPtr.Cmp(this.ValPtr) < 0 && this.ValPtr.Cmp(high.ValPtr) < 0
}

// (a,b]
func (this *Identifier) InRightClosure(low, high *Identifier) bool {
	return high.ValPtr.Cmp(this.ValPtr) == 0 || this.In(low, high)
}

// [a,b)
func (this *Identifier) InLeftClosure(low, high *Identifier) bool {
	return low.ValPtr.Cmp(this.ValPtr) == 0 || this.In(low, high)
}

func (this *Identifier) PlusTwoPower(exp uint) Identifier {
	y := big.NewInt(1)
	y.Lsh(y,exp)
	y.Add(y,this.ValPtr)
	if y.BitLen()> BIT_WIDTH {
		tmp:=big.NewInt(1)
		y.Mod(y,tmp.Lsh(tmp,uint(BIT_WIDTH)))
	}
	return Identifier{y}
}
func (this *Identifier)CopyFrom(id *Identifier)  {
	if id.ValPtr==nil{
		this.ValPtr=nil
		return
	}
	if this.ValPtr==nil {
		this.ValPtr = new(big.Int)
	}
	this.ValPtr.Set(id.ValPtr)
}


func IDlize(s string) (ret Identifier) {
	ret.ValPtr =new(big.Int)
	tmp := sha1.Sum([]byte(s))
	ret.ValPtr.SetBytes(tmp[:])
	return
}

type Address struct {
	Addr string // "host:Port"
	Port int
	Id   Identifier
	//lock *sync.Mutex
}

func (this *Address) isNil() bool {
	return this.Addr == ""
}
func (this *Address) Nullify()  {
	this.Addr =""
	this.Port=0
	this.Id.ValPtr=nil
	//this.lock=new(sync.Mutex)
}

func (this *Address)CopyFrom(addr *Address)  {
	//this.lock.Lock()
	//addr.lock.Lock()
	this.Addr=addr.Addr
	this.Port=addr.Port
	this.Id.CopyFrom(&addr.Id)
	//this.lock.Unlock()
	//addr.lock.Unlock()
	//log.Trace("COPY trace:",this.Port," with sha1 ",this.Id)
	//log.Trace("COPY trace original :",addr.Port," with sha1 ",addr.Id)
}

func NewAddress(addr string)(reply Address)  {
	reply.Addr =addr
	reply.Id =IDlize(addr)
	reply.Port, _ =strconv.Atoi(strings.Split(addr,":")[1])
	//reply.lock=new(sync.Mutex)
	return
}

type AddressWithBoolean struct {
	Addr Address
	Stat bool
}

func NewAddressWithBoolean(addr *Address,stat bool)AddressWithBoolean  {
	var stru AddressWithBoolean
	stru.Addr.Nullify()
	stru.Addr.CopyFrom(addr)
	stru.Stat=stat
	return stru
}
