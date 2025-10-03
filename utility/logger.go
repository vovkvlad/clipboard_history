package utility

import (
	"io"
	"log"
	"os"
	"path"
)

func InitLogger() *os.File {
	// TODO: Get log file path based on the OS
	logFile, err := os.OpenFile(path.Join(".", "tmp", "app.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	multiWriter := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multiWriter)

	return logFile
}
