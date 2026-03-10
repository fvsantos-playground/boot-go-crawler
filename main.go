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
	} else if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	rawURL := args[0]
	cfg := configure(rawURL, 10)

	fmt.Printf("starting crawl of: %v\n", rawURL)
	// cfg.wg.Add(1)
	// go cfg.crawlPage(rawURL)

	cfg.wg.Go(func() {
		cfg.crawlPage(rawURL)
	})

	cfg.wg.Wait()
	fmt.Printf("finished crawling %v\n", rawURL)
}
