package handler

import (
	"testing"
)

func TestNormalizeSEOHost(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		expected string
	}{
		{
			name:     "web3main.pro exact",
			raw:      "web3main.pro",
			expected: "web3main.pro",
		},
		{
			name:     "web3main.pro uppercase",
			raw:      "WEB3MAIN.PRO",
			expected: "web3main.pro",
		},
		{
			name:     "web3main.pro mixed case",
			raw:      "Web3Main.Pro",
			expected: "web3main.pro",
		},
		{
			name:     "www.web3main.pro",
			raw:      "www.web3main.pro",
			expected: "web3main.pro",
		},
		{
			name:     "WWW.WEB3MAIN.PRO",
			raw:      "WWW.WEB3MAIN.PRO",
			expected: "web3main.pro",
		},
		{
			name:     "web3main.pro with port 443",
			raw:      "web3main.pro:443",
			expected: "web3main.pro",
		},
		{
			name:     "www.web3main.pro with port 8080",
			raw:      "www.web3main.pro:8080",
			expected: "web3main.pro",
		},
		{
			name:     "www.web3main.pro with port 443",
			raw:      "www.web3main.pro:443",
			expected: "web3main.pro",
		},
		{
			name:     "phpqa.ru",
			raw:      "phpqa.ru",
			expected: "phpqa.ru",
		},
		{
			name:     "PHPQA.RU",
			raw:      "PHPQA.RU",
			expected: "phpqa.ru",
		},
		{
			name:     "www.phpqa.ru",
			raw:      "www.phpqa.ru",
			expected: "phpqa.ru",
		},
		{
			name:     "phpqa.ru with port 443",
			raw:      "phpqa.ru:443",
			expected: "phpqa.ru",
		},
		{
			name:     "www.phpqa.ru with port 8080",
			raw:      "www.phpqa.ru:8080",
			expected: "phpqa.ru",
		},
		{
			name:     "evil.com falls back to default",
			raw:      "evil.com",
			expected: defaultSEOHost,
		},
		{
			name:     "evil.com with port falls back to default",
			raw:      "evil.com:443",
			expected: defaultSEOHost,
		},
		{
			name:     "www.evil.com falls back to default",
			raw:      "www.evil.com",
			expected: defaultSEOHost,
		},
		{
			name:     "empty string falls back to default",
			raw:      "",
			expected: defaultSEOHost,
		},
		{
			name:     "127.0.0.1 falls back to default",
			raw:      "127.0.0.1",
			expected: defaultSEOHost,
		},
		{
			name:     "localhost falls back to default",
			raw:      "localhost",
			expected: defaultSEOHost,
		},
		{
			name:     "arbitrary.domain falls back to default",
			raw:      "arbitrary.domain",
			expected: defaultSEOHost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeSEOHost(tt.raw)
			if result != tt.expected {
				t.Errorf("normalizeSEOHost(%q) = %q, want %q", tt.raw, result, tt.expected)
			}
		})
	}
}
