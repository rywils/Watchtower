## Watchtower

Watchtower is a passive local network watcher.
It monitors your LAN and reports meaningful changes — best when used on untrusted networks to monitor for evil twin attacks, etc

No scanning. No probing. No noise.

### What it does

* Detects new devices
* Detects devices leaving
* Detects MAC address changes
* Persists a baseline across runs
* Ignores broadcast, gateway, and junk entries
* It only tells you what actually matters.

---

### Usage

```bash
go build
sudo ./watchtower
```

ARP access requires elevated privileges.

### Output example

```less
[*] Watchtower running
[+] New device 192.168.0.114 ca:3f:33:7c:a4:b4
[!] MAC changed 192.168.0.42 00:11:22:33:44:55 → aa:bb:cc:dd:ee:ff
```
---

### State

Watchtower stores state locally:

```bash
~/.watchtower/state.json
```

Delete this file to reset the baseline.

