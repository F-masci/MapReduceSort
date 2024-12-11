package utils

import "log"

// CheckError Checks if the request has produced an error
func CheckError(err error) {
	if err != nil {
		log.Fatal(err) // Logs the error and terminates the program if an error occurs
	}
}
