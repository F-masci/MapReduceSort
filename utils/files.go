package utils

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Scrive un array ordinato su un file
func WriteResultToFile(arr []int, filename string) error {
	// Ordina l'array
	sort.Ints(arr)

	// Controlla se la directory "result" esiste, altrimenti la crea
	_, err := os.Stat("result")
	if os.IsNotExist(err) {
		// Crea la directory se non esiste
		err := os.Mkdir("result", 0755)
		if err != nil {
			fmt.Println("Errore nella creazione della directory:", err)
			return err
		}
	}

	// Ora crea il file nella directory result
	file, err := os.Create("result/" + filename)
	if err != nil {
		fmt.Println("Errore nella creazione del file:", err)
		return err
	}
	defer file.Close() // Assicurati di chiudere il file quando finito

	// Converti ogni numero in stringa e scrivi sul file
	var sb strings.Builder
	for _, num := range arr {
		sb.WriteString(strconv.Itoa(num) + "\n")
	}

	_, err = file.WriteString(sb.String())
	if err != nil {
		return fmt.Errorf("errore durante la scrittura nel file: %v", err)
	}

	return nil
}
