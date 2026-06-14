<!-- <p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/license-MIT-blue?style=flat-square" alt="License">
  <img src="https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-green?style=flat-square" alt="Platform">
  <img src="https://img.shields.io/github/v/release/netkit-tool/netkit?style=flat-square" alt="Release">
</p> -->

# NetKit

A multi-functional network reconnaissance and security toolkit written in Go.

面向网络管理员日常运维、安全工程师资产梳理、开发者学习网络编程。

---

## Features

- **Host Discovery** - ICMP/TCP ping, CIDR & file input, concurrent scanning
- **Port Scanning** - TCP Connect scan, Top100/Top1000/full/custom ports
- **CIDR Toolkit** - Expand, info, and aggregate CIDR blocks
- **DNS Enumeration** - Record query, subdomain brute force, DNSSEC detection *(Phase 2)*
- **Service Fingerprinting** - Banner grab, TLS fingerprint, protocol probe *(Phase 2)*
- **Web Recon** - HTTP probe, tech stack detection *(Phase 2)*
- **SSL/TLS Audit** - Certificate check, weak cipher detection, security grading *(Phase 3)*
- **Route Tracing** - MTR-style traceroute with latency stats *(Phase 3)*
- **Multi-format Output** - Table, JSON, CSV out of the box
- **Pipeline Friendly** - Chain commands via stdin/stdout *(Phase 2)*
- **Cross-platform** - Linux, macOS, Windows, zero dependency deployment

## Installation

### Pre-built Binaries

Download from [GitHub Releases](https://github.com/SiYue-ZO/netkit/releases).

### Build from Source

```bash
git clone https://github.com/SiYue-ZO/netkit.git
cd netkit
make build
```

### Go Install

```bash
go install github.com/SiYue-ZO/netkit@latest
```

## Quick Start

```bash
# Show version
netkit version

# Host discovery
netkit ping -t 192.168.1.0/24
netkit ping -t 10.0.0.1 --method icmp
netkit ping -l targets.txt

# Port scanning
netkit portscan -t 192.168.1.1
netkit portscan -t 192.168.1.1 -p 22,80,443,8080
netkit portscan -t 192.168.1.0/24 -p top1000 -c 200
netkit portscan -t 192.168.1.1 -p 1-65535 --skip-ping
netkit portscan -l hosts.txt -o results.json -f json

# CIDR toolkit
netkit cidr info 192.168.1.0/24
netkit cidr expand 10.0.0.0/29
netkit cidr aggregate 192.168.1.0/24 192.168.2.0/24

# Output formats
netkit portscan -t target -f json
netkit portscan -t target -f csv -o scan.csv
netkit portscan -t target -f table
```

## Usage

### Global Flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--config` | `-C` | | Config file path |
| `--timeout` | `-T` | 5 | Timeout in seconds |
| `--threads` | `-c` | 50 | Concurrent threads |
| `--verbose` | `-v` | false | Verbose output |
| `--output` | `-o` | | Output file path |
| `--format` | `-f` | table | Output format (table/json/csv) |
| `--no-color` | | false | Disable colored output |

### ping - Host Discovery

```bash
netkit ping -t <targets> [flags]

Flags:
  -t, --target string[]   Targets (comma-separated, supports CIDR)
  -l, --list string       Target list file
      --method string     Probe method: icmp|tcp|auto (default "auto")
      --count int         Probe count per target (default 1)
```

### portscan - Port Scanner

```bash
netkit portscan -t <targets> [flags]

Flags:
  -t, --target string[]   Targets (comma-separated, supports CIDR)
  -l, --list string       Target list file
  -p, --ports string      Port range: 80,443,1-1000,top100,top1000,full (default "top100")
      --proto string      Protocol: tcp|udp (default "tcp")
      --skip-ping         Skip host discovery, scan directly
```

### cidr - CIDR Toolkit

```bash
netkit cidr <subcommand>

Subcommands:
  expand <cidr>           Expand CIDR to IP list
  info <cidr>             Show CIDR block info
  aggregate <cidr...>     Aggregate multiple CIDRs
```

## Configuration

NetKit supports YAML config files. Place `netkit.yaml` in `~/.netkit/` or the current directory.

```yaml
# ~/.netkit/netkit.yaml
timeout: 5
threads: 50
format: table
verbose: false

portscan:
  default_ports: top100
  proto: tcp
  skip_ping: false

ping:
  method: auto
  count: 1
```

## Project Structure

```
netkit/
├── main.go                     # Entry point
├── cmd/                        # Cobra commands
│   ├── root.go                 # Root command + global flags + Viper
│   ├── version.go              # Version info
│   ├── ping.go                 # Ping command
│   ├── portscan.go             # Port scan command
│   └── cidr.go                 # CIDR command
├── internal/                   # Internal packages
│   ├── ping/                   # Host discovery engine
│   ├── scanner/                # Port scanning engine
│   ├── cidrutil/               # CIDR utilities
│   └── output/                 # Output formatting
├── pkg/                        # Public packages
│   ├── types/                  # Shared type definitions
│   └── pool/                   # Goroutine pool
├── configs/                    # Default configs
├── .goreleaser.yml             # Cross-platform release
├── .golangci.yml               # Linter config
├── .github/workflows/ci.yml   # CI pipeline
└── Makefile                    # Build automation
```

## Development

```bash
# Build
make build

# Run
make run ARGS="version"

# Test
make test

# Lint
make lint

# Cross-compile all platforms
make cross-build

# Release snapshot
make snapshot
```

## Roadmap

See [ROADMAP.md](ROADMAP.md) for the full development plan.

| Phase | Version | Focus |
|---|---|---|
| Phase 1 | v0.1.0 | CLI framework + ping + portscan + cidr |
| Phase 2 | v0.2.0 | dns + fingerprint + web + whois + pipeline |
| Phase 3 | v0.3.0 | ssl + trace + SYN/UDP scan |
| Phase 4 | v0.4.0 | Reports + shell completion + test coverage |
| Phase 5 | v0.5.0+ | Fingerprint DB + vuln detection + API |

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This tool is intended for authorized network security assessment and educational purposes only. Users are responsible for ensuring they have proper authorization before scanning any networks or systems. The developers assume no liability and are not responsible for any misuse or damage.
