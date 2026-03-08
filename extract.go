package main

import (
	"fmt"
	"io"
	"net/http"
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

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	reader := strings.NewReader(htmlBody)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	var urls []string
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			parsedUrl, err := url.Parse(src)
			if err == nil {
				resolvedUrl := baseURL.ResolveReference(parsedUrl)
				urls = append(urls, resolvedUrl.String())
			}
		}
	})

	return urls, nil
}

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "BootCrawler/1.0")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("Invalid status code: %v", res.StatusCode)
	}

	content_type := res.Header.Get("Content-Type")
	if content_type != "text/html" {
		return "", fmt.Errorf("Invalid content type: %v", content_type)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
