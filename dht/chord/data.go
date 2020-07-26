package chord

type KVPair struct {
	k, v string
	kID  Identifier
}

type Data struct {
	dat map[string]string
}
/*
func (this *Data)splitBy(id Identifier)(reply Data)  {
	// any key before Id will be split out
	for x,y:=range this.dat{
		if IDlize(x).ValPtr.Cmp(id.ValPtr)<=0{
			reply.dat[x]=y
			delete(this.dat,x)
		}
	}
	return
}

 */

