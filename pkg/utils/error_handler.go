package utils

import (
	"fmt"
	"log"
	"os"
)

func HandleError(err error, message string) error {
	errorLogger := log.New(os.Stderr, "\nERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger.Println(message, err)
	return fmt.Errorf(message)
}
