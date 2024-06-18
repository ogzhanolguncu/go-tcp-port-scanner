package main

import (
	"flag"
	"fmt"
	"math"
	"net"
	"sync"
	"time"
)

const DEFAULT_TIMEOUT = time.Millisecond * 300
const MAX_PORT = math.MaxUint16
const DEFAULT_WORKER_COUNT = 100

type Payload struct {
	port        uint
	host        string
	timeout     time.Duration
	concurrency uint
}

var (
	port        = flag.Uint("port", 0, "Target port")
	host        = flag.String("host", "", "Target host")
	timeout     = flag.Int("timeout", 300, "Timeout in millis")
	concurrency = flag.Uint("concurrency", DEFAULT_WORKER_COUNT, "Number of concurrent workers")
)

func parseCLIArgs() (Payload, error) {
	flag.Parse()

	if *host == "" {
		return Payload{}, fmt.Errorf("Host is required!")
	}

	return Payload{
		port:        *port,
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

	if *port == 0 {
		fmt.Println("No port provided, performing a scan from 1 to 65535")
		portScanner(args.host)
	} else {
		port, err := checkPort(args)
		if err == nil {
			fmt.Printf("Port: %d is open\n", port)
		} else {
			fmt.Println("Error:", err)
		}
	}
}

func portScanner(host string) {
	var wg sync.WaitGroup
	ports := make(chan uint, MAX_PORT)

	for i := 0; i < DEFAULT_WORKER_COUNT; i++ {
		wg.Add(1)
		go worker(host, ports, &wg)
	}

	for port := uint(1); port <= MAX_PORT; port++ {
		ports <- port
	}
	close(ports)

	wg.Wait()
}

func worker(host string, ports chan uint, wg *sync.WaitGroup) {
	defer wg.Done()
	for port := range ports {
		if _, err := checkPort(Payload{
			port: port,
			host: host,
		}); err == nil {
			fmt.Printf("Port: %d is open\n", port)
		}
	}
}

func checkPort(args Payload) (uint, error) {
	address := fmt.Sprintf("%s:%d", args.host, args.port)
	conn, err := net.DialTimeout("tcp", address, DEFAULT_TIMEOUT)
	if err != nil {
		return 0, err
	}
	conn.Close()
	return args.port, nil
}
