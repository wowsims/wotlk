package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/tools"
)

func ReadAtlasLootData(inputsDir string) *WowDatabase {
	db := NewWowDatabase()

	// Read these in reverse order, because some items are listed in multiple expansions
	// and we want to overwrite with the earliest value.
	readAtlasLootSourceData(db, proto.Expansion_ExpansionWotlk, "https://raw.githubusercontent.com/Hoizame/AtlasLootClassic/master/AtlasLootClassic_Data/source-wrath.lua")
	readAtlasLootSourceData(db, proto.Expansion_ExpansionTbc, "https://raw.githubusercontent.com/Hoizame/AtlasLootClassic/master/AtlasLootClassic_Data/source-tbc.lua")
	readAtlasLootSourceData(db, proto.Expansion_ExpansionVanilla, "https://raw.githubusercontent.com/Hoizame/AtlasLootClassic/master/AtlasLootClassic_Data/source.lua")

	readAtlasLootDungeonData(db, proto.Expansion_ExpansionVanilla, "https://raw.githubusercontent.com/Hoizame/AtlasLootClassic/master/AtlasLootClassic_DungeonsAndRaids/data.lua")
	readAtlasLootDungeonData(db, proto.Expansion_ExpansionTbc, "https://raw.githubusercontent.com/Hoizame/AtlasLootClassic/master/AtlasLootClassic_DungeonsAndRaids/data-tbc.lua")
	readAtlasLootDungeonData(db, proto.Expansion_ExpansionWotlk, "https://raw.githubusercontent.com/Hoizame/AtlasLootClassic/master/AtlasLootClassic_DungeonsAndRaids/data-wrath.lua")

	readAtlasLootFactionData(db, proto.Expansion_ExpansionWotlk, "https://raw.githubusercontent.com/Hoizame/AtlasLootClassic/master/AtlasLootClassic_Factions/data-wrath.lua")

	readZoneData(db)
	readFactionData(db, inputsDir)

	return db
}

func readAtlasLootSourceData(db *WowDatabase, expansion proto.Expansion, srcUrl string) {
	srcTxt, err := tools.ReadWeb(srcUrl)
	if err != nil {
		log.Fatalf("Error reading atlasloot file %s", err)
	}

	itemPattern := regexp.MustCompile(`^\[([0-9]+)\] = {(.*)},$`)
	typePattern := regexp.MustCompile(`\[3\] = (\d+),.*\[4\] = (\d+)`)
	lines := strings.Split(srcTxt, "\n")
	for _, line := range lines {
		match := itemPattern.FindStringSubmatch(line)
		if match != nil {
			idStr := match[1]
			id, _ := strconv.Atoi(idStr)
			item := &proto.UIItem{Id: int32(id), Expansion: expansion}
			if _, ok := db.Items[item.Id]; ok {
				continue
			}

			paramsStr := match[2]
			typeMatch := typePattern.FindStringSubmatch(paramsStr)
			if typeMatch != nil {
				itemType, _ := strconv.Atoi(typeMatch[1])
				spellID, _ := strconv.Atoi(typeMatch[2])
				if prof, ok := AtlasLootProfessionIDs[itemType]; ok {
					item.Sources = append(item.Sources, &proto.UIItemSource{
						Source: &proto.UIItemSource_Crafted{
							Crafted: &proto.CraftedSource{
								Profession: prof,
								SpellId:    int32(spellID),
							},
						},
					})
				}
			}

			db.MergeItem(item)
		}
	}
}

func readAtlasLootDungeonData(db *WowDatabase, expansion proto.Expansion, srcUrl string) {
	srcTxt, err := tools.ReadWeb(srcUrl)
	if err != nil {
		log.Fatalf("Error reading atlasloot file %s", err)
	}

	// Convert newline to '@@@' so we can do regexes on the whole file as 1 line.
	srcTxt = strings.ReplaceAll(srcTxt, "\n", "@@@")

	dungeonPattern := regexp.MustCompile(`data\["([^"]+)"] = {.*?\sMapID = (\d+),.*?items = {(.*?)@@@}@@@`)
	npcNameAndIDPattern := regexp.MustCompile(`^[^@]*?AL\["(.*?)"\]\)?,(.*?(@@@\s*npcID = {?(\d+),))?`)
	diffItemsPattern := regexp.MustCompile(`\[([A-Z0-9]+_DIFF)\] = (({.*?@@@\s*},?@@@)|(.*?@@@\s*\),?@@@))`)
	itemsPattern := regexp.MustCompile(`@@@\s+{(.*?)},`)
	itemParamPattern := regexp.MustCompile(`AL\["(.*?)"\]`)
	for _, dungeonMatch := range dungeonPattern.FindAllStringSubmatch(srcTxt, -1) {
		fmt.Printf("Zone: %s\n", dungeonMatch[1])
		zoneID, _ := strconv.Atoi(dungeonMatch[2])
		db.MergeZone(&proto.UIZone{
			Id:        int32(zoneID),
			Expansion: expansion,
		})

		npcSplits := strings.Split(dungeonMatch[3], "name = ")[1:]
		for _, npcSplit := range npcSplits {
			npcSplit = strings.ReplaceAll(npcSplit, "AtlasLoot:GetRetByFaction(", "")
			npcMatch := npcNameAndIDPattern.FindStringSubmatch(npcSplit)
			if npcMatch == nil {
				panic("No npc match: " + npcSplit)
			}
			npcName := npcMatch[1]
			npcID := 0
			if len(npcMatch) > 3 {
				npcID, _ = strconv.Atoi(npcMatch[4])
			}
			if npcName == "Onyxia" { // AtlasLoot uses 15956 for some reason, which is the ID for Anub'Rekan.
				npcID = 10184
			} else if npcName == "Yogg-Saron" { // AtlasLoot uses 33271 for some reason, which is the ID for General Vezax.
				npcID = 33288
			}
			fmt.Printf("NPC: %s/%d\n", npcName, npcID)
			if npcID != 0 {
				db.MergeNpc(&proto.UINPC{
					Id:     int32(npcID),
					ZoneId: int32(zoneID),
					Name:   npcName,
				})
			}

			for _, difficultyMatch := range diffItemsPattern.FindAllStringSubmatch(npcSplit, -1) {
				difficulty, ok := AtlasLootDifficulties[difficultyMatch[1]]
				if !ok {
					log.Fatalf("Invalid difficulty for NPC %s: %s", npcName, difficultyMatch[1])
				}

				curCategory := ""
				curLocation := 0

				for _, itemMatch := range itemsPattern.FindAllStringSubmatch(difficultyMatch[0], -1) {
					itemParams := core.MapSlice(strings.Split(itemMatch[1], ","), strings.TrimSpace)
					location, _ := strconv.Atoi(itemParams[0]) // Location within AtlasLoot's menu.

					idStr := itemParams[1]
					if idStr[0] == 'n' || idStr[0] == '"' { // nil or "xxx"
						if len(itemParams) > 3 {
							if paramMatch := itemParamPattern.FindStringSubmatch(itemParams[3]); paramMatch != nil {
								curCategory = paramMatch[1]
								curLocation = location
							}
						}
						if len(itemParams) > 4 {
							if paramMatch := itemParamPattern.FindStringSubmatch(itemParams[4]); paramMatch != nil {
								curCategory = paramMatch[1]
								curLocation = location
							}
						}
					} else { // item ID
						itemID, _ := strconv.Atoi(idStr)
						//fmt.Printf("Item: %d\n", itemID)
						dropSource := &proto.DropSource{
							Difficulty: difficulty,
							ZoneId:     int32(zoneID),
						}
						if npcID == 0 {
							dropSource.OtherName = npcName
						} else {
							dropSource.NpcId = int32(npcID)
						}

						if curCategory != "" && location == curLocation+1 {
							curLocation = location
							dropSource.Category = curCategory
						}

						item := &proto.UIItem{Id: int32(itemID), Sources: []*proto.UIItemSource{{
							Source: &proto.UIItemSource_Drop{
								Drop: dropSource,
							},
						}}}
						db.MergeItem(item)
					}
				}
			}
		}
	}
}

func readAtlasLootFactionData(db *WowDatabase, expansion proto.Expansion, srcUrl string) {
	srcTxt, err := tools.ReadWeb(srcUrl)
	if err != nil {
		log.Fatalf("Error reading atlasloot file %s", err)
	}

	// Convert newline to '@@@' so we can do regexes on the whole file as 1 line.
	regex := regexp.MustCompile(`\r?\n`)
	srcTxt = regex.ReplaceAllString(srcTxt, "@@@")
	srcTxt = strings.ReplaceAll(srcTxt, "Updated in SoD", "")

	factionpattern := regexp.MustCompile(`data\["([^"]+)"] = {.*?\sFactionID = (\d+),.*?items = {(.*?)@@@}@@@`)
	repLevelPattern := regexp.MustCompile(`{ -- (Friendly|Honored|Revered|Exalted)\s?[\d]?@@@\s+name =(.*?@@@\s+},?@@@\s+},?)`)
	itemsPattern := regexp.MustCompile(`@@@\s+{(.*?)},`)

	for _, factionMatch := range factionpattern.FindAllStringSubmatch(srcTxt, -1) {
		factionID, err := strconv.Atoi(factionMatch[2])
		if err != nil {
			fmt.Printf("Error reading faction %s\n", factionMatch[1])
			return
		}
		fmt.Printf("Faction: %s\n", factionMatch[1])

		db.MergeFaction(&proto.UIFaction{
			Id:        int32(factionID),
			Expansion: expansion,
		})

		for _, repLevelMatch := range repLevelPattern.FindAllStringSubmatch(factionMatch[3], -1) {
			repLevel := repLevelMatch[1]
			fmt.Printf("Reputation: %s\n", repLevel)

			for _, itemMatch := range itemsPattern.FindAllStringSubmatch(repLevelMatch[2], -1) {
				itemParams := core.MapSlice(strings.Split(itemMatch[1], ","), strings.TrimSpace)

				idStr := itemParams[1]
				itemID, _ := strconv.Atoi(idStr)

				if itemID != 0 {
					// fmt.Printf("Item: %d\n", itemID)
					repSource := &proto.RepSource{
						RepFactionId: int32(factionID),
						RepLevel:     AtlasLootRepLevels[repLevel],
					}

					item := &proto.UIItem{Id: int32(itemID)}
					item.Sources = append(item.Sources, &proto.UIItemSource{
						Source: &proto.UIItemSource_Rep{
							Rep: repSource,
						},
					})

					db.MergeItem(item)
				}
			}
		}
	}
}

func readZoneData(db *WowDatabase) {
	zoneIDs := make([]int32, 0, len(db.Zones))
	for zoneID := range db.Zones {
		zoneIDs = append(zoneIDs, zoneID)
	}
	zoneIDStrs := core.MapSlice(zoneIDs, func(zoneID int32) string { return strconv.Itoa(int(zoneID)) })

	zoneTM := &WowheadTooltipManager{
		TooltipManager{
			FilePath:   "",
			UrlPattern: "https://nether.wowhead.com/wotlk/tooltip/zone/%s",
		},
	}
	zoneTooltips := zoneTM.FetchFromWeb(zoneIDStrs)

	tooltipPattern := regexp.MustCompile(`{"name":"(.*?)",`)
	for i, zoneID := range zoneIDs {
		tooltip := zoneTooltips[zoneIDStrs[i]]
		match := tooltipPattern.FindStringSubmatch(tooltip)
		if match == nil {
			log.Fatalf("Error parsing zone tooltip %s", tooltip)
		}
		db.Zones[zoneID].Name = match[1]
	}
}

type FactionConfig struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

func readFactionData(db *WowDatabase, inputsDir string) {
	data, err := os.ReadFile(fmt.Sprintf("%s/factions.json", inputsDir))
	if err != nil {
		log.Fatalf("failed to load talent json file: %s", err)
	}

	var buf bytes.Buffer
	err = json.Compact(&buf, []byte(data))
	if err != nil {
		log.Fatalf("failed to compact json: %s", err)
	}

	var jsonFactions []FactionConfig

	err = json.Unmarshal(buf.Bytes(), &jsonFactions)
	if err != nil {
		log.Fatalf("failed to parse talent to json %s", err)
	}

	for _, factionConfig := range jsonFactions {
		if db.Factions[factionConfig.Id] != nil {
			db.Factions[factionConfig.Id].Name = factionConfig.Name
		}
	}
}

var AtlasLootProfessionIDs = map[int]proto.Profession{
	//4: proto.Profession_FirstAid,
	5: proto.Profession_Blacksmithing,
	6: proto.Profession_Leatherworking,
	7: proto.Profession_Alchemy,
	//9: proto.Profession_Cooking,
	10: proto.Profession_Mining,
	11: proto.Profession_Tailoring,
	12: proto.Profession_Engineering,
	13: proto.Profession_Enchanting,
	17: proto.Profession_Jewelcrafting,
	18: proto.Profession_Inscription,
}
var AtlasLootDifficulties = map[string]proto.DungeonDifficulty{
	"NORMAL_DIFF":  proto.DungeonDifficulty_DifficultyNormal,
	"HEROIC_DIFF":  proto.DungeonDifficulty_DifficultyHeroic,
	"ALPHA_DIFF":   proto.DungeonDifficulty_DifficultyTitanRuneAlpha,
	"BETA_DIFF":    proto.DungeonDifficulty_DifficultyTitanRuneBeta,
	"RAID10_DIFF":  proto.DungeonDifficulty_DifficultyRaid10,
	"RAID10H_DIFF": proto.DungeonDifficulty_DifficultyRaid10H,
	"RAID25_DIFF":  proto.DungeonDifficulty_DifficultyRaid25,
	"RAID25H_DIFF": proto.DungeonDifficulty_DifficultyRaid25H,
}

var AtlasLootRepLevels = map[string]proto.RepLevel{
	"Hated":      proto.RepLevel_RepLevelHated,
	"Hostile":    proto.RepLevel_RepLevelHostile,
	"Unfriendly": proto.RepLevel_RepLevelUnfriendly,
	"Neutral":    proto.RepLevel_RepLevelNeutral,
	"Friendly":   proto.RepLevel_RepLevelFriendly,
	"Honored":    proto.RepLevel_RepLevelHonored,
	"Revered":    proto.RepLevel_RepLevelRevered,
	"Exalted":    proto.RepLevel_RepLevelExalted,
}
