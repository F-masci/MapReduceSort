package structs

import "time"

// SortMapRequest RPC Map request
type SortMapRequest struct {
	Client    string
	Timestamp time.Time
	Request   []int
}

// SortReduceRequest RPC Reduce request
type SortReduceRequest struct {
	Client    string
	Timestamp time.Time
	Request   []int
}
