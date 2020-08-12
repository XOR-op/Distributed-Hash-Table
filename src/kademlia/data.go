package kademlia

import (
	//"sync"
	"time"
	"github.com/sasha-s/go-deadlock"
)

type Data struct {
	Key, Value string
}

type Storage struct {
	data           map[string]string
	expireTable    map[string]time.Time
	duplicateTable map[string]time.Time
	republishTable map[string]time.Time
	lock           deadlock.RWMutex
}

func (self *Storage) Has(key string) (ok bool) {
	self.lock.RLock()
	defer self.lock.RUnlock()
	_, ok = self.data[key]
	return
}

func (self *Storage) Get(key string) (value string, ok bool) {
	self.lock.RLock()
	defer self.lock.RUnlock()
	value, ok = self.data[key]
	return
}

func (self *Storage) Store(key, value string, fromOriginal bool) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if _,ok:=self.republishTable[key];ok{
		// self is owner
		return
	}
	if fromOriginal {
		self.expireTable[key] = time.Now().Add(tExpire)
		self.duplicateTable[key] = time.Now().Add(tDuplicate)
	} else if _,ok:=self.data[key];!ok{
		//self.expireTable[key]=time.Now().Add(tExpire)
		self.expireTable[key].Add(tDuplicate)
	}
	self.data[key] = value
}

func (self *Storage)OwnStore(key,value string){
	self.lock.Lock()
	defer self.lock.Unlock()
	self.data[key] = value
	self.duplicateTable[key] = time.Now().Add(tDuplicate)
	self.republishTable[key]=time.Now().Add(tRepublish)
}

func (self *Storage) Expire() {
	self.lock.Lock()
	defer self.lock.Unlock()
	toDelete := make([]string, 0)
	for k, t := range self.expireTable {
		if time.Now().After(t) {
			toDelete = append(toDelete, k)
		}
	}
	for _, k := range toDelete {
		delete(self.data, k)
		delete(self.expireTable, k)
		delete(self.duplicateTable, k)
	}
}

func (self *Storage)NeedDuplicate() (reply map[string]string) {
	self.lock.RLock()
	defer self.lock.RUnlock()
	reply=make(map[string]string)
	for k,t:=range self.duplicateTable{
		if time.Now().After(t){
			reply[k]=self.data[k]
			self.duplicateTable[k]=time.Now().Add(tDuplicate)
		}
	}
	return
}

func (self *Storage)NeedRepublish() (reply map[string]string) {
	self.lock.RLock()
	defer self.lock.RUnlock()
	reply=make(map[string]string)
	for k,t:=range self.republishTable{
		if time.Now().After(t){
			reply[k]=self.data[k]
			self.republishTable[k]=time.Now().Add(tRepublish)
		}
	}
	return
}

func NewStorage() (reply *Storage) {
	reply = new(Storage)
	reply.data = make(map[string]string)
	reply.expireTable = make(map[string]time.Time)
	reply.republishTable=make(map[string]time.Time)
	reply.duplicateTable=make(map[string]time.Time)
	return
}
