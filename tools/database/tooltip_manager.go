package database

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/tools"
)

// Generic class for fetching tooltip info from the web.
type TooltipManager struct {
	FilePath   string
	UrlPattern string
}

func (tm *TooltipManager) Read() map[int32]string {
	strDB := tools.ReadMapOrNil(tm.FilePath)

	db := core.MapMap(strDB, func(k, v string) (int32, string) {
		itemID, err := strconv.Atoi(k)
		if err != nil {
			log.Fatal("Invalid item ID: " + k)
		}

		return int32(itemID), v
	})

	fmt.Printf("\n--\nTOOLTIPS LOADED: %d\n--\n", len(db))
	return db
}

func (tm *TooltipManager) FetchFromWeb(idsToFetch []string) map[string]string {
	newTooltips := tools.ReadWebMultiMap(idsToFetch, func(id string) string {
		return fmt.Sprintf(tm.UrlPattern, id)
	})

	newTooltips = core.FilterMap(newTooltips, func(id string, body string) bool {
		if len(body) < 2 {
			fmt.Printf("Missing tooltip data for %s", id)
			return false
		}
		if strings.Contains(body, "\"error\":") {
			// fmt.Printf("Error in tooltip for %s: %s\n", id, body)
			return false
		}
		return true
	})

	return newTooltips
}

func (tm *TooltipManager) Fetch(minId, maxId int32, otherIds []string) {
	strDB := tools.ReadMapOrNil(tm.FilePath)

	var idsToFetch []string
	for i := minId; i <= maxId; i++ {
		id := strconv.Itoa(int(i))
		// Don't fetch tooltips already in the DB.
		if _, ok := strDB[id]; !ok {
			idsToFetch = append(idsToFetch, id)
		}
	}
	idsToFetch = append(idsToFetch, otherIds...)

	newTooltips := tm.FetchFromWeb(idsToFetch)

	for k, v := range newTooltips {
		strDB[k] = v
	}

	tools.WriteMapSortByIntKey(tm.FilePath, strDB)
}
