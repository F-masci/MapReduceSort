package config

import "fmt"

type WorkerAddress struct {
	Host  string
	Port  int
	Proto string
}

// Address Calcola l'indirizzo del mapper dinamicamente
func (m WorkerAddress) Address() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}

var MapperAddresses = []WorkerAddress{
	{Host: "localhost", Port: 45980, Proto: "tcp"},
	{Host: "localhost", Port: 45981, Proto: "tcp"},
	{Host: "localhost", Port: 45982, Proto: "tcp"},
	{Host: "localhost", Port: 45983, Proto: "tcp"},
}

var ReducerAddresses = []WorkerAddress{
	{Host: "localhost", Port: 45980, Proto: "tcp"},
	{Host: "localhost", Port: 45981, Proto: "tcp"},
	{Host: "localhost", Port: 45982, Proto: "tcp"},
	{Host: "localhost", Port: 45983, Proto: "tcp"},
}

const (
	CasualNumbersDim   = 1000
	CasualNumbersRange = 10000
)
