package watcher

import (
	"bytes"
	"net"
	"os/exec"
	"strings"
	"time"
)

func ReadARP(ignoredIPs map[string]struct{}) *State {
	out, err := exec.Command("arp", "-an").Output()
	if err != nil {
		return NewState()
	}

	state := NewState()
	now := time.Now().Unix()

	lines := bytes.Split(out, []byte("\n"))
	for _, l := range lines {
		fields := strings.Fields(string(l))
		if len(fields) < 4 {
			continue
		}

		ip := strings.Trim(fields[1], "()")
		mac := fields[3]

		// --- FILTERS ---

		// Ignore broadcast MAC
		if mac == "ff:ff:ff:ff:ff:ff" {
			continue
		}

		// Ignore invalid MACs
		if _, err := net.ParseMAC(mac); err != nil {
			continue
		}

		// Ignore broadcast IPs
		if strings.HasSuffix(ip, ".255") {
			continue
		}

		// Ignore non-private IP space 
		parsedIP := net.ParseIP(ip)
		if parsedIP == nil || !parsedIP.IsPrivate() {
			continue
		}

		if ignoredIPs != nil {
			if _, ok := ignoredIPs[ip]; ok {
				continue
			}
		}

		state.Devices[ip] = Device{
			IP:       ip,
			MAC:      mac,
			LastSeen: now,
		}
	}

	return state
}
