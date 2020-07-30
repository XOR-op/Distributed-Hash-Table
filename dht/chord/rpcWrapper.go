package chord

type RPCWrapper struct {
	Node *ChordNode
}

func (this *RPCWrapper) Successor(_ int,reply *Address)(err error)  {
	defer this.Node.RecoverErr(&err)
	return this.Node.RPCSuccessor(reply)
}

func (this *RPCWrapper) Predecessor(_ int,reply *Address)(err error)  {
	defer this.Node.RecoverErr(&err)
	return this.Node.RPCPredecessor(reply)
}

func (this *RPCWrapper) FindIDSuccessor(id Identifier,reply *AddressWithBoolean)(err error){
	defer this.Node.RecoverErr(&err)
	err=this.Node.RPCFindIDSuccessor(id,reply)
	return
}
func (this *RPCWrapper) FindIDSuccessorWithValidation(id Identifier,reply *AddressWithBoolean)(err error) {
	defer this.Node.RecoverErr(&err)
	err = this.Node.RPCFindIDSuccessorWithValidation(id, reply)
	return
}

func (this *RPCWrapper) Notify(callee Address,_ *int) (err error)  {
	defer this.Node.RecoverErr(&err)
	return this.Node.RPCNotify(callee)
}

func (this *RPCWrapper)CopyList(_ int,reply *[ALTERNATIVE_SIZE]Address) (err error)  {
	defer this.Node.RecoverErr(&err)
	return this.Node.RPCCopyList(reply)
}
func (this *RPCWrapper)LocalFindIDSuccessor(id Identifier,reply *Address) (err error)  {
	defer this.Node.RecoverErr(&err)
	return this.Node.FindIdSuccessor(id,reply)
}

func (this *RPCWrapper)Put(kv KVPair,stat *bool) (err error)  {
	defer this.Node.RecoverErr(&err)
	return this.Node.RPCPut(kv,stat)
}

func (this *RPCWrapper)Get(key string,reply *StringWithBoolean)(err error){
	defer this.Node.RecoverErr(&err)
	return this.Node.RPCGet(key,reply)
}

func (this *RPCWrapper)Delete(key string,stat *bool)(err error){
	defer this.Node.RecoverErr(&err)
	return this.Node.RPCDelete(key,stat)
}

func (this *RPCWrapper)MoveData(caller Address,reply *map[string]string)(err error)  {
	defer this.Node.RecoverErr(&err)
	return this.Node.SplitBy(caller,reply)
}