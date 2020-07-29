package main

import (
	log "github.com/sirupsen/logrus"
	easy_formatter "github.com/t-tomalak/logrus-easy-formatter"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	NodeSize = 70
	DataSize = 3000
	NodeLeft = 70
)

func NodeAnswerDumpAll(node *[NodeSize]dhtNode) {
	for i, n := range node {
		if i < NodeLeft {
			n.AnswerDump()
		}
	}
}

func DumpNodeAll(node *[NodeSize]dhtNode) {
	log.Println("**********DUMP BEGIN************")
	for _, n := range node {
		n.Dump(2)
	}
	log.Println("***********DUMP END*************")
}

func GetData() map[string]string {
	m := make(map[string]string)
	for i := 0; i < DataSize; i++ {
		m[strconv.Itoa(i)] = "<#" + strconv.Itoa(i) + "#>abc"
	}
	return m
}

func BigProcedure() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("A Panic", err)
			panic(err)
		}
	}()
	log.SetFormatter(&easy_formatter.Formatter{
			TimestampFormat: "2006-01-02 15:04:05.000",
			LogFormat:       "[%lvl%]: %time% - %msg%\n",
		})
	log.SetLevel(log.InfoLevel)
	//log.SetLevel(log.TraceLevel)
	//log.SetLevel(log.DebugLevel)
	if true {
		//if false{
		f, err := os.OpenFile("log/log_"+strings.ReplaceAll(time.Now().Format(time.Stamp), " ", "-")+".log", os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		log.SetOutput(f)
	}
	log.Info("Current Program with ", NodeSize, " Node")
	//log.SetReportCaller(true)
	node := [NodeSize]dhtNode{}
	data := GetData()
	startPort := 13301
	for i, port := 0, startPort; i < NodeSize; i, port = i+1, port+1 {
		node[i] = NewNode(port)
		node[i].Run()
	}
	for i := 0; i < NodeSize/2; i++ {
		if i == 0 {
			node[i].Create()
		} else {
			node[i].Join("localhost:" + strconv.Itoa(startPort))
			time.Sleep(500 * time.Millisecond)
		}
		//DumpNodeAll(&node)
	}
	time.Sleep(2 * time.Second)

	for i := 0; i < DataSize; i++ {
		s := strconv.Itoa(i)
		if stat := node[0].Put(s, data[s]); !stat {
			log.Warning("[KEY INSERT FAILED] key:", s)
		}
		time.Sleep(5 * time.Millisecond)
	}
	time.Sleep(500 * time.Millisecond)
	for i := NodeSize / 2; i < NodeSize; i++ {
		node[i].Join("localhost:" + strconv.Itoa(startPort))
		time.Sleep(500 * time.Millisecond)
		//DumpNodeAll(&node)
	}

	time.Sleep(2 * time.Second)
	log.Info("=== All nodes joined ===")
	for i := 0; i < DataSize; i++ {
		s := strconv.Itoa(i)
		stat, ans := node[0].Get(s)
		if !stat || ans != data[s] {
			log.Warning("[KEY FIND FAILED] ", stat, " key:", s, " want:", data[s], " wrong:", ans)
		}
	}
	for i:=0;i<NodeSize;i++{
		node[i].DataDump()
	}
	NodeAnswerDumpAll(&node)
	for i := 0; i < NodeSize; i++ {
		node[i].Quit()
	}
	time.Sleep(3 * time.Second)
	log.Info("All process finished")
}
