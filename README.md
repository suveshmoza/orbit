# Orbit

Find the best DNS server for your network.

Orbit is a simple terminal tool that tests popular DNS resolvers — like Cloudflare, Google, and Quad9 — and shows you which one is best for your network.

<img width="1756" height="1080" alt="orbit-demo" src="https://github.com/user-attachments/assets/75738d39-fc29-4cfe-a278-ae2e76ac3b94" />


## Requirements

- Go 1.26+

## Usage

```bash
go run .
```

Build a binary:

```bash
go build -o orbit .
./orbit
```

Press `q` or `Ctrl+C` to quit at any time. In-flight queries are cancelled cleanly.

## Configuration

Orbit reads `config.json` from the working directory. Example:

```json
{
  "config": {
    "samples": 5,
    "timeout_ms": 2000,
    "test_domains": ["www.google.com", "amazon.com"]
  },
  "servers": [
    { "name": "Cloudflare", "address": "1.1.1.1", "port": 53 },
    { "name": "Google", "address": "8.8.8.8", "port": 53 }
  ]
}
```

| Field | Description |
|-------|-------------|
| `samples` | Queries per domain, per server |
| `timeout_ms` | Per-query timeout in milliseconds |
| `test_domains` | Domains to resolve |
| `servers` | Resolvers to compare (`name`, `address`, `port`) |

The default config benchmarks Cloudflare, Google, OpenDNS, and Quad9 against popular domains.

## Output

While running, Orbit shows the current query, a progress bar, and a live results table. Each row reports how many queries passed and the mean, median, and lowest RTT. When the benchmark finishes, the server with the best median is highlighted in green.

## How it works

1. Load `config.json`
2. Spawn a goroutine per DNS server
3. For each server, query every domain `samples` times over UDP
4. Stream results to the TUI and update stats in real time
5. Pick the fastest server by median RTT

## Todo

- [ ] Add CLI flags for samples, timeout and opening config file in editor
- [ ] Add more DNS servers and domains to the default config
- [ ] Better error handling and logging
- [ ] Write tests
