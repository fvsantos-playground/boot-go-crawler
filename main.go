package main

import (
	"fmt"
	"net/url"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseUrl, err := url.Parse(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("starting crawl of: %v\n", baseUrl)

	// data, err := getHTML(baseUrl.String())
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// fmt.Println(data)

	crawlPage(baseUrl.String(), baseUrl.String(), map[string]int{})
}
