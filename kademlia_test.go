package main

import (
	"DHT/src/kademlia"
	"time"
)

const portStart = 13301
func main()  {
	kademlia.DefaultInitialize()
	defer kademlia.DefaultClose()
	arr:=make([]*kademlia.Node,5)
	for i:=0;i<5;i++{
		arr[i]=kademlia.NewNode(portStart+i)
		if i!=0{
			arr[i].Join(portStart)
		}
	}
	time.Sleep(5*time.Second)
}
