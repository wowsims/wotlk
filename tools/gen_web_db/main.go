package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
)

// Example usage:
// go run ./tools/gen_web_db.go -infile=./assets/item_data/all_item_tooltips.csv -outfile=./assets/item_data/all_items_db.json
// go run ./tools/gen_web_db.go -infile=./assets/spell_data/all_spell_tooltips.csv -outfile=./assets/spell_data/all_spells_db.json

func main() {
	infile := flag.String("infile", "", "Path to input .csv file for tooltips data.")
	outfile := flag.String("outfile", "db.json", "Path to output file for generated json db.")
	flag.Parse()

	if *infile == "" {
		panic("infile flag is required!")
	}

	lines := readLines(*infile)
	// Ignore first line
	lines = lines[1:]

	pattern := regexp.MustCompile(`(\d+),\s*.*{"name":\s*"(.+?)".*"icon":\s*"(.+?)"`)

	var items []ItemData
	for _, line := range lines {
		matches := pattern.FindSubmatch([]byte(line))
		if len(matches) < 4 || matches[1] == nil || matches[2] == nil || matches[3] == nil {
			continue
		}

		id, err := strconv.Atoi(string(matches[1]))
		if err != nil {
			log.Fatal("Invalid ID: " + string(matches[1]))
		}

		items = append(items, ItemData{
			ID:   id,
			Name: string(matches[2]),
			Icon: string(matches[3]),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	file, _ := json.Marshal(items)
	_ = ioutil.WriteFile(*outfile, file, 0644)
}

type ItemData struct {
	ID   int
	Name string
	Icon string
}

func readLines(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return lines
}
