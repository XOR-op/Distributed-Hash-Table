package kademlia

func Triple(condition bool, a, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

type FindStat int

const (
	FindFail      = 0
	FindEnough    = 1
	FindNotEnough = 2
	FindValue     = 3
)

type FindNodeRequest struct {
	Target Identifier
	Auth   *Contact
}

type FindNodeResponse struct {
	Stat   FindStat
	Auth   *Contact
	KNodes [K]*Contact
	Amount int
	err    error
}

type FindValueRequest struct {
	Key  string
	ID   Identifier
	Auth *Contact
}

type FindValueResponse struct {
	Stat   FindStat
	Auth   *Contact
	KNodes *[K]*Contact
	Amount int
	Value  *string
	err    error
}

type PingRequest struct {
	Auth *Contact
}

type StoreRequest struct {
	Key, Value string
	Auth       *Contact
}
