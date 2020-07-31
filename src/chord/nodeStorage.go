package chord

import (
	log "github.com/sirupsen/logrus"
	"time"
)

const RETRY_TIMEOUT = 20*time.Millisecond

func (this *ChordNode) Put(key, value string) (ret bool) {
	var addr Address
	var err error
	for i:=0;i<3;i++ {
		err=this.FindIdSuccessor(IDlize(key), &addr)
		if err!=nil{
			time.Sleep(RETRY_TIMEOUT)
			continue
		}
		if err = RemoteCall(addr, "RPCWrapper.Put", KVPair{key, value}, &ret); err != nil {
			this.validateSuccessor(false)
			continue
		}else {
			return
		}
	}
	log.Warning(this.addr.Port, " Put pair {", key, ":", value, "} failed:",err)
	return false
}

func (this *ChordNode) Get(key string) (ret bool, value string) {
	var addr Address
	var reply StringWithBoolean
	var err error
	for i:=0;i<3;i++ {
		err=this.FindIdSuccessor(IDlize(key), &addr)
		if err!=nil{
			time.Sleep(RETRY_TIMEOUT)
			continue
		}
		if err = RemoteCall(addr, "RPCWrapper.Get", key, &reply); err != nil {
			this.validateSuccessor(false)
			continue
		}else {
			return reply.Stat, reply.Str
		}
	}
	log.Warning(this.addr.Port, " Get Key{", key, "} failed:",err)
	return false, err.Error()
}

func (this *ChordNode) Delete(key string) (ret bool) {
	var addr Address
	var err error
	for i:=0;i<3;i++ {
		err=this.FindIdSuccessor(IDlize(key), &addr)
		if err!=nil{
			time.Sleep(RETRY_TIMEOUT)
			continue
		}
		if err = RemoteCall(addr, "RPCWrapper.Delete", key, &ret); err != nil {
			this.validateSuccessor(false)
			continue
		}else {
			return
		}
	}
	log.Warning(this.addr.Port, " Delete Key{", key, "} failed:",err)
	return false
}

func (this *ChordNode) MoveData(caller Address, reply *map[string]string) (err error) {
	// return data whose id < caller.id
	log.Trace(this.addr.Port, " [RPC] invoked of MoveData(", caller.Port, ")")
	if this.nodePredecessor.isNil() {
		log.Debug(this.addr.Port, " cannot move data")
		return nil
	}
	client, err := Dial("tcp", this.nodePredecessor.Addr)
	if  err != nil {
		this.validateSuccessor(true)
		client, err = Dial("tcp", this.nodePredecessor.Addr)
	}
	// update backup
	_ = client.Call("RPCWrapper.DropPartialBackup", caller, nil)
	*reply = this.storage.SplitBy(&caller.Id, &this.addr.Id)
	return nil
}

func (this *ChordNode) GetStorage(reply *map[string]string) error {
	// return my storage
	log.Trace(this.addr.Port, " [RPC] invoked of GetStorage()")
	this.storage.lock.RLock()
	defer this.storage.lock.RUnlock()
	*reply = this.storage.Storage
	return nil
}

func (this *ChordNode) UpdateBackup(data *map[string]string) error {
	// add data to backup
	log.Trace(this.addr.Port, " [RPC] invoked of UpdateBackup()")
	this.succStorageBackup.Merge(data)
	return nil
}

func (this *ChordNode) UpdateStorage(data *map[string]string) error {
	// add data to backup
	log.Trace(this.addr.Port, " [RPC] invoked of UpdateStorage()")
	this.storage.Merge(data)
	return nil
}

func (this *ChordNode) DropPartialBackup(caller Address) error {
	log.Trace(this.addr.Port, " [RPC] invoked of DropPartialBackup(", caller.Port, ")")
	newMap:= this.succStorageBackup.SplitBy(&caller.Id, &this.nodeSuccessor.Id)
	this.succStorageBackup.lock.Lock()
	defer this.succStorageBackup.lock.Unlock()
	this.succStorageBackup.Storage=newMap
	return nil
}
