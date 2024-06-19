package main

import (
	"flag"
	"fmt"
	"log/slog"
	"math"
	"net"
	"sync"
	"time"

	"github.com/ogzhanolguncu/go-port-scanner/port_parser"
)

const DEFAULT_TIMEOUT = time.Millisecond * 300
const MAX_PORT = math.MaxUint16
const DEFAULT_WORKER_COUNT = 100

type Payload struct {
	port        port_parser.Port
	host        string
	timeout     time.Duration
	concurrency uint
}

var (
	/*
		--port=443
		--port=443,80
		--port=20-443
	*/
	port = flag.String("port", "0", "Target port --port=443\nSelected ports --port=443,80 (splitted by ',')\nRange --port=20-443 (splitted by '-') or port range")
	/*
		--host=127.0.0.1
	*/
	host        = flag.String("host", "", "Target host")
	timeout     = flag.Int("timeout", 300, "Timeout in millis")
	concurrency = flag.Uint("concurrency", DEFAULT_WORKER_COUNT, "Number of concurrent workers")
)

func parseCLIArgs() (Payload, error) {
	flag.Parse()

	var p port_parser.Port
	err := p.Parse(*port)
	if err != nil {
		return Payload{}, err
	}

	if *host == "" {
		return Payload{}, fmt.Errorf("host is required")
	}

	return Payload{
		port:        p,
		host:        *host,
		timeout:     time.Millisecond * time.Duration(*timeout),
		concurrency: *concurrency,
	}, nil
}

func main() {
	args, err := parseCLIArgs()
	if err != nil {
		fmt.Println("Error parsing arguments:", err)
		flag.Usage()
		return
	}

	switch args.port.Type {
	case port_parser.FullScan:
		fmt.Println("No port provided, performing a scan from 1 to 65535")
		slog.Info("Running a full sweep")
		portListScanner(args.host, nil)

	case port_parser.MultiplePorts:
		slog.Info("Running a sweep for given ports")
		portListScanner(args.host, args.port.PortList)

	case port_parser.SinglePort:
		slog.Info("Running a sweep for single port")
		if args.port.Type == port_parser.SinglePort {
			port, err := checkPort(args)
			if err == nil {
				fmt.Printf("Port: %d is open\n", port)
			} else {
				fmt.Println("Error:", err)
			}
		}
	}

}

func portListScanner(host string, portList []int) {
	var wg sync.WaitGroup
	ports := make(chan uint, MAX_PORT)

	for i := 0; i < DEFAULT_WORKER_COUNT; i++ {
		wg.Add(1)
		go worker(host, ports, &wg)
	}

	if portList != nil {
		for _, port := range portList {
			ports <- uint(port)
		}
	} else {
		for port := uint(1); port <= MAX_PORT; port++ {
			ports <- port
		}
	}
	close(ports)

	wg.Wait()
}

func worker(host string, ports chan uint, wg *sync.WaitGroup) {
	defer wg.Done()
	for port := range ports {
		if _, err := checkPort(Payload{
			port: port_parser.Port{SinglePort: int(port)},
			host: host,
		}); err == nil {
			fmt.Printf("Port: %d is open\n", port)
		}
	}
}

func checkPort(args Payload) (uint, error) {
	address := fmt.Sprintf("%s:%d", args.host, args.port.SinglePort)
	conn, err := net.DialTimeout("tcp", address, DEFAULT_TIMEOUT)
	if err != nil {
		return 0, err
	}
	conn.Close()
	return uint(args.port.SinglePort), nil
}
