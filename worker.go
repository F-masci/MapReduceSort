package main

import (
	"MapReduceSort/config"
	"MapReduceSort/structs"
	"MapReduceSort/utils"
	"errors"
	"flag"
	"log"
	"net"
	"net/rpc"
	"sort"
	"strconv"
	"sync"
	"time"
)

// WorkerHandler Handler for exposing the worker's RPC methods
type WorkerHandler struct {
	mapService    bool                       // Indicates whether the map service is enabled
	reduceService bool                       // Indicates whether the reduce service is enabled
	reduceQueue   map[ReduceQueueKey][][]int // Queue for storing reduce requests
	mutex         sync.Mutex                 // Mutex for synchronizing access to the reduce queue
}

// ReduceQueueKey Struct to uniquely identify reduce requests by client and timestamp
type ReduceQueueKey struct {
	client    string
	timestamp time.Time
}

// Map Handles the map operation
func (wh *WorkerHandler) Map(request structs.SortMapRequest, _ *structs.SortMapResponse) error {

	// Ensure the map service is enabled
	if !wh.mapService {
		return errors.New("map service not enabled")
	}

	// Sort the input data
	sort.Ints(request.Request)
	log.Println("Mapped")

	// Retrieve the addresses of reducers
	_, _, reducerAddresses := config.ParseConfig()
	nReducer := len(reducerAddresses)
	distribution := utils.DistributeNumbersByValue(request.Request, nReducer)

	// Send sorted chunks to reducers
	for i := range nReducer {

		chunk := distribution[i]
		reducerConfig := reducerAddresses[i]

		// Connect to the reducer
		client, err := rpc.Dial(reducerConfig.Proto, reducerConfig.Address())
		utils.CheckError(err)

		// Prepare the request for the reducer
		request := structs.SortReduceRequest{Client: request.Client, Timestamp: request.Timestamp, Request: chunk, ReducerIdx: i}
		var reply structs.SortReduceResponse

		// Call the reduce method on the worker
		err = client.Call("WorkerHandler.Reduce", request, &reply)
		utils.CheckError(err)

	}

	return nil
}

// Reduce Handles the reduce operation
func (wh *WorkerHandler) Reduce(request structs.SortReduceRequest, _ *structs.SortMapResponse) error {

	// Ensure the reduce service is enabled
	if !wh.reduceService {
		return errors.New("reduce service not enabled")
	}

	// Generate a unique key for the reduce queue
	reduceQueueKey := ReduceQueueKey{client: request.Client, timestamp: request.Timestamp}
	wh.mutex.Lock()
	// Append the request to the reduce queue
	wh.reduceQueue[reduceQueueKey] = append(wh.reduceQueue[reduceQueueKey], request.Request)
	wh.mutex.Unlock()

	// Retrieve the addresses of the mappers
	_, mapperAddresses, _ := config.ParseConfig()
	// If all the mappers' requests have been received, start the reduce operation
	if len(wh.reduceQueue[reduceQueueKey]) == len(mapperAddresses) {
		// Create a filename for the result
		filename := request.Client + "_" + request.Timestamp.Format("2006_01_02_15_04_05") + "__" + strconv.Itoa(request.ReducerIdx)
		// Merge and sort the arrays from the reduce queue
		res := utils.MergeSortedArrays(wh.reduceQueue[reduceQueueKey])
		// Write the result to a file
		err := utils.WriteResultToFile(res, filename)
		utils.CheckError(err)
		log.Println(reduceQueueKey, "Reduced")

		// Clean up the reduce queue for the current key
		wh.mutex.Lock()
		delete(wh.reduceQueue, reduceQueueKey)
		wh.mutex.Unlock()
	}

	return nil
}

func main() {

	/* --- ARGUMENTS --- */

	// Specify the parameters to contact the worker
	address := flag.String("address", "localhost", "Specifies the address to contact the worker")
	port := flag.Int("port", 0, "Specifies the port to contact the worker")
	proto := flag.String("proto", "tcp", "Specifies the protocol to contact the worker")
	mapService := flag.Bool("map", false, "Specifies whether to enable the map service on the worker")
	reduceService := flag.Bool("reduce", false, "Specifies whether to enable the reduce service on the worker")
	flag.Parse()

	if *port == 0 {
		log.Fatal("Please specify the port number")
	}

	// Allocate the worker address configuration
	workerAddress := structs.WorkerAddress{
		Host:  *address,
		Port:  *port,
		Proto: *proto,
	}

	// Allocate the worker handler and pass a pointer
	workerHandler := &WorkerHandler{
		mapService:    *mapService,
		reduceService: *reduceService,
		reduceQueue:   make(map[ReduceQueueKey][][]int),
	}

	// Create the RPC server
	server := rpc.NewServer()
	err := server.Register(workerHandler)
	utils.CheckError(err)

	// Bind the server to the specified address
	lis, err := net.Listen(workerAddress.Proto, workerAddress.Address())
	utils.CheckError(err)
	log.Printf("RPC worker listens on port %d", workerAddress.Port)

	// Start the worker to listen for incoming connections
	server.Accept(lis)
}
