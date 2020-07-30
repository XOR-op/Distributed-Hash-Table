package chord

import (
	log "github.com/sirupsen/logrus"
	"net/rpc"
	"runtime"
	"time"
)

func Must(err error) bool {
	if err != nil {
		pc, _, _, _ := runtime.Caller(1)
		log.Warning(runtime.FuncForPC(pc).Name(), ":", err.Error())
		panic(err)
	}
	return true
}

type TemporaryError interface {
	Temporary() bool

}

func RemoteCall(addr Address, method string, arg, ret interface{}) (err error) {
	defer func() {
		if t := recover(); t != nil {
			pc, _, _, _ := runtime.Caller(3)
			log.Warning("[ERROR] RemoteCall ", method,
				" from ", runtime.FuncForPC(pc).Name(), " to ", addr.Addr, " fail:", t)
			err = t.(error)
		}
	}()
	//addr.Validate(false,121)
	client, err := Dial("tcp", addr.Addr)
	//addr.Validate(false,122)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	if err = client.Call(method, arg, ret); err != nil {
		panic(err)
	}
	log.Trace(method)
	//addr.Validate(false,123)
	return
}

func Dial(method, addr string) (*rpc.Client, error) {
	client, err := rpc.Dial(method, addr)
	for retry_cnt := 1; err != nil && retry_cnt <= 4; retry_cnt++ {
		// avoid "resource temporarily unavailable" error
		if terr, ok := err.(TemporaryError); ok && terr.Temporary() {
			pc, _, _, _ := runtime.Caller(1)
			log.Warning("Dial ", method,
				" from ", runtime.FuncForPC(pc).Name(), " to ", addr, " resource unavailable")
			time.Sleep(time.Duration(retry_cnt) * 50 * time.Millisecond)
			client, err = rpc.Dial("tcp", addr)
		} else {
			break
		}
	}
	return client, err
}

func (this *ChordNode)RecoverErr(errPtr *error) {
	if t := recover(); t != nil {
		*errPtr = t.(error)
		log.Warn(this.addr.Port, " Catch panic:", *errPtr)
	}else {
		*errPtr=nil
	}
}
