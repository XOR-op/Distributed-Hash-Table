package main

import (
	"DHT/src/kademlia"
	"strconv"
	"time"
)

const (
	portStart = 13301
	kadNodeN=25
	kadDataN=80
)

func gen(i int)(string,string)  {
	return strconv.Itoa(i),"@"+strconv.Itoa(i)+"#"
}
func main()  {
	kademlia.DefaultInitialize()
	defer kademlia.DefaultClose()
	arr:=make([]*kademlia.Node,kadNodeN)
	for i:=0;i<kadNodeN;i++{
		arr[i]=kademlia.NewNode(portStart+i)
		if i!=0{
			arr[i].Join(portStart)
			time.Sleep(time.Second)
			kademlia.DefaultLogger.Info(portStart+i,"Joined")
		}
	}
	kademlia.DefaultLogger.Info("Join done")
	nodeCur:=7
	for i:=0;i<kadDataN;i++{
		nodeCur=(nodeCur+i^2)%10
		arr[nodeCur].Store(gen(i))
		kademlia.DefaultLogger.Debug("store",i)
	}
	kademlia.DefaultLogger.Info("Store done")
	time.Sleep(2*time.Second)
	for i:=0;i<kadDataN;i++{
		k,v:=gen(i)
		nodeCur=(nodeCur+i^2)%kadNodeN
		for !arr[nodeCur].Status(){
			nodeCur=(nodeCur+i^2)%kadNodeN
		}
		val,ok:=arr[nodeCur].Get(k)
		if !ok{
			kademlia.DefaultLogger.Error("Missing:",k)
		}else {
			if v!=val{
				kademlia.DefaultLogger.Error("not compatible:[",k,":",v,"],wrong:",val)
			}
		}
		time.Sleep(50*time.Millisecond)
	}
	for i:=10;i<15;i++{
		arr[i].Quit()
		time.Sleep(200*time.Millisecond)
	}
	kademlia.DefaultLogger.Info("Get after quit start")
	time.Sleep(20*time.Second)
	for i:=kadDataN-1;i>=0;i--{
		k,v:=gen(i)
		nodeCur=(nodeCur+i^2)%kadNodeN
		for !arr[nodeCur].Status(){
			nodeCur=(nodeCur+i^2)%kadNodeN
		}
		val,ok:=arr[nodeCur].Get(k)
		if !ok{
			kademlia.DefaultLogger.Error("Missing:",k)
		}else {
			if v!=val{
				kademlia.DefaultLogger.Error("not compatible:[",k,":",v,"],wrong:",val)
			}
		}
		time.Sleep(50*time.Millisecond)
	}
	kademlia.DefaultLogger.Info("Get done")
	for i:=0;i<kadNodeN;i++{
		arr[i].Quit()
	}
	kademlia.DefaultLogger.Info("Done")
}
