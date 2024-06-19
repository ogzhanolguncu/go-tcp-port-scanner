package port_parser

import (
	"reflect"
	"testing"
)

func TestPortParse(t *testing.T) {
	tests := []struct {
		input       string
		expected    Port
		expectError bool
	}{
		{"443", Port{Type: SinglePort, SinglePort: 443}, false},
		{"443,80", Port{Type: MultiplePorts, PortList: []int{443, 80}}, false},
		{"20-25", Port{Type: MultiplePorts, PortList: []int{20, 21, 22, 23, 24, 25}}, false},
		{"0", Port{Type: FullScan}, false},
		{"65535", Port{Type: SinglePort, SinglePort: 65535}, false},
		{"443,invalid", Port{}, true},
		{"invalid", Port{}, true},
		{"20-20", Port{Type: MultiplePorts, PortList: []int{20}}, false},
		{"-20", Port{}, true},
		{"20-", Port{}, true},
		{"20-10", Port{Type: MultiplePorts, PortList: []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}}, false},
	}

	for _, test := range tests {
		var p Port
		err := p.Parse(test.input)
		if test.expectError {
			if err == nil {
				t.Errorf("expected error for input %s, but got none", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for input %s: %v", test.input, err)
			}
			if !reflect.DeepEqual(p, test.expected) {
				t.Errorf("for input %s, expected %v, got %v", test.input, test.expected, p)
			}
		}
	}
}
