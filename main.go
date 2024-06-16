package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {
	portPtr := flag.Uint("port", 42, "Target port")
	hostPtr := flag.String("host", "", "Target host")

	flag.Parse()

	checkPort(*hostPtr, *portPtr)

}

func checkPort(host string, port uint) {
	// Concatenated address: localhost:3333
	address := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Something went wrong when openning a conection to ", address)
	}
	result := fmt.Sprintf("Fort: %d is open", port)
	fmt.Println(result)
	conn.Close()
}
