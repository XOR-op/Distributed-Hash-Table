package chord

import log "github.com/sirupsen/logrus"

type RPCWrapper struct {
	Node *ChordNode
}

func (this *RPCWrapper) Successor(_ int,reply *Address)(err error)  {
	defer func() {
		if t:=recover();t!=nil{
			err=t.(error)
			log.Warn(this.Node.addr.Port," Catch panic:",err)
		}
	}()
	return this.Node.RPCSuccessor(reply)
}

func (this *RPCWrapper) Predecessor(_ int,reply *Address)(err error)  {
	defer func() {
		if t:=recover();t!=nil{
			err=t.(error)
			log.Warn(this.Node.addr.Port," Catch panic:",err)
		}
	}()
	return this.Node.RPCPredecessor(reply)
}

func (this *RPCWrapper) FindIDSuccessor(id Identifier,reply *AddressWithBoolean)(err error){
	defer func() {
		if t:=recover();t!=nil{
			err=t.(error)
			log.Warn(this.Node.addr.Port," Catch panic:",err)
		}
	}()
	err=this.Node.RPCFindIDSuccessor(id,reply)
	return
}

func (this *RPCWrapper) Notify(callee Address,_ *int) (err error)  {
	defer func() {
		if t:=recover();t!=nil{
			err=t.(error)
			log.Warn(this.Node.addr.Port," Catch panic:",err)
		}
	}()
	return this.Node.RPCNotify(callee)
}

func (this *RPCWrapper)CopyList(_ int,reply *[ALTERNATIVE_SIZE]Address) (err error)  {
	defer func() {
		if t:=recover();t!=nil{
			err=t.(error)
			log.Warn(this.Node.addr.Port," Catch panic:",err)
		}
	}()
	return this.Node.RPCCopyList(reply)
}
func (this *RPCWrapper)LocalFindIDSuccessor(id Identifier,reply *Address) (err error)  {
	defer func() {
		if t:=recover();t!=nil{
			err=t.(error)
			log.Warn(this.Node.addr.Port," Catch panic:",err)
		}
	}()
	return this.Node.FindIdSuccessor(id,reply)
}

