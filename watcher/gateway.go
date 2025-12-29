package watcher

import (
	"net"
	"os/exec"
	"strings"
)

func DefaultGatewayIPs() map[string]struct{} {
	gateways := map[string]struct{}{}

	if addGatewayFromIPRoute(gateways) {
		return gateways
	}
	if addGatewayFromRouteGet(gateways) {
		return gateways
	}
	_ = addGatewayFromNetstat(gateways)

	return gateways
}

func addGatewayFromIPRoute(gateways map[string]struct{}) bool {
	out, err := exec.Command("ip", "route", "show", "default").Output()
	if err != nil {
		return false
	}

	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		for i := 0; i+1 < len(fields); i++ {
			if fields[i] == "via" {
				if addGateway(gateways, fields[i+1]) {
					return true
				}
			}
		}
	}

	return len(gateways) > 0
}

func addGatewayFromRouteGet(gateways map[string]struct{}) bool {
	out, err := exec.Command("route", "-n", "get", "default").Output()
	if err != nil {
		return false
	}

	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "gateway:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 && addGateway(gateways, fields[1]) {
				return true
			}
		}
	}

	return len(gateways) > 0
}

func addGatewayFromNetstat(gateways map[string]struct{}) bool {
	out, err := exec.Command("netstat", "-rn").Output()
	if err != nil {
		return false
	}

	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		if fields[0] == "default" || fields[0] == "0.0.0.0" {
			if addGateway(gateways, fields[1]) {
				return true
			}
		}
	}

	return len(gateways) > 0
}

func addGateway(gateways map[string]struct{}, candidate string) bool {
	ip := net.ParseIP(candidate)
	if ip == nil {
		return false
	}
	ip = ip.To4()
	if ip == nil {
		return false
	}

	gateways[ip.String()] = struct{}{}
	return true
}
