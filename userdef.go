package main

/* In this file, you should implement function "NewNode" and
 * a struct which implements the interface "dhtNode".
 */
import "DHT/dht/chord"

func NewNode(port int)(reply dhtNode) {
	var node chord.Node
	node.Init(port)
	reply=&node
	return
}

// Todo: implement a struct which implements the interface "dhtNode".