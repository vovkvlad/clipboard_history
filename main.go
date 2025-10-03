package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vovkvlad/clipboard_history/core/clipboard"
	"github.com/vovkvlad/clipboard_history/core/storage"
	"github.com/vovkvlad/clipboard_history/utility"
)

func main() {
	// Initialize logger
	logFile := utility.InitLogger()
	defer func() {
		if err := logFile.Close(); err != nil {
			log.Printf("Error closing log file: %v", err)
		}
	}()

	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Logger initialized")

	// Initialize database
	_, err := storage.InitDb()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Println("Database initialized")

	// Create watcher and add listener
	watcher := clipboard.NewWatcher()
	watcher.Add_listener(func(data []byte) {
		log.Println("Clipboard changed:", string(data))
	})

	// Set up signal handling for Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a channel to signal when the goroutine is done
	done := make(chan error, 1)

	// Start clipboard watcher in a goroutine
	go func() {
		done <- clipboard.Start(ctx, watcher)
	}()

	log.Println("Press Ctrl+C to exit.")

	// Wait for either Ctrl+C or the goroutine to finish
	select {
	case <-sigChan:
		log.Println("\nReceived interrupt signal, cancelling...")
		cancel()
		// Wait for the goroutine to finish
		if err := <-done; err != nil {
			log.Fatalf("Clipboard watcher error: %v\n", err)
		}
	case err := <-done:
		if err != nil {
			log.Fatalf("Clipboard watcher error: %v\n", err)
		} else {
			log.Println("Clipboard watcher finished")
		}
	}

	fmt.Println("Program terminated.")
}
