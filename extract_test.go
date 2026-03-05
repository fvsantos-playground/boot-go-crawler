package main

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

func TestGetHeadingFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "h1 heading",
			input:    "<html><body><h1>Welcome to Boot.dev</h1></body></html>",
			expected: "Welcome to Boot.dev",
		},
		{
			name:     "h2 heading",
			input:    "<html><body><h2>Welcome to Boot.dev</h2></body></html>",
			expected: "Welcome to Boot.dev",
		},
		{
			name:     "h3 heading",
			input:    "<html><body><h3>Welcome to Boot.dev</h3></body></html>",
			expected: "",
		},
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := getHeadingFromHTML(test.input)
			if err != nil {
				t.Errorf("Test %v - %s FAIL: unexpected error: %v", i, test.name, err)
				return
			}
			if actual != test.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, test.name, test.expected, actual)
			}
		})
	}
}

func TestParagraphInsideMainTag(t *testing.T) {
	expected := "Main paragraph."
	input := fmt.Sprintf("<html><body><p>Outside paragraph.</p><main><p>%v</p></main></body></html>", expected)

	actual, err := getFirstParagraphFromHTML(input)
	if err != nil {
		t.Errorf("expected %q, got %q", expected, actual)
	} else if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTML(t *testing.T) {
	expected := "Outside paragraph."
	input := fmt.Sprintf("<html><body><p>%v</p><p>Main paragraph.</p></body></html>", expected)

	actual, err := getFirstParagraphFromHTML(input)
	if err != nil {
		t.Errorf("expected %q, got %q", expected, actual)
	} else if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestNoParagraphs(t *testing.T) {
	input := "<html><body><h1>Welcome to Boot.dev</h1></body></html>"

	actual, err := getFirstParagraphFromHTML(input)
	if err != nil {
		t.Errorf("expected %q, got %q", "", actual)
	} else if actual != "" {
		t.Errorf("expected %q, got %q", "", actual)
	}
}

func TestGetURLsFromHTMLAbsolute(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body><a href="https://crawler-test.com"><span>Boot.dev</span></a></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getURLsFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://crawler-test.com"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetURLsFromHTMLRelative(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body><a href="/relative-path"><span>Boot.dev</span></a></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getURLsFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://crawler-test.com/relative-path"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetMultipleURLsFromHTMLMixed(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body>
	<a href="/relative-path"><span>Boot.dev</span></a>
	<a href="https://crawler-test.com/absolute-path"><span>Boot.dev</span></a>
	<a href="https://external.com"><span>External</span></a>
	<a href="/relative-path"><span>Boot.dev</span></a>
	</body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getURLsFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://crawler-test.com/relative-path", "https://crawler-test.com/absolute-path", "https://external.com", "https://crawler-test.com/relative-path"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetImagesFromHTMLRelative(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body><img src="/logo.png" alt="Logo"></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getImagesFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://crawler-test.com/logo.png"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestNoImageFromHTML(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body><p>No image</p</body></html>`

	baseUrl, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getImagesFromHTML(inputBody, baseUrl)
	var expected []string
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetMultipleImagesFromHTMLMixed(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body>
	<img src="/relative-path"><span>Boot.dev</span></img></img>
	<img src="https://crawler-test.com/absolute-path"><span>Boot.dev</span></img>
	<img src="https://external.com"><span>External</span></img>
	<img src="/relative-path"><span>Boot.dev</span></img>
	</body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getImagesFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://crawler-test.com/relative-path", "https://crawler-test.com/absolute-path", "https://external.com", "https://crawler-test.com/relative-path"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
