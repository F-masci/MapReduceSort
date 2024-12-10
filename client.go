package main

import (
	"MapReduceSort/config"
	"MapReduceSort/structs"
	"MapReduceSort/utils"
	"flag"
	"fmt"
	"math/rand"
	"net/rpc"
)

func main() {

	masterAddresses, _, _ := config.ParseConfig()

	// Recuperare gli argomenti
	client := flag.String("client", "guest", "Identificatore del client")
	flag.Parse()

	numbers := make([]int, config.CasualNumbersDim)

	// Genera N numeri casuali
	for i := 0; i < config.CasualNumbersDim; i++ {
		numbers[i] = rand.Intn(config.CasualNumbersRange) // Numeri interi casuali nell'intervallo [0, CasualNumbersRange]
	}

	// Stampa i numeri casuali
	fmt.Printf("Numeri generati: %v\n", numbers)

	idx := rand.Intn(len(masterAddresses))
	master, err := rpc.Dial(masterAddresses[idx].Proto, masterAddresses[idx].Address())
	utils.CheckError(err)

	request := structs.SortRequest{Request: numbers, Client: *client}

	master.Go("MasterHandler.Sort", request, nil, nil)
	utils.CheckError(err)

}
