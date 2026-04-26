package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"watchtower/watcher"
)

func main() {
	interval := flag.Duration("interval", 3*time.Second, "polling interval (e.g. 3s, 10s)")
	noDNS := flag.Bool("no-dns", false, "disable reverse DNS lookups in human-readable output")
	jsonOutput := flag.Bool("json", false, "emit JSON events (one per line)")
	ignoreCSV := flag.String("ignore", "", "comma-separated IPs to ignore")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	w, err := watcher.New(watcher.Config{
		PollInterval:    *interval,
		ResolveDNS:      !*noDNS,
		JSONOutput:      *jsonOutput,
		StaticIgnoredIP: watcher.ParseIgnoredIPs(*ignoreCSV),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("[*] Watchtower running")
	w.Run(ctx)
	fmt.Println("\n[*] Watchtower stopped")
}
