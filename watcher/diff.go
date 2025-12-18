package watcher

import (
	"net"
	"strings"
	"time"
)

type EventType string

const (
	EventNewDevice  EventType = "new_device"
	EventGoneDevice EventType = "device_gone"
	EventMACChange  EventType = "mac_change"
)

type Event struct {
	Type      EventType
	IP        string
	OldMAC    string
	NewMAC    string
	Timestamp int64
}

func isIgnorable(ip, mac string) bool {
	// broadcast MAC
	if mac == "ff:ff:ff:ff:ff:ff" {
		return true
	}

	// broadcast IP
	if strings.HasSuffix(ip, ".255") {
		return true
	}

	parsed := net.ParseIP(ip)
	if parsed == nil || !parsed.IsPrivate() {
		return true
	}

	// IGNORE GATEWAY
	if strings.HasSuffix(ip, ".1") {
		return true
	}

	return false
}

func Diff(prev, curr *State) []Event {
	now := time.Now().Unix()
	events := []Event{}

	// new / changed
	for ip, d := range curr.Devices {
		if isIgnorable(ip, d.MAC) {
			continue
		}

		if old, ok := prev.Devices[ip]; !ok {
			events = append(events, Event{
				Type:      EventNewDevice,
				IP:        ip,
				NewMAC:    d.MAC,
				Timestamp: now,
			})
		} else if old.MAC != d.MAC {
			events = append(events, Event{
				Type:      EventMACChange,
				IP:        ip,
				OldMAC:    old.MAC,
				NewMAC:    d.MAC,
				Timestamp: now,
			})
		}
	}

	// gone
	for ip, old := range prev.Devices {
		if isIgnorable(ip, old.MAC) {
			continue
		}

		if _, ok := curr.Devices[ip]; !ok {
			events = append(events, Event{
				Type:      EventGoneDevice,
				IP:        ip,
				OldMAC:    old.MAC,
				Timestamp: now,
			})
		}
	}

	return events
}

func (e Event) Print() {
	switch e.Type {
	case EventNewDevice:
		println("[+] New device", e.IP, e.NewMAC)
	case EventGoneDevice:
		println("[-] Device left", e.IP, e.OldMAC)
	case EventMACChange:
		println("[!] MAC changed", e.IP, e.OldMAC, "→", e.NewMAC)
	}
}

