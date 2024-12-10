package config

import (
	"MapReduceSort/structs"
	"MapReduceSort/utils"
	"encoding/json"
	"os"
)

const (
	CasualNumbersDim   = 1000
	CasualNumbersRange = 10000
)

// Funzione per leggere il file JSON e fare il parsing
func loadAddresses(filename string, target interface{}) error {
	// Leggi il file JSON
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decodifica il contenuto JSON
	decoder := json.NewDecoder(file)
	err = decoder.Decode(target)
	if err != nil {
		return err
	}

	return nil
}

func ParseConfig() (masterAddresses []structs.MasterAddress, mapperAddresses []structs.WorkerAddress, reducerAddresses []structs.WorkerAddress) {

	err := loadAddresses("config/master.json", &masterAddresses)
	utils.CheckError(err)

	err = loadAddresses("config/mapper.json", &mapperAddresses)
	utils.CheckError(err)

	err = loadAddresses("config/reducer.json", &reducerAddresses)
	utils.CheckError(err)

	return
}
