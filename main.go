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
	port       uint
	host       string
	timeout    uint
	concurreny uint
}

func main() {
	port := flag.Uint("port", 0, "Target port")
	host := flag.String("host", "", "Target host")
	// timeout := flag.Uint("timeout", 300, "Timeout in millis")
	// concurrency := flag.Uint("concurrency", DEFAULT_WORKER_COUNT, "Number of concurrent workers")

	flag.Parse()
	// timeoutDuration := time.Millisecond * time.Duration(*timeout)

	if *host == "" {
		fmt.Println("Host is required.")
		return
	}

	if *port == 0 {
		fmt.Println("No port provided, performing a scan from 1 to 65535")
		portScanner(*host)
	} else {
		port, err := checkPort(*host, *port)
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

	// Launch worker goroutines
	for i := 0; i < DEFAULT_WORKER_COUNT; i++ {
		wg.Add(1)
		go worker(host, ports, &wg)
	}

	// Send ports to the workers
	for port := uint(1); port <= MAX_PORT; port++ {
		ports <- port
	}
	close(ports)

	// Wait for all workers to finish
	wg.Wait()
}

func worker(host string, ports chan uint, wg *sync.WaitGroup) {
	defer wg.Done()
	for port := range ports {
		if _, err := checkPort(host, port); err == nil {
			fmt.Printf("Port: %d is open\n", port)
		}
	}
}

func checkPort(host string, port uint) (uint, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, DEFAULT_TIMEOUT)
	if err != nil {
		return 0, err
	}
	conn.Close()
	return port, nil
}
