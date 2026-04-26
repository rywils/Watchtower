package watcher

import (
	"encoding/json"
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

func (e Event) JSON() string {
	b, err := json.Marshal(e)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func (e Event) Text(resolveDNS bool) string {
	ipHost := formatIPHostname(e.IP, resolveDNS)
	switch e.Type {
	case EventNewDevice:
		return fmt.Sprintf("[+] %s has joined the network.", ipHost)
	case EventGoneDevice:
		return fmt.Sprintf("[-] %s has left the network.", ipHost)
	case EventMACChange:
		return fmt.Sprintf("[!] %s MAC address has changed to %s (was %s).", ipHost, e.NewMAC, e.OldMAC)
	}
	return fmt.Sprintf("[?] Unknown event for %s", e.IP)
}

func formatIPHostname(ip string, resolveDNS bool) string {
	if !resolveDNS {
		return ip
	}
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
