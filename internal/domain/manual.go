package domain

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ManualSource struct{}

func NewManualSource() *ManualSource {
	return &ManualSource{}
}

func (m *ManualSource) GetDomains() ([]string, error) {
	fmt.Print("Enter domains (one per line, empty line to finish):\n")

	scanner := bufio.NewScanner(os.Stdin)
	var domains []string

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			break
		}

		if isValidDomain(line) {
			domains = append(domains, line)
		} else {
			fmt.Printf("Invalid domain format: %s\n", line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	if len(domains) == 0 {
		return nil, fmt.Errorf("no valid domains provided")
	}

	return domains, nil
}

func (m *ManualSource) Name() string {
	return "Manual Input"
}

func isValidDomain(domain string) bool {
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}

	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return false
	}

	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return false
	}

	for i, part := range parts {
		if len(part) == 0 || len(part) > 63 {
			return false
		}

		for j, c := range part {
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-') {
				return false
			}
			if c == '-' && (j == 0 || j == len(part)-1) {
				return false
			}
		}

		if i == len(parts)-1 {
			if len(part) < 2 {
				return false
			}
		}
	}

	return true
}
