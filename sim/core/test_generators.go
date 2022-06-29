package core

import (
	"fmt"
	"strings"

	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/proto"
	googleProto "google.golang.org/protobuf/proto"
)

type SingleCharacterStatsTestGenerator struct {
	Name    string
	Request *proto.ComputeStatsRequest
}

func (generator *SingleCharacterStatsTestGenerator) NumTests() int {
	return 1
}
func (generator *SingleCharacterStatsTestGenerator) GetTest(testIdx int) (string, *proto.ComputeStatsRequest, *proto.StatWeightsRequest, *proto.RaidSimRequest) {
	return generator.Name, generator.Request, nil, nil
}

type SingleStatWeightsTestGenerator struct {
	Name    string
	Request *proto.StatWeightsRequest
}

func (generator *SingleStatWeightsTestGenerator) NumTests() int {
	return 1
}
func (generator *SingleStatWeightsTestGenerator) GetTest(testIdx int) (string, *proto.ComputeStatsRequest, *proto.StatWeightsRequest, *proto.RaidSimRequest) {
	return generator.Name, nil, generator.Request, nil
}

type SingleDpsTestGenerator struct {
	Name    string
	Request *proto.RaidSimRequest
}

func (generator *SingleDpsTestGenerator) NumTests() int {
	return 1
}
func (generator *SingleDpsTestGenerator) GetTest(testIdx int) (string, *proto.ComputeStatsRequest, *proto.StatWeightsRequest, *proto.RaidSimRequest) {
	return generator.Name, nil, nil, generator.Request
}

type GearSetCombo struct {
	Label   string
	GearSet *proto.EquipmentSpec
}
type SpecOptionsCombo struct {
	Label       string
	SpecOptions interface{}
}
type BuffsCombo struct {
	Label    string
	Raid     *proto.RaidBuffs
	Party    *proto.PartyBuffs
	Debuffs  *proto.Debuffs
	Player   *proto.IndividualBuffs
	Consumes *proto.Consumes
}
type EncounterCombo struct {
	Label     string
	Encounter *proto.Encounter
}
type SettingsCombos struct {
	Class       proto.Class
	Races       []proto.Race
	GearSets    []GearSetCombo
	SpecOptions []SpecOptionsCombo
	Buffs       []BuffsCombo
	Encounters  []EncounterCombo
	SimOptions  *proto.SimOptions
}

func (combos *SettingsCombos) NumTests() int {
	return len(combos.Races) * len(combos.GearSets) * len(combos.SpecOptions) * len(combos.Buffs) * len(combos.Encounters)
}

func (combos *SettingsCombos) GetTest(testIdx int) (string, *proto.ComputeStatsRequest, *proto.StatWeightsRequest, *proto.RaidSimRequest) {
	testNameParts := []string{}

	raceIdx := testIdx % len(combos.Races)
	testIdx /= len(combos.Races)
	race := combos.Races[raceIdx]
	testNameParts = append(testNameParts, race.String()[4:])

	gearSetIdx := testIdx % len(combos.GearSets)
	testIdx /= len(combos.GearSets)
	gearSetCombo := combos.GearSets[gearSetIdx]
	testNameParts = append(testNameParts, gearSetCombo.Label)

	specOptionsIdx := testIdx % len(combos.SpecOptions)
	testIdx /= len(combos.SpecOptions)
	specOptionsCombo := combos.SpecOptions[specOptionsIdx]
	testNameParts = append(testNameParts, specOptionsCombo.Label)

	buffsIdx := testIdx % len(combos.Buffs)
	testIdx /= len(combos.Buffs)
	buffsCombo := combos.Buffs[buffsIdx]
	testNameParts = append(testNameParts, buffsCombo.Label)

	encounterIdx := testIdx % len(combos.Encounters)
	testIdx /= len(combos.Encounters)
	encounterCombo := combos.Encounters[encounterIdx]
	testNameParts = append(testNameParts, encounterCombo.Label)

	rsr := &proto.RaidSimRequest{
		Raid: SinglePlayerRaidProto(
			WithSpec(&proto.Player{
				Race:      race,
				Class:     combos.Class,
				Equipment: gearSetCombo.GearSet,
				Consumes:  buffsCombo.Consumes,
				Buffs:     buffsCombo.Player,
				// TODO: Allow cooldowns in tests
				//Cooldowns: &proto.Cooldowns{
				//	Cooldowns: []*proto.Cooldown{
				//		&proto.Cooldown{
				//			Id: &proto.ActionID{
				//				RawId: &proto.ActionID_SpellId{
				//					SpellId: 12043,
				//				},
				//			},
				//			Timings: []float64{
				//				5,
				//			},
				//		},
				//	},
				//},
			}, specOptionsCombo.SpecOptions),
			buffsCombo.Party,
			buffsCombo.Raid,
			buffsCombo.Debuffs),
		Encounter:  encounterCombo.Encounter,
		SimOptions: combos.SimOptions,
	}

	return strings.Join(testNameParts, "-"), nil, nil, rsr
}

// Returns all items that meet the given conditions.
type ItemFilter struct {
	// If set to ClassUnknown, any class is fine.
	Class proto.Class

	ArmorType proto.ArmorType

	// Blank list allows any value. Otherwise item must match 1 value from the list.
	WeaponTypes       []proto.WeaponType
	HandTypes         []proto.HandType
	RangedWeaponTypes []proto.RangedWeaponType

	// Item IDs to ignore.
	IDBlacklist []int32
}

// Returns whether the given item matches the conditions of this filter.
//
// If equipChecksOnly is true, will only check conditions related to whether
// the item is equippable.
func (filter *ItemFilter) Matches(item items.Item, equipChecksOnly bool) bool {
	if filter.Class != proto.Class_ClassUnknown && len(item.ClassAllowlist) > 0 {
		found := false
		for _, class := range item.ClassAllowlist {
			if class == filter.Class {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if item.Type == proto.ItemType_ItemTypeWeapon {
		if len(filter.WeaponTypes) > 0 {
			found := false
			for _, weaponType := range filter.WeaponTypes {
				if weaponType == item.WeaponType {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}

		if len(filter.HandTypes) > 0 {
			found := false
			for _, handType := range filter.HandTypes {
				if handType == item.HandType {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	} else if item.Type == proto.ItemType_ItemTypeRanged {
		if len(filter.RangedWeaponTypes) > 0 {
			found := false
			for _, rangedWeaponType := range filter.RangedWeaponTypes {
				if rangedWeaponType == item.RangedWeaponType {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	} else {
		if filter.ArmorType != proto.ArmorType_ArmorTypeUnknown {
			if item.ArmorType > filter.ArmorType {
				return false
			}
		}
	}

	if !equipChecksOnly {
		if !HasItemEffect(item.ID) {
			return false
		}

		if len(filter.IDBlacklist) > 0 {
			for _, itemID := range filter.IDBlacklist {
				if itemID == item.ID {
					return false
				}
			}
		}
	}

	return true
}

func (filter *ItemFilter) FindAllItems() []items.Item {
	filteredItems := []items.Item{}

	for _, item := range items.ByID {
		if filter.Matches(item, false) {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems
}

func (filter *ItemFilter) FindAllSets() []*ItemSet {
	filteredSets := []*ItemSet{}

	for _, set := range GetAllItemSets() {
		firstItem := items.ByID[set.ItemIDs()[0]]
		if filter.Matches(firstItem, true) {
			filteredSets = append(filteredSets, set)
		}
	}

	return filteredSets
}

func (filter *ItemFilter) FindAllMetaGems() []items.Gem {
	filteredGems := []items.Gem{}

	for _, gem := range items.GemsByID {
		if gem.Color == proto.GemColor_GemColorMeta {
			filteredGems = append(filteredGems, gem)
		}
	}

	return filteredGems
}

type ItemsTestGenerator struct {
	// Fields describing the base API request.
	Player     *proto.Player
	PartyBuffs *proto.PartyBuffs
	RaidBuffs  *proto.RaidBuffs
	Debuffs    *proto.Debuffs
	Encounter  *proto.Encounter
	SimOptions *proto.SimOptions

	// Some fields are populated automatically.
	ItemFilter ItemFilter

	initialized bool

	items []items.Item
	sets  []*ItemSet

	metagems []items.Gem

	metaSocketIdx int
}

func (generator *ItemsTestGenerator) init() {
	if generator.initialized {
		return
	}
	generator.initialized = true

	generator.ItemFilter.Class = generator.Player.Class
	if generator.ItemFilter.IDBlacklist == nil {
		generator.ItemFilter.IDBlacklist = []int32{}
	}
	for _, itemSpec := range generator.Player.Equipment.Items {
		generator.ItemFilter.IDBlacklist = append(generator.ItemFilter.IDBlacklist, itemSpec.Id)
	}

	generator.items = generator.ItemFilter.FindAllItems()
	generator.sets = generator.ItemFilter.FindAllSets()

	baseEquipment := items.ProtoToEquipment(*generator.Player.Equipment)
	generator.metaSocketIdx = -1
	for i, socketColor := range baseEquipment[proto.ItemSlot_ItemSlotHead].GemSockets {
		if socketColor == proto.GemColor_GemColorMeta {
			generator.metaSocketIdx = i
			break
		}
	}
	if generator.metaSocketIdx == -1 {
		return
	}
	generator.metagems = generator.ItemFilter.FindAllMetaGems()
}

func (generator *ItemsTestGenerator) NumTests() int {
	generator.init()
	return len(generator.items) + len(generator.sets) + len(generator.metagems)
}

func (generator *ItemsTestGenerator) GetTest(testIdx int) (string, *proto.ComputeStatsRequest, *proto.StatWeightsRequest, *proto.RaidSimRequest) {
	generator.init()
	label := ""

	playerCopy := googleProto.Clone(generator.Player).(*proto.Player)
	equipment := items.ProtoToEquipment(*playerCopy.Equipment)
	if testIdx < len(generator.items) {
		testItem := generator.items[testIdx]
		equipment.EquipItem(generator.items[testIdx])
		label = fmt.Sprintf("%s-%d", strings.ReplaceAll(testItem.Name, " ", ""), testItem.ID)
	} else if testIdx < len(generator.items)+len(generator.sets) {
		testSet := generator.sets[testIdx-len(generator.items)]
		for _, itemID := range testSet.ItemIDs() {
			setItem := items.ByID[itemID]
			equipment.EquipItem(setItem)
		}
		label = strings.ReplaceAll(testSet.Name, " ", "")
	} else {
		testMetaGem := generator.metagems[testIdx-len(generator.items)-len(generator.sets)]
		headItem := &equipment[proto.ItemSlot_ItemSlotHead]
		for len(headItem.Gems) <= generator.metaSocketIdx {
			headItem.Gems = append(headItem.Gems, items.Gem{})
		}
		headItem.Gems[generator.metaSocketIdx] = testMetaGem
		label = strings.ReplaceAll(testMetaGem.Name, " ", "")
	}
	playerCopy.Equipment = equipment.ToEquipmentSpecProto()

	rsr := &proto.RaidSimRequest{
		Raid: SinglePlayerRaidProto(
			playerCopy,
			generator.PartyBuffs,
			generator.RaidBuffs,
			generator.Debuffs),
		Encounter:  generator.Encounter,
		SimOptions: generator.SimOptions,
	}

	return label, nil, nil, rsr
}

type SubGenerator struct {
	name      string
	generator TestGenerator
}

type CombinedTestGenerator struct {
	subgenerators []SubGenerator
}

func (generator *CombinedTestGenerator) NumTests() int {
	total := 0
	for _, child := range generator.subgenerators {
		total += child.generator.NumTests()
	}
	return total
}

func (generator *CombinedTestGenerator) GetTest(testIdx int) (string, *proto.ComputeStatsRequest, *proto.StatWeightsRequest, *proto.RaidSimRequest) {
	remaining := testIdx
	for _, child := range generator.subgenerators {
		numTests := child.generator.NumTests()
		if remaining < numTests {
			testName, csr, swr, rsr := child.generator.GetTest(remaining)
			return child.name + "-" + testName, csr, swr, rsr
		}
		remaining -= numTests
	}

	panic("Invalid testIdx")
}

type CharacterSuiteConfig struct {
	Class proto.Class

	Race        proto.Race
	GearSet     GearSetCombo
	SpecOptions SpecOptionsCombo

	RaidBuffs   *proto.RaidBuffs
	PartyBuffs  *proto.PartyBuffs
	PlayerBuffs *proto.IndividualBuffs
	Consumes    *proto.Consumes
	Debuffs     *proto.Debuffs

	IsTank          bool
	InFrontOfTarget bool

	OtherRaces       []proto.Race
	OtherGearSets    []GearSetCombo
	OtherSpecOptions []SpecOptionsCombo

	ItemFilter ItemFilter

	StatsToWeigh    []proto.Stat
	EPReferenceStat proto.Stat
}

func FullCharacterTestSuiteGenerator(config CharacterSuiteConfig) TestGenerator {
	allRaces := append(config.OtherRaces, config.Race)
	allGearSets := append(config.OtherGearSets, config.GearSet)
	allSpecOptions := append(config.OtherSpecOptions, config.SpecOptions)

	defaultPlayer := WithSpec(
		&proto.Player{
			Class:     config.Class,
			Race:      config.Race,
			Equipment: config.GearSet.GearSet,
			Consumes:  config.Consumes,
			Buffs:     config.PlayerBuffs,

			InFrontOfTarget: config.InFrontOfTarget,
		},
		config.SpecOptions.SpecOptions)

	defaultRaid := SinglePlayerRaidProto(defaultPlayer, config.PartyBuffs, config.RaidBuffs, config.Debuffs)
	if config.IsTank {
		defaultRaid.Tanks = append(defaultRaid.Tanks, &proto.RaidTarget{TargetIndex: 0})
	}

	generator := &CombinedTestGenerator{
		subgenerators: []SubGenerator{
			SubGenerator{
				name: "CharacterStats",
				generator: &SingleCharacterStatsTestGenerator{
					Name: "Default",
					Request: &proto.ComputeStatsRequest{
						Raid: defaultRaid,
					},
				},
			},
			SubGenerator{
				name: "Settings",
				generator: &SettingsCombos{
					Class:       config.Class,
					Races:       allRaces,
					GearSets:    allGearSets,
					SpecOptions: allSpecOptions,
					Buffs: []BuffsCombo{
						BuffsCombo{
							Label: "NoBuffs",
						},
						BuffsCombo{
							Label:    "FullBuffs",
							Raid:     config.RaidBuffs,
							Party:    config.PartyBuffs,
							Debuffs:  config.Debuffs,
							Player:   config.PlayerBuffs,
							Consumes: config.Consumes,
						},
					},
					Encounters: MakeDefaultEncounterCombos(config.Debuffs),
					SimOptions: DefaultSimTestOptions,
				},
			},
			SubGenerator{
				name: "AllItems",
				generator: &ItemsTestGenerator{
					Player:     defaultPlayer,
					RaidBuffs:  config.RaidBuffs,
					PartyBuffs: config.PartyBuffs,
					Debuffs:    config.Debuffs,
					Encounter:  MakeSingleTargetEncounter(0),
					SimOptions: DefaultSimTestOptions,
					ItemFilter: config.ItemFilter,
				},
			},
		},
	}

	newRaid := googleProto.Clone(defaultRaid).(*proto.Raid)
	newRaid.Parties[0].Players[0].InFrontOfTarget = !newRaid.Parties[0].Players[0].InFrontOfTarget

	generator.subgenerators = append(generator.subgenerators, SubGenerator{
		name: "SwitchInFrontOfTarget",
		generator: &SingleDpsTestGenerator{
			Name: "Default",
			Request: &proto.RaidSimRequest{
				Raid:       newRaid,
				Encounter:  MakeSingleTargetEncounter(0),
				SimOptions: DefaultSimTestOptions,
			},
		},
	})

	if len(config.StatsToWeigh) > 0 {
		generator.subgenerators = append(generator.subgenerators, SubGenerator{
			name: "StatWeights",
			generator: &SingleStatWeightsTestGenerator{
				Name: "Default",
				Request: &proto.StatWeightsRequest{
					Player:     defaultPlayer,
					RaidBuffs:  config.RaidBuffs,
					PartyBuffs: config.PartyBuffs,
					Debuffs:    config.Debuffs,
					Encounter:  MakeSingleTargetEncounter(0),
					SimOptions: StatWeightsDefaultSimTestOptions,
					Tanks:      defaultRaid.Tanks,

					StatsToWeigh:    config.StatsToWeigh,
					EpReferenceStat: config.EPReferenceStat,
				},
			},
		})
	}

	if config.Consumes.Drums == proto.Drums_DrumsUnknown {
		newRaid := googleProto.Clone(defaultRaid).(*proto.Raid)
		newRaid.Parties[0].Players[0].Consumes.Drums = proto.Drums_DrumsOfBattle

		generator.subgenerators = append(generator.subgenerators, SubGenerator{
			name: "SelfDrums",
			generator: &SingleDpsTestGenerator{
				Name: "DPS",
				Request: &proto.RaidSimRequest{
					Raid:       newRaid,
					Encounter:  MakeSingleTargetEncounter(0),
					SimOptions: DefaultSimTestOptions,
				},
			},
		})
	}

	// Add this separately so it's always last, which makes it easy to find in the
	// displayed test results.
	generator.subgenerators = append(generator.subgenerators, SubGenerator{
		name: "Average",
		generator: &SingleDpsTestGenerator{
			Name: "Default",
			Request: &proto.RaidSimRequest{
				Raid:       defaultRaid,
				Encounter:  MakeSingleTargetEncounter(5),
				SimOptions: AverageDefaultSimTestOptions,
			},
		},
	})

	return generator
}
