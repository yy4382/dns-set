package ip

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type APIDetector struct {
	client *http.Client
}

func NewAPIDetector() *APIDetector {
	return &APIDetector{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (a *APIDetector) GetIPv4() (net.IP, error) {
	resp, err := a.client.Get("https://api-ipv4.ip.sb/ip")
	if err != nil {
		return nil, fmt.Errorf("failed to get IPv4 from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	ipStr := strings.TrimSpace(string(body))
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address received: %s", ipStr)
	}

	if ip.To4() == nil {
		return nil, fmt.Errorf("received non-IPv4 address: %s", ipStr)
	}

	return ip, nil
}

func (a *APIDetector) GetIPv6() (net.IP, error) {
	resp, err := a.client.Get("https://api-ipv6.ip.sb/ip")
	if err != nil {
		return nil, fmt.Errorf("failed to get IPv6 from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	ipStr := strings.TrimSpace(string(body))
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address received: %s", ipStr)
	}

	if ip.To4() != nil {
		return nil, fmt.Errorf("received IPv4 address instead of IPv6: %s", ipStr)
	}

	return ip, nil
}

func (a *APIDetector) Name() string {
	return "External API (ip.sb)"
}