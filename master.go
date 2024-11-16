package main

import (
	"MapReduceSort/config"
	"MapReduceSort/structs"
	"MapReduceSort/utils"
	"fmt"
	"math/rand"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial(config.MapperProto, config.MapperAddr)
	utils.CheckError(err)

	numbers := make([]int, 100)

	// Genera 100 numeri casuali
	for i := 0; i < 100; i++ {
		numbers[i] = rand.Intn(100) // Numeri interi casuali nell'intervallo [0, 99]
	}

	request := structs.SortMapRequest{Request: numbers}
	var reply structs.SortMapResponse

	err = client.Call("SortHandler.Sort", request, &reply)
	utils.CheckError(err)

	fmt.Println("Risultato ordinato dal server:", reply)
}
