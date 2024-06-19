package port_parser

import (
	"fmt"
	"strconv"
	"strings"
)

type PortType int

const (
	SinglePort PortType = iota
	MultiplePorts
	FullScan
)

type Port struct {
	Type       PortType
	SinglePort int
	PortList   []int
}

func (p *Port) String() string {
	switch p.Type {
	case SinglePort:
		return fmt.Sprintf("%d", p.SinglePort)
	case MultiplePorts:
		return fmt.Sprintf("%v", p.PortList)
	case FullScan:
		return "0"
	default:
		return ""
	}
}

func (p *Port) Parse(value string) error {
	switch {
	case value == "0":
		p.Type = FullScan
		return nil

	case strings.Contains(value, ","):
		ports := strings.Split(value, ",")
		var portList []int

		for _, portStr := range ports {
			port, err := strconv.Atoi(portStr)
			if err != nil {
				return fmt.Errorf("invalid port number: %s", portStr)
			}
			portList = append(portList, port)
		}
		p.Type = MultiplePorts
		p.PortList = portList
		return nil

	case strings.Contains(value, "-"):
		ports := strings.Split(value, "-")
		if len(ports) != 2 {
			return fmt.Errorf("invalid port range: %s", value)
		}

		startPort, err := strconv.Atoi(ports[0])
		if err != nil {
			return fmt.Errorf("invalid start port: %s", ports[0])
		}

		endPort, err := strconv.Atoi(ports[1])
		if err != nil {
			return fmt.Errorf("invalid end port: %s", ports[1])
		}

		if startPort == endPort {
			p.PortList = []int{startPort}
		} else {
			for i := min(startPort, endPort); i <= max(startPort, endPort); i++ {
				p.PortList = append(p.PortList, i)
			}
		}

		p.Type = MultiplePorts
		return nil

	default:
		port, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid port: %s", value)
		}
		p.Type = SinglePort
		p.SinglePort = port
		return nil
	}
}
