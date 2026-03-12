package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"
)

func writeJSONReport(pages map[string]PageData, filename string) error {
	fmt.Printf("generating report for %v pages...\n", len(pages))
	sortedKeys := slices.Sorted(maps.Keys(pages))
	sortedPageData := make([]PageData, len(sortedKeys))
	for i, data := range sortedKeys {
		sortedPageData[i] = pages[data]
	}

	data, err := json.MarshalIndent(sortedPageData, "", "  ")
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("writing report to %v...\n", filename)
	os.WriteFile(filename, data, 0644)
	return nil

}
