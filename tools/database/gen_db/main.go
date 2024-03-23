package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/wowsims/wotlk/sim"
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
// go run ./tools/database/gen_db -outDir=assets -gen=wago-db2-items
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
		database.NewWowheadItemTooltipManager(fmt.Sprintf("%s/wowhead_item_tooltips.csv", inputsDir)).Fetch(int32(*minId), int32(*maxId), database.OtherItemIdsToFetch)
		return
	} else if *genAsset == "wowhead-spells" {
		database.NewWowheadSpellTooltipManager(fmt.Sprintf("%s/wowhead_spell_tooltips.csv", inputsDir)).Fetch(int32(*minId), int32(*maxId), []string{})
		return
	} else if *genAsset == "wowhead-gearplannerdb" {
		tools.WriteFile(fmt.Sprintf("%s/wowhead_gearplannerdb.txt", inputsDir), tools.ReadWebRequired("https://nether.wowhead.com/wotlk/data/gear-planner?dv=100"))
		return
	} else if *genAsset == "wotlk-items" {
		database.NewWotlkItemTooltipManager(fmt.Sprintf("%s/wotlk_items_tooltips.csv", inputsDir)).Fetch(int32(*minId), int32(*maxId), []string{})
		return
	} else if *genAsset == "wago-db2-items" {
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
			// Only included items that are in wowheads gearplanner db
			// Wowhead doesn't seem to have a field/flag to signify 'not available / in game' but their gearplanner db has them filtered
			item := response.ToItemProto()
			if _, ok := wowheadDB.Items[strconv.Itoa(int(item.Id))]; ok {
				db.MergeItem(item)
			}
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

	for _, spellIds := range GetAllRotationSpellIds() {
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

	// There is an 'unavailable' version of many t9 set pieces, e.g. https://www.wowhead.com/wotlk/item=48842/thralls-hauberk
	triumphItems := core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		return strings.HasSuffix(item.Name, "of Triumph")
	})
	db.Items = core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		nameToMatch := item.Name + " of Triumph"
		for _, item := range triumphItems {
			if item.Name == nameToMatch {
				return false
			}
		}
		return true
	})

	// Theres an invalid 251 t10 set for every class
	// The invalid set has a higher item id than the 'correct' ones
	t10invalidItems := core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		return item.SetName != "" && item.Ilvl == 251
	})
	db.Items = core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		for _, t10item := range t10invalidItems {
			if t10item.Name == item.Name && item.Ilvl == t10item.Ilvl && item.Id > t10item.Id {
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

	db.Gems = core.FilterMap(db.Gems, func(_ int32, gem *proto.UIGem) bool {
		if strings.HasSuffix(gem.Name, "Stormjewel") {
			gem.Unique = false
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
	} else if item.Quality == proto.ItemQuality_ItemQualityArtifact {
		return false
	} else if item.Quality > proto.ItemQuality_ItemQualityHeirloom {
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
		if item.Quality != proto.ItemQuality_ItemQualityHeirloom && item.Ilvl < 140 {
			return false
		}
	}
	if item.Ilvl == 0 {
		fmt.Printf("Missing ilvl: %s\n", item.Name)
	}

	return true
}
func simmableGemFilter(_ int32, gem *proto.UIGem) bool {
	if _, ok := database.GemAllowList[gem.Id]; ok {
		return true
	}

	// Allow all meta gems
	if gem.Color == proto.GemColor_GemColorMeta {
		return true
	}

	// Arbitrary to filter out old gems
	if gem.Id < 39900 {
		return false
	}

	return gem.Quality >= proto.ItemQuality_ItemQualityUncommon
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

func CreateTempAgent(r *proto.Raid) core.Agent {
	encounter := core.MakeSingleTargetEncounter(0.0)
	env, _, _ := core.NewEnvironment(r, encounter, false)
	return env.Raid.Parties[0].Players[0]
}

type RotContainer struct {
	Name string
	Raid *proto.Raid
}

func GetAllRotationSpellIds() map[string][]int32 {
	sim.RegisterAll()

	rotMapping := []RotContainer{
		{Name: "feral", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassDruid,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "-503202132322010053120230310511-205503012",
		}, &proto.Player_FeralDruid{FeralDruid: &proto.FeralDruid{Options: &proto.FeralDruid_Options{}, Rotation: &proto.FeralDruid_Rotation{}}}), nil, nil, nil)},
		{Name: "balance", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassDruid,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "5102233115331303213305311031--205003002",
		}, &proto.Player_BalanceDruid{BalanceDruid: &proto.BalanceDruid{Options: &proto.BalanceDruid_Options{}}}), nil, nil, nil)},
		{Name: "guardian", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassDruid,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "-503232132322010353120300313511-20350001",
		}, &proto.Player_FeralTankDruid{FeralTankDruid: &proto.FeralTankDruid{Options: &proto.FeralTankDruid_Options{}}}), nil, nil, nil)},
		{Name: "restodruid", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassDruid,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "05320031103--230023312131502331050313051",
		}, &proto.Player_RestorationDruid{RestorationDruid: &proto.RestorationDruid{Options: &proto.RestorationDruid_Options{}}}), nil, nil, nil)},
		{Name: "elemental", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassShaman,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "0532001523212351322301351-005052031",
		}, &proto.Player_ElementalShaman{ElementalShaman: &proto.ElementalShaman{Options: &proto.ElementalShaman_Options{}}}), nil, nil, nil)},
		{Name: "enhance", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassShaman,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "053030152-30405003105021333031131031051",
		}, &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{}}}), nil, nil, nil)},
		{Name: "restosham", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassShaman,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "-3020503-50005331335310501122331251",
		}, &proto.Player_RestorationShaman{RestorationShaman: &proto.RestorationShaman{Options: &proto.RestorationShaman_Options{}}}), nil, nil, nil)},
		{Name: "hunter", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassHunter,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "-015305101-5000032500033330532135301311",
		}, &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{}}}), nil, nil, nil)},
		{Name: "mage", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassMage,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "23000513310033015032310250532-03-023303001",
		}, &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{}}}), nil, nil, nil)},
		{Name: "healingpriest", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassPriest,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "0503203130300512301313231251-2351010303",
		}, &proto.Player_HealingPriest{HealingPriest: &proto.HealingPriest{Options: &proto.HealingPriest_Options{}}}), nil, nil, nil)},
		{Name: "shadow", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassPriest,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "05032031--325023051223010323151301351",
		}, &proto.Player_ShadowPriest{ShadowPriest: &proto.ShadowPriest{Options: &proto.ShadowPriest_Options{}}}), nil, nil, nil)},
		{Name: "rogue", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassRogue,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "00532000523-0252051050035010223100501251",
		}, &proto.Player_Rogue{Rogue: &proto.Rogue{Options: &proto.Rogue_Options{}}}), nil, nil, nil)},
		{Name: "warrior", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassWarrior,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "302023102331-305053000520310053120500351",
		}, &proto.Player_Warrior{Warrior: &proto.Warrior{Options: &proto.Warrior_Options{}}}), nil, nil, nil)},
		{Name: "protwarrior", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassWarrior,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "2500030023-302-053351225000012521030113321",
		}, &proto.Player_ProtectionWarrior{ProtectionWarrior: &proto.ProtectionWarrior{Options: &proto.ProtectionWarrior_Options{}}}), nil, nil, nil)},
		{Name: "holypally", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassPaladin,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "50350151020013053100515221-50023131203",
		}, &proto.Player_HolyPaladin{HolyPaladin: &proto.HolyPaladin{Options: &proto.HolyPaladin_Options{}}}), nil, nil, nil)},
		{Name: "protpally", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassPaladin,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "-05005135200132311333312321-511302012003",
		}, &proto.Player_ProtectionPaladin{ProtectionPaladin: &proto.ProtectionPaladin{Options: &proto.ProtectionPaladin_Options{}}}), nil, nil, nil)},
		{Name: "ret", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassPaladin,
			Race:          proto.Race_RaceBloodElf,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "050501-05-05232051203331302133231331",
		}, &proto.Player_RetributionPaladin{RetributionPaladin: &proto.RetributionPaladin{Options: &proto.RetributionPaladin_Options{}}}), nil, nil, nil)},
		{Name: "warlock", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassWarlock,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "2350002030023510253500331151--550000051",
		}, &proto.Player_Warlock{Warlock: &proto.Warlock{Options: &proto.Warlock_Options{}}}), nil, nil, nil)},
		{Name: "dk", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassDeathknight,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "-320043500002-2300303050032152000150013133051",
		}, &proto.Player_Deathknight{Deathknight: &proto.Deathknight{Options: &proto.Deathknight_Options{}}}), nil, nil, nil)},
		{Name: "tankdk", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassDeathknight,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "005510153330330220102013-3050505100023101-002",
		}, &proto.Player_TankDeathknight{TankDeathknight: &proto.TankDeathknight{Options: &proto.TankDeathknight_Options{}, Rotation: &proto.TankDeathknight_Rotation{}}}), nil, nil, nil)},
	}

	ret_db := make(map[string][]int32, 0)

	for _, r := range rotMapping {
		f := CreateTempAgent(r.Raid).GetCharacter()

		spells := make([]int32, 0, len(f.Spellbook))

		for _, s := range f.Spellbook {
			if s.SpellID != 0 {
				spells = append(spells, s.SpellID)
			}
		}

		for _, s := range f.GetAuras() {
			if s.ActionID.SpellID != 0 {
				spells = append(spells, s.ActionID.SpellID)
			}
		}

		ret_db[r.Name] = spells
	}
	return ret_db
}
