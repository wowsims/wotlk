package database

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/tailscale/hujson"
	"github.com/wowsims/sod/sim/core/proto"
)

// Example db input file: https://nether.wowhead.com/classic/data/gear-planner?dv=100

func ParseWowheadDB(dbContents string) WowheadDatabase {
	var wowheadDB WowheadDatabase

	// Each part looks like 'WH.setPageData("wow.gearPlanner.some.name", {......});'
	parts := strings.Split(dbContents, "WH.setPageData(")

	for _, dbPart := range parts {
		// fmt.Printf("Part len: %d\n", len(dbPart))
		if len(dbPart) < 10 {
			continue
		}
		dbPart = strings.TrimSpace(dbPart)
		dbPart = strings.TrimRight(dbPart, ");")

		if dbPart[0] != '"' {
			continue
		}
		secondQuoteIdx := strings.Index(dbPart[1:], "\"")
		if secondQuoteIdx == -1 {
			continue
		}
		dbName := dbPart[1 : secondQuoteIdx+1]
		// fmt.Printf("DB name: %s\n", dbName)

		commaIdx := strings.Index(dbPart, ",")
		dbContents := dbPart[commaIdx+1:]
		if dbName == "wow.gearPlanner.classic.item" {
			standardized, err := hujson.Standardize([]byte(dbContents)) // Removes invalid JSON, such as trailing commas
			if err != nil {
				log.Fatalf("Failed to standardize json %s\n\n%s\n\n%s", err, dbContents[0:30], dbContents[len(dbContents)-30:])
			}

			err = json.Unmarshal(standardized, &wowheadDB.Items)
			if err != nil {
				log.Fatalf("failed to parse wowhead item db to json %s\n\n%s", err, dbContents[0:30])
			}
		}

		if dbName == "wow.gearPlanner.classic.randomEnchant" {
			standardized, err := hujson.Standardize([]byte(dbContents)) // Removes invalid JSON, such as trailing commas
			if err != nil {
				log.Fatalf("Failed to standardize json %s\n\n%s\n\n%s", err, dbContents[0:30], dbContents[len(dbContents)-30:])
			}

			err = json.Unmarshal(standardized, &wowheadDB.RandomSuffixes)
			if err != nil {
				log.Fatalf("failed to parse wowhead random suffix db to json %s\n\n%s", err, dbContents[0:30])
			}
		}
	}

	fmt.Printf("\n--\nWowhead DB items loaded: %d\n--\n", len(wowheadDB.Items))
	fmt.Printf("\n--\nWowhead DB random suffixes loaded: %d\n--\n", len(wowheadDB.RandomSuffixes))

	return wowheadDB
}

type WowheadDatabase struct {
	Items          map[string]WowheadItem
	RandomSuffixes map[string]WowheadRandomSuffix
}

type WowheadRandomSuffix struct {
	ID    int32                    `json:"id"`
	Name  string                   `json:"name"`
	Stats WowheadRandomSuffixStats `json:"stats"`
}

type WowheadRandomSuffixStats struct {
	Armor             int32 `json:"armor"`
	Strength          int32 `json:"str"`
	Agility           int32 `json:"agi"`
	Stamina           int32 `json:"sta"`
	Intellect         int32 `json:"int"`
	Spirit            int32 `json:"spi"`
	SpellPower        int32 `json:"spldmg"`
	ArcanePower       int32 `json:"arcsplpwr"`
	FirePower         int32 `json:"firsplpwr"`
	FrostPower        int32 `json:"frosplpwr"`
	HolyPower         int32 `json:"holsplpwr"`
	NaturePower       int32 `json:"natsplpwr"`
	ShadowPower       int32 `json:"shasplpwr"`
	MeleeCrit         int32 `json:"mlecritstrkpct"`
	MP5               int32 `json:"manargn"`
	AttackPower       int32 `json:"mleatkpwr"`
	RangedAttackPower int32 `json:"rgdatkpwr"`
	Defense           int32 `json:"def"`
	Block             int32 `json:"blockpct"`
	Dodge             int32 `json:"dodgepct"`
	ArcaneResistance  int32 `json:"arcres"`
	FireResistance    int32 `json:"firres"`
	FrostResistance   int32 `json:"frores"`
	NatureResistance  int32 `json:"natres"`
	ShadowResistance  int32 `json:"shares"`
	Healing           int32 `json:"splheal"`
}

func (wrs WowheadRandomSuffix) ToProto() *proto.ItemRandomSuffix {
	stats := Stats{
		proto.Stat_StatArmor:             float64(wrs.Stats.Armor),
		proto.Stat_StatStrength:          float64(wrs.Stats.Strength),
		proto.Stat_StatAgility:           float64(wrs.Stats.Agility),
		proto.Stat_StatStamina:           float64(wrs.Stats.Stamina),
		proto.Stat_StatIntellect:         float64(wrs.Stats.Intellect),
		proto.Stat_StatSpirit:            float64(wrs.Stats.Spirit),
		proto.Stat_StatSpellPower:        float64(wrs.Stats.SpellPower),
		proto.Stat_StatArcanePower:       float64(wrs.Stats.ArcanePower),
		proto.Stat_StatFirePower:         float64(wrs.Stats.FirePower),
		proto.Stat_StatFrostPower:        float64(wrs.Stats.FrostPower),
		proto.Stat_StatHolyPower:         float64(wrs.Stats.HolyPower),
		proto.Stat_StatNaturePower:       float64(wrs.Stats.NaturePower),
		proto.Stat_StatShadowPower:       float64(wrs.Stats.ShadowPower),
		proto.Stat_StatMeleeCrit:         float64(wrs.Stats.MeleeCrit),
		proto.Stat_StatMP5:               float64(wrs.Stats.MP5),
		proto.Stat_StatAttackPower:       float64(wrs.Stats.AttackPower),
		proto.Stat_StatRangedAttackPower: float64(wrs.Stats.RangedAttackPower),
		proto.Stat_StatDefense:           float64(wrs.Stats.Defense),
		proto.Stat_StatBlock:             float64(wrs.Stats.Block),
		proto.Stat_StatDodge:             float64(wrs.Stats.Dodge),
		proto.Stat_StatArcaneResistance:  float64(wrs.Stats.ArcaneResistance),
		proto.Stat_StatFireResistance:    float64(wrs.Stats.FireResistance),
		proto.Stat_StatFrostResistance:   float64(wrs.Stats.FrostResistance),
		proto.Stat_StatNatureResistance:  float64(wrs.Stats.NatureResistance),
		proto.Stat_StatShadowResistance:  float64(wrs.Stats.ShadowResistance),
		proto.Stat_StatHealing:           float64(wrs.Stats.Healing),
	}

	return &proto.ItemRandomSuffix{
		Id:    wrs.ID,
		Name:  wrs.Name,
		Stats: toSlice(stats),
	}
}

type WowheadItem struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`

	Quality       int32 `json:"quality"`
	Ilvl          int32 `json:"itemLevel"`
	Phase         int32 `json:"contentPhase"`
	RequiresLevel int32 `json:"reqlevel"`

	Stats               WowheadItemStats `json:"stats"`
	RandomSuffixOptions []int32          `json:"randomEnchants"`

	SourceTypes   []int32             `json:"source"` // 1 = Crafted, 2 = Dropped by, 3 = sold by zone vendor? barely used, 4 = Quest, 5 = Sold by
	SourceDetails []WowheadItemSource `json:"sourcemore"`
}
type WowheadItemStats struct {
	Armor int32 `json:"armor"`
}
type WowheadItemSource struct {
	C        int32  `json:"c"`
	Name     string `json:"n"`    // Name of crafting spell
	Icon     string `json:"icon"` // Icon corresponding to the named entity
	EntityID int32  `json:"ti"`   // Crafting Spell ID / NPC ID / ?? / Quest ID
	ZoneID   int32  `json:"z"`    // Only for drop / sold by sources
}

func (wi WowheadItem) ToProto() *proto.UIItem {
	var sources []*proto.UIItemSource
	for i, details := range wi.SourceDetails {
		switch wi.SourceTypes[i] {
		case 1: // Crafted
			// We'll get this from AtlasLoot instead because it can also tell us the profession.
			//sources = append(sources, &proto.UIItemSource{
			//	Source: &proto.UIItemSource_Crafted{
			//		Crafted: &proto.CraftedSource{
			//			SpellId: details.EntityID,
			//		},
			//	},
			//})
		case 2: // Dropped by
			// Do nothing, we'll get this from AtlasLoot.
		case 3: // Sold by zone vendor? barely used
		case 4: // Quest
			sources = append(sources, &proto.UIItemSource{
				Source: &proto.UIItemSource_Quest{
					Quest: &proto.QuestSource{
						Id:   details.EntityID,
						Name: details.Name,
					},
				},
			})
		case 5: // Sold by
			sources = append(sources, &proto.UIItemSource{
				Source: &proto.UIItemSource_SoldBy{
					SoldBy: &proto.SoldBySource{
						NpcId:   details.EntityID,
						NpcName: details.Name,
						ZoneId:  details.ZoneID,
					},
				},
			})
		}
	}

	return &proto.UIItem{
		Id:                  wi.ID,
		Name:                wi.Name,
		Icon:                wi.Icon,
		Ilvl:                wi.Ilvl,
		Phase:               wi.Phase,
		RequiresLevel:       wi.RequiresLevel,
		Sources:             sources,
		RandomSuffixOptions: wi.RandomSuffixOptions,
	}
}
