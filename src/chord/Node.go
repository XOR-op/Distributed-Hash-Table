package chord

import (
	log "github.com/sirupsen/logrus"
)

type Node struct {
	N *ChordNode
}

func (this *Node)Init(port int)  {
	this.N =new(ChordNode)
	this.N.Init(port)
}

func (this Node) Run() {
	this.N.Run()
}

func (this Node) Create() {
	this.N.Create()
	this.N.RunDaemon()
}

func (this Node) Join(addr string)bool {
	if err:=this.N.Join(NewAddress(addr));err!=nil{
		log.Fatal(err)
	}
	this.N.RunDaemon()
	return true
}

func (this Node) Quit() {
	// todo fix data

	this.N.Quit()
}

func (this Node) ForceQuit() {
	this.N.ForceQuit()
}

func (this Node)Ping(addr string)(ok bool)  {
	defer func() {
		if t:=recover();t!=nil{
			if t.(error)==exitSig{
				ok=false
			}else {
				panic(t)
			}
		}
	}()
	return this.N.Ping(addr)
}


func (this Node) Put(key, value string)(ok bool) {
	defer func() {
		if t:=recover();t!=nil{
			if t.(error)==exitSig{
				ok=false
			}else {
				panic(t)
			}
		}
	}()
	return this.N.Put(key,value)
}

func (this Node) Get(key string) (ok bool, val string) {
	defer func() {
		if t:=recover();t!=nil{
			if t.(error)==exitSig{
				ok=false
			}else {
				panic(t)
			}
		}
	}()
	return this.N.Get(key)
}

func (this Node) Delete(key string)(ok bool) {
	defer func() {
		if t:=recover();t!=nil{
			if t.(error)==exitSig{
				ok=false
			}else {
				panic(t)
			}
		}
	}()
	return this.N.Delete(key)
}

func (this Node) Dump(verbose int) {
	this.N.Dump(verbose)
}
func (this Node) AnswerDump() {
	this.N.AnswerDump()
}

func (this Node) DataDump() {
	this.N.DataDump()
}
