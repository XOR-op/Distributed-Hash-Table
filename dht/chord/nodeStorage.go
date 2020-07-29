package chord

import log "github.com/sirupsen/logrus"

func (this *ChordNode)Put(key,value string)(ret bool ) {
	var addr Address
	Must(this.FindIdSuccessor(IDlize(key),&addr))
	if err:=RemoteCall(addr,"RPCWrapper.Put", KVPair{key,value},&ret);err!=nil{
		log.Warning(this.addr.Port," Put pair {",key,":",value,"} failed")
		ret=false
	}
	return
}

func (this *ChordNode)Get(key string)(ret bool,value string)  {
	var addr Address
	var reply StringWithBoolean
	Must(this.FindIdSuccessor(IDlize(key),&addr))
	if err:=RemoteCall(addr,"RPCWrapper.Get",key,&reply);err!=nil{
		log.Warning(this.addr.Port," Get Key{",key,"} failed")
		return false,"ERROR OCCURRED"
	}
	return reply.Stat,reply.Str
}

func (this *ChordNode)Delete(key string)(ret bool) {
	var addr Address
	Must(this.FindIdSuccessor(IDlize(key),&addr))
	if err:=RemoteCall(addr,"RPCWrapper.Delete",key,&ret);err!=nil{
		log.Warning(this.addr.Port," Delete Key{",key,"} failed")
		ret=false
	}
	return
}

func (this *ChordNode)SplitBy(caller Address,reply *map[string]string)error  {
	log.Trace(this.addr.Port," [RPC] invoked of SplitBy()")
	*reply=this.storage.SplitBy(&caller.Id,&this.addr.Id)
	return nil
}