package main

import (
	"MapReduceSort/config"
	"MapReduceSort/structs"
	"MapReduceSort/utils"
	"fmt"
	"math/rand"
	"net/rpc"
	"sync"
	"time"
)

func main() {

	numbers := make([]int, config.CasualNumbersDim)

	// Genera N numeri casuali
	for i := 0; i < config.CasualNumbersDim; i++ {
		numbers[i] = rand.Intn(config.CasualNumbersRange) // Numeri interi casuali nell'intervallo [0, CasualNumbersRange]
	}

	// Stampa i numeri casuali
	fmt.Printf("Numeri generati: %v\n", numbers)

	nMapper := len(config.MapperAddresses)
	chunkSize := len(numbers) / nMapper
	remainder := len(numbers) % nMapper

	// Suddivido lo slice
	chunks := make([][]int, nMapper)
	start := 0
	for i := 0; i < nMapper; i++ {
		end := start + chunkSize
		if i < remainder {
			end++ // Aggiungi un elemento extra ai primi "remainder" mapper
		}
		chunks[i] = numbers[start:end]
		start = end
	}

	// Ordina i chunk
	var wg sync.WaitGroup
	for i := range nMapper {

		chunk := chunks[i]
		mapperConfig := config.MapperAddresses[i]

		mapper, err := rpc.Dial(mapperConfig.Proto, mapperConfig.Address())
		utils.CheckError(err)

		request := structs.SortMapRequest{Client: "localhost", Timestamp: time.Now(), Request: chunk}

		wg.Add(1)

		go func() {
			defer wg.Done() // Decrementa il counter del WaitGroup quando la chiamata è finita
			mapper.Go("WorkerHandler.Map", request, nil, nil)
			utils.CheckError(err)
		}()
	}
	wg.Wait()
}
