package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yy4382/dns-set/internal/config"
	"github.com/yy4382/dns-set/internal/dns"
	"github.com/yy4382/dns-set/internal/ui"
)

var rootCmd = &cobra.Command{
	Use:   "dns-set",
	Short: "A tool for managing DNS records on DNS providers",
	Long: `dns-set is a command-line tool for automatically managing DNS records.
It can read domains from multiple sources (manual input, Caddyfile),
detect IP addresses through various methods (network interface, API, manual),
and update DNS records on supported providers (currently Cloudflare).`,
	RunE: runDNSSet,
}

func init() {
	rootCmd.Flags().StringP("config", "c", "", "Config file path")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runDNSSet(cmd *cobra.Command, args []string) error {
	configPath, _ := cmd.Flags().GetString("config")

	var cfg *config.Config
	var err error

	if configPath != "" {
		cfg, err = config.LoadWithConfigPath(configPath)
	} else {
		cfg, err = config.Load()
	}

	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	apiToken := cfg.Cloudflare.APIToken

	if apiToken == "" {
		tempCLI := ui.NewCLI(cfg, nil)

		promptedToken, err := tempCLI.PromptAndSaveAPIToken(configPath)
		if err != nil {
			return fmt.Errorf("failed to configure API token: %w", err)
		}

		apiToken = promptedToken

		if configPath != "" {
			cfg, err = config.LoadWithConfigPath(configPath)
		} else {
			cfg, err = config.Load()
		}

		if err != nil {
			return fmt.Errorf("failed to reload configuration: %w", err)
		}
	}

	provider, err := dns.NewCloudflareProvider(apiToken)
	if err != nil {
		return fmt.Errorf("failed to initialize Cloudflare provider: %w", err)
	}

	cli := ui.NewCLI(cfg, provider)
	return cli.Run()
}
