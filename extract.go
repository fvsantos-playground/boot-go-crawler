package main

import (
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
