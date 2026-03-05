package main

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHTML(html string) (string, error) {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", err
	}

	selector := doc.Find("h1")
	val, err := selector.First().Html()
	if err != nil {
		return "", err
	}

	if val == "" {
		selector = doc.Find("h2")
		val, err = selector.First().Html()
		if err != nil {
			return "", err
		}
	}

	return strings.TrimSpace(val), nil
}

func getFirstParagraphFromHTML(html string) (string, error) {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", err
	}

	selector := doc.Find("main p")
	val, err := selector.First().Html()
	if err != nil {
		return "", err
	}

	if val == "" {
		selector = doc.Find("p")
		val, err = selector.First().Html()
		if err != nil {
			return "", err
		}
	}

	return strings.TrimSpace(val), nil
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	reader := strings.NewReader(htmlBody)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	var urls []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			parsedURL, err := url.Parse(href)
			if err == nil {
				resolvedURL := baseURL.ResolveReference(parsedURL)
				urls = append(urls, resolvedURL.String())
			}
		}
	})

	return urls, nil
}

