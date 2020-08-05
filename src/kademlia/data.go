package kademlia

import (
	"sync"
	"time"
)

type Data struct {
	Key, Value string
}

type Storage struct {
	data        map[string]string
	expireTable map[string]time.Time
	lock        sync.RWMutex
}

func (self *Storage) Has(key string) (ok bool) {
	self.lock.RLock()
	defer self.lock.RUnlock()
	_, ok = self.data[key]
	return
}

func (self *Storage) Get(key string) (value string,ok bool) {
	self.lock.RLock()
	defer self.lock.RUnlock()
	value,ok=self.data[key]
	return
}

func (self *Storage) Store(key, value string, expireT time.Time) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.data[key]=value
	self.expireTable[key]=expireT
}

func (self *Storage)Expire(){
	self.lock.Lock()
	defer self.lock.Unlock()
	toDelete:=make([]string,0)
	for k,t:=range self.expireTable{
		if time.Now().After(t){
			toDelete=append(toDelete, k)
		}
	}
	for _,k:=range toDelete{
		delete(self.data,k)
		delete(self.expireTable,k)
	}
}
