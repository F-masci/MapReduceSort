package structs

import "time"

// SortRequest RPC Master request
type SortRequest struct {
	Request []int
	Client  string
}

// SortMapRequest RPC Map request
type SortMapRequest struct {
	Client    string
	Timestamp time.Time
	Request   []int
	MapperIdx int
}

// SortReduceRequest RPC Reduce request
type SortReduceRequest struct {
	Client     string
	Timestamp  time.Time
	Request    []int
	ReducerIdx int
}
