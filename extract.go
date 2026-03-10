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
	if !strings.Contains(strings.ToLower(content_type), "text/html") {
		return "", fmt.Errorf("Invalid content type: %v", content_type)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	baseUrl, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	if currentURL.Host != baseUrl.Host {
		fmt.Printf("Skipping %v because it is not on the same domain as %v\n", rawCurrentURL, rawBaseURL)
		return
	}

	normalizedURL, err := normalizeUrl(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	count, ok := pages[normalizedURL]
	if ok {
		fmt.Printf("Already crawled %v %v times\n", rawCurrentURL, count)
		pages[normalizedURL] = count + 1
	} else {
		pages[normalizedURL] = 1
		fmt.Printf("Crawling %v...\n", rawCurrentURL)
		html, err := getHTML(rawCurrentURL)
		if err != nil {
			fmt.Println(err)
			return
		}

		pageData, err := extractPageData(html, rawCurrentURL)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Page title: %v\n", pageData.Heading)
		fmt.Printf("Page outgoing links count: %v\n", len(pageData.OutgoingLinks))
		for _, link := range pageData.OutgoingLinks {
			crawlPage(rawBaseURL, link, pages)
		}
	}
}
