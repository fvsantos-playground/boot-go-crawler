package main

import (
	"fmt"
	"net/url"
)

type PageData struct {
	URL            string
	Heading        string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func extractPageData(html, pageURL string) (PageData, error) {
	fmt.Printf("Extracting data from %v...\n", pageURL)
	pageData := PageData{
		URL: pageURL,
	}

	heading, err := getHeadingFromHTML(html)
	if err != nil {
		return PageData{}, err
	}
	pageData.Heading = heading

	firstParagraph, err := getFirstParagraphFromHTML(html)
	if err != nil {
		return PageData{}, err
	}
	pageData.FirstParagraph = firstParagraph

	baseUrl, err := url.Parse(pageURL)
	if err != nil {
		return PageData{}, err
	}

	outgoingLinks, err := getURLsFromHTML(html, baseUrl)
	if err != nil {
		return PageData{}, err
	}
	pageData.OutgoingLinks = outgoingLinks

	imageURLs, err := getImagesFromHTML(html, baseUrl)
	if err != nil {
		return PageData{}, err
	}
	pageData.ImageURLs = imageURLs

	return pageData, nil
}
