package main

import (
	"MapReduceSort/config"
	"MapReduceSort/structs"
	"MapReduceSort/utils"
	"flag"
	"log"
	"math/rand"
	"net/rpc"
)

func main() {

	/* --- SETUP --- */

	// Specify the client identifier connecting to the master
	client := flag.String("client", "guest", "Client identifier")
	// Select the master's config index from the file
	masterIdx := flag.Int("master-idx", -1, "Specifies the index of the master to use from the config file")
	// Provide the master's configuration manually
	masterAddress := flag.String("master-address", "", "Specifies the master's address to use")
	masterPort := flag.Int("master-port", 0, "Specifies the master's port to use")
	masterProto := flag.String("master-proto", "", "Specifies the master's protocol to use")
	flag.Parse()

	// Retrieve master's addresses from the config file
	masterAddresses, _, _ := config.ParseConfig()

	// Load the master's configuration
	var masterConfig = structs.MasterAddress{}
	if *masterIdx > -1 && *masterIdx < len(masterAddresses) { // Select config from the file
		masterConfig = masterAddresses[*masterIdx]
	} else if *masterAddress != "" {
		masterConfig = structs.MasterAddress{ // Select manually provided config
			Host:  *masterAddress,
			Port:  *masterPort,
			Proto: *masterProto,
		}
	} else {
		log.Fatalln("Please select a master to use")
	}

	/* --- RANDOM GENERATION --- */

	// Generate N random numbers
	numbers := make([]int, config.CasualNumbersDim)
	for i := 0; i < config.CasualNumbersDim; i++ {
		numbers[i] = rand.Intn(config.CasualNumbersRange) // [0, config.CasualNumbersRange)
	}
	log.Println("Generated numbers:", numbers)

	/* --- RPC --- */

	// Create the connection to the master
	master, err := rpc.Dial(masterConfig.Proto, masterConfig.Address())
	utils.CheckError(err)

	// Send the Sort request
	request := structs.SortRequest{Request: numbers, Client: *client}
	master.Go("MasterHandler.Sort", request, nil, nil)
	utils.CheckError(err)
}
