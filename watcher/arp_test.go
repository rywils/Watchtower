package watcher

import "testing"

func TestParseARPFiltersAndParsesDevices(t *testing.T) {
	out := []byte(`? (192.168.1.10) at aa:bb:cc:dd:ee:ff on en0 ifscope [ethernet]
? (192.168.1.255) at ff:ff:ff:ff:ff:ff on en0 ifscope [ethernet]
? (8.8.8.8) at 11:22:33:44:55:66 on en0 ifscope [ethernet]
? (192.168.1.1) at 00:11:22:33:44:55 on en0 ifscope [ethernet]
`)
	ignored := map[string]struct{}{
		"192.168.1.1": {},
	}

	state := ParseARP(out, ignored)
	if len(state.Devices) != 1 {
		t.Fatalf("expected exactly 1 device, got %d", len(state.Devices))
	}
	if _, ok := state.Devices["192.168.1.10"]; !ok {
		t.Fatalf("expected 192.168.1.10 to be present")
	}
}
