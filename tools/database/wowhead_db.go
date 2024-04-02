package database

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/tailscale/hujson"
	"github.com/wowsims/wotlk/sim/core/proto"
)

// Example db input file: https://nether.wowhead.com/wotlk/data/gear-planner?dv=100

func ParseWowheadDB(dbContents string) WowheadDatabase {
	var wowheadDB WowheadDatabase

	// Each part looks like 'WH.setPageData("wow.gearPlanner.some.name", {......});'
	parts := strings.Split(dbContents, "WH.setPageData(")

	for _, dbPart := range parts {
		//fmt.Printf("Part len: %d\n", len(dbPart))
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
		//fmt.Printf("DB name: %s\n", dbName)

		commaIdx := strings.Index(dbPart, ",")
		dbContents := dbPart[commaIdx+1:]
		if dbName == "wow.gearPlanner.wrath.item" {
			standardized, err := hujson.Standardize([]byte(dbContents)) // Removes invalid JSON, such as trailing commas
			if err != nil {
				log.Fatalf("Failed to standardize json %s\n\n%s\n\n%s", err, dbContents[0:30], dbContents[len(dbContents)-30:])
			}

			err = json.Unmarshal(standardized, &wowheadDB.Items)
			if err != nil {
				log.Fatalf("failed to parse wowhead item db to json %s\n\n%s", err, dbContents[0:30])
			}
		}
	}

	fmt.Printf("\n--\nWowhead DB items loaded: %d\n--\n", len(wowheadDB.Items))

	return wowheadDB
}

type WowheadDatabase struct {
	Items map[string]WowheadItem
}

type WowheadItem struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`

	Quality int32 `json:"quality"`
	Ilvl    int32 `json:"itemLevel"`
	Phase   int32 `json:"contentPhase"`

	RaceMask  uint16 `json:"raceMask"`
	ClassMask uint16 `json:"classMask"`

	Stats WowheadItemStats `json:"stats"`

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
		Id:                 wi.ID,
		Name:               wi.Name,
		Icon:               wi.Icon,
		Ilvl:               wi.Ilvl,
		Phase:              wi.Phase,
		FactionRestriction: wi.getFactionRstriction(),
		ClassAllowlist:     wi.getClassRestriction(),
		Sources:            sources,
	}
}

func (wi WowheadItem) getFactionRstriction() proto.UIItem_FactionRestriction {
	if wi.RaceMask == 1101 {
		return proto.UIItem_FACTION_RESTRICTION_ALLIANCE_ONLY
	} else if wi.RaceMask == 690 {
		return proto.UIItem_FACTION_RESTRICTION_HORDE_ONLY
	} else {
		return proto.UIItem_FACTION_RESTRICTION_UNSPECIFIED
	}
}

type ClassMask uint16

const (
	ClassMaskWarrior     ClassMask = 1 << iota
	ClassMaskPaladin               // 2
	ClassMaskHunter                // 4
	ClassMaskRogue                 // 8
	ClassMaskPriest                // 16
	ClassMaskDeathKnight           // 32
	ClassMaskShaman                // 64
	ClassMaskMage                  // 128
	ClassMaskWarlock               // 256
	ClassMaskUnknown               // 512 seemingly unused?
	ClassMaskDruid                 // 1024
)

func (wi WowheadItem) getClassRestriction() []proto.Class {
	classAllowlist := []proto.Class{}
	if wi.ClassMask&uint16(ClassMaskWarrior) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassWarrior)
	}
	if wi.ClassMask&uint16(ClassMaskPaladin) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassPaladin)
	}
	if wi.ClassMask&uint16(ClassMaskHunter) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassHunter)
	}
	if wi.ClassMask&uint16(ClassMaskRogue) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassRogue)
	}
	if wi.ClassMask&uint16(ClassMaskPriest) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassPriest)
	}
	if wi.ClassMask&uint16(ClassMaskDruid) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassDruid)
	}
	if wi.ClassMask&uint16(ClassMaskShaman) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassShaman)
	}
	if wi.ClassMask&uint16(ClassMaskMage) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassMage)
	}
	if wi.ClassMask&uint16(ClassMaskWarlock) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassWarlock)
	}
	if wi.ClassMask&uint16(ClassMaskDeathKnight) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassDeathknight)
	}

	return classAllowlist
}
