package main

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
)

const N int =120
const LEFT int=70

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
			time.Sleep(500*time.Millisecond)
		}
		DumpAll(&node)
	}
	time.Sleep(5*time.Second)
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

func Logg(err *error)  {
	if t:=recover();t!=nil{
		fmt.Println(t.(error))
		*err=t.(error)
	}else {
		*err = nil
	}
}

func small()(err error)  {
	defer Logg(&err)
	fmt.Println("1")
	fmt.Println("2")
	panic(errors.New("Hello"))
	fmt.Println("3")
	fmt.Println("4")
	return nil
}

func MinorTest()  {
	fmt.Println(small())
}

func main()  {
	//Procedure()
	BigProcedure()
	//MinorTest()
}