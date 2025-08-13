package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Cloudflare  CloudflareConfig `mapstructure:"cloudflare"`
	Preferences PreferencesConfig `mapstructure:"preferences"`
}

type CloudflareConfig struct {
	APIToken string `mapstructure:"api_token"`
}

type PreferencesConfig struct {
	CaddyfilePath  string `mapstructure:"caddyfile_path"`
	DefaultTTL     int    `mapstructure:"default_ttl"`
}

func Load() (*Config, error) {
	if err := loadEnvFiles(); err != nil {
		return nil, fmt.Errorf("failed to load .env files: %w", err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	
	configDir, err := getConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}
	viper.AddConfigPath(configDir)

	viper.SetEnvPrefix("DNS_SET")
	viper.AutomaticEnv()

	viper.BindEnv("cloudflare.api_token", "CLOUDFLARE_API_TOKEN")
	viper.BindEnv("preferences.caddyfile_path", "DNS_SET_CADDYFILE_PATH")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func Save(config *Config) error {
	configDir, err := getConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	viper.Set("cloudflare", config.Cloudflare)
	viper.Set("preferences", config.Preferences)

	configPath := filepath.Join(configDir, "config.yaml")
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func getConfigDir() (string, error) {
	if customConfigDir := os.Getenv("DNS_SET_CONFIG_DIR"); customConfigDir != "" {
		return customConfigDir, nil
	}

	if configDir := os.Getenv("XDG_CONFIG_HOME"); configDir != "" {
		return filepath.Join(configDir, "dns-set"), nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".config", "dns-set"), nil
}

func loadEnvFiles() error {
	envPaths := []string{
		".env",
		".env.local",
	}

	configDir, err := getConfigDir()
	if err == nil {
		envPaths = append(envPaths, 
			filepath.Join(configDir, ".env"),
			filepath.Join(configDir, ".env.local"),
		)
	}

	for _, path := range envPaths {
		if _, err := os.Stat(path); err == nil {
			if loadErr := godotenv.Load(path); loadErr != nil {
				return fmt.Errorf("failed to load %s: %w", path, loadErr)
			}
		}
	}

	return nil
}

func setDefaults() {
	viper.SetDefault("preferences.caddyfile_path", "/etc/caddy/Caddyfile")
	viper.SetDefault("preferences.default_ttl", 300)
}