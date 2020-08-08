package chord

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"net"
	"net/rpc"
	"strconv"
	"strings"
	"time"
)

var exitSig = errors.New("EXIT")

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
			log.Debug(this.addr.Port, " RPC Normal Quit")
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
			if t := recover(); t != nil {
				if t.(error) == exitSig {
					log.Debug(this.addr.Port, " Daemon Normal Quit")
				} else {
					log.Warning(this.addr.Port, " ", t.(error))
				}
			}
			this.lis.Close()
			log.Debug(this.addr.Port, " Daemon Quited.")
		}()
		counter := 0
		for {
			counter = (counter + 1) % 4
			time.Sleep(time.Duration(UPDATE_INTERVAL) * time.Millisecond)
			select {
			case <-this.daemonContext.Done():
				log.Debug(this.addr.Port, " Daemon Normal Quit")
				return
			default:
				log.Debug(this.addr.Port, " Daemon invoked")
				this.CheckPredecessor()
				this.MayFatal()
				Must(this.Stabilize())
				this.MayFatal()
				Must(this.FixFingers())
				this.MayFatal()
				if counter == 0 {
					this.UpdateAlternativeSuccessor()
				}
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
	for i := 0; i < 3; i++ {
		chanArrived := make(chan error)
		go func() {
			client, err := Dial("tcp", addr)
			if err == nil {
				client.Close()
			}
			chanArrived <- err
		}()
		select {
		case <-this.daemonContext.Done():
			log.Debug(this.addr.Port, " Ping exit")
			panic(exitSig)
			return false
		case err := <-chanArrived:
			if err == nil {
				return true
			} else {
				if terr, ok := err.(TemporaryError); ok && terr.Temporary() {
					time.Sleep(time.Duration(i+1) * 30 * time.Millisecond)
					break // continue
				}
				if !strings.Contains(err.Error(),"refused") {
					// dirty handling
					log.Warning(this.addr.Port, " ping with err:", err)
				}
				return false
			}
		case <-time.After(time.Duration(300+(i*500)) * time.Millisecond):
			break
		}
	}
	log.Warning(this.addr.Port, " ping ", addr, " timeout FAILED")
	return false
}

func (this *ChordNode) Run() {
	this.server = rpc.NewServer()
	_ = this.server.Register(&RPCWrapper{this})
	l, err := net.Listen("tcp", ":"+strconv.Itoa(this.addr.Port))
	if err != nil {
		log.Fatal("Listen error in", this.addr.Addr,":",err)
		panic(err)
	}
	this.lis = l
	this.quitRPC = make(chan bool, 1)
	this.daemonContext, this.quitDaemonInvoker = context.WithCancel(context.Background())
	go this.RPCServe(l)
}

func (this *ChordNode) ForceQuit() {
	// rpc server should be down
	if !this.running{
		return
	}
	this.running=false
	log.Debug(this.addr.Port, " try quit")
	this.quitRPC <- true
	log.Debug(this.addr.Port, " quit half")
	this.quitDaemonInvoker()
	log.Info(this.addr.Port, " has been forced quited")
}

func (this *ChordNode) Quit() {
	this.ForceQuit()
}
