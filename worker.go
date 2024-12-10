package main

import (
	"MapReduceSort/config"
	"MapReduceSort/structs"
	"MapReduceSort/utils"
	"flag"
	"log"
	"net"
	"net/rpc"
	"sort"
	"strconv"
	"sync"
	"time"
)

// WorkerHandler Handler per esporre il metodo RPC del worker
type WorkerHandler struct {
	ReduceQueue map[ReduceQueueKey][][]int
	Mutex       sync.Mutex
}

type ReduceQueueKey struct {
	Client    string
	Timestamp time.Time
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

func (wh *WorkerHandler) Map(request structs.SortMapRequest, _ *structs.SortMapResponse) error {

	sort.Ints(request.Request)
	log.Println("Mapped")

	_, _, reducerAddresses := config.ParseConfig()
	nReducer := len(reducerAddresses)
	distribution := distributeNumbersByValue(request.Request, config.CasualNumbersRange, nReducer)

	// Ordina i chunk
	for i := range nReducer {

		chunk := distribution[i]
		reducerConfig := reducerAddresses[i]

		client, err := rpc.Dial(reducerConfig.Proto, reducerConfig.Address())
		utils.CheckError(err)

		request := structs.SortReduceRequest{Client: request.Client, Timestamp: request.Timestamp, Request: chunk, ReducerIdx: i}
		var reply structs.SortReduceResponse

		err = client.Call("WorkerHandler.Reduce", request, &reply)
		utils.CheckError(err)

	}

	return nil
}

func (wh *WorkerHandler) Reduce(request structs.SortReduceRequest, _ *structs.SortMapResponse) error {

	reduceQueueKey := ReduceQueueKey{Client: request.Client, Timestamp: request.Timestamp}
	wh.Mutex.Lock()
	wh.ReduceQueue[reduceQueueKey] = append(wh.ReduceQueue[reduceQueueKey], request.Request)
	wh.Mutex.Unlock()

	_, mapperAddresses, _ := config.ParseConfig()
	if len(wh.ReduceQueue[reduceQueueKey]) == len(mapperAddresses) {
		filename := request.Client + "_" + request.Timestamp.Format("2006_01_02_15_04_05") + "__" + strconv.Itoa(request.ReducerIdx)
		res := mergeSortedArrays(wh.ReduceQueue[reduceQueueKey])
		err := utils.WriteResultToFile(res, filename)
		utils.CheckError(err)
		log.Println(reduceQueueKey, "Reduced")
		wh.Mutex.Lock()
		delete(wh.ReduceQueue, reduceQueueKey)
		wh.Mutex.Unlock()
	}

	return nil
}

func main() {

	// Recuperare gli argomenti
	address := flag.String("address", "localhost", "Specifica l'indirizzo su cui contattare il master")
	port := flag.Int("port", 0, "Specifica la porta su cui contattare il master")
	proto := flag.String("proto", "tcp", "Specifica il protocollo con cui contattare il master")
	_ = flag.String("mode", "both", "Specifica se utilizzare il worker come mapper, reducer o entrambi")
	flag.Parse()

	if *port == 0 {
		log.Fatal("Specificare il numero di porta")
	}

	workerAddress := structs.WorkerAddress{
		Host:  *address,
		Port:  *port,
		Proto: *proto,
	}

	// Alloco la funzione del worker e passo un puntatore
	sortHandler := &WorkerHandler{ReduceQueue: make(map[ReduceQueueKey][][]int)}

	server := rpc.NewServer()
	err := server.Register(sortHandler)
	utils.CheckError(err)

	// Faccio il bind con gli indirizzi specificati nelle config
	lis, err := net.Listen(workerAddress.Proto, workerAddress.Address())
	utils.CheckError(err)
	log.Printf("RPC worker listens on port %d", workerAddress.Port)

	// Metto il worker in ascolto
	server.Accept(lis)

}
