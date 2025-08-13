package domain

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCaddyfileSource_GetDomains(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expected    []string
		expectError bool
	}{
		{
			name: "simple domains",
			content: `example.com {
	root * /var/www
}

api.example.com {
	reverse_proxy localhost:3000
}`,
			expected:    []string{"example.com", "api.example.com", "localhost"},
			expectError: false,
		},
		{
			name: "domains with ports",
			content: `example.com:80 {
	root * /var/www
}

api.example.com:443 {
	reverse_proxy localhost:3000
}`,
			expected:    []string{"example.com", "api.example.com", "localhost"},
			expectError: false,
		},
		{
			name: "mixed content",
			content: `# Comment line
example.com {
	root * /var/www
}

:8080 {
	respond "Hello"
}

localhost {
	respond "Local"
}

sub.example.org {
	proxy / backend:9000
}`,
			expected:    []string{"example.com", "localhost", "sub.example.org"},
			expectError: false,
		},
		{
			name: "no valid domains",
			content: `:8080 {
	respond "Hello"
}

:9000 {
	respond "World"
}`,
			expected:    nil,
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			caddyfilePath := filepath.Join(tmpDir, "Caddyfile")
			
			err := os.WriteFile(caddyfilePath, []byte(test.content), 0644)
			require.NoError(t, err)

			source := NewCaddyfileSource(caddyfilePath)
			domains, err := source.GetDomains()

			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, test.expected, domains)
			}
		})
	}
}

func TestCaddyfileSource_GetDomains_FileNotFound(t *testing.T) {
	source := NewCaddyfileSource("/nonexistent/Caddyfile")
	_, err := source.GetDomains()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open Caddyfile")
}

func TestCaddyfileSourceName(t *testing.T) {
	path := "/etc/caddy/Caddyfile"
	source := NewCaddyfileSource(path)
	assert.Equal(t, "Caddyfile (/etc/caddy/Caddyfile)", source.Name())
}