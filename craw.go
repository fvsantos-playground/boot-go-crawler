package main

import (
	"fmt"
	"net/url"
)

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() { <-cfg.concurrencyControl }()

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	if currentURL.Host != cfg.baseURL.Host {
		fmt.Printf("Skipping %v because it is not on the same domain as %v\n", rawCurrentURL, cfg.baseURL.String())
		return
	}

	normalizedURL, err := normalizeUrl(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	isFirstVisit := cfg.addPageVisit(normalizedURL)
	if !isFirstVisit {
		fmt.Printf("Already crawled %v\n", rawCurrentURL)
		return
	}

	fmt.Printf("************* Crawling %v *************\n", rawCurrentURL)
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
	fmt.Printf("************* Data extracted: %v *************\n", rawCurrentURL)
	cfg.pages[normalizedURL] = pageData

	for _, link := range pageData.OutgoingLinks {
		cfg.wg.Go(func() {
			cfg.crawlPage(link)
		})
	}

}
