package chord

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net"
	"net/rpc"
	"runtime"
	"strconv"
	"sync"
)

const (
	BIT_WIDTH        int = 160
	ALTERNATIVE_SIZE int = 5
	UPDATE_INTERVAL  int = 200
)

type ChordNode struct {
	addr                   Address
	finger                 [BIT_WIDTH]Address
	server                 *rpc.Server
	lis                    net.Listener
	quitRPC                chan bool
	DaemonContext          context.Context
	quitDaemonInvoker      context.CancelFunc
	nodeSuccessor          *Address
	nodePredecessor        Address
	data                   Data
	fingerUpdateIndex      int
	alternativeSuccessors  [ALTERNATIVE_SIZE]Address
	fingerAndSuccessorLock sync.RWMutex
	predecessorLock        sync.RWMutex
	alternativeListLock    sync.RWMutex
	dataLock               sync.Mutex
	validateSuccessorLock  sync.Mutex
	notifyLock             sync.Mutex
	joinLock               sync.Mutex
}

func (this *ChordNode) Init(port int) {
	this.addr.Addr = "localhost:" + strconv.Itoa(port)
	this.addr.Port = port
	this.addr.Id = IDlize(this.addr.Addr)
	this.nodeSuccessor = &this.finger[0]
	for i, _ := range this.finger {
		this.finger[i].Nullify()
	}
	for i, _ := range this.alternativeSuccessors {
		this.alternativeSuccessors[i].Nullify()
	}
}

func (this *ChordNode) validateSuccessor(ignoreCurrent bool) {
	// guarantee alternativeSuccessor valid
	this.validateSuccessorLock.Lock()
	defer this.validateSuccessorLock.Unlock()
	this.MayFatal()
	log.Trace(this.addr.Port, " is validating successor")
	if !ignoreCurrent {
		this.fingerAndSuccessorLock.RLock()
		if this.Ping(this.nodeSuccessor.Addr) {
			this.fingerAndSuccessorLock.RUnlock()
			log.Trace(this.addr.Port, " validating successor passed")
			return
		}
		this.fingerAndSuccessorLock.RUnlock()
	}
	log.Trace(this.addr.Port, " validateSuccessor another branch")
	ref := this.nodeSuccessor.Addr
	var addr Address
	this.alternativeListLock.RLock()
	for _, addr = range this.alternativeSuccessors {
		if addr.Addr == ref {
			continue
		} else {
			ref = addr.Addr
		}
		addr.Validate(true, this.addr.Port)
		log.Trace(this.addr.Port, " go")
		if this.Ping(addr.Addr) {
			goto done
		}
	}
	this.alternativeListLock.RUnlock()
	log.Fatal(this.addr.Port, " NO VALID SUCCESSOR!!!",this.alternativeSuccessors)
	return
	//panic(errors.New("NO VALID SUCCESSOR!!!"))
done:
	log.Trace(this.addr.Port, " may done")
	this.MayFatal()
	for _, v := range this.alternativeSuccessors {
		v.Validate(false, this.addr.Port)
	}
	this.alternativeListLock.RUnlock()
	log.Trace(this.addr.Port, " alternative passed")
	log.Debug(this.addr.Port, " validate set successor:", addr.Port)
	this.fingerAndSuccessorLock.Lock()
	this.nodeSuccessor.CopyFrom(&addr)
	this.fingerAndSuccessorLock.Unlock()
	log.Trace(this.addr.Port, " copy done")
	this.MayFatal()
	// now alternativeSuccessor can be modified due to no use of addr
	log.Trace(this.addr.Port, " alternative copy start")
	pc, _, _, _ := runtime.Caller(1)
	callerName := runtime.FuncForPC(pc).Name()
	log.Trace(this.addr.Port, " father ", callerName)
	this.UpdateAlternativeSuccessor()
	log.Trace(this.addr.Port, " alternative copy done")
	this.MayFatal()
	log.Trace(this.addr.Port, " validating return")
	return
}

func (this *ChordNode) closestPrecedingNode(id Identifier) *Address {
	// todo also utilize alternativeSuccessors
	defer log.Trace(this.addr.Port, " Exit find finger")
	log.Trace(this.addr.Port, " Enter find finger")
	this.fingerAndSuccessorLock.RLock()
	defer this.fingerAndSuccessorLock.RUnlock()
	hasTested:=make(map[string]bool)
	for i := BIT_WIDTH - 1; i >= 0; i -= 1 {
		if !this.finger[i].isNil() {
			if _,ok:=hasTested[this.finger[i].Addr];!ok {
				if this.finger[i].Id.In(&this.addr.Id, &id) && this.Ping(this.finger[i].Addr) {
					return &this.finger[i]
				}else {
					hasTested[this.finger[i].Addr]=true
				}
			}
		}
	}
	return &this.finger[0]
}

func (this *ChordNode) FindIdSuccessor(id Identifier, reply *Address) (err error) {
	defer func() {
		if t := recover(); t != nil {
			err = t.(error)
			log.Warning(this.addr.Port, " ", err)
			//log.Fatal(err)
		}
	}()
	this.MayFatal()
	if id.InRightClosure(&this.addr.Id, &this.nodeSuccessor.Id) {
		log.Trace(GOid(), this.addr.Port, " First branch")
		this.validateSuccessor(false)
		this.fingerAndSuccessorLock.RLock()
		log.Trace(GOid(), this.addr.Port, " locally get ID ", id, "'s successor:", this.nodeSuccessor.Port)
		reply.CopyFrom(this.nodeSuccessor)
		this.fingerAndSuccessorLock.RUnlock()
		return nil
	} else {
		log.Trace(GOid(), this.addr.Port, " Second branch")
		stru := NewAddressWithBoolean(this.closestPrecedingNode(id), false)
		var backup Address
		backup.Nullify()
		backup.CopyFrom(this.nodeSuccessor)
		for !stru.Stat {
			//todo retry
			log.Trace(this.addr.Port, " Sub loop")
			// copy for debug only
			var tmpAddr Address
			tmpAddr.CopyFrom(&stru.Addr)
			if cerr := RemoteCall(tmpAddr, "RPCWrapper.FindIDSuccessor", id, &stru); cerr != nil {
				if stru.Addr.Addr==backup.Addr{
					this.validateSuccessor(false)
					backup.CopyFrom(this.nodeSuccessor)
				}
				log.Warning(this.addr.Port," fail ",stru.Addr.Port,", has to retry ",backup.Port)
				stru.Addr.CopyFrom(&backup)
			}
			log.Trace(GOid(), "cur stru.addr:", stru.Addr.Port)
		}
		log.Trace(GOid(), this.addr.Port, " remotely get ID ", id, "'s successor:", stru.Addr.Port)
		reply.CopyFrom(&stru.Addr)
		this.MayFatal()
		return nil
	}
}

func (this *ChordNode) Stabilize() (err error) {
	defer func() {
		if t := recover(); t != nil {
			err = t.(error)
			log.Warning(this.addr.Port, " ", err)
		}
	}()
	defer log.Trace(this.addr.Port," stabilize done")
	log.Trace(this.addr.Port, " stabilize. Cur suc:", this.nodeSuccessor.Port)
	this.MayFatal()
	client, err := Dial("tcp", this.nodeSuccessor.Addr)
	for err != nil {
		// retry and reconnect
		this.validateSuccessor(false)
		client, err = Dial("tcp", this.nodeSuccessor.Addr)
	}
	defer client.Close()
	var x Address
	if err := client.Call("RPCWrapper.Predecessor", 0, &x); err != nil {
		log.Warning(this.addr.Port," find successor.predecessor failed")
		this.validateSuccessor(false)
		return nil // ignore explicitly
	}
	if !x.isNil() && x.Id.In(&this.addr.Id, &this.nodeSuccessor.Id) {
		xclient, erra := Dial("tcp", x.Addr)
		if erra == nil {
			defer xclient.Close()
			this.fingerAndSuccessorLock.Lock()
			log.Debug(this.addr.Port, " stabilize original:", this.nodeSuccessor.Port, " set successor:", x.Port)
			this.nodeSuccessor.CopyFrom(&x)
			this.fingerAndSuccessorLock.Unlock()
			Must(xclient.Call("RPCWrapper.Notify", this.addr, nil))
			//x.Validate(false, this.addr.Port)
			this.UpdateAlternativeSuccessor()
			x.Validate(false, this.addr.Port)
			return erra
		}
	}
	Must(client.Call("RPCWrapper.Notify", this.addr, nil))
	return

}

func (this *ChordNode) FixFingers() error {
	log.Trace(this.addr.Port, " Fix fingers:", this.fingerUpdateIndex)
	this.fingerUpdateIndex = (this.fingerUpdateIndex + 1) % BIT_WIDTH
	return this.FindIdSuccessor(this.addr.Id.PlusTwoPower(uint(this.fingerUpdateIndex)), &this.finger[this.fingerUpdateIndex])
}

func (this *ChordNode) CheckPredecessor() {
	log.Trace(GOid(), this.addr.Port, " check predecessor. Cur pre:", this.nodePredecessor.Port)
	this.predecessorLock.Lock()
	defer this.predecessorLock.Unlock()
	if !this.Ping(this.nodePredecessor.Addr) {
		if !this.nodePredecessor.isNil() {
			log.Trace(this.addr.Port, " clear predecessor")
		}
		this.nodePredecessor.Nullify()
	}
}

func (this *ChordNode) UpdateAlternativeSuccessor() {
	log.Trace(GOid(), this.addr.Port, " check alternative successor")
	var warehouse [ALTERNATIVE_SIZE]Address

	if this.addr.Addr != this.nodeSuccessor.Addr {
		if RemoteCall(*this.nodeSuccessor, "RPCWrapper.CopyList", 0, &warehouse)!=nil{
			return
		}
	}
	this.alternativeListLock.Lock()
	defer this.alternativeListLock.Unlock()
	for i:=0;i<ALTERNATIVE_SIZE;i++{
		this.alternativeSuccessors[i].CopyFrom(&warehouse[i])
	}
}

func (this *ChordNode) Join(addr Address) (err error) {
	defer func() {
		if t := recover(); t != nil {
			log.Warning(t.(error))
			panic(t)
		}
	}()
	this.joinLock.Lock()
	defer this.joinLock.Unlock()
	log.Trace(this.addr.Port, " Joined from ", addr.Port)
	this.nodePredecessor.Nullify()
	defer func() {
		if t := recover(); t != nil {
			err = t.(error)
		}
	}()
	// find successor by addr
	client, err := Dial("tcp", addr.Addr)
	Must(err)
	Must(client.Call("RPCWrapper.LocalFindIDSuccessor", this.addr.Id, this.nodeSuccessor))
	_ = client.Close()
	// call to successor
	client, err = Dial("tcp", this.nodeSuccessor.Addr)
	Must(err)
	Must(client.Call("RPCWrapper.CopyList", 0, &this.alternativeSuccessors))
	Must(client.Call("RPCWrapper.Notify", this.addr, nil))
	log.Debug(this.addr.Port, " after join. Suc:", this.nodeSuccessor.Port)
	return
}
