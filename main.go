package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(args) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	var maxConcurrency, maxPages int = 10, 50
	if len(args) > 1 {
		fmt.Sscanf(args[1], "%d", &maxConcurrency)
	}

	if len(args) > 2 {
		fmt.Sscanf(args[2], "%d", &maxPages)
	}

	rawURL := args[0]
	cfg := configure(rawURL, maxConcurrency, maxPages)

	fmt.Printf("starting crawl of: %v\n", rawURL)
	// cfg.wg.Add(1)
	// go cfg.crawlPage(rawURL)

	cfg.wg.Go(func() {
		cfg.crawlPage(rawURL)
	})

	cfg.wg.Wait()
	fmt.Printf("finished crawling %v\n", rawURL)

	writeJSONReport(cfg.pages, "report.json")
}
