package ui

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/yy4382/dns-set/internal/config"
	"github.com/yy4382/dns-set/internal/dns"
	"github.com/yy4382/dns-set/internal/domain"
	"github.com/yy4382/dns-set/internal/ip"
	"golang.org/x/term"
)

type CLI struct {
	config   *config.Config
	provider dns.DNSProvider
	scanner  *bufio.Scanner
}

func NewCLI(cfg *config.Config, provider dns.DNSProvider) *CLI {
	return &CLI{
		config:   cfg,
		provider: provider,
		scanner:  bufio.NewScanner(os.Stdin),
	}
}

func (c *CLI) Run() error {
	fmt.Printf("=== DNS Setter - %s Provider ===\n\n", c.provider.Name())

	domainSource, err := c.selectDomainSource()
	if err != nil {
		return fmt.Errorf("failed to select domain source: %w", err)
	}

	domains, err := domainSource.GetDomains()
	if err != nil {
		return fmt.Errorf("failed to get domains: %w", err)
	}

	selectedDomains, err := c.selectDomains(domains)
	if err != nil {
		return fmt.Errorf("failed to select domains: %w", err)
	}

	ipDetector, err := c.selectIPDetector()
	if err != nil {
		return fmt.Errorf("failed to select IP detector: %w", err)
	}

	recordTypes, err := c.selectRecordTypes()
	if err != nil {
		return fmt.Errorf("failed to select record types: %w", err)
	}

	proxied, err := c.selectProxyStatus()
	if err != nil {
		return fmt.Errorf("failed to select proxy status: %w", err)
	}

	for _, recordType := range recordTypes {
		var ip net.IP
		var err error

		if recordType == dns.RecordTypeA {
			ip, err = ipDetector.GetIPv4()
		} else {
			ip, err = ipDetector.GetIPv6()
		}

		if err != nil {
			fmt.Printf("Failed to get %s address: %v\n", recordType, err)
			continue
		}

		fmt.Printf("\nDetected %s address: %s\n", recordType, ip.String())

		for _, domain := range selectedDomains {
			fmt.Printf("Updating %s record for %s...\n", recordType, domain)

			err = c.provider.UpdateRecord(domain, recordType, ip, c.config.Preferences.DefaultTTL, proxied)
			if err != nil {
				fmt.Printf("Failed to update %s record for %s: %v\n", recordType, domain, err)
			} else {
				proxyStatus := "DNS only"
				if proxied {
					proxyStatus = "Proxied"
				}
				fmt.Printf("Successfully updated %s record for %s (%s)\n", recordType, domain, proxyStatus)
			}
		}
	}

	fmt.Println("\nDNS update completed!")
	return nil
}

func (c *CLI) selectDomainSource() (domain.DomainSource, error) {
	fmt.Println("Select domain source:")
	fmt.Println("1. Manual input")
	fmt.Println("2. Caddyfile")

	choice, err := c.promptChoice("Enter choice (1-2): ", 1, 2)
	if err != nil {
		return nil, err
	}

	switch choice {
	case 1:
		return domain.NewManualSource(), nil
	case 2:
		caddyfilePath, err := c.promptCaddyfilePath(c.config.Preferences.CaddyfilePath)
		if err != nil {
			return nil, err
		}
		return domain.NewCaddyfileSource(caddyfilePath), nil
	default:
		return nil, fmt.Errorf("invalid choice")
	}
}

func (c *CLI) selectDomains(domains []string) ([]string, error) {
	if len(domains) == 1 {
		fmt.Printf("Found domain: %s\n", domains[0])
		return domains, nil
	}

	fmt.Printf("\nFound %d domains:\n", len(domains))
	for i, domain := range domains {
		fmt.Printf("%d. %s\n", i+1, domain)
	}

	fmt.Println("Select domains to update (comma-separated numbers, or 'all'):")
	c.scanner.Scan()
	input := strings.TrimSpace(c.scanner.Text())

	if input == "all" {
		return domains, nil
	}

	var selected []string
	parts := strings.Split(input, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		index, err := strconv.Atoi(part)
		if err != nil || index < 1 || index > len(domains) {
			return nil, fmt.Errorf("invalid selection: %s", part)
		}
		selected = append(selected, domains[index-1])
	}

	if len(selected) == 0 {
		return nil, fmt.Errorf("no domains selected")
	}

	return selected, nil
}

func (c *CLI) selectIPDetector() (ip.IPDetector, error) {
	fmt.Println("\nSelect IP detection method:")
	fmt.Println("1. Network interface")
	fmt.Println("2. External API (ip.sb)")
	fmt.Println("3. Manual input")

	choice, err := c.promptChoice("Enter choice (1-3): ", 1, 3)
	if err != nil {
		return nil, err
	}

	switch choice {
	case 1:
		return ip.NewInterfaceDetector(), nil
	case 2:
		return ip.NewAPIDetector(), nil
	case 3:
		return ip.NewManualDetector(), nil
	default:
		return nil, fmt.Errorf("invalid choice")
	}
}

func (c *CLI) selectRecordTypes() ([]dns.RecordType, error) {
	fmt.Println("\nSelect record types to update:")
	fmt.Println("1. IPv4 (A) only")
	fmt.Println("2. IPv6 (AAAA) only")
	fmt.Println("3. Both IPv4 and IPv6")

	choice, err := c.promptChoice("Enter choice (1-3): ", 1, 3)
	if err != nil {
		return nil, err
	}

	switch choice {
	case 1:
		return []dns.RecordType{dns.RecordTypeA}, nil
	case 2:
		return []dns.RecordType{dns.RecordTypeAAAA}, nil
	case 3:
		return []dns.RecordType{dns.RecordTypeA, dns.RecordTypeAAAA}, nil
	default:
		return nil, fmt.Errorf("invalid choice")
	}
}

func (c *CLI) selectProxyStatus() (bool, error) {
	fmt.Println("\nSelect Cloudflare proxy status:")
	fmt.Println("1. DNS only (grey cloud)")
	fmt.Println("2. Proxied (yellow cloud)")

	choice, err := c.promptChoice("Enter choice (1-2): ", 1, 2)
	if err != nil {
		return false, err
	}

	return choice == 2, nil
}

func (c *CLI) promptChoice(prompt string, min, max int) (int, error) {
	for {
		fmt.Print(prompt)
		if !c.scanner.Scan() {
			return 0, fmt.Errorf("failed to read input")
		}

		input := strings.TrimSpace(c.scanner.Text())
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("Please enter a number between %d and %d\n", min, max)
			continue
		}

		if choice < min || choice > max {
			fmt.Printf("Please enter a number between %d and %d\n", min, max)
			continue
		}

		return choice, nil
	}
}

func (c *CLI) promptCaddyfilePath(defaultPath string) (string, error) {
	if _, err := os.Stat(defaultPath); err == nil {
		return defaultPath, nil
	}

	fmt.Printf("Caddyfile not found at %s\n", defaultPath)
	fmt.Print("Please enter Caddyfile path (absolute or relative to current directory): ")

	if !c.scanner.Scan() {
		return "", fmt.Errorf("failed to read input")
	}

	userPath := strings.TrimSpace(c.scanner.Text())
	if userPath == "" {
		return "", fmt.Errorf("no path provided")
	}

	var resolvedPath string
	if filepath.IsAbs(userPath) {
		resolvedPath = userPath
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current working directory: %w", err)
		}
		resolvedPath = filepath.Join(cwd, userPath)
	}

	if _, err := os.Stat(resolvedPath); err != nil {
		return "", fmt.Errorf("Caddyfile not found at %s: %w", resolvedPath, err)
	}

	fmt.Printf("Using Caddyfile: %s\n", resolvedPath)
	return resolvedPath, nil
}

func (c *CLI) PromptAndSaveAPIToken(configPath string) (string, error) {
	fmt.Println("\n=== Cloudflare API Token Required ===")
	fmt.Println("To use dns-set with Cloudflare, you need to provide an API token.")
	fmt.Println("You can create one at: https://dash.cloudflare.com/profile/api-tokens")
	fmt.Println("Make sure the token has the following permissions:")
	fmt.Println("  - Zone:Read")
	fmt.Println("  - DNS:Edit")
	fmt.Print("\nPlease enter your Cloudflare API token (input will be hidden): ")

	tokenBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("failed to read API token: %w", err)
	}
	fmt.Println()

	token := strings.TrimSpace(string(tokenBytes))
	if token == "" {
		return "", fmt.Errorf("API token cannot be empty")
	}

	if len(token) < 40 {
		fmt.Println("Warning: The entered token seems too short. Cloudflare API tokens are typically 40+ characters.")
		fmt.Print("Do you want to continue anyway? (y/N): ")

		if !c.scanner.Scan() {
			return "", fmt.Errorf("failed to read confirmation")
		}

		confirmation := strings.TrimSpace(strings.ToLower(c.scanner.Text()))
		if confirmation != "y" && confirmation != "yes" {
			return "", fmt.Errorf("API token setup cancelled")
		}
	}

	newConfig := &config.Config{
		Cloudflare: config.CloudflareConfig{
			APIToken: token,
		},
		Preferences: c.config.Preferences,
	}

	if err := config.SaveToPath(newConfig, configPath); err != nil {
		return "", fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("✓ API token saved to configuration file\n")
	return token, nil
}
