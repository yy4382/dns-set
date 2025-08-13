package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidDomain(t *testing.T) {
	tests := []struct {
		domain   string
		expected bool
	}{
		{"example.com", true},
		{"sub.example.com", true},
		{"my-site.example.org", true},
		{"test123.co.uk", true},
		{"", false},
		{".", false},
		{".example.com", false},
		{"example.com.", false},
		{"example", false},
		{"ex..ample.com", false},
		{"example-.com", false},
		{"-example.com", false},
		{"example.c", false},
		{"verylongdomainnamethatisgreaterthan63charactersandshouldreaallyfail.com", false},
	}

	for _, test := range tests {
		t.Run(test.domain, func(t *testing.T) {
			result := isValidDomain(test.domain)
			assert.Equal(t, test.expected, result, "Domain: %s", test.domain)
		})
	}
}

func TestManualSourceName(t *testing.T) {
	source := NewManualSource()
	assert.Equal(t, "Manual Input", source.Name())
}
