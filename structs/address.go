package structs

import "fmt"

type MasterAddress struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	Proto string `json:"proto"`
}

type WorkerAddress struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	Proto string `json:"proto"`
}

// Address Calcola l'indirizzo del worker dinamicamente
func (m WorkerAddress) Address() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}

// Address Calcola l'indirizzo del mapper dinamicamente
func (m MasterAddress) Address() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}
