package core

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type CharacterBuildPhase uint8

func (cbp CharacterBuildPhase) Matches(other CharacterBuildPhase) bool {
	return (cbp & other) != 0
}

const (
	CharacterBuildPhaseNone CharacterBuildPhase = 0
	CharacterBuildPhaseBase CharacterBuildPhase = 1 << iota
	CharacterBuildPhaseGear
	CharacterBuildPhaseTalents
	CharacterBuildPhaseBuffs
	CharacterBuildPhaseConsumes
)

const CharacterBuildPhaseAll = CharacterBuildPhaseBase | CharacterBuildPhaseGear | CharacterBuildPhaseTalents | CharacterBuildPhaseBuffs | CharacterBuildPhaseConsumes

// Character is a data structure to hold all the shared values that all
// class logic shares.
// All players have stats, equipment, auras, etc
type Character struct {
	Unit

	Name  string // Different from Label, needed for returned results.
	Race  proto.Race
	Class proto.Class
	Spec  proto.Spec

	// Current gear.
	Equipment
	//Item Swap Handler
	ItemSwap ItemSwap

	// Consumables this Character will be using.
	Consumes *proto.Consumes

	// Base stats for this Character.
	baseStats stats.Stats

	// Handles scaling that only affects stats from items
	itemStatMultipliers stats.Stats
	// Used to track if we need to separately apply multipliers, because
	// equipment was already applied
	equipStatsApplied bool

	// Bonus stats for this Character, specified in the UI and/or EP
	// calculator
	bonusStats     stats.Stats
	bonusMHDps     float64
	bonusOHDps     float64
	bonusRangedDps float64

	professions [2]proto.Profession

	glyphs            [6]int32
	PrimaryTalentTree uint8

	// Provides major cooldown management behavior.
	majorCooldownManager

	// Up reference to this Character's Party.
	Party *Party

	// This character's index within its party [0-4].
	PartyIndex int

	defensiveTrinketCD *Timer
	offensiveTrinketCD *Timer
	conjuredCD         *Timer

	Pets []*Pet // cached in AddPet, for advance()
}

func NewCharacter(party *Party, partyIndex int, player *proto.Player) Character {
	if player.Database != nil {
		addToDatabase(player.Database)
	}

	character := Character{
		Unit: Unit{
			Type:        PlayerUnit,
			Index:       int32(party.Index*5 + partyIndex),
			Level:       CharacterLevel,
			auraTracker: newAuraTracker(),
			PseudoStats: stats.NewPseudoStats(),
			Metrics:     NewUnitMetrics(),

			StatDependencyManager: stats.NewStatDependencyManager(),

			ReactionTime:         max(0, time.Duration(player.ReactionTimeMs)*time.Millisecond),
			ChannelClipDelay:     max(0, time.Duration(player.ChannelClipDelayMs)*time.Millisecond),
			DistanceFromTarget:   player.DistanceFromTarget,
			NibelungAverageCasts: player.NibelungAverageCasts,
		},

		Name:  player.Name,
		Race:  player.Race,
		Class: player.Class,
		Spec:  PlayerProtoToSpec(player),

		Equipment: ProtoToEquipment(player.Equipment),

		professions: [2]proto.Profession{
			player.Profession1,
			player.Profession2,
		},

		Party:      party,
		PartyIndex: partyIndex,

		majorCooldownManager: newMajorCooldownManager(player.Cooldowns),
	}

	character.GCD = character.NewTimer()

	character.Label = fmt.Sprintf("%s (#%d)", character.Name, character.Index+1)

	if player.Glyphs != nil {
		character.glyphs = [6]int32{
			player.Glyphs.Major1,
			player.Glyphs.Major2,
			player.Glyphs.Major3,
			player.Glyphs.Minor1,
			player.Glyphs.Minor2,
			player.Glyphs.Minor3,
		}
	}
	character.PrimaryTalentTree = GetPrimaryTalentTreeIndex(player.TalentsString)

	character.Consumes = &proto.Consumes{}
	if player.Consumes != nil {
		character.Consumes = player.Consumes
	}

	character.baseStats = BaseStats[BaseStatsKey{Race: character.Race, Class: character.Class}]

	character.AddStats(character.baseStats)
	character.addUniversalStatDependencies()
	for i := range character.itemStatMultipliers {
		character.itemStatMultipliers[i] = 1
	}

	if player.BonusStats != nil {
		if player.BonusStats.Stats != nil {
			character.bonusStats = stats.FromFloatArray(player.BonusStats.Stats)
		}
		if player.BonusStats.PseudoStats != nil {
			ps := player.BonusStats.PseudoStats
			character.bonusMHDps = ps[proto.PseudoStat_PseudoStatMainHandDps]
			character.bonusOHDps = ps[proto.PseudoStat_PseudoStatOffHandDps]
			character.bonusRangedDps = ps[proto.PseudoStat_PseudoStatRangedDps]
			character.PseudoStats.BonusMHDps += character.bonusMHDps
			character.PseudoStats.BonusOHDps += character.bonusOHDps
			character.PseudoStats.BonusRangedDps += character.bonusRangedDps
		}
	}

	if weapon := character.OffHand(); weapon.ID != 0 {
		if weapon.WeaponType == proto.WeaponType_WeaponTypeShield {
			character.PseudoStats.CanBlock = true
		}
	}
	character.PseudoStats.InFrontOfTarget = player.InFrontOfTarget

	if player.EnableItemSwap && player.ItemSwap != nil {
		character.enableItemSwap(player.ItemSwap, character.DefaultMeleeCritMultiplier(), character.DefaultMeleeCritMultiplier(), 0)
	}

	return character
}

func (character *Character) applyEquipScaling(stat stats.Stat, multiplier float64) float64 {
	var oldValue = character.EquipStats()[stat]
	character.itemStatMultipliers[stat] *= multiplier
	var newValue = character.EquipStats()[stat]
	return newValue - oldValue
}

func (character *Character) ApplyEquipScaling(stat stats.Stat, multiplier float64) {
	var statDiff stats.Stats
	statDiff[stat] = character.applyEquipScaling(stat, multiplier)
	// Equipment stats already applied, so need to manually at the bonus to
	// the character now to ensure correct values
	if character.equipStatsApplied {
		character.AddStats(statDiff)
	}
}

func (character *Character) ApplyDynamicEquipScaling(sim *Simulation, stat stats.Stat, multiplier float64) {
	statDiff := character.applyEquipScaling(stat, multiplier)
	character.AddStatDynamic(sim, stat, statDiff)
}

func (character *Character) RemoveEquipScaling(stat stats.Stat, multiplier float64) {
	var statDiff stats.Stats
	statDiff[stat] = character.applyEquipScaling(stat, 1/multiplier)
	// Equipment stats already applied, so need to manually at the bonus to
	// the character now to ensure correct values
	if character.equipStatsApplied {
		character.AddStats(statDiff)
	}
}

func (character *Character) RemoveDynamicEquipScaling(sim *Simulation, stat stats.Stat, multiplier float64) {
	statDiff := character.applyEquipScaling(stat, 1/multiplier)
	character.AddStatDynamic(sim, stat, statDiff)
}

func (character *Character) EquipStats() stats.Stats {
	var baseEquipStats = character.Equipment.Stats()
	var bonusEquipStats = baseEquipStats.Add(character.bonusStats)
	return bonusEquipStats.DotProduct(character.itemStatMultipliers)
}

func (character *Character) applyEquipment() {
	if character.equipStatsApplied {
		panic("Equipment stats already applied to character!")
	}
	character.AddStats(character.EquipStats())
	character.equipStatsApplied = true
}

func (character *Character) addUniversalStatDependencies() {
	character.AddStat(stats.Health, 20-10*20)
	character.AddStatDependency(stats.Stamina, stats.Health, 10)
	character.AddStatDependency(stats.Agility, stats.Armor, 2)
}

// Returns a partially-filled PlayerStats proto for use in the CharacterStats api call.
func (character *Character) applyAllEffects(agent Agent, raidBuffs *proto.RaidBuffs, partyBuffs *proto.PartyBuffs, individualBuffs *proto.IndividualBuffs) *proto.PlayerStats {
	playerStats := &proto.PlayerStats{}

	measureStats := func() *proto.UnitStats {
		return &proto.UnitStats{
			Stats:       character.SortAndApplyStatDependencies(character.stats).ToFloatArray(),
			PseudoStats: character.GetPseudoStatsProto(),
		}
	}

	applyRaceEffects(agent)
	character.applyProfessionEffects()
	character.applyBuildPhaseAuras(CharacterBuildPhaseBase)
	playerStats.BaseStats = measureStats()

	character.applyEquipment()
	character.applyItemEffects(agent)
	character.applyItemSetBonusEffects(agent)
	character.applyBuildPhaseAuras(CharacterBuildPhaseGear)
	playerStats.GearStats = measureStats()

	agent.ApplyTalents()
	character.applyBuildPhaseAuras(CharacterBuildPhaseTalents)
	playerStats.TalentsStats = measureStats()

	applyBuffEffects(agent, raidBuffs, partyBuffs, individualBuffs)
	character.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)
	playerStats.BuffsStats = measureStats()

	applyConsumeEffects(agent)
	character.applyBuildPhaseAuras(CharacterBuildPhaseConsumes)
	playerStats.ConsumesStats = measureStats()
	character.clearBuildPhaseAuras(CharacterBuildPhaseAll)

	for _, petAgent := range character.PetAgents {
		applyPetBuffEffects(petAgent, raidBuffs, partyBuffs, individualBuffs)
	}

	return playerStats
}
func (character *Character) applyBuildPhaseAuras(phase CharacterBuildPhase) {
	sim := Simulation{}
	character.Env.MeasuringStats = true
	for _, aura := range character.auras {
		if aura.BuildPhase.Matches(phase) {
			aura.Activate(&sim)
		}
	}
	character.Env.MeasuringStats = false
}
func (character *Character) clearBuildPhaseAuras(phase CharacterBuildPhase) {
	sim := Simulation{}
	character.Env.MeasuringStats = true
	for _, aura := range character.auras {
		if aura.BuildPhase.Matches(phase) {
			aura.Deactivate(&sim)
		}
	}
	character.Env.MeasuringStats = false
}

// Apply effects from all equipped core.
func (character *Character) applyItemEffects(agent Agent) {
	for slot, eq := range character.Equipment {
		if applyItemEffect, ok := itemEffects[eq.ID]; ok {
			applyItemEffect(agent)
		}

		for _, g := range eq.Gems {
			if applyGemEffect, ok := itemEffects[g.ID]; ok {
				applyGemEffect(agent)
			}
		}

		if applyEnchantEffect, ok := enchantEffects[eq.Enchant.EffectID]; ok {
			applyEnchantEffect(agent)
		}

		if applyWeaponEffect, ok := weaponEffects[eq.Enchant.EffectID]; ok {
			applyWeaponEffect(agent, proto.ItemSlot(slot))
		}
	}

	if character.ItemSwap.IsEnabled() {
		offset := int(proto.ItemSlot_ItemSlotMainHand)
		for i, item := range character.ItemSwap.unEquippedItems {
			if applyEnchantEffect, ok := enchantEffects[item.Enchant.EffectID]; ok {
				applyEnchantEffect(agent)
			}

			if applyWeaponEffect, ok := weaponEffects[item.Enchant.EffectID]; ok {
				applyWeaponEffect(agent, proto.ItemSlot(offset+i))
			}
		}
	}
}

func (character *Character) AddPet(pet PetAgent) {
	if character.Env != nil {
		panic("Pets must be added during construction!")
	}

	character.PetAgents = append(character.PetAgents, pet)
	character.Pets = append(character.Pets, pet.GetPet())
}

func (character *Character) GetBaseStats() stats.Stats {
	return character.baseStats
}

// Returns the crit multiplier for a spell.
// https://web.archive.org/web/20081014064638/http://elitistjerks.com/f31/t12595-relentless_earthstorm_diamond_-_melee_only/p4/
// https://github.com/TheGroxEmpire/TBC_DPS_Warrior_Sim/issues/30
// TODO "primaryModifiers" could be modelled as a PseudoStat, since they're unit-specific. "secondaryModifiers" apply to a specific set of spells.
func (character *Character) calculateCritMultiplier(normalCritDamage float64, primaryModifiers float64, secondaryModifiers float64) float64 {
	if character.HasMetaGemEquipped(34220) ||
		character.HasMetaGemEquipped(32409) ||
		character.HasMetaGemEquipped(41285) ||
		character.HasMetaGemEquipped(41398) {
		primaryModifiers *= 1.03
	}
	return 1.0 + (normalCritDamage*primaryModifiers-1.0)*(1.0+secondaryModifiers)
}
func (character *Character) calculateHealingCritMultiplier(normalCritDamage float64, primaryModifiers float64, secondaryModifiers float64) float64 {
	if character.HasMetaGemEquipped(41376) {
		primaryModifiers *= 1.03
	}
	return 1.0 + (normalCritDamage*primaryModifiers-1.0)*(1.0+secondaryModifiers)
}
func (character *Character) SpellCritMultiplier(primaryModifiers float64, secondaryModifiers float64) float64 {
	return character.calculateCritMultiplier(1.5, primaryModifiers, secondaryModifiers)
}
func (character *Character) MeleeCritMultiplier(primaryModifiers float64, secondaryModifiers float64) float64 {
	return character.calculateCritMultiplier(2.0, primaryModifiers, secondaryModifiers)
}
func (character *Character) HealingCritMultiplier(primaryModifiers float64, secondaryModifiers float64) float64 {
	return character.calculateHealingCritMultiplier(1.5, primaryModifiers, secondaryModifiers)
}
func (character *Character) DefaultSpellCritMultiplier() float64 {
	return character.SpellCritMultiplier(1, 0)
}
func (character *Character) DefaultMeleeCritMultiplier() float64 {
	return character.MeleeCritMultiplier(1, 0)
}
func (character *Character) DefaultHealingCritMultiplier() float64 {
	return character.HealingCritMultiplier(1, 0)
}

func (character *Character) AddRaidBuffs(_ *proto.RaidBuffs) {
}
func (character *Character) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
	if character.Race == proto.Race_RaceDraenei {
		partyBuffs.HeroicPresence = true
	}

	switch character.MainHand().ID {
	case ItemIDAtieshMage:
		partyBuffs.AtieshMage += 1
	case ItemIDAtieshWarlock:
		partyBuffs.AtieshWarlock += 1
	}

	switch character.Neck().ID {
	case ItemIDBraidedEterniumChain:
		partyBuffs.BraidedEterniumChain = true
	case ItemIDChainOfTheTwilightOwl:
		partyBuffs.ChainOfTheTwilightOwl = true
	case ItemIDEyeOfTheNight:
		partyBuffs.EyeOfTheNight = true
	}
}

func (character *Character) initialize(agent Agent) {
	character.majorCooldownManager.initialize(character)
	character.ItemSwap.initialize(character)

	character.gcdAction = &PendingAction{
		Priority: ActionPriorityGCD,
		OnAction: func(sim *Simulation) {
			if hc := &character.Hardcast; hc.Expires != startingCDTime && hc.Expires <= sim.CurrentTime {
				hc.Expires = startingCDTime
				if hc.OnComplete != nil {
					hc.OnComplete(sim, hc.Target)
				}
			}

			if sim.CurrentTime < 0 {
				return
			}

			if sim.Options.Interactive {
				if character.GCD.IsReady(sim) {
					sim.NeedsInput = true
				}
				return
			}

			if character.Rotation != nil {
				character.Rotation.DoNextAction(sim)
				return
			}
		},
	}
}

func (character *Character) Finalize() {
	if character.Env.IsFinalized() {
		return
	}

	character.PseudoStats.ParryHaste = character.PseudoStats.CanParry

	character.Unit.finalize()

	// For now, restrict this optimization to rogues only. Ferals will require
	// some extra logic to handle their ExcessEnergy() calc.
	if character.Class == proto.Class_ClassRogue {
		character.Env.RegisterPostFinalizeEffect(func() {
			character.energyBar.setupEnergyThresholds()
		})
	}

	character.majorCooldownManager.finalize()
}

func (character *Character) FillPlayerStats(playerStats *proto.PlayerStats) {
	if playerStats == nil {
		return
	}

	character.applyBuildPhaseAuras(CharacterBuildPhaseAll)
	playerStats.FinalStats = &proto.UnitStats{
		Stats:       character.GetStats().ToFloatArray(),
		PseudoStats: character.GetPseudoStatsProto(),
	}
	character.clearBuildPhaseAuras(CharacterBuildPhaseAll)
	playerStats.Sets = character.GetActiveSetBonusNames()

	playerStats.Metadata = character.GetMetadata()
	for _, pet := range character.Pets {
		playerStats.Pets = append(playerStats.Pets, &proto.PetStats{
			Metadata: pet.GetMetadata(),
		})
	}

	if character.Rotation != nil {
		playerStats.RotationStats = character.Rotation.getStats()
	}
}

func (character *Character) reset(sim *Simulation, agent Agent) {
	character.Unit.reset(sim, agent)
	character.majorCooldownManager.reset(sim)
	character.ItemSwap.reset(sim)
	character.CurrentTarget = character.defaultTarget

	agent.Reset(sim)

	for _, petAgent := range character.PetAgents {
		petAgent.GetPet().reset(sim, petAgent)
	}
}

func (character *Character) HasProfession(prof proto.Profession) bool {
	return prof == character.professions[0] || prof == character.professions[1]
}

func (character *Character) HasGlyph(glyphID int32) bool {
	for _, g := range character.glyphs {
		if g == glyphID {
			return true
		}
	}
	return false
}

func (character *Character) HasTrinketEquipped(itemID int32) bool {
	return character.Trinket1().ID == itemID ||
		character.Trinket2().ID == itemID
}

func (character *Character) HasRingEquipped(itemID int32) bool {
	return character.Finger1().ID == itemID || character.Finger2().ID == itemID
}

func (character *Character) HasMetaGemEquipped(gemID int32) bool {
	for _, gem := range character.Head().Gems {
		if gem.ID == gemID {
			return true
		}
	}
	return false
}

// Returns the MH weapon if one is equipped, and null otherwise.
func (character *Character) GetMHWeapon() *Item {
	weapon := character.MainHand()
	if weapon.ID == 0 {
		return nil
	}
	return weapon
}
func (character *Character) HasMHWeapon() bool {
	return character.GetMHWeapon() != nil
}

// Returns the OH weapon if one is equipped, and null otherwise. Note that
// shields / Held-in-off-hand items are NOT counted as weapons in this function.
func (character *Character) GetOHWeapon() *Item {
	weapon := character.OffHand()
	if weapon.ID == 0 ||
		weapon.WeaponType == proto.WeaponType_WeaponTypeShield ||
		weapon.WeaponType == proto.WeaponType_WeaponTypeOffHand {
		return nil
	} else {
		return weapon
	}
}
func (character *Character) HasOHWeapon() bool {
	return character.GetOHWeapon() != nil
}

// Returns the ranged weapon if one is equipped, and null otherwise.
func (character *Character) GetRangedWeapon() *Item {
	weapon := character.Ranged()
	if weapon.ID == 0 ||
		weapon.RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeIdol ||
		weapon.RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeLibram ||
		weapon.RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeTotem {
		return nil
	} else {
		return weapon
	}
}
func (character *Character) HasRangedWeapon() bool {
	return character.GetRangedWeapon() != nil
}

func (character *Character) GetProcMaskForEnchant(effectID int32) ProcMask {
	return character.getProcMaskFor(func(weapon *Item) bool {
		return weapon.Enchant.EffectID == effectID
	})
}

func (character *Character) GetProcMaskForItem(itemID int32) ProcMask {
	return character.getProcMaskFor(func(weapon *Item) bool {
		return weapon.ID == itemID
	})
}

func (character *Character) GetProcMaskForTypes(weaponTypes ...proto.WeaponType) ProcMask {
	return character.getProcMaskFor(func(weapon *Item) bool {
		return weapon == nil || slices.Contains(weaponTypes, weapon.WeaponType)
	})
}

func (character *Character) getProcMaskFor(pred func(weapon *Item) bool) ProcMask {
	mask := ProcMaskUnknown
	if pred(character.MainHand()) {
		mask |= ProcMaskMeleeMH
	}
	if pred(character.OffHand()) {
		mask |= ProcMaskMeleeOH
	}
	return mask
}

func (character *Character) doneIteration(sim *Simulation) {
	// Need to do pets first, so we can add their results to the owners.
	for _, pet := range character.Pets {
		pet.doneIteration(sim)
		character.Metrics.AddFinalPetMetrics(&pet.Metrics)
	}

	character.Unit.doneIteration(sim)
}

func (character *Character) GetPseudoStatsProto() []float64 {
	return []float64{
		proto.PseudoStat_PseudoStatMainHandDps:          character.AutoAttacks.MH().DPS(),
		proto.PseudoStat_PseudoStatOffHandDps:           character.AutoAttacks.OH().DPS(),
		proto.PseudoStat_PseudoStatRangedDps:            character.AutoAttacks.Ranged().DPS(),
		proto.PseudoStat_PseudoStatBlockValueMultiplier: character.PseudoStats.BlockValueMultiplier,
		// Base values are modified by Enemy attackTables, but we display for LVL 80 enemy as paperdoll default
		proto.PseudoStat_PseudoStatDodge: character.PseudoStats.BaseDodge + character.GetDiminishedDodgeChance(),
		proto.PseudoStat_PseudoStatParry: character.PseudoStats.BaseParry + character.GetDiminishedParryChance(),
	}
}

func (character *Character) GetMetricsProto() *proto.UnitMetrics {
	metrics := character.Metrics.ToProto()
	metrics.Name = character.Name
	metrics.UnitIndex = character.UnitIndex
	metrics.Auras = character.auraTracker.GetMetricsProto()

	metrics.Pets = make([]*proto.UnitMetrics, len(character.Pets))
	for i, pet := range character.Pets {
		metrics.Pets[i] = pet.GetMetricsProto()
	}

	return metrics
}

func (character *Character) GetDefensiveTrinketCD() *Timer {
	return character.GetOrInitTimer(&character.defensiveTrinketCD)
}
func (character *Character) GetOffensiveTrinketCD() *Timer {
	return character.GetOrInitTimer(&character.offensiveTrinketCD)
}
func (character *Character) GetConjuredCD() *Timer {
	return character.GetOrInitTimer(&character.conjuredCD)
}

// Returns the talent tree (0, 1, or 2) of the tree with the most points.
//
// talentStr is expected to be a wowhead-formatted talent string, e.g.
// "12123131-123123123-123123213"
func GetPrimaryTalentTreeIndex(talentStr string) uint8 {
	trees := strings.Split(talentStr, "-")
	bestTree := 0
	bestTreePoints := 0

	for treeIdx, treeStr := range trees {
		points := 0
		for talentIdx := 0; talentIdx < len(treeStr); talentIdx++ {
			v, _ := strconv.Atoi(string(treeStr[talentIdx]))
			points += v
		}

		if points > bestTreePoints {
			bestTreePoints = points
			bestTree = treeIdx
		}
	}

	return uint8(bestTree)
}

// Uses proto reflection to set fields in a talents proto (e.g. MageTalents,
// WarriorTalents) based on a talentsStr. treeSizes should contain the number
// of talents in each tree, usually around 30. This is needed because talent
// strings truncate 0's at the end of each tree, so we can't infer the start index
// of the tree from the string.
func FillTalentsProto(data protoreflect.Message, talentsStr string, treeSizes [3]int) {
	treeStrs := strings.Split(talentsStr, "-")
	fieldDescriptors := data.Descriptor().Fields()

	var offset int
	for treeIdx, treeStr := range treeStrs {
		for talentIdx, talentValStr := range treeStr {
			talentVal, _ := strconv.Atoi(string(talentValStr))
			talentOffset := offset + talentIdx + 1
			fd := fieldDescriptors.ByNumber(protowire.Number(talentOffset))
			if fd == nil {
				panic(fmt.Sprintf("Couldn't find proto field for talent #%d, full string: %s", talentOffset, talentsStr))
			}
			if fd.Kind() == protoreflect.BoolKind {
				data.Set(fd, protoreflect.ValueOfBool(talentVal == 1))
			} else { // Int32Kind
				data.Set(fd, protoreflect.ValueOfInt32(int32(talentVal)))
			}
		}
		offset += treeSizes[treeIdx]
	}
}
