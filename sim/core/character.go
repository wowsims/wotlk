package core

import (
	"fmt"
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

	// Current gear.
	Equip Equipment
	//Item Swap Handler
	ItemSwap ItemSwap

	// Consumables this Character will be using.
	Consumes *proto.Consumes

	// Base stats for this Character.
	baseStats stats.Stats

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

			DistanceFromTarget: player.DistanceFromTarget,
		},

		Name:  player.Name,
		Race:  player.Race,
		Class: player.Class,
		Equip: ProtoToEquipment(player.Equipment),
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

	if player.BonusStats != nil {
		if player.BonusStats.Stats != nil {
			character.AddStats(stats.FromFloatArray(player.BonusStats.Stats))
		}
		if player.BonusStats.PseudoStats != nil {
			ps := player.BonusStats.PseudoStats
			character.PseudoStats.BonusMHDps += ps[proto.PseudoStat_PseudoStatMainHandDps]
			character.PseudoStats.BonusOHDps += ps[proto.PseudoStat_PseudoStatOffHandDps]
			character.PseudoStats.BonusRangedDps += ps[proto.PseudoStat_PseudoStatRangedDps]
		}
	}

	if weapon := character.Equip[proto.ItemSlot_ItemSlotOffHand]; weapon.ID != 0 {
		if weapon.WeaponType == proto.WeaponType_WeaponTypeShield {
			character.PseudoStats.CanBlock = true
		}
	}
	character.PseudoStats.InFrontOfTarget = player.InFrontOfTarget

	return character
}

func (character *Character) addUniversalStatDependencies() {
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

	character.AddStats(character.Equip.Stats())
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

	for _, petAgent := range character.Pets {
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
	for slot, eq := range character.Equip {
		if applyItemEffect, ok := itemEffects[eq.ID]; ok {
			applyItemEffect(agent)
		}

		for _, g := range eq.Gems {
			if applyGemEffect, ok := itemEffects[g.ID]; ok {
				applyGemEffect(agent)
			}
		}

		// TODO: should we use eq.Enchant.EffectID because some enchants use a spellID instead of itemID?
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

	character.Pets = append(character.Pets, pet)
}

func (character *Character) GetPet(name string) PetAgent {
	for _, petAgent := range character.Pets {
		if petAgent.GetPet().Name == name {
			return petAgent
		}
	}
	panic(character.Name + " has no pet with name " + name)
}

func (character *Character) MultiplyMeleeSpeed(sim *Simulation, amount float64) {
	character.Unit.MultiplyMeleeSpeed(sim, amount)

	if len(character.Pets) > 0 {
		for _, petAgent := range character.Pets {
			petAgent.OwnerAttackSpeedChanged(sim)
		}
	}
}

func (character *Character) MultiplyRangedSpeed(sim *Simulation, amount float64) {
	character.Unit.MultiplyRangedSpeed(sim, amount)

	if len(character.Pets) > 0 {
		for _, petAgent := range character.Pets {
			petAgent.OwnerAttackSpeedChanged(sim)
		}
	}
}

func (character *Character) MultiplyAttackSpeed(sim *Simulation, amount float64) {
	character.Unit.MultiplyAttackSpeed(sim, amount)

	if len(character.Pets) > 0 {
		for _, petAgent := range character.Pets {
			petAgent.OwnerAttackSpeedChanged(sim)
		}
	}
}

func (character *Character) GetBaseStats() stats.Stats {
	return character.baseStats
}

// Returns the crit multiplier for a spell.
// https://web.archive.org/web/20081014064638/http://elitistjerks.com/f31/t12595-relentless_earthstorm_diamond_-_melee_only/p4/
// https://github.com/TheGroxEmpire/TBC_DPS_Warrior_Sim/issues/30
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

func (character *Character) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
}
func (character *Character) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
	if character.Race == proto.Race_RaceDraenei {
		partyBuffs.HeroicPresence = true
	}

	if character.Equip[ItemSlotMainHand].ID == ItemIDAtieshMage {
		partyBuffs.AtieshMage += 1
	}
	if character.Equip[ItemSlotMainHand].ID == ItemIDAtieshWarlock {
		partyBuffs.AtieshWarlock += 1
	}

	if character.Equip[ItemSlotNeck].ID == ItemIDBraidedEterniumChain {
		partyBuffs.BraidedEterniumChain = true
	}
	if character.Equip[ItemSlotNeck].ID == ItemIDChainOfTheTwilightOwl {
		partyBuffs.ChainOfTheTwilightOwl = true
	}
	if character.Equip[ItemSlotNeck].ID == ItemIDEyeOfTheNight {
		partyBuffs.EyeOfTheNight = true
	}
}

func (character *Character) initialize(agent Agent) {
	character.majorCooldownManager.initialize(character)

	character.gcdAction = &PendingAction{
		Priority: ActionPriorityGCD,
		OnAction: func(sim *Simulation) {
			if sim.CurrentTime < 0 {
				return
			}

			if character.Rotation != nil {
				character.Rotation.DoNextAction(sim)
				return
			}

			character.TryUseCooldowns(sim)
			if character.GCD.IsReady(sim) {
				agent.OnGCDReady(sim)

				if !character.doNothing && character.GCD.IsReady(sim) && (!character.IsWaiting() && !character.IsWaitingForMana()) {
					msg := fmt.Sprintf("Character `%s` did not perform any actions. Either this is a bug or agent should use 'WaitUntil' or 'WaitForMana' to explicitly wait.\n\tIf character has no action to perform use 'DoNothing'.", character.Label)
					panic(msg)
				}
				character.doNothing = false
			}
		},
	}
}

func (character *Character) Finalize(playerStats *proto.PlayerStats) {
	if character.Env.IsFinalized() {
		return
	}

	character.PseudoStats.ParryHaste = character.PseudoStats.CanParry

	character.Unit.finalize()

	character.majorCooldownManager.finalize()
	character.ItemSwap.finalize()

	if playerStats != nil {
		character.applyBuildPhaseAuras(CharacterBuildPhaseAll)
		playerStats.FinalStats = &proto.UnitStats{
			Stats:       character.GetStats().ToFloatArray(),
			PseudoStats: character.GetPseudoStatsProto(),
		}
		character.clearBuildPhaseAuras(CharacterBuildPhaseAll)
		playerStats.Sets = character.GetActiveSetBonusNames()
		playerStats.Cooldowns = character.GetMajorCooldownIDs()
	}
}

func (character *Character) init(sim *Simulation, agent Agent) {
	character.Unit.init(sim)
}

func (character *Character) reset(sim *Simulation, agent Agent) {
	character.Unit.reset(sim, agent)
	character.majorCooldownManager.reset(sim)
	character.ItemSwap.reset(sim)
	character.CurrentTarget = character.defaultTarget

	agent.Reset(sim)

	for _, petAgent := range character.Pets {
		petAgent.GetPet().reset(sim, petAgent)
	}
}

// Advance moves time forward counting down auras, CDs, mana regen, etc
func (character *Character) advance(sim *Simulation, elapsedTime time.Duration) {
	character.Unit.advance(sim, elapsedTime)

	if len(character.Pets) > 0 {
		for _, petAgent := range character.Pets {
			if !petAgent.GetPet().enabled {
				continue
			}
			petAgent.GetPet().advance(sim, elapsedTime)
		}
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
	return character.Equip[ItemSlotTrinket1].ID == itemID ||
		character.Equip[ItemSlotTrinket2].ID == itemID
}

func (character *Character) HasRingEquipped(itemID int32) bool {
	return character.Equip[ItemSlotFinger1].ID == itemID ||
		character.Equip[ItemSlotFinger2].ID == itemID
}

func (character *Character) HasMetaGemEquipped(gemID int32) bool {
	for _, gem := range character.Equip[ItemSlotHead].Gems {
		if gem.ID == gemID {
			return true
		}
	}
	return false
}

// Returns the MH weapon if one is equipped, and null otherwise.
func (character *Character) GetMHWeapon() *Item {
	weapon := &character.Equip[proto.ItemSlot_ItemSlotMainHand]
	if weapon.ID == 0 {
		return nil
	} else {
		return weapon
	}
}
func (character *Character) HasMHWeapon() bool {
	return character.GetMHWeapon() != nil
}

// Returns the OH weapon if one is equipped, and null otherwise. Note that
// shields / Held-in-off-hand items are NOT counted as weapons in this function.
func (character *Character) GetOHWeapon() *Item {
	weapon := &character.Equip[proto.ItemSlot_ItemSlotOffHand]
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
	weapon := &character.Equip[proto.ItemSlot_ItemSlotRanged]
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

// Returns the hands that the item is equipped in, as (MH, OH).
func (character *Character) GetWeaponHands(itemID int32) (bool, bool) {
	mh := false
	oh := false
	if weapon := character.GetMHWeapon(); weapon != nil && weapon.ID == itemID {
		mh = true
	}
	if weapon := character.GetOHWeapon(); weapon != nil && weapon.ID == itemID {
		oh = true
	}
	return mh, oh
}

func (character *Character) doneIteration(sim *Simulation) {
	// Need to do pets first so we can add their results to the owners.
	if len(character.Pets) > 0 {
		for _, petAgent := range character.Pets {
			pet := petAgent.GetPet()
			pet.doneIteration(sim)
			character.Metrics.AddFinalPetMetrics(&pet.Metrics)
		}
	}

	character.Unit.doneIteration(sim)
}

func (character *Character) GetPseudoStatsProto() []float64 {
	vals := make([]float64, stats.PseudoStatsLen)
	vals[proto.PseudoStat_PseudoStatMainHandDps] = character.WeaponFromMainHand(0).DPS()
	vals[proto.PseudoStat_PseudoStatOffHandDps] = character.WeaponFromOffHand(0).DPS()
	vals[proto.PseudoStat_PseudoStatRangedDps] = character.WeaponFromRanged(0).DPS()
	vals[proto.PseudoStat_PseudoStatBlockValueMultiplier] = character.PseudoStats.BlockValueMultiplier
	// Base values are modified by Enemy attackTables, but we display for LVL 80 enemy as paperdoll default
	vals[proto.PseudoStat_PseudoStatDodge] = character.PseudoStats.BaseDodge + character.GetDiminishedDodgeChance()
	vals[proto.PseudoStat_PseudoStatParry] = character.PseudoStats.BaseParry + character.GetDiminishedParryChance()
	//vals[proto.PseudoStat_PseudoStatMiss] = 0.05 + character.GetDiminishedMissChance() + character.PseudoStats.ReducedPhysicalHitTakenChance
	return vals
}

func (character *Character) GetMetricsProto() *proto.UnitMetrics {
	metrics := character.Metrics.ToProto()
	metrics.Name = character.Name
	metrics.UnitIndex = character.UnitIndex
	metrics.Auras = character.auraTracker.GetMetricsProto()

	metrics.Pets = []*proto.UnitMetrics{}
	for _, petAgent := range character.Pets {
		metrics.Pets = append(metrics.Pets, petAgent.GetPet().GetMetricsProto())
	}

	return metrics
}

type BaseStatsKey struct {
	Race  proto.Race
	Class proto.Class
}

var BaseStats = map[BaseStatsKey]stats.Stats{}

// To calculate base stats, get a naked level 70 of the race/class you want, ideally without any talents to mess up base stats.
//  Basic stats are as-shown (str/agi/stm/int/spirit)

// Base Spell Crit is calculated by
//   1. Take as-shown value (troll shaman have 3.5%)
//   2. Calculate the bonus from int (for troll shaman that would be 104/78.1=1.331% crit)
//   3. Subtract as-shown from int bouns (3.5-1.331=2.169)
//   4. 2.169*22.08 (rating per crit percent) = 47.89 crit rating.

// Base mana can be looked up here: https://wowwiki-archive.fandom.com/wiki/Base_mana

// I assume a similar processes can be applied for other stats.

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
// strings truncate 0's at the end of each tree so we can't infer the start index
// of the tree from the string.
func FillTalentsProto(data protoreflect.Message, talentsStr string, treeSizes [3]int) {
	treeStrs := strings.Split(talentsStr, "-")
	fieldDescriptors := data.Descriptor().Fields()

	var offset int
	for treeIdx, treeStr := range treeStrs {
		for talentIdx, talentValStr := range treeStr {
			talentVal, _ := strconv.Atoi(string(talentValStr))
			fd := fieldDescriptors.ByNumber(protowire.Number(offset + talentIdx + 1))
			if fd.Kind() == protoreflect.BoolKind {
				data.Set(fd, protoreflect.ValueOfBool(talentVal == 1))
			} else { // Int32Kind
				data.Set(fd, protoreflect.ValueOfInt32(int32(talentVal)))
			}
		}
		offset += treeSizes[treeIdx]
	}
}
