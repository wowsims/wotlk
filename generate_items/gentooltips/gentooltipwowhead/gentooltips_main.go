package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	f, err := os.Create("all_item_tooltips.csv")
	if err != nil {
		log.Fatalf("failed to open file to write: %s", err)
	}
	type result struct {
		id    int
		value string
	}
	results := make(chan result, 10)
	go func() {
		for i := 1; i < 52205; i++ {
			url := fmt.Sprintf("https://wowhead.com/wotlk/tooltip/item/%d?json", i)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("Error fetching %d: %s\n", i, err)
				continue
			}
			body, _ := ioutil.ReadAll(resp.Body)
			bstr := string(body)
			if strings.Contains(bstr, "\"error\":") {
				fmt.Printf("Error in tooltip for %d: %s\n", i, bstr)
				continue
			}
			fmt.Printf("Found tooltip for %d\n", i)
			results <- result{id: i, value: bstr}
		}
		close(results)
	}()

	for res := range results {
		url := fmt.Sprintf("https://wowhead.com/wotlk/tooltip/item/%d?json", res.id)
		f.WriteString(fmt.Sprintf("%d, %s, %s\n", res.id, url, res.value))
	}
}
