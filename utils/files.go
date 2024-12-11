package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// WriteResultToFile Writes an array to a file
func WriteResultToFile(arr []int, filename string) error {

	const folder string = ".output/"

	// Check if the "result" directory exists, otherwise create it
	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		// Create the directory if it does not exist
		err := os.Mkdir(folder, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return err
		}
	}

	// Now create the file in the directory
	file, err := os.Create(folder + filename)
	if err != nil {
		return fmt.Errorf("error creating the file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			_ = fmt.Errorf("error closing the file: %v", err)
		}
	}(file)

	// Convert each number to a string and write it to the file
	var sb strings.Builder
	for _, num := range arr {
		sb.WriteString(strconv.Itoa(num) + "\n")
	}

	// Write the entire string to the file
	_, err = file.WriteString(sb.String())
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}
