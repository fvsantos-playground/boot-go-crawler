package main

import "testing"

func TestNormalizeUrl(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "https with trailing slash",
			input:    "https://www.boot.dev/blog/path/",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "https without trailing slash",
			input:    "https://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "http with trailing slash",
			input:    "http://www.boot.dev/blog/path/",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "http without trailing slash",
			input:    "http://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "https with query params",
			input:    "https://www.boot.dev/blog/path?query=param",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "https with query params and trailing slash",
			input:    "https://www.boot.dev/blog/path?query=param/",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "https with query params and trailing slash",
			input:    "https://www.boot.dev/blog/path?query=param",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "without path",
			input:    "https://www.boot.dev",
			expected: "www.boot.dev",
		},
		{
			name:     "without path and trailing slash",
			input:    "https://www.boot.dev/",
			expected: "www.boot.dev",
		},
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := normalizeUrl(test.input)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, test.name, err)
				return
			}
			if actual != test.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, test.name, test.expected, actual)
			}
		})
	}
}
