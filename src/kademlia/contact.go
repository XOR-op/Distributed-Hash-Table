package kademlia

type Contact struct {
	address string
	port    int
	id      Identifier
}

func (this *Contact)Duplicate()(reply *Contact)  {
	reply=new(Contact)
	reply.address=this.address
	reply.port=this.port
	reply.id.As(&this.id)
	return
}

func (this *Contact)TestConn() bool {
	return Ping(this)
}


