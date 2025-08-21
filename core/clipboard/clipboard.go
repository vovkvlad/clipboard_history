package clipboard

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.design/x/clipboard"
)

func Start() {
	err := clipboard.Init()

	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	change := clipboard.Watch(ctx, clipboard.FmtText)

	fmt.Println("Waiting for clipboard changes... Press Ctrl+C to exit.")

	// Set up signal handling for Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a channel to signal when the goroutine is done
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("\nShutting down...")
				return
			case data := <-change:
				fmt.Println("Clipboard changed:", string(data))
			}
		}
	}()

	// Wait for either Ctrl+C or the goroutine to finish
	select {
	case <-sigChan:
		fmt.Println("\nReceived interrupt signal, cancelling...")
		cancel()
	case <-done:
		fmt.Println("Goroutine finished")
	}

	// Wait for the goroutine to actually finish
	<-done
	fmt.Println("Program terminated.")
}
