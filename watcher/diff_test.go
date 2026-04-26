package watcher

import "testing"

func TestDiffDetectsNewGoneAndMACChange(t *testing.T) {
	prev := &State{
		Devices: map[string]Device{
			"192.168.1.20": {IP: "192.168.1.20", MAC: "aa:aa:aa:aa:aa:aa"},
			"192.168.1.30": {IP: "192.168.1.30", MAC: "bb:bb:bb:bb:bb:bb"},
		},
	}
	curr := &State{
		Devices: map[string]Device{
			"192.168.1.20": {IP: "192.168.1.20", MAC: "cc:cc:cc:cc:cc:cc"},
			"192.168.1.40": {IP: "192.168.1.40", MAC: "dd:dd:dd:dd:dd:dd"},
		},
	}

	events := Diff(prev, curr, nil)
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
}

func TestIsIgnorableUsesExplicitGatewayListNotDotOne(t *testing.T) {
	ignored := map[string]struct{}{
		"192.168.1.1": {},
	}
	if !isIgnorable("192.168.1.1", "00:11:22:33:44:55", ignored) {
		t.Fatalf("expected explicit gateway IP to be ignored")
	}
	if isIgnorable("192.168.1.21", "00:11:22:33:44:55", ignored) {
		t.Fatalf("did not expect normal private host to be ignored")
	}
}
