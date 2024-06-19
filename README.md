# TCP Port Scanner

Simple TCP Port scanner written in Go.

## Features

- TCP scan for detecting open ports
- Option to specify a range of ports or a single port

## Usage

```bash
Usage: go-port-scanner [OPTIONS] --host <IP> --port <NUMBER> --concurrency <NUMBER>

Options:
  -h, --host             Target host
  -p, --port             Target port or port range
  -c  --concurrency      Number of concurrent requests
  -t  --timeout          Timeout for TCP calls


 go-port-scanner --host 45.33.32.156
```
