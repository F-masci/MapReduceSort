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

	// Select the config from the file
	idx := flag.Int("idx", -1, "Specifies the index of the config to use from the file")
	// Provide the configuration manually
	address := flag.String("address", "localhost", "Specifies the address to contact the master")
	port := flag.Int("port", 0, "Specifies the port to contact the master")
	proto := flag.String("proto", "tcp", "Specifies the protocol to contact the master")
	flag.Parse()

	// Retrieve master's addresses from the config file
	masterAddresses, _, _ := config.ParseConfig()

	// Load the master's configuration
	var masterConfig = structs.MasterAddress{}
	if *idx > -1 && *idx < len(masterAddresses) { // Select config from the file
		masterConfig = masterAddresses[*idx]
	} else if *address != "" && *port != 0 && *proto != "" {
		masterConfig = structs.MasterAddress{ // Select manually provided config
			Host:  *address,
			Port:  *port,
			Proto: *proto,
		}
	} else {
		log.Fatalln("Please provide a config to use")
	}

	/* --- RPC --- */

	// Allocate the master's structure
	masterHandler := &MasterHandler{}

	// Allocate the server
	server := rpc.NewServer()
	err := server.Register(masterHandler)
	utils.CheckError(err)

	// Bind to the specified address
	lis, err := net.Listen(masterConfig.Proto, masterConfig.Address())
	utils.CheckError(err)
	log.Printf("RPC master listens on %s:%s:%d", masterConfig.Proto, masterConfig.Host, masterConfig.Port)

	// Start the master listening for new connections
	server.Accept(lis)
}
