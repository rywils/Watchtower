package watcher

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const pollInterval = 3 * time.Second

const gatewayRefreshInterval = 2 * time.Minute

type Config struct {
	PollInterval    time.Duration
	ResolveDNS      bool
	JSONOutput      bool
	StaticIgnoredIP map[string]struct{}
}

type Watcher struct {
	state              *State
	cfg                Config
	cachedGatewayIPs   map[string]struct{}
	lastGatewayRefresh time.Time
}

func New(cfg Config) (*Watcher, error) {
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = pollInterval
	}
	if cfg.StaticIgnoredIP == nil {
		cfg.StaticIgnoredIP = map[string]struct{}{}
	}

	state, err := LoadState()
	if err != nil {
		return nil, err
	}
	if state == nil {
		state = NewState()
		println("[*] Baseline created")
	}
	return &Watcher{
		state:            state,
		cfg:              cfg,
		cachedGatewayIPs: map[string]struct{}{},
	}, nil
}

func (w *Watcher) Run(ctx context.Context) {
	ticker := time.NewTicker(w.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if err := SaveState(w.state); err != nil {
				fmt.Printf("[!] failed to save state: %v\n", err)
			}
			return
		case <-ticker.C:
			ignored := w.ignoredIPs()
			curr := ReadARP(ignored)
			events := Diff(w.state, curr, ignored)
			for _, e := range events {
				if w.cfg.JSONOutput {
					fmt.Println(e.JSON())
					continue
				}
				fmt.Println(e.Text(w.cfg.ResolveDNS))
			}
			w.state = curr
		}
	}
}

func (w *Watcher) ignoredIPs() map[string]struct{} {
	combined := map[string]struct{}{}
	for ip := range w.cfg.StaticIgnoredIP {
		combined[ip] = struct{}{}
	}

	now := time.Now()
	if now.Sub(w.lastGatewayRefresh) > gatewayRefreshInterval || len(w.cachedGatewayIPs) == 0 {
		w.cachedGatewayIPs = DefaultGatewayIPs()
		w.lastGatewayRefresh = now
	}
	for ip := range w.cachedGatewayIPs {
		combined[ip] = struct{}{}
	}
	return combined
}

func ParseIgnoredIPs(csv string) map[string]struct{} {
	out := map[string]struct{}{}
	for _, part := range strings.Split(csv, ",") {
		p := strings.TrimSpace(part)
		if p == "" {
			continue
		}
		out[p] = struct{}{}
	}
	return out
}
