package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_DefaultValues(t *testing.T) {
	config, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "api", config.Preferences.IPSource)
	assert.Equal(t, "/etc/caddy/Caddyfile", config.Preferences.CaddyfilePath)
	assert.Equal(t, 300, config.Preferences.DefaultTTL)
}

func TestLoad_WithConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".config", "dns-set")
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	configContent := `cloudflare:
  api_token: "test-token"
preferences:
  ip_source: "interface"
  caddyfile_path: "/custom/Caddyfile"
  default_ttl: 600`

	configPath := filepath.Join(configDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	config, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "test-token", config.Cloudflare.APIToken)
	assert.Equal(t, "interface", config.Preferences.IPSource)
	assert.Equal(t, "/custom/Caddyfile", config.Preferences.CaddyfilePath)
	assert.Equal(t, 600, config.Preferences.DefaultTTL)
}

func TestLoad_WithEnvironmentVariables(t *testing.T) {
	oldToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	oldIPSource := os.Getenv("DNS_SET_IP_SOURCE")
	defer func() {
		os.Setenv("CLOUDFLARE_API_TOKEN", oldToken)
		os.Setenv("DNS_SET_IP_SOURCE", oldIPSource)
	}()

	os.Setenv("CLOUDFLARE_API_TOKEN", "env-token")
	os.Setenv("DNS_SET_IP_SOURCE", "manual")

	config, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "env-token", config.Cloudflare.APIToken)
	assert.Equal(t, "manual", config.Preferences.IPSource)
}

func TestLoad_WithDotEnvFile(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	oldToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	oldIPSource := os.Getenv("DNS_SET_IP_SOURCE")
	defer func() {
		os.Setenv("CLOUDFLARE_API_TOKEN", oldToken)
		os.Setenv("DNS_SET_IP_SOURCE", oldIPSource)
	}()
	
	os.Unsetenv("CLOUDFLARE_API_TOKEN")
	os.Unsetenv("DNS_SET_IP_SOURCE")

	envContent := `CLOUDFLARE_API_TOKEN=dotenv-token
DNS_SET_IP_SOURCE=interface`

	err := os.WriteFile(".env", []byte(envContent), 0644)
	require.NoError(t, err)

	config, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "dotenv-token", config.Cloudflare.APIToken)
	assert.Equal(t, "interface", config.Preferences.IPSource)
}

func TestLoad_WithCustomConfigDir(t *testing.T) {
	tmpDir := t.TempDir()
	customConfigDir := filepath.Join(tmpDir, "custom-config")
	err := os.MkdirAll(customConfigDir, 0755)
	require.NoError(t, err)

	configContent := `cloudflare:
  api_token: "custom-dir-token"`

	configPath := filepath.Join(customConfigDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	oldConfigDir := os.Getenv("DNS_SET_CONFIG_DIR")
	defer os.Setenv("DNS_SET_CONFIG_DIR", oldConfigDir)
	os.Setenv("DNS_SET_CONFIG_DIR", customConfigDir)

	config, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "custom-dir-token", config.Cloudflare.APIToken)
}

func TestSave(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	config := &Config{
		Cloudflare: CloudflareConfig{
			APIToken: "test-token",
		},
		Preferences: PreferencesConfig{
			IPSource:      "interface",
			CaddyfilePath: "/test/Caddyfile",
			DefaultTTL:    600,
		},
	}

	err := Save(config)
	require.NoError(t, err)

	configPath := filepath.Join(tmpDir, ".config", "dns-set", "config.yaml")
	assert.FileExists(t, configPath)

	loadedConfig, err := Load()
	require.NoError(t, err)

	assert.Equal(t, config.Cloudflare.APIToken, loadedConfig.Cloudflare.APIToken)
	assert.Equal(t, config.Preferences.IPSource, loadedConfig.Preferences.IPSource)
	assert.Equal(t, config.Preferences.CaddyfilePath, loadedConfig.Preferences.CaddyfilePath)
	assert.Equal(t, config.Preferences.DefaultTTL, loadedConfig.Preferences.DefaultTTL)
}