package main

import (
	"DHT/dht/chord"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const N int =120
const LEFT int=50

func DumpAll(node *[N]dhtNode)  {
	log.Println("**********DUMP BEGIN************")
	for _,n:=range node {
		n.Dump(2)
	}
	log.Println("***********DUMP END*************")
}
func AnswerDumpAll(node *[N]dhtNode)  {
	for i,n:=range node {
		if i<LEFT {
			n.AnswerDump()
		}
	}
}
func Procedure()  {
	defer func() {
		if err:=recover();err!=nil{
			log.Println("A Panic",err)
			panic(err)
		}
	}()
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	//log.SetLevel(log.InfoLevel)
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&log.TextFormatter{})
	if true{
	//if false{
		f, err := os.OpenFile("log/log_"+strings.ReplaceAll(time.Now().Format(time.Stamp)," ","-")+".log", os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		log.SetOutput(f)
	}
	log.Info("Current Program with ",N, " Node")
	//log.SetReportCaller(true)
	node:=[N]dhtNode{}
	startPort:=13301
	for i,port:=0,startPort;i<N;i,port=i+1,port+1{
		node[i]=NewNode(port)
		node[i].Run()
	}
	for i,n:=range node{
		if i==0{
			n.Create()
		}else {
			n.Join("localhost:"+strconv.Itoa(startPort))
			time.Sleep(1000*time.Millisecond)
		}
		DumpAll(&node)
	}
	time.Sleep(1*time.Second)
	log.Info("=== All nodes joined ===")
	for i:=N-1;i>=LEFT;i--{
		log.Debug(i,"?")
		node[i].ForceQuit()
		log.Debug(i, " Quit")
		time.Sleep(1000*time.Millisecond)
		DumpAll(&node)
	}
	time.Sleep(10*time.Second)
	DumpAll(&node)
	log.Info("=== Partial nodes quited ===")
	AnswerDumpAll(&node)
	for i,n:=range node{
		if i<LEFT {
			n.Quit()
		}
	}
	time.Sleep(3*time.Second)
	log.Info("All process finished")
}

func MinorTest()  {
	a:=chord.NewAddress("localhost:1234")
	b:=chord.NewAddress("localhost:1235")
	//c:=chord.NewAddress("localhost:1236")
	d:=new(chord.Address)
	d.CopyFrom(&a)
	a.CopyFrom(&b)
	a.Validate(false,1)
	b.Validate(false,2)
	d.Validate(false,3)
	fmt.Println(runtime.GOMAXPROCS(6))


}

func main()  {
	Procedure()
	//MinorTest()
}