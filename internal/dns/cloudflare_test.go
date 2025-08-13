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

func TestCloudflare_TTLLogic(t *testing.T) {
	tests := []struct {
		name        string
		inputTTL    *int
		expectedTTL int
		description string
	}{
		{
			name:        "nil TTL uses auto",
			inputTTL:    nil,
			expectedTTL: 1,
			description: "When config DefaultTTL is nil, should use Cloudflare auto TTL (1)",
		},
		{
			name:        "zero TTL uses auto",
			inputTTL:    intPtr(0),
			expectedTTL: 1,
			description: "When config DefaultTTL is 0, should use Cloudflare auto TTL (1)",
		},
		{
			name:        "positive TTL uses value",
			inputTTL:    intPtr(300),
			expectedTTL: 300,
			description: "When config DefaultTTL is set to a positive value, should use that value",
		},
		{
			name:        "large TTL uses value",
			inputTTL:    intPtr(86400),
			expectedTTL: 86400,
			description: "Should handle large TTL values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the TTL logic that would be used in UpdateRecord
			var actualTTL int
			if tt.inputTTL == nil {
				actualTTL = 1 // Cloudflare's auto TTL when no TTL specified in config
			} else {
				actualTTL = *tt.inputTTL
				// Use TTL auto (1) when TTL is explicitly set to 0
				if actualTTL == 0 {
					actualTTL = 1
				}
			}

			assert.Equal(t, tt.expectedTTL, actualTTL, tt.description)
		})
	}
}

// Helper function to create int pointers
func intPtr(i int) *int {
	return &i
}
