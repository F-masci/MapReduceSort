package structs

import "fmt"

// MasterAddress Defines the address of the master with host, port, and protocol
type MasterAddress struct {
	Host  string `json:"host"`  // Hostname or IP address of the master
	Port  int    `json:"port"`  // Port number of the master
	Proto string `json:"proto"` // Protocol used to contact the master (e.g., "tcp")
}

// WorkerAddress Defines the address of the worker with host, port, and protocol
type WorkerAddress struct {
	Host  string `json:"host"`  // Hostname or IP address of the worker
	Port  int    `json:"port"`  // Port number of the worker
	Proto string `json:"proto"` // Protocol used to contact the worker (e.g., "tcp")
}

// Address Returns the address of the worker in the format "host:port"
func (m WorkerAddress) Address() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port) // Constructs the address string
}

// Address Returns the address of the master in the format "host:port"
func (m MasterAddress) Address() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port) // Constructs the address string
}
