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
	reply.Validate(false,"RPCPredecessor")
	return nil
}

func (this *ChordNode) RPCFindIDSuccessor(id Identifier,reply *AddressWithBoolean) error{
	log.Trace(this.addr.Port," [RPC] invoked of RPCFindIDSuccessor(",id,")")
	this.MayFatal()
	this.validateSuccessor(false)
	this.fingerAndSuccessorLock.RLock()
	rep:=id.InRightClosure(&this.addr.Id,&this.nodeSuccessor.Id)
	this.fingerAndSuccessorLock.RUnlock()
	if rep{
		*reply=NewAddressWithBoolean(this.nodeSuccessor,true)
		reply.Addr.Validate(false,"RPCFindIDSuccessor true branch")
		return nil
	}else {
		*reply=NewAddressWithBoolean(this.closestPrecedingNode(id),false)
		log.Trace(this.addr.Port," [RPC] invoked of RPCFindIDSuccessor AND Next:",reply.Addr.Port," sha1:",reply.Addr.Id.ValPtr)
		this.addr.Validate(false,"RPCFindIDSuccessor false branch 1")
		reply.Addr.Validate(false,"RPCFindIDSuccessor false branch 2")
		return nil
	}
}

func (this *ChordNode) RPCNotify(callee Address)error  {
	this.notifyLock.Lock()
	defer this.notifyLock.Unlock()
	log.Trace(this.addr.Port," [RPC] invoked of RPCNotify(",callee.Addr,")")
	this.predecessorLock.Lock()
	defer this.predecessorLock.Unlock()
	if this.nodePredecessor.isNil() || callee.Id.In(&this.nodePredecessor.Id,&this.addr.Id){
		log.Debug(this.addr.Port,": postPre:",this.nodePredecessor.Port," new:",callee)
		this.nodePredecessor=callee
	}
	callee.Validate(false,"RPCNotify")
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

