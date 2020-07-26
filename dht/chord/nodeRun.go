package chord

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net"
	"net/rpc"
	"strconv"
	"time"
)

func (this *ChordNode) Create() {
	this.nodePredecessor.Nullify()
	this.nodeSuccessor.CopyFrom(&this.addr)
	for i, _ := range this.alternativeSuccessors {
		this.alternativeSuccessors[i].CopyFrom(&this.addr)
	}
}

func (this *ChordNode) RPCServe(l net.Listener) {
	defer l.Close()
	defer log.Debug(this.addr.Port, " RPC server Quited.")
	for {
		conn, err := l.Accept()
		select {
		case <-this.quitRPC:
			return
		default:
			if err != nil {
				log.Warning("rpc.Serve: accept:", err.Error())
				return
			}
			go this.server.ServeConn(conn)
		}
	}
}

func (this *ChordNode) RunDaemon() {
	go func() {
		defer func() {
			recover()
			this.lis.Close()
			log.Debug(this.addr.Port, " Daemon Quited.")
		}()
		for {
			time.Sleep(time.Duration(UPDATE_INTERVAL) * time.Millisecond)
			select {
			case <-this.DaemonContext.Done():
				return
			default:
				log.Debug(this.addr.Port, " Daemon invoked")
				this.CheckPredecessor()
				Must(this.Stabilize())
				Must(this.FixFingers())
				log.Debug(this.addr.Port, " Daemon Slept")
				//this.Dump(2)
			}
		}
	}()
}

func (this *ChordNode) Ping(addr string) bool {
	if addr == "" {
		return false
	}
	log.Trace(this.addr.Port, " ping ", addr)
	for i := 0; i < 2; i++ {
		chanArrived:=make(chan bool)
		go func() {
			client, err := net.Dial("tcp", addr)
			if err==nil{
				client.Close()
			}
			chanArrived<-err==nil
		}()
		select {
		case <-this.DaemonContext.Done():
			panic("EXIT")
			return false
		case ok:=<-chanArrived:
			if ok{
				return true
			}
		case <-time.After(300*time.Millisecond):
			break
		}
	}
	log.Warning(this.addr.Port, " ping ", addr, " FAILED")
	return false
}

func (this *ChordNode) Run() {
	this.server = rpc.NewServer()
	_ = this.server.Register(&RPCWrapper{this})
	l, err := net.Listen("tcp", ":"+strconv.Itoa(this.addr.Port))
	if err != nil {
		log.Println("Listen error in", this.addr.Addr)
		panic(err)
	}
	this.lis = l
	this.quitRPC = make(chan bool, 1)
	this.DaemonContext,this.quitDaemonInvoker = context.WithCancel(context.Background())
	go this.RPCServe(l)
}

func (this *ChordNode) ForceQuit() {
	// rpc server should be down
	this.quitRPC <- true
	this.quitDaemonInvoker()
}

func (this *ChordNode) Quit() {
	this.ForceQuit()
}
