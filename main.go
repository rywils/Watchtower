package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"watchtower/watcher"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	w, err := watcher.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("[*] Watchtower running")
	w.Run(ctx)
	fmt.Println("\n[*] Watchtower stopped")
}

