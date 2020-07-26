package chord

import (
	log "github.com/sirupsen/logrus"
	"net/rpc"
	"runtime"
)

func Must(err error) bool {
	if err != nil {
		pc, _, _, _ := runtime.Caller(1)
		log.Warning(runtime.FuncForPC(pc).Name(), ":", err.Error())
		panic(err)
	}
	return true
}


func RemoteCall(addr Address, method string, arg, ret interface{}) (err error) {
	defer func() {
		if t := recover(); t != nil {
			pc, _, _, _ := runtime.Caller(3)
			log.Warning("[ERROR] RemoteCall ", method,
				" from ", runtime.FuncForPC(pc).Name()," to ", addr.Addr, " fail:", t)
			err = t.(error)
		}
	}()
	addr.Validate(false,121)
	client, err := rpc.Dial("tcp", addr.Addr)
	addr.Validate(false,122)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	if err = client.Call(method, arg, ret); err != nil {
		panic(err)
	}
	log.Trace(method)
	addr.Validate(false,123)
	return
}
