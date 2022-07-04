package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
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

	maxID := 52205
	threads := 8
	numPer := maxID / threads

	minID := 1

	wg := &sync.WaitGroup{}
	for crawler := 0; crawler < threads; crawler++ {
		wg.Add(1)
		go func(min, max int) {
			client := http.Client{}
			for i := min; i < max; i++ {
				url := fmt.Sprintf("https://wotlkdb.com/?item=%d&power", i)
				resp, err := client.Get(url)
				if err != nil {
					fmt.Printf("Error fetching %d: %s\n", i, err)
					continue
				}
				body, _ := ioutil.ReadAll(resp.Body)
				bstr := string(body)
				bstr = strings.Replace(bstr, fmt.Sprintf("$WowheadPower.registerItem('%d', 0, ", i), "", 1)
				bstr = strings.TrimSuffix(bstr, ";")
				bstr = strings.TrimSuffix(bstr, ")")
				bstr = strings.ReplaceAll(bstr, "\n", "")
				bstr = strings.ReplaceAll(bstr, "\t", "")
				bstr = strings.Replace(bstr, "name_enus: '", "\"name\": \"", 1)
				bstr = strings.Replace(bstr, "quality:", "\"quality\":", 1)
				bstr = strings.Replace(bstr, "icon: '", "\"icon\": \"", 1)
				bstr = strings.Replace(bstr, "tooltip_enus: '", "\"tooltip\": \"", 1)
				bstr = strings.ReplaceAll(bstr, "',", "\",")
				bstr = strings.ReplaceAll(bstr, "\\'", "'")
				// replace the '} with "}
				if strings.HasSuffix(bstr, "'}") {
					bstr = bstr[:len(bstr)-2] + "\"}"
				}

				// fmt.Printf("Found tooltip for %d\n%s\n", i, bstr)
				results <- result{id: i, value: bstr}
			}
		}(minID, minID+numPer)
		fmt.Printf("Started thread %d, %d to %d\n", crawler, minID, minID+numPer-1)
		minID += numPer
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	total := 0
	var lastUpdate time.Time
	for res := range results {
		total++

		if time.Since(lastUpdate).Seconds() > 2 {
			lastUpdate = time.Now()
			fmt.Printf("Tooltips %d/%d complete\n", total, maxID)
		}

		if strings.Contains(res.value, "\"error\":") {
			// fmt.Printf("Error in tooltip for %d: %s\n", i, bstr)
			continue
		}

		url := fmt.Sprintf("https://wotlkdb.com/?item=%d&power", res.id)
		f.WriteString(fmt.Sprintf("%d, %s, %s\n", res.id, url, res.value))
	}
}
