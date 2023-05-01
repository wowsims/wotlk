package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	_ "github.com/wowsims/wotlk/sim/encounters" // Needed for preset encounters.
	"github.com/wowsims/wotlk/tools"
	"github.com/wowsims/wotlk/tools/database"
)

// To do a full re-scrape, delete the previous output file first.
// go run ./tools/database/gen_db -outDir=assets -gen=atlasloot
// go run ./tools/database/gen_db -outDir=assets -gen=wowhead-items
// go run ./tools/database/gen_db -outDir=assets -gen=wowhead-spells -maxid=75000
// go run ./tools/database/gen_db -outDir=assets -gen=wowhead-gearplannerdb
// go run ./tools/database/gen_db -outDir=assets -gen=wotlk-items
// go run ./tools/database/gen_db -outDir=assets -gen=wago-db2-sparse
// go run ./tools/database/gen_db -outDir=assets -gen=db

var minId = flag.Int("minid", 1, "Minimum ID to scan for")
var maxId = flag.Int("maxid", 57000, "Maximum ID to scan for")
var outDir = flag.String("outDir", "assets", "Path to output directory for writing generated .go files.")
var genAsset = flag.String("gen", "", "Asset to generate. Valid values are 'db', 'atlasloot', 'wowhead-items', 'wowhead-spells', 'wowhead-itemdb', 'wotlk-items', and 'wago-db2-items'")

func main() {
	flag.Parse()
	if *outDir == "" {
		panic("outDir flag is required!")
	}

	dbDir := fmt.Sprintf("%s/database", *outDir)
	inputsDir := fmt.Sprintf("%s/db_inputs", *outDir)

	if *genAsset == "atlasloot" {
		db := database.ReadAtlasLootData()
		db.WriteJson(fmt.Sprintf("%s/atlasloot_db.json", inputsDir))
		return
	} else if *genAsset == "wowhead-items" {
		database.NewWowheadItemTooltipManager(fmt.Sprintf("%s/wowhead_item_tooltips.csv", inputsDir)).Fetch(int32(*minId), int32(*maxId))
		return
	} else if *genAsset == "wowhead-spells" {
		database.NewWowheadSpellTooltipManager(fmt.Sprintf("%s/wowhead_spell_tooltips.csv", inputsDir)).Fetch(int32(*minId), int32(*maxId))
		return
	} else if *genAsset == "wowhead-gearplannerdb" {
		tools.WriteFile(fmt.Sprintf("%s/wowhead_gearplannerdb.txt", inputsDir), tools.ReadWebRequired("https://nether.wowhead.com/wotlk/data/gear-planner?dv=100"))
		return
	} else if *genAsset == "wotlk-items" {
		database.NewWotlkItemTooltipManager(fmt.Sprintf("%s/wotlk_items_tooltips.csv", inputsDir)).Fetch(int32(*minId), int32(*maxId))
		return
	} else if *genAsset == "wago-db2-sparse" {
		tools.WriteFile(fmt.Sprintf("%s/wago_db2_items.csv", inputsDir), tools.ReadWebRequired("https://wago.tools/db2/ItemSparse/csv?build=3.4.2.49311"))
		return
	} else if *genAsset != "db" {
		panic("Invalid gen value")
	}

	itemTooltips := database.NewWowheadItemTooltipManager(fmt.Sprintf("%s/wowhead_item_tooltips.csv", inputsDir)).Read()
	spellTooltips := database.NewWowheadSpellTooltipManager(fmt.Sprintf("%s/wowhead_spell_tooltips.csv", inputsDir)).Read()
	wowheadDB := database.ParseWowheadDB(tools.ReadFile(fmt.Sprintf("%s/wowhead_gearplannerdb.txt", inputsDir)))
	atlaslootDB := database.ReadDatabaseFromJson(tools.ReadFile(fmt.Sprintf("%s/atlasloot_db.json", inputsDir)))
	factionRestrictions := database.ParseItemFactionRestrictionsFromWagoDB(tools.ReadFile(fmt.Sprintf("%s/wago_db2_items.csv", inputsDir)))

	db := database.NewWowDatabase()
	db.Encounters = core.PresetEncounters
	db.GlyphIDs = getGlyphIDsFromJson(fmt.Sprintf("%s/glyph_id_map.json", inputsDir))

	for _, response := range itemTooltips {
		if response.IsEquippable() {
			db.MergeItem(response.ToItemProto())
		} else if response.IsGem() {
			db.MergeGem(response.ToGemProto())
		}
	}

	for _, wowheadItem := range wowheadDB.Items {
		item := wowheadItem.ToProto()
		if _, ok := db.Items[item.Id]; ok {
			db.MergeItem(item)
		}
	}
	for _, item := range atlaslootDB.Items {
		if _, ok := db.Items[item.Id]; ok {
			db.MergeItem(item)
		}
	}

	db.MergeItems(database.ItemOverrides)
	db.MergeGems(database.GemOverrides)
	db.MergeEnchants(database.EnchantOverrides)
	ApplyGlobalFilters(db)
	AttachFactionInformation(db, factionRestrictions)

	leftovers := db.Clone()
	ApplyNonSimmableFilters(leftovers)
	leftovers.WriteBinaryAndJson(fmt.Sprintf("%s/leftover_db.bin", dbDir), fmt.Sprintf("%s/leftover_db.json", dbDir))

	ApplySimmableFilters(db)
	for _, enchant := range db.Enchants {
		if enchant.ItemId != 0 {
			db.AddItemIcon(enchant.ItemId, itemTooltips)
		}
		if enchant.SpellId != 0 {
			db.AddSpellIcon(enchant.SpellId, spellTooltips)
		}
	}

	for _, itemID := range database.ExtraItemIcons {
		db.AddItemIcon(itemID, itemTooltips)
	}

	for _, item := range db.Items {
		for _, source := range item.Sources {
			if crafted := source.GetCrafted(); crafted != nil {
				db.AddSpellIcon(crafted.SpellId, spellTooltips)
			}
		}
	}

	for _, spellId := range database.SharedSpellsIcons {
		db.AddSpellIcon(spellId, spellTooltips)
	}

	for _, spellIds := range GetAllTalentSpellIds(&inputsDir) {
		for _, spellId := range spellIds {
			db.AddSpellIcon(spellId, spellTooltips)
		}
	}

	atlasDBProto := atlaslootDB.ToUIProto()
	db.MergeZones(atlasDBProto.Zones)
	db.MergeNpcs(atlasDBProto.Npcs)

	db.WriteBinaryAndJson(fmt.Sprintf("%s/db.bin", dbDir), fmt.Sprintf("%s/db.json", dbDir))
}

// Filters out entities which shouldn't be included anywhere.
func ApplyGlobalFilters(db *database.WowDatabase) {
	db.Items = core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		if _, ok := database.ItemDenyList[item.Id]; ok {
			return false
		}

		for _, pattern := range database.DenyListNameRegexes {
			if pattern.MatchString(item.Name) {
				return false
			}
		}
		return true
	})

	// There is an 'unavailable' version of every naxx set, e.g. https://www.wowhead.com/wotlk/item=43728/bonescythe-gauntlets
	heroesItems := core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		return strings.HasPrefix(item.Name, "Heroes' ")
	})
	db.Items = core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		nameToMatch := "Heroes' " + item.Name
		for _, heroItem := range heroesItems {
			if heroItem.Name == nameToMatch {
				return false
			}
		}
		return true
	})

	// There is an 'unavailable' version of many t8 set pieces, e.g. https://www.wowhead.com/wotlk/item=46235/darkruned-gauntlets
	valorousItems := core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		return strings.HasPrefix(item.Name, "Valorous ")
	})
	db.Items = core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		nameToMatch := "Valorous " + item.Name
		for _, item := range valorousItems {
			if item.Name == nameToMatch {
				return false
			}
		}
		return true
	})

	db.Gems = core.FilterMap(db.Gems, func(_ int32, gem *proto.UIGem) bool {
		if _, ok := database.GemDenyList[gem.Id]; ok {
			return false
		}

		for _, pattern := range database.DenyListNameRegexes {
			if pattern.MatchString(gem.Name) {
				return false
			}
		}
		return true
	})

	db.ItemIcons = core.FilterMap(db.ItemIcons, func(_ int32, icon *proto.IconData) bool {
		return icon.Name != "" && icon.Icon != ""
	})
	db.SpellIcons = core.FilterMap(db.SpellIcons, func(_ int32, icon *proto.IconData) bool {
		return icon.Name != "" && icon.Icon != ""
	})
}

// AttachFactionInformation attaches faction information (faction restrictions) to the DB items.
func AttachFactionInformation(db *database.WowDatabase, factionRestrictions map[int32]proto.UIItem_FactionRestriction) {
	for _, item := range db.Items {
		item.FactionRestriction = factionRestrictions[item.Id]
	}
}

// Filters out entities which shouldn't be included in the sim.
func ApplySimmableFilters(db *database.WowDatabase) {
	db.Items = core.FilterMap(db.Items, simmableItemFilter)
	db.Gems = core.FilterMap(db.Gems, simmableGemFilter)
}
func ApplyNonSimmableFilters(db *database.WowDatabase) {
	db.Items = core.FilterMap(db.Items, func(id int32, item *proto.UIItem) bool {
		return !simmableItemFilter(id, item)
	})
	db.Gems = core.FilterMap(db.Gems, func(id int32, gem *proto.UIGem) bool {
		return !simmableGemFilter(id, gem)
	})
}
func simmableItemFilter(_ int32, item *proto.UIItem) bool {
	if _, ok := database.ItemAllowList[item.Id]; ok {
		return true
	}

	if item.Quality < proto.ItemQuality_ItemQualityUncommon {
		return false
	} else if item.Quality > proto.ItemQuality_ItemQualityLegendary {
		return false
	} else if item.Quality < proto.ItemQuality_ItemQualityEpic {
		if item.Ilvl < 145 {
			return false
		}
		if item.Ilvl < 149 && item.SetName == "" {
			return false
		}
	} else {
		// Epic and legendary items might come from classic, so use a lower ilvl threshold.
		if item.Ilvl < 140 {
			return false
		}
	}
	if item.Ilvl == 0 {
		fmt.Printf("Missing ilvl: %s\n", item.Name)
	}

	return true
}
func simmableGemFilter(_ int32, gem *proto.UIGem) bool {
	if gem.Quality < proto.ItemQuality_ItemQualityUncommon {
		return false
	}
	return true
}

type TalentConfig struct {
	FieldName string `json:"fieldName"`
	// Spell ID for each rank of this talent.
	// Omitted ranks will be inferred by incrementing from the last provided rank.
	SpellIds  []int32 `json:"spellIds"`
	MaxPoints int32   `json:"maxPoints"`
}

type TalentTreeConfig struct {
	Name          string         `json:"name"`
	BackgroundUrl string         `json:"backgroundUrl"`
	Talents       []TalentConfig `json:"talents"`
}

func getSpellIdsFromTalentJson(infile *string) []int32 {
	data, err := os.ReadFile(*infile)
	if err != nil {
		log.Fatalf("failed to load talent json file: %s", err)
	}

	var buf bytes.Buffer
	err = json.Compact(&buf, []byte(data))
	if err != nil {
		log.Fatalf("failed to compact json: %s", err)
	}

	var talents []TalentTreeConfig

	err = json.Unmarshal(buf.Bytes(), &talents)
	if err != nil {
		log.Fatalf("failed to parse talent to json %s", err)
	}

	spellIds := make([]int32, 0)

	for _, tree := range talents {
		for _, talent := range tree.Talents {
			spellIds = append(spellIds, talent.SpellIds...)

			// Infer omitted spell IDs.
			if len(talent.SpellIds) < int(talent.MaxPoints) {
				curSpellId := talent.SpellIds[len(talent.SpellIds)-1]
				for i := len(talent.SpellIds); i < int(talent.MaxPoints); i++ {
					curSpellId++
					spellIds = append(spellIds, curSpellId)
				}
			}
		}
	}

	return spellIds
}

func GetAllTalentSpellIds(inputsDir *string) map[string][]int32 {
	talentsDir := fmt.Sprintf("%s/../../ui/core/talents/trees", *inputsDir)
	specFiles := []string{
		"deathknight.json",
		"druid.json",
		"hunter.json",
		"hunter_cunning.json",
		"hunter_ferocity.json",
		"hunter_tenacity.json",
		"mage.json",
		"paladin.json",
		"priest.json",
		"rogue.json",
		"shaman.json",
		"warlock.json",
		"warrior.json",
	}

	ret_db := make(map[string][]int32, 0)

	for _, specFile := range specFiles {
		specPath := fmt.Sprintf("%s/%s", talentsDir, specFile)
		ret_db[specFile[:len(specFile)-5]] = getSpellIdsFromTalentJson(&specPath)
	}

	return ret_db

}

type GlyphID struct {
	ItemID  int32 `json:"itemId"`
	SpellID int32 `json:"spellId"`
}

func getGlyphIDsFromJson(infile string) []*proto.GlyphID {
	data, err := os.ReadFile(infile)
	if err != nil {
		log.Fatalf("failed to load glyph json file: %s", err)
	}

	var buf bytes.Buffer
	err = json.Compact(&buf, []byte(data))
	if err != nil {
		log.Fatalf("failed to compact json: %s", err)
	}

	var glyphIDs []GlyphID

	err = json.Unmarshal(buf.Bytes(), &glyphIDs)
	if err != nil {
		log.Fatalf("failed to parse glyph IDs to json %s", err)
	}

	return core.MapSlice(glyphIDs, func(gid GlyphID) *proto.GlyphID {
		return &proto.GlyphID{
			ItemId:  gid.ItemID,
			SpellId: gid.SpellID,
		}
	})
}
