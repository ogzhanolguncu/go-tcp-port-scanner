package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

const TIMEOUT = 5
const MAX_PORT = 65535

func main() {
	portPtr := flag.Uint("port", 0, "Target port")
	hostPtr := flag.String("host", "", "Target host")

	flag.Parse()

	if *hostPtr == "" {
		fmt.Println("Host is required.")
	}

	if *portPtr == 0 {
		fmt.Println("No port provided, performing a scan from 1 to 65535")

		for port := uint(1); port <= MAX_PORT; port++ {
			port, err := checkPort(*hostPtr, port)
			if err == nil {
				fmt.Printf("Port: %d is open\n", port)
			}
		}
	} else {
		port, err := checkPort(*hostPtr, *portPtr)
		if err == nil {
			fmt.Printf("Port: %d is open\n", port)
		} else {
			fmt.Println("Error:", err)
		}
	}

}

/*
Checks open ports for given host and port. If port is missing does a vanilla scan from 1-65536
*/
func checkPort(host string, port uint) (uint, error) {
	// Concatenated address: localhost:3333
	address := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.DialTimeout("tcp", address, time.Second*TIMEOUT)

	if err != nil {
		return 0, err
	}

	conn.Close()
	return port, nil
}
