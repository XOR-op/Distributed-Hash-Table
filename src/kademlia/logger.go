package kademlia

import (
	"github.com/sirupsen/logrus"
	easy_formatter "github.com/t-tomalak/logrus-easy-formatter"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Logger struct {
	file *os.File
	log  *logrus.Logger
	info string
}

const (
	defaultSkip  = 2
	//defaultLevel = logrus.TraceLevel
	//outputToFile = true
	defaultLevel = logrus.InfoLevel
	outputToFile = false
)

var DefaultLogger *Logger

func NewLogger(info string) (reply *Logger) {
	reply = new(Logger)
	reply.log = DefaultLogger.log
	reply.file = nil
	reply.info = info
	return
}

func DefaultInitialize() {
	DefaultLogger = new(Logger)
	DefaultLogger.file = nil
	if outputToFile {
		DefaultLogger.log = &logrus.Logger{}
		DefaultLogger.log.SetFormatter(&easy_formatter.Formatter{
			TimestampFormat: "2006-01-02 15:04:05.000",
			LogFormat:       "[%lvl%]: %time% - %msg%\n",
		})
		f, err := os.OpenFile("log/log_"+strings.ReplaceAll(time.Now().Format(time.Stamp), " ", "-")+".log", os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		DefaultLogger.file = f
		DefaultLogger.log.SetOutput(f)
	}else {
		DefaultLogger.log=logrus.StandardLogger()
	}
	DefaultLogger.info = "DEFAULT"
	DefaultLogger.log.SetLevel(defaultLevel)
}

func DefaultClose() {
	if outputToFile {
		_ = DefaultLogger.file.Close()
	}
}

func (this *Logger) Prefix() string {
	pc, _, line, _ := runtime.Caller(defaultSkip)
	return "[" + this.info + "] @" + runtime.FuncForPC(pc).Name() + "(line " + strconv.Itoa(line) + "): "
}

func (this *Logger) Info(args ...interface{}) {
	this.log.Info(this.Prefix(), args)
}

func (this *Logger) Warning(args ...interface{}) {
	this.log.Warning(this.Prefix(), args)
}

func (this *Logger) Fatal(args ...interface{}) {
	this.log.Fatal(this.Prefix(), args)
}

func (this *Logger) Debug(args ...interface{}) {
	this.log.Debug(this.Prefix(), args)
}

func (this *Logger) Trace(args ...interface{}) {
	this.log.Trace(this.Prefix(), args)
}
func (this *Logger) Error(args ...interface{}) {
	this.log.Error(this.Prefix(), args)
}

func GetMyName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
