## Watchtower

Watchtower is a passive local network watcher.
It monitors your LAN and reports meaningful changes - best when used on untrusted networks to monitor for evil twin attacks, etc.

No scanning. No probing. No noise.

### What it does

* Detects new devices
* Detects devices leaving
* Detects MAC address changes
* Persists a baseline across runs
* Ignores broadcast, detected gateway, and junk entries
* Supports human-readable or JSON-line output
* Supports custom polling interval and custom ignored IPs

---

### Usage

```bash
go build
sudo ./watchtower
```

ARP access requires elevated privileges.

### Flags

```bash
./watchtower -interval 5s -ignore "192.168.1.2,192.168.1.3"
./watchtower -no-dns
./watchtower -json
```

* `-interval` sets poll cadence (`time.Duration` format).
* `-no-dns` disables reverse-DNS lookups for faster, non-blocking output.
* `-json` emits one JSON event per line.
* `-ignore` adds a comma-separated list of IPs to ignore.

### Output example

```less
[*] Watchtower running
[+] 192.168.0.114 (my-phone.local) has joined the network.
[!] 192.168.0.42 (unknown) MAC address has changed to aa:bb:cc:dd:ee:ff (was 00:11:22:33:44:55).
```

JSON mode:

```json
{"Type":"new_device","IP":"192.168.0.114","OldMAC":"","NewMAC":"ca:3f:33:7c:a4:b4","Timestamp":1714100000}
```
---

### State

Watchtower stores state locally:

```bash
~/.watchtower/state.json
```

Delete this file to reset the baseline.

If state cannot be read or parsed, Watchtower now exits with a clear error instead of silently resetting baseline.

