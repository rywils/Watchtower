package watcher

import (
	"fmt"
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

func isIgnorable(ip, mac string, ignoredIPs map[string]struct{}) bool {
	// broadcast MAC
	if mac == "ff:ff:ff:ff:ff:ff" {
		return true
	}

	if ignoredIPs != nil {
		if _, ok := ignoredIPs[ip]; ok {
			return true
		}
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

func Diff(prev, curr *State, ignoredIPs map[string]struct{}) []Event {
	now := time.Now().Unix()
	events := []Event{}

	// new / changed
	for ip, d := range curr.Devices {
		if isIgnorable(ip, d.MAC, ignoredIPs) {
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
		if isIgnorable(ip, old.MAC, ignoredIPs) {
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
	ipHost := formatIPHostname(e.IP)
	switch e.Type {
	case EventNewDevice:
		fmt.Printf("[+] %s has joined the network.\n", ipHost)
	case EventGoneDevice:
		fmt.Printf("[-] %s has left the network.\n", ipHost)
	case EventMACChange:
		fmt.Printf("[!] %s MAC address has changed to %s (was %s).\n", ipHost, e.NewMAC, e.OldMAC)
	}
}

func formatIPHostname(ip string) string {
	hostname := hostnameForIP(ip)
	return fmt.Sprintf("%s (%s)", ip, hostname)
}

func hostnameForIP(ip string) string {
	names, err := net.LookupAddr(ip)
	if err != nil || len(names) == 0 {
		return "unknown"
	}
	return strings.TrimSuffix(names[0], ".")
}
