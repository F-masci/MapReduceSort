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
)

// SortHandler Handler per esporre il metodo RPC
type SortHandler struct{}

func (SortHandler) Sort(request structs.SortMapRequest, reply *structs.SortMapResponse) error {

	reply.Response = make([]int, len(request.Request))
	copy(reply.Response, request.Request)
	sort.Ints(reply.Response)
	fmt.Println("Richiesta ricevuta: ", request.Request, "\nNumeri ordinati: ", reply.Response)

	return nil
}

func main() {

	sortHandler := new(SortHandler)
	server := rpc.NewServer() // create a server
	err := server.Register(sortHandler)
	utils.CheckError(err)

	lis, err := net.Listen(config.MapperProto, config.MapperAddr) // create a listener that handles RPCs
	utils.CheckError(err)
	log.Printf("RPC server listens on port %s", config.MapperPort)

	go func() {
		for {
			server.Accept(lis) // register the listener and accept inbound RPCs
		}
	}()
	select {}

}
