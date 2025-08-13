package ip

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type ManualDetector struct{}

func NewManualDetector() *ManualDetector {
	return &ManualDetector{}
}

func (m *ManualDetector) GetIPv4() (net.IP, error) {
	return m.promptForIP("IPv4")
}

func (m *ManualDetector) GetIPv6() (net.IP, error) {
	return m.promptForIP("IPv6")
}

func (m *ManualDetector) promptForIP(ipType string) (net.IP, error) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("Enter %s address: ", ipType)
		if !scanner.Scan() {
			return nil, fmt.Errorf("failed to read input")
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			return nil, fmt.Errorf("no IP address provided")
		}

		ip := net.ParseIP(input)
		if ip == nil {
			fmt.Printf("Invalid IP address format: %s\n", input)
			continue
		}

		if ipType == "IPv4" && ip.To4() == nil {
			fmt.Printf("Please enter a valid IPv4 address, got: %s\n", input)
			continue
		}

		if ipType == "IPv6" && ip.To4() != nil {
			fmt.Printf("Please enter a valid IPv6 address, got: %s\n", input)
			continue
		}

		return ip, nil
	}
}

func (m *ManualDetector) Name() string {
	return "Manual Input"
}
