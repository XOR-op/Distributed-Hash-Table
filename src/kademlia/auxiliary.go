package kademlia

func Triple(condition bool, a, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

type FindStat int

const (
	FindFail      FindStat = 0
	FindEnough    FindStat = 1
	FindNotEnough FindStat = 2
	FindValue     FindStat = 3
)

type FindNodeRequest struct {
	Target Identifier
	Auth   *Contact
}

type FindNodeResponse struct {
	Stat   FindStat
	Auth   *Contact
	KNodes []*Contact
	Amount int
	Err    error
}

type FindValueRequest struct {
	Key  string
	ID   Identifier
	Auth *Contact
}

type FindValueResponse struct {
	Stat   FindStat
	Auth   *Contact
	KNodes []*Contact
	Amount int
	Value  string
	Err    error
}

type PingRequest struct {
	Auth *Contact
}

type StoreRequest struct {
	Key, Value string
	Auth       *Contact
	original   bool
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
