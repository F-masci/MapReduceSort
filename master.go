package main

import (
	"MapReduceSort/config"
	"MapReduceSort/structs"
	"MapReduceSort/utils"
	"flag"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

type MasterHandler struct{}

func (MasterHandler) Sort(request structs.SortRequest, reply *structs.SortMapResponse) error {

	log.Println("New request from", request.Client)

	_, mapperAddresses, _ := config.ParseConfig()
	numbers := request.Request

	nMapper := len(mapperAddresses)
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
	timestamp := time.Now()
	for i := range nMapper {

		chunk := chunks[i]
		mapperConfig := mapperAddresses[i]

		mapper, err := rpc.Dial(mapperConfig.Proto, mapperConfig.Address())
		utils.CheckError(err)

		request := structs.SortMapRequest{Client: request.Client, Timestamp: timestamp, Request: chunk, MapperIdx: i}

		wg.Add(1)

		go func() {
			defer wg.Done()
			mapper.Go("WorkerHandler.Map", request, nil, nil)
			utils.CheckError(err)
		}()
	}
	wg.Wait()

	return nil

}

func main() {

	// Recuperare gli argomenti
	address := flag.String("address", "localhost", "Specifica l'indirizzo su cui contattare il master")
	port := flag.Int("port", 0, "Specifica la porta su cui contattare il master")
	proto := flag.String("proto", "tcp", "Specifica il protocollo con cui contattare il master")
	flag.Parse()

	if *port == 0 {
		log.Fatal("Specificare il numero di porta")
	}

	masterAddress := structs.MasterAddress{
		Host:  *address,
		Port:  *port,
		Proto: *proto,
	}

	// Alloco la funzione del worker e passo un puntatore
	masterHandler := MasterHandler{}

	server := rpc.NewServer()
	err := server.Register(masterHandler)
	utils.CheckError(err)

	// Faccio il bind con gli indirizzi specificati nelle config
	lis, err := net.Listen(masterAddress.Proto, masterAddress.Address())
	utils.CheckError(err)
	log.Printf("RPC master listens on port %d", masterAddress.Port)

	// Metto il master in ascolto
	server.Accept(lis)

}
