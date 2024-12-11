package config

import (
	"MapReduceSort/structs"
	"MapReduceSort/utils"
	"encoding/json"
	"fmt"
	"os"
)

const (
	CasualNumbersDim   = 1000  // Dimension of the random numbers array
	CasualNumbersRange = 10000 // Maximum range for the random numbers
)

// loadAddresses Reads the JSON file and performs parsing
func loadAddresses(filename string, target interface{}) error {
	// Open the JSON file
	file, err := os.Open(filename)
	if err != nil {
		return err // Return error if file cannot be opened
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			_ = fmt.Errorf("error closing the file: %v", err)
		}
	}(file)

	// Decode the JSON content
	decoder := json.NewDecoder(file)
	err = decoder.Decode(target)
	if err != nil {
		return err
	}

	return nil
}

// ParseConfig Parses the configuration files for master, mapper, and reducer addresses
func ParseConfig() (masterAddresses []structs.MasterAddress, mapperAddresses []structs.WorkerAddress, reducerAddresses []structs.WorkerAddress) {

	// Load master addresses from JSON file
	err := loadAddresses("config/master.json", &masterAddresses)
	utils.CheckError(err)

	// Load mapper addresses from JSON file
	err = loadAddresses("config/mapper.json", &mapperAddresses)
	utils.CheckError(err)

	// Load reducer addresses from JSON file
	err = loadAddresses("config/reducer.json", &reducerAddresses)
	utils.CheckError(err)

	return
}
