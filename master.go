package main

import (
	"MapReduceSort/config"
	"MapReduceSort/structs"
	"MapReduceSort/utils"
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"
)

// MasterHandler Structure of the master
type MasterHandler struct{}

// Sort Executes the sorting of the input values
func (MasterHandler) Sort(request structs.SortRequest, _ *structs.SortMapResponse) error {

	log.Println("New request from", request.Client)

	/* --- SETUP --- */

	// Retrieve the mappers' addresses from the config
	_, mapperAddresses, _ := config.ParseConfig()
	numbers := request.Request

	// Calculate chunk distribution
	nMapper := len(mapperAddresses)
	chunkSize := len(numbers) / nMapper
	remainder := len(numbers) % nMapper

	/* --- CHUNKS --- */

	// Split the input into chunks
	chunks := make([][]int, nMapper)
	start := 0
	for i := 0; i < nMapper; i++ {
		end := start + chunkSize
		if i < remainder {
			end++ // Add an extra element to the first "remainder" mappers
		}
		chunks[i] = numbers[start:end]
		start = end
	}

	/* --- RPC --- */

	// Send chunks to the mappers
	timestamp := time.Now()
	for i := range mapperAddresses {
		// Load the chunk assigned to the mapper
		chunk := chunks[i]
		mapperConfig := mapperAddresses[i]

		// Create the connection to the mapper
		conn, err := rpc.Dial(mapperConfig.Proto, mapperConfig.Address())
		if err != nil {
			return fmt.Errorf("failed to connect to mapper at %s: %w", mapperConfig.Address(), err)
		}

		// Prepare the request
		req := structs.SortMapRequest{Client: request.Client, Timestamp: timestamp, Request: chunk, MapperIdx: i}

		// Perform the RPC call in a goroutine
		go func() {
			err := conn.Go("WorkerHandler.Map", req, nil, nil).Error
			utils.CheckError(err)
		}()
	}

	return nil
}

func main() {

	/* --- SETUP --- */

	// Specify the parameters to contact the master
	address := flag.String("address", "localhost", "Specifies the address to contact the master")
	port := flag.Int("port", 0, "Specifies the port to contact the master")
	proto := flag.String("proto", "tcp", "Specifies the protocol to contact the master")

	flag.Parse()

	if *port == 0 {
		log.Fatal("Please specify the port number")
	}

	// Allocate the master's configuration to bind to
	masterAddress := structs.MasterAddress{
		Host:  *address,
		Port:  *port,
		Proto: *proto,
	}

	/* --- RPC --- */

	// Allocate the master's structure
	masterHandler := &MasterHandler{}

	// Allocate the server
	server := rpc.NewServer()
	err := server.Register(masterHandler)
	utils.CheckError(err)

	// Bind to the specified address
	lis, err := net.Listen(masterAddress.Proto, masterAddress.Address())
	utils.CheckError(err)
	log.Printf("RPC master listens on port %d", masterAddress.Port)

	// Start the master listening for new connections
	server.Accept(lis)
}
