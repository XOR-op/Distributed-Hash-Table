package kademlia

type KBucket struct {
	indicator bucketNode
	size      int
	maxSize   int
}

func (this *KBucket) Add(addr *Contact) {
	if this.size<this.maxSize{
		n:=newBucketNode(addr)
		n.attachAfter(this.VirtualNode())
	}else {
		n:=this.Tail().Detach()
		if n.element.TestConn(){
			n.attachAfter(this.VirtualNode())
		}else {
			n=newBucketNode(addr)
			n.attachAfter(this.VirtualNode())
		}
	}
}

func NewKBucket(maxSize int) (reply *KBucket) {
	reply = new(KBucket)
	reply.indicator = bucketNode{nil, nil, nil}
	reply.indicator.prev = &reply.indicator
	reply.indicator.next = &reply.indicator
	reply.size = 0
	reply.maxSize = maxSize
	return
}

func (this *KBucket) Head() *bucketNode {
	return this.indicator.next
}

func (this *KBucket) Tail() *bucketNode {
	return this.indicator.prev
}

func (this *KBucket) VirtualNode() *bucketNode {
	return &this.indicator
}

type bucketNode struct {
	next    *bucketNode
	prev    *bucketNode
	element *Contact
}

func newBucketNode(val *Contact) (reply *bucketNode) {
	reply=new(bucketNode)
	reply.next = nil
	reply.prev = nil
	reply.element = val.Duplicate()
	return
}

func (this *bucketNode) attachAfter(n *bucketNode) {
	this.next = n.next
	this.prev = n
	n.next.prev = this
	n.next = this
}

func (this *bucketNode) Detach()*bucketNode {
	this.next.prev = this.prev
	this.prev.next = this.next
	this.next = nil
	this.prev = nil
	return this
}


