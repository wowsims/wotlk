package core

import (
	"fmt"
	"slices"
	"strings"

	"github.com/wowsims/wotlk/sim/core/proto"
	googleProto "google.golang.org/protobuf/proto"
)

type SingleCharacterStatsTestGenerator struct {
	Name    string
	Request *proto.ComputeStatsRequest
}

func (generator *SingleCharacterStatsTestGenerator) NumTests() int {
	return 1
}
func (generator *SingleCharacterStatsTestGenerator) GetTest(_ int) (string, *proto.ComputeStatsRequest, *proto.StatWeightsRequest, *proto.RaidSimRequest) {
	return generator.Name, generator.Request, nil, nil
}

type SingleStatWeightsTestGenerator struct {
	Name    string
	Request *proto.StatWeightsRequest
}

func (generator *SingleStatWeightsTestGenerator) NumTests() int {
	return 1
}
func (generator *SingleStatWeightsTestGenerator) GetTest(_ int) (string, *proto.ComputeStatsRequest, *proto.StatWeightsRequest, *proto.RaidSimRequest) {
	return generator.Name, nil, generator.Request, nil
}

type SingleDpsTestGenerator struct {
	Name    string
	Request *proto.RaidSimRequest
}

func (generator *SingleDpsTestGenerator) NumTests() int {
	return 1
}
func (generator *SingleDpsTestGenerator) GetTest(_ int) (string, *proto.ComputeStatsRequest, *proto.StatWeightsRequest, *proto.RaidSimRequest) {
	return generator.Name, nil, nil, generator.Request
}

type RotationCastsTestGenerator struct {
	SpecOptions []SpecOptionsCombo
	PartyBuffs  *proto.PartyBuffs
	RaidBuffs   *proto.RaidBuffs
	Debuffs     *proto.Debuffs
	Player      *proto.Player
	Encounter   *proto.Encounter
	SimOptions  *proto.SimOptions
}

func (generator *RotationCastsTestGenerator) NumTests() int {
	return len(generator.SpecOptions)
}

func (generator *RotationCastsTestGenerator) GetTest(testIdx int) (string, *proto.ComputeStatsRequest, *proto.StatWeightsRequest, *proto.RaidSimRequest) {
	rsr := &proto.RaidSimRequest{
		Raid: SinglePlayerRaidProto(
			WithSpec(generator.Player, generator.SpecOptions[testIdx].SpecOptions),
			generator.PartyBuffs,
			generator.RaidBuffs,
			generator.Debuffs,
		),
		Encounter:  generator.Encounter,
		SimOptions: generator.SimOptions,
	}
	return generator.SpecOptions[testIdx].Label, nil, nil, rsr
}

type GearSetCombo struct {
	Label   string
	GearSet *proto.EquipmentSpec
}
type TalentsCombo struct {
	Label   string
	Talents string
	Glyphs  *proto.Glyphs
}
type SpecOptionsCombo struct {
	Label       string
	SpecOptions interface{}
}
type RotationCombo struct {
	Label    string
	Rotation *proto.APLRotation
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
	TalentSets  []TalentsCombo
	SpecOptions []SpecOptionsCombo
	Rotations   []RotationCombo
	Buffs       []BuffsCombo
	Encounters  []EncounterCombo
	SimOptions  *proto.SimOptions
	IsHealer    bool
	Cooldowns   *proto.Cooldowns
}

func (combos *SettingsCombos) NumTests() int {
	return len(combos.Races) * len(combos.GearSets) * len(combos.TalentSets) * len(combos.SpecOptions) * len(combos.Buffs) * len(combos.Encounters) * max(1, len(combos.Rotations))
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

	talentSetIdx := testIdx % len(combos.TalentSets)
	testIdx /= len(combos.TalentSets)
	talentSetCombo := combos.TalentSets[talentSetIdx]
	// We never use more than 1 talent combo, so it just makes the names longer.
	//testNameParts = append(testNameParts, talentSetCombo.Label)

	specOptionsIdx := testIdx % len(combos.SpecOptions)
	testIdx /= len(combos.SpecOptions)
	specOptionsCombo := combos.SpecOptions[specOptionsIdx]
	testNameParts = append(testNameParts, specOptionsCombo.Label)

	rotationsCombo := RotationCombo{Label: "None", Rotation: &proto.APLRotation{}}
	if len(combos.Rotations) > 0 {
		rotationsIdx := testIdx % len(combos.Rotations)
		testIdx /= len(combos.Rotations)
		rotationsCombo = combos.Rotations[rotationsIdx]
		testNameParts = append(testNameParts, rotationsCombo.Label)
	}

	buffsIdx := testIdx % len(combos.Buffs)
	testIdx /= len(combos.Buffs)
	buffsCombo := combos.Buffs[buffsIdx]
	testNameParts = append(testNameParts, buffsCombo.Label)

	encounterIdx := testIdx % len(combos.Encounters)
	encounterCombo := combos.Encounters[encounterIdx]
	testNameParts = append(testNameParts, encounterCombo.Label)

	rsr := &proto.RaidSimRequest{
		Raid: SinglePlayerRaidProto(
			WithSpec(&proto.Player{
				Race:               race,
				Class:              combos.Class,
				Equipment:          gearSetCombo.GearSet,
				TalentsString:      talentSetCombo.Talents,
				Glyphs:             talentSetCombo.Glyphs,
				Consumes:           buffsCombo.Consumes,
				Buffs:              buffsCombo.Player,
				Profession1:        proto.Profession_Engineering,
				Cooldowns:          combos.Cooldowns,
				Rotation:           rotationsCombo.Rotation,
				DistanceFromTarget: 30,
				ReactionTimeMs:     150,
				ChannelClipDelayMs: 50,
			}, specOptionsCombo.SpecOptions),
			buffsCombo.Party,
			buffsCombo.Raid,
			buffsCombo.Debuffs),
		Encounter:  encounterCombo.Encounter,
		SimOptions: combos.SimOptions,
	}
	if combos.IsHealer {
		rsr.Raid.TargetDummies = 1
	}

	return strings.Join(testNameParts, "-"), nil, nil, rsr
}

// Returns all items that meet the given conditions.
type ItemFilter struct {
	// If set to ClassUnknown, any class is fine.
	Class proto.Class

	ArmorType proto.ArmorType

	// Empty lists allows any value. Otherwise, item must match a value from the list.
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
func (filter *ItemFilter) Matches(item Item, equipChecksOnly bool) bool {
	if item.Type == proto.ItemType_ItemTypeWeapon {
		if len(filter.WeaponTypes) > 0 && !slices.Contains(filter.WeaponTypes, item.WeaponType) {
			return false
		}
		if len(filter.HandTypes) > 0 && !slices.Contains(filter.HandTypes, item.HandType) {
			return false
		}
	} else if item.Type == proto.ItemType_ItemTypeRanged {
		if len(filter.RangedWeaponTypes) > 0 && !slices.Contains(filter.RangedWeaponTypes, item.RangedWeaponType) {
			return false
		}
	} else {
		if filter.ArmorType != proto.ArmorType_ArmorTypeUnknown {
			if item.ArmorType > filter.ArmorType {
				return false
			}
		}
	}

	if !equipChecksOnly {
		if !HasItemEffectForTest(item.ID) {
			return false
		}

		if slices.Contains(filter.IDBlacklist, item.ID) {
			return false
		}
	}

	return true
}

func (filter *ItemFilter) FindAllItems() []Item {
	var filteredItems []Item

	for _, item := range ItemsByID {
		if filter.Matches(item, false) {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems
}

func (filter *ItemFilter) FindAllSets() []*ItemSet {
	var filteredSets []*ItemSet

	for _, set := range sets {
		if setItems := set.Items(); len(setItems) > 0 {
			if filter.Matches(setItems[0], true) {
				filteredSets = append(filteredSets, set)
			}
		}
	}

	return filteredSets
}

func (filter *ItemFilter) FindAllMetaGems() []Gem {
	var filteredGems []Gem

	for _, gem := range GemsByID {
		if gem.Color == proto.GemColor_GemColorMeta {
			if !strings.Contains(gem.Name, "Skyfire") &&
				!strings.Contains(gem.Name, "Earthstorm") &&
				!strings.Contains(gem.Name, "Starfire") &&
				!strings.Contains(gem.Name, "Unstable") {
				filteredGems = append(filteredGems, gem)
			}
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
	IsHealer   bool

	// Some fields are populated automatically.
	ItemFilter ItemFilter

	initialized bool

	items []Item
	sets  []*ItemSet

	metagems []Gem

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

	baseEquipment := ProtoToEquipment(generator.Player.Equipment)
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
	equipment := ProtoToEquipment(playerCopy.Equipment)
	if testIdx < len(generator.items) {
		testItem := generator.items[testIdx]
		equipment.EquipItem(generator.items[testIdx])
		label = fmt.Sprintf("%s-%d", strings.ReplaceAll(testItem.Name, " ", ""), testItem.ID)
	} else if testIdx < len(generator.items)+len(generator.sets) {
		testSet := generator.sets[testIdx-len(generator.items)]
		for _, setItem := range testSet.Items() {
			equipment.EquipItem(setItem)
		}
		label = strings.ReplaceAll(testSet.Name, " ", "")
	} else {
		testMetaGem := generator.metagems[testIdx-len(generator.items)-len(generator.sets)]
		headItem := &equipment[proto.ItemSlot_ItemSlotHead]
		for len(headItem.Gems) <= generator.metaSocketIdx {
			headItem.Gems = append(headItem.Gems, Gem{})
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
	if generator.IsHealer {
		rsr.Raid.TargetDummies = 1
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
	Glyphs      *proto.Glyphs
	Talents     string
	Rotation    RotationCombo

	Consumes *proto.Consumes

	IsHealer        bool
	IsTank          bool
	InFrontOfTarget bool

	OtherRaces       []proto.Race
	OtherGearSets    []GearSetCombo
	OtherSpecOptions []SpecOptionsCombo
	OtherRotations   []RotationCombo

	ItemFilter ItemFilter

	StatsToWeigh    []proto.Stat
	EPReferenceStat proto.Stat

	Cooldowns *proto.Cooldowns
}

func FullCharacterTestSuiteGenerator(config CharacterSuiteConfig) TestGenerator {
	allRaces := append(config.OtherRaces, config.Race)
	allGearSets := append(config.OtherGearSets, config.GearSet)
	allTalentSets := []TalentsCombo{{
		Label:   "Talents",
		Talents: config.Talents,
		Glyphs:  config.Glyphs,
	}}
	allSpecOptions := append(config.OtherSpecOptions, config.SpecOptions)
	allRotations := append(config.OtherRotations, config.Rotation)

	defaultPlayer := WithSpec(
		&proto.Player{
			Class:         config.Class,
			Race:          config.Race,
			Equipment:     config.GearSet.GearSet,
			Consumes:      config.Consumes,
			Buffs:         FullIndividualBuffs,
			TalentsString: config.Talents,
			Glyphs:        config.Glyphs,
			Profession1:   proto.Profession_Engineering,
			Rotation:      config.Rotation.Rotation,
			Cooldowns:     config.Cooldowns,

			InFrontOfTarget:    config.InFrontOfTarget,
			DistanceFromTarget: 30,
			ReactionTimeMs:     150,
			ChannelClipDelayMs: 50,
		},
		config.SpecOptions.SpecOptions)

	defaultRaid := SinglePlayerRaidProto(defaultPlayer, FullPartyBuffs, FullRaidBuffs, FullDebuffs)
	if config.IsTank {
		defaultRaid.Tanks = append(defaultRaid.Tanks, &proto.UnitReference{Type: proto.UnitReference_Player, Index: 0})
	}
	if config.IsHealer {
		defaultRaid.TargetDummies = 1
	}

	generator := &CombinedTestGenerator{
		subgenerators: []SubGenerator{
			{
				name: "CharacterStats",
				generator: &SingleCharacterStatsTestGenerator{
					Name: "Default",
					Request: &proto.ComputeStatsRequest{
						Raid: defaultRaid,
					},
				},
			},
			{
				name: "Settings",
				generator: &SettingsCombos{
					Class:       config.Class,
					Races:       allRaces,
					GearSets:    allGearSets,
					TalentSets:  allTalentSets,
					SpecOptions: allSpecOptions,
					Rotations:   allRotations,
					Buffs: []BuffsCombo{
						{
							Label: "NoBuffs",
						},
						{
							Label:    "FullBuffs",
							Raid:     FullRaidBuffs,
							Party:    FullPartyBuffs,
							Debuffs:  FullDebuffs,
							Player:   FullIndividualBuffs,
							Consumes: config.Consumes,
						},
					},
					IsHealer:   config.IsHealer,
					Encounters: MakeDefaultEncounterCombos(),
					SimOptions: DefaultSimTestOptions,
					Cooldowns:  config.Cooldowns,
				},
			},
			{
				name: "AllItems",
				generator: &ItemsTestGenerator{
					Player:     defaultPlayer,
					RaidBuffs:  FullRaidBuffs,
					PartyBuffs: FullPartyBuffs,
					Debuffs:    FullDebuffs,
					Encounter:  MakeSingleTargetEncounter(0),
					SimOptions: DefaultSimTestOptions,
					ItemFilter: config.ItemFilter,
					IsHealer:   config.IsHealer,
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
					RaidBuffs:  FullRaidBuffs,
					PartyBuffs: FullPartyBuffs,
					Debuffs:    FullDebuffs,
					Encounter:  MakeSingleTargetEncounter(0),
					SimOptions: StatWeightsDefaultSimTestOptions,
					Tanks:      defaultRaid.Tanks,

					StatsToWeigh:    config.StatsToWeigh,
					EpReferenceStat: config.EPReferenceStat,
				},
			},
		})
	}

	// Add this separately, so it's always last, which makes it easy to find in the
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
