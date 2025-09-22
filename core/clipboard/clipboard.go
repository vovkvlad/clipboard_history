package clipboard

import (
	"context"
	"fmt"

	"golang.design/x/clipboard"
)

func Start(ctx context.Context, watcher *Watcher) error {
	err := clipboard.Init()
	if err != nil {
		return err
	}

	change := clipboard.Watch(ctx, clipboard.FmtText)

	fmt.Println("Waiting for clipboard changes...")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Clipboard watcher shutting down...")
			return nil
		case data := <-change:
			watcher.Notify(data)
		}
	}
}
