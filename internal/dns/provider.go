package dns

import "net"

type RecordType string

const (
	RecordTypeA    RecordType = "A"
	RecordTypeAAAA RecordType = "AAAA"
)

type Record struct {
	ID      string
	Name    string
	Type    RecordType
	Content string
	TTL     int
	Proxied bool
}

type DNSProvider interface {
	UpdateRecord(domain string, recordType RecordType, ip net.IP, ttl *int, proxied bool) error
	ListRecords(domain string) ([]Record, error)
	Name() string
}
