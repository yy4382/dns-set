package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractRootDomain(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		expected string
	}{
		{
			name:     "subdomain",
			domain:   "test.yyang.dev",
			expected: "yyang.dev",
		},
		{
			name:     "deep subdomain",
			domain:   "api.test.yyang.dev",
			expected: "yyang.dev",
		},
		{
			name:     "root domain",
			domain:   "yyang.dev",
			expected: "yyang.dev",
		},
		{
			name:     "com domain",
			domain:   "example.com",
			expected: "example.com",
		},
		{
			name:     "subdomain with com",
			domain:   "www.example.com",
			expected: "example.com",
		},
		{
			name:     "single part",
			domain:   "localhost",
			expected: "localhost",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractRootDomain(tt.domain)
			assert.Equal(t, tt.expected, result)
		})
	}
}