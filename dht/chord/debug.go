package chord

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

func (this *ChordNode) Dump(verbose int) {
	//s := make([]int, 10)
	//for i, _ := range s {
	//	s[i] = this.finger[i].Port
	//}
	log.WithFields(log.Fields{
		"Addr":        this.addr.Port,
		"Successor":   this.nodeSuccessor.Port,
		"Predecessor": this.nodePredecessor.Port,
		//"Finger":      s,
	}).Info("[DUMP]\n")
	//if this.nodeSuccessor.Port!=0 {
	//	this.MayFatal()
	//}
	//this.addr.Validate(false,"self")
	//this.nodeSuccessor.Validate(false,"succ")
	//this.nodePredecessor.Validate(false,"pred")
}
func (this *ChordNode) AnswerDump() {
	//s := make([]int, 10)
	//for i, _ := range s {
	//	s[i] = this.finger[i].Port
	//}
	fmt.Printf("Addr=%d Predecessor=%d Successor=%d\n",
		this.addr.Port,this.nodePredecessor.Port,this.nodeSuccessor.Port)
}

func (this *ChordNode)MayFatal()  {
	pc, _, _, _ := runtime.Caller(1)
	callerName:=runtime.FuncForPC(pc).Name()
	if this.nodeSuccessor.isNil(){
		log.Warning(this.addr.Port," successor is nil!")
	}
	this.fingerAndSuccessorLock.RLock()
	if reSha1:=IDlize(this.nodeSuccessor.Addr);reSha1.ValPtr.Cmp(this.nodeSuccessor.Id.ValPtr)!=0{
		log.Fatal(callerName," this:",this.addr.Port," Succ:",this.nodeSuccessor.Port," Correct:",reSha1," Wrong:",this.nodeSuccessor.Id)
		panic(errors.New("WRONG1"))
	}
	log.Debug(this.addr.Port," Successor:",this.nodeSuccessor.Port, " sha1 check passed")
	/*
	for i,x:=range this.finger{
		if !x.isNil() {
			whoami:=strconv.Itoa(this.addr.Port)+"^^^"
			if x.isNil(){
				return
			}
			if reSha1:=IDlize(x.Addr);reSha1.ValPtr.Cmp(x.Id.ValPtr)!=0{
				log.Fatal(callerName, " from ",whoami,":this FINGER ",i,":",x.Port," Correct:",reSha1," Wrong:",x.Id)
			}
		}
	}

	 */
	this.fingerAndSuccessorLock.RUnlock()
	this.predecessorLock.RLock()
	defer this.predecessorLock.RUnlock()
	if this.nodePredecessor.isNil(){
		return
	}
	if reSha1:=IDlize(this.nodePredecessor.Addr);reSha1.ValPtr.Cmp(this.nodePredecessor.Id.ValPtr)!=0{
		log.Fatal(callerName," this:",this.addr.Port," Prede:",this.nodePredecessor.Port," Correct:",reSha1," Wrong:",this.nodePredecessor.Id)
		panic(errors.New("WRONG2"))
	}
	log.Debug(this.addr.Port," Predecessor:",this.nodePredecessor.Port, " sha1 check passed")
}


func (this *Address)Validate(willLog bool,whoami interface{})  {
	if this.isNil(){
		return
	}
	if reSha1:=IDlize(this.Addr);reSha1.ValPtr.Cmp(this.Id.ValPtr)!=0{
		pc, _, _, _ := runtime.Caller(1)
		log.Fatal(runtime.FuncForPC(pc).Name(), " from ",whoami,":this Address ",this.Port," Correct:",reSha1," Wrong:",this.Id)
	}else {
		if willLog {
			log.Trace("from ",whoami," this Address ", this.Port, " sha1 check passed")
		}
	}
}

func (this *ChordNode)DataDump(){
	this.storage.lock.Lock()
	defer this.storage.lock.Unlock()
	log.Info("From {",this.addr.Port,"}'s data:")
	for k,v:=range this.storage.Storage {
		log.Info("Key:",k," ,value:",v)
	}
	log.Info("{",this.addr.Port,"} done")
}



func GOid() string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	//id, err := strconv.Atoi(idField)
	//if err != nil {
	//	panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	//}
	return "[threadid:"+idField+"] "
}