package chord

import "sync"


type KVPair struct {
	Key   string
	Value string
}

func (this *KVPair)isNil()bool  {
	return this.Key ==""
}

func NewNilChordKV() KVPair {
	return KVPair{"",""}
}


type Table struct {
	Storage map[string]string
	lock    sync.RWMutex
}

func (this *Table)Get(key string,valuePtr *string) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if val,ok:=this.Storage[key];ok{
		*valuePtr=val
		return true
	}
	return false
}

func (this *Table)Put(key,value string) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.Storage[key]=value
	return true
}

func (this *Table)Delete(key string)bool  {
	this.lock.Lock()
	defer this.lock.Unlock()
	if _,ok:=this.Storage[key];ok{
		delete(this.Storage,key)
		return true
	}
	return false
}

func (this *Table)SplitBy(predecessorID,thisID *Identifier)(reply map[string]string){
	this.lock.Lock()
	defer this.lock.Unlock()
	reply=make(map[string]string)
	for k,v:=range this.Storage {
		sha1sum:=IDlize(k)
		if !sha1sum.InRightClosure(predecessorID,thisID){
			reply[k]=v
		}
	}
	for k,_:=range reply {
		delete(this.Storage,k)
	}
	return reply
}

func (this *Table)Merge(rhs *map[string]string)  {
	this.lock.Lock()
	defer this.lock.Unlock()
	for k,v:=range *rhs{
		this.Storage[k]=v
	}
	*rhs=make(map[string]string)
}
