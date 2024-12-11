package structs

import "time"

// SortRequest Defines the RPC request sent to the master for sorting
type SortRequest struct {
	Request []int  // Array of integers to be sorted
	Client  string // Identifier for the client making the request
}

// SortMapRequest Defines the RPC request sent by the master to the mapper for sorting
type SortMapRequest struct {
	Client    string    // Identifier for the client making the request
	Timestamp time.Time // Timestamp of the request
	Request   []int     // Array of integers to be processed by the mapper
	MapperIdx int       // Index of the mapper handling the request
}

// SortReduceRequest Defines the RPC request sent by the mapper to the reducer for sorting
type SortReduceRequest struct {
	Client     string    // Identifier for the client making the request
	Timestamp  time.Time // Timestamp of the request
	Request    []int     // Array of integers to be processed by the reducer
	ReducerIdx int       // Index of the reducer handling the request
}
