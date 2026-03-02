package main

import (
	"errors"
	"fmt"
	uri "net/url"
	"strings"
)

func normalizeUrl(url string) (string, error) {
	if url == "" {
		return "", errors.New("url is required")
	}

	u, err := uri.Parse(url)
	if err != nil {
		return "", err
	}

	host, path := strings.TrimSpace(u.Host), strings.TrimSpace(strings.TrimSuffix(u.Path, "/"))
	if len(path) == 0 {
		return host, nil
	}

	return fmt.Sprintf("%s%s", host, path), nil
}
