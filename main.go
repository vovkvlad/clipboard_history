package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/vovkvlad/clipboard_history/core/clipboard"
)

func main() {
	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create watcher and add listener
	watcher := clipboard.NewWatcher()
	watcher.Add_listener(func(data []byte) {
		fmt.Println("Clipboard changed:", string(data))
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

	fmt.Println("Press Ctrl+C to exit.")

	// Wait for either Ctrl+C or the goroutine to finish
	select {
	case <-sigChan:
		fmt.Println("\nReceived interrupt signal, cancelling...")
		cancel()
		// Wait for the goroutine to finish
		if err := <-done; err != nil {
			fmt.Printf("Clipboard watcher error: %v\n", err)
		}
	case err := <-done:
		if err != nil {
			fmt.Printf("Clipboard watcher error: %v\n", err)
		} else {
			fmt.Println("Clipboard watcher finished")
		}
	}

	fmt.Println("Program terminated.")
}
