package watcher

import (
	"context"
	"time"
)

const pollInterval = 3 * time.Second

type Watcher struct {
	state *State
}

func New() (*Watcher, error) {
	state, err := LoadState()
	if err != nil {
		return nil, err
	}
	if state == nil {
		state = NewState()
		println("[*] Baseline created")
	}
	return &Watcher{state: state}, nil
}

func (w *Watcher) Run(ctx context.Context) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			SaveState(w.state)
			return
		case <-ticker.C:
			ignored := DefaultGatewayIPs()
			curr := ReadARP(ignored)
			events := Diff(w.state, curr, ignored)
			for _, e := range events {
				e.Print()
			}
			w.state = curr
		}
	}
}
