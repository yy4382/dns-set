package ip

import (
	"fmt"
	"net"
)

type InterfaceDetector struct{}

func NewInterfaceDetector() *InterfaceDetector {
	return &InterfaceDetector{}
}

func (i *InterfaceDetector) GetIPv4() (net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			ip := ipNet.IP
			if ip.IsLoopback() || ip.IsPrivate() {
				continue
			}

			if ip.To4() != nil {
				return ip, nil
			}
		}
	}

	return nil, fmt.Errorf("no public IPv4 address found")
}

func (i *InterfaceDetector) GetIPv6() (net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			ip := ipNet.IP
			if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() {
				continue
			}

			if ip.To4() == nil {
				return ip, nil
			}
		}
	}

	return nil, fmt.Errorf("no public IPv6 address found")
}

func (i *InterfaceDetector) Name() string {
	return "Network Interface"
}
