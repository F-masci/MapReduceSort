package main

import (
	"MapReduceSort/config"
	"MapReduceSort/structs"
	"MapReduceSort/utils"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sort"
	"strconv"
)

// WorkerHandler Handler per esporre il metodo RPC del worker
type WorkerHandler struct {
	ID          int
	ReduceQueue [][]int
}

func distributeNumbersByValue(input []int, maxValue, nodes int) [][]int {
	// Calcola il range per ciascun nodo
	rangeSize := (maxValue + 1) / nodes

	// Inizializza il risultato
	result := make([][]int, nodes)

	// Distribuisci i numeri nei nodi
	for _, num := range input {
		if num > maxValue {
			continue // Salta i numeri che superano il massimo valore
		}
		// Calcola il nodo a cui appartiene il numero
		node := num / rangeSize
		if node >= nodes {
			node = nodes - 1 // Assicura che il valore massimo vada nell'ultimo nodo
		}
		result[node] = append(result[node], num)
	}

	return result
}

func mergeSortedArrays(input [][]int) []int {
	var result []int
	// Unisci tutti gli array
	for _, array := range input {
		result = append(result, array...)
	}
	// Ordina l'array risultante
	sort.Ints(result)
	return result
}

func (WorkerHandler) Map(request structs.SortMapRequest, reply *structs.SortMapResponse) error {

	reply.Response = make([]int, len(request.Request))
	copy(reply.Response, request.Request)
	sort.Ints(reply.Response)
	fmt.Println("Richiesta ricevuta: ", request.Request, "\nNumeri ordinati: ", reply.Response)

	nReducer := len(config.ReducerAddresses)
	distribution := distributeNumbersByValue(reply.Response, config.CasualNumbersRange, nReducer)

	// Ordina i chunk
	for i := range nReducer {

		chunk := distribution[i]
		reducerConfig := config.ReducerAddresses[i]

		client, err := rpc.Dial(reducerConfig.Proto, reducerConfig.Address())
		utils.CheckError(err)

		request := structs.SortReduceRequest{Client: request.Client, Timestamp: request.Timestamp, Request: chunk}
		var reply structs.SortReduceResponse

		err = client.Call("WorkerHandler.Reduce", request, &reply)
		utils.CheckError(err)

	}

	return nil
}

func (wh *WorkerHandler) Reduce(request structs.SortReduceRequest, _ *structs.SortMapResponse) error {

	wh.ReduceQueue = append(wh.ReduceQueue, request.Request)
	fmt.Println("[", wh.ID, "] Richiesta di reduce ricevuta da ", request.Client, "-", request.Timestamp)

	if len(wh.ReduceQueue) == len(config.MapperAddresses) {
		res := mergeSortedArrays(wh.ReduceQueue)
		utils.WriteResultToFile(res, request.Client+"_"+request.Timestamp.Format("2006_01_02_15_04_05")+"__"+strconv.Itoa(wh.ID))
	}

	return nil
}

func main() {

	workerAddresses := utils.MergeAndRemoveDuplicates(config.MapperAddresses, config.ReducerAddresses)

	// Avvio i listener dei workers
	for idx, workerConfig := range workerAddresses {

		// Alloco la funzione del worker e passo un puntatore
		sortHandler := &WorkerHandler{ID: idx, ReduceQueue: make([][]int, 0)}

		server := rpc.NewServer()
		err := server.Register(sortHandler)
		utils.CheckError(err)

		// Faccio il bind con gli indirizzi specificati nelle config
		lis, err := net.Listen(workerConfig.Proto, workerConfig.Address())
		utils.CheckError(err)
		log.Printf("RPC worker listens on port %d", workerConfig.Port)

		// Metto il mapper in ascolto su una goroutine differente
		go server.Accept(lis)
	}
	// Impedisco alla routine principale di terminare
	select {}

}
