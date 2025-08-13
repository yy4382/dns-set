package domain

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type CaddyfileSource struct {
	path string
}

func NewCaddyfileSource(path string) *CaddyfileSource {
	return &CaddyfileSource{path: path}
}

func (c *CaddyfileSource) GetDomains() ([]string, error) {
	file, err := os.Open(c.path)
	if err != nil {
		return nil, fmt.Errorf("failed to open Caddyfile at %s: %w", c.path, err)
	}
	defer file.Close()

	domains := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	
	domainRegex := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		if strings.Contains(line, "{") {
			line = strings.Split(line, "{")[0]
			line = strings.TrimSpace(line)
		}
		
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		
		for _, part := range parts {
			part = strings.TrimSpace(part)
			
			if part == "" || strings.HasPrefix(part, ":") {
				continue
			}
			
			if strings.Contains(part, ":") {
				part = strings.Split(part, ":")[0]
			}
			
			if domainRegex.MatchString(part) && !strings.Contains(part, "*") && !isDirective(part) {
				domains[part] = true
			}
		}
	}
	
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading Caddyfile: %w", err)
	}
	
	result := make([]string, 0, len(domains))
	for domain := range domains {
		result = append(result, domain)
	}
	
	if len(result) == 0 {
		return nil, fmt.Errorf("no valid domains found in Caddyfile")
	}
	
	return result, nil
}

func (c *CaddyfileSource) Name() string {
	return fmt.Sprintf("Caddyfile (%s)", c.path)
}

func isDirective(word string) bool {
	commonDirectives := map[string]bool{
		"root":            true,
		"respond":         true,
		"reverse_proxy":   true,
		"proxy":          true,
		"file_server":    true,
		"encode":         true,
		"header":         true,
		"rewrite":        true,
		"uri":            true,
		"try_files":      true,
		"basicauth":      true,
		"request_header": true,
		"import":         true,
		"log":            true,
		"tls":            true,
		"backend":        true,
	}
	
	if isDirective, exists := commonDirectives[word]; exists {
		return isDirective
	}
	
	return false
}