package utils

import "log"

// CheckError Controlla se la richiesta ha prodotto un errore
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
