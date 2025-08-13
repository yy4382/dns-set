package ip

import "net"

type IPDetector interface {
	GetIPv4() (net.IP, error)
	GetIPv6() (net.IP, error)
	Name() string
}
