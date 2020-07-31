package chord

import log "github.com/sirupsen/logrus"

func (this *ChordNode) RPCSuccessor(reply *Address)error  {
	log.Trace(this.addr.Port," [RPC] invoked of RPCSuccessor()")
	this.MayFatal()
	this.fingerAndSuccessorLock.RLock()
	reply.CopyFrom(this.nodeSuccessor)
	this.fingerAndSuccessorLock.RUnlock()
	reply.Validate(false,"RPCSuccessor")
	return nil
}

func (this *ChordNode) RPCPredecessor(reply *Address)error  {
	log.Trace(this.addr.Port," [RPC] invoked of RPCPredecessor()")
	this.CheckPredecessor()
	this.predecessorLock.RLock()
	reply.CopyFrom(&this.nodePredecessor)
	this.predecessorLock.RUnlock()
	//reply.Validate(false,"RPCPredecessor")
	return nil
}

func (this *ChordNode) RPCFindIDSuccessorWithValidation(id Identifier,reply *AddressWithBoolean)error  {
	this.validateSuccessor(false)
	return this.RPCFindIDSuccessor(id,reply)
}

func (this *ChordNode) RPCFindIDSuccessor(id Identifier,reply *AddressWithBoolean) error{
	log.Trace(this.addr.Port," [RPC] invoked of RPCFindIDSuccessor(",id,")")
	//this.MayFatal()
	this.validateSuccessor(false)
	this.fingerAndSuccessorLock.RLock()
	rep:=id.InRightClosure(&this.addr.Id,&this.nodeSuccessor.Id)
	this.fingerAndSuccessorLock.RUnlock()
	if rep{
		*reply=NewAddressWithBoolean(this.nodeSuccessor,true)
		//reply.Addr.Validate(false,"RPCFindIDSuccessor true branch")
		return nil
	}else {
		*reply=NewAddressWithBoolean(this.closestPrecedingNode(id),false)
		log.Trace(this.addr.Port," [RPC] invoked of RPCFindIDSuccessor AND Next:",reply.Addr.Port," sha1:",reply.Addr.Id.ValPtr)
		//this.addr.Validate(false,"RPCFindIDSuccessor false branch 1")
		//reply.Addr.Validate(false,"RPCFindIDSuccessor false branch 2")
		return nil
	}
}

func (this *ChordNode) RPCNotify(caller Address)error  {
	this.notifyLock.Lock()
	defer this.notifyLock.Unlock()
	log.Trace(this.addr.Port," [RPC] invoked of RPCNotify(", caller.Addr,")")
	this.predecessorLock.Lock()
	defer this.predecessorLock.Unlock()
	if this.nodePredecessor.isNil() || caller.Id.In(&this.nodePredecessor.Id,&this.addr.Id){
		log.Debug(this.addr.Port,": postPre:",this.nodePredecessor.Port," new:", caller)
		this.nodePredecessor.CopyFrom(&caller)
	}
	caller.Validate(false,"RPCNotify")
	return nil
}

func (this *ChordNode)RPCCopyList(reply *[ALTERNATIVE_SIZE]Address)error  {
	log.Trace(this.addr.Port," [RPC] invoked of RPCCopyList()")
	this.alternativeListLock.RLock()
	for i:=1;i<ALTERNATIVE_SIZE;i++{
		reply[i].CopyFrom(&this.alternativeSuccessors[i-1])
		reply[i].Validate(false,this.addr.Port)
		this.alternativeSuccessors[i-1].Validate(false,this.addr.Port)
	}
	//copy(reply[1:],this.alternativeSuccessors[:ALTERNATIVE_SIZE-1])
	this.alternativeListLock.RUnlock()
	reply[0].CopyFrom(this.nodeSuccessor)
	reply[0].Validate(true,this.addr.Port)
	this.addr.Validate(true,this.addr.Port)
	this.MayFatal()
	return nil
}

func (this *ChordNode)RPCPut(kv KVPair,stat *bool)(err error) {
	log.Trace(this.addr.Port, " [RPC] invoked of Put(", kv.Key, ",", kv.Value, ")")
	*stat = this.storage.Put(kv.Key, kv.Value)
	var ignore bool
	if !this.nodePredecessor.isNil() {
		if errn := RemoteCall(this.nodePredecessor, "RPCWrapper.BackupPut", kv, &ignore); errn != nil {
			log.Warning(this.addr.Port, " Put backup failed:", err)
		}
	} else {
		log.Warning(this.addr.Port, " Cannot backup put:", kv.Key)
	}
	return
}

func (this *ChordNode)RPCGet(key string,reply *StringWithBoolean)error{
	log.Trace(this.addr.Port," [RPC] invoked of Get(",key,")")
	reply.Stat=this.storage.Get(key,&reply.Str)
	return nil
}

func (this *ChordNode)RPCDelete(key string,stat *bool)(err error){
	log.Trace(this.addr.Port," [RPC] invoked of Delete(",key,")")
	*stat=this.storage.Delete(key)
	var ignore bool
	if !this.nodePredecessor.isNil() {
		if errn := RemoteCall(this.nodePredecessor, "RPCWrapper.BackupDelete", key, &ignore); errn != nil {
			log.Warning(this.addr.Port, " Delete backup failed:", err)
		}
	} else {
		log.Warning(this.addr.Port, " Cannot backup delete:", key)
	}
	return nil
}

func (this *ChordNode)RPCBackupPut(kv KVPair,stat *bool)error{
	log.Trace(this.addr.Port," [RPC] invoked of BackupPut(",kv.Key,",",kv.Value,")")
	*stat=this.succStorageBackup.Put(kv.Key,kv.Value)
	return nil
}

func (this *ChordNode)RPCBackupDelete(key string,stat *bool)error {
	log.Trace(this.addr.Port," [RPC] invoked of BackupDelete(",key,")")
	*stat = this.succStorageBackup.Delete(key)
	return nil
}

func (this *ChordNode)RPCCopyFingers(_ int,reply *[BIT_WIDTH]Address)error{
	// todo
	return nil
}

func (this *ChordNode)RPCUpdatePredecessor(addr [2]Address,_ *int)error  {
	// {caller,caller.predecessor}
	this.predecessorLock.Lock()
	defer this.predecessorLock.Unlock()
	if addr[0].Addr==this.nodePredecessor.Addr{
		log.Debug(this.addr.Port,": postPre:",this.nodePredecessor.Port," new:", addr[1])
		this.nodePredecessor.CopyFrom(&addr[1])
	}
	return nil
}

func (this *ChordNode)RPCUpdateSuccessor(addr [2]Address,_ *int)error {
	// {caller,caller.predecessor}
	this.fingerAndSuccessorLock.Lock()
	defer this.fingerAndSuccessorLock.Unlock()
	if addr[0].Addr == this.nodeSuccessor.Addr {
		log.Debug(this.addr.Port, ": postSuc:", this.nodeSuccessor.Port, " new:", addr[1])
		this.nodeSuccessor.CopyFrom(&addr[1])
	}
	return nil
}