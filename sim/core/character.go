package core

import (
	"fmt"
	"time"

	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

// Character is a data structure to hold all the shared values that all
// class logic shares.
// All players have stats, equipment, auras, etc
type Character struct {
	Unit

	Name         string // Different from Label, needed for returned results.
	Race         proto.Race
	ShattFaction proto.ShattrathFaction
	Class        proto.Class

	// Current gear.
	Equip items.Equipment

	// Pets owned by this Character.
	Pets []PetAgent

	// Consumables this Character will be using.
	Consumes proto.Consumes

	// Base stats for this Character.
	baseStats stats.Stats

	// Provides stat dependency management behavior.
	stats.StatDependencyManager

	// Provides major cooldown management behavior.
	majorCooldownManager

	// Up reference to this Character's Party.
	Party *Party

	// This character's index within its party [0-4].
	PartyIndex int

	// Total amount of remaining additional mana expected for the current sim iteration,
	// beyond this Character's mana pool. This should include mana potions / runes /
	// innervates / etc.
	ExpectedBonusMana float64

	// Hack for ensuring we don't apply windfury totem aura if there's already
	// a MH imbue.
	// TODO: Figure out a cleaner way to do this.
	HasMHWeaponImbue bool

	defensiveTrinketCD *Timer
	offensiveTrinketCD *Timer
	conjuredCD         *Timer
}

func NewCharacter(party *Party, partyIndex int, player proto.Player) Character {
	character := Character{
		Unit: Unit{
			Type:        PlayerUnit,
			Index:       int32(party.Index*5 + partyIndex),
			Level:       CharacterLevel,
			auraTracker: newAuraTracker(),
			PseudoStats: stats.NewPseudoStats(),
			Metrics:     NewUnitMetrics(),
		},

		Name:         player.Name,
		Race:         player.Race,
		ShattFaction: player.ShattFaction,
		Class:        player.Class,
		Equip:        items.ProtoToEquipment(*player.Equipment),

		Party:      party,
		PartyIndex: partyIndex,

		majorCooldownManager: newMajorCooldownManager(player.Cooldowns),
	}

	character.GCD = character.NewTimer()

	character.Label = fmt.Sprintf("%s (#%d)", character.Name, character.Index+1)

	if player.Consumes != nil {
		character.Consumes = *player.Consumes
	}

	character.baseStats = BaseStats[BaseStatsKey{Race: character.Race, Class: character.Class}]

	bonusStats := stats.Stats{}
	if player.BonusStats != nil {
		copy(bonusStats[:], player.BonusStats[:])
	}

	character.AddStats(character.baseStats)
	character.AddStats(bonusStats)
	character.addUniversalStatDependencies()

	if weapon := character.Equip[proto.ItemSlot_ItemSlotOffHand]; weapon.ID != 0 {
		if weapon.WeaponType == proto.WeaponType_WeaponTypeShield {
			character.PseudoStats.CanBlock = true
		}
	}
	character.PseudoStats.InFrontOfTarget = player.InFrontOfTarget
	character.addEffectPets()

	return character
}

func (character *Character) addUniversalStatDependencies() {
	character.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Stamina,
		ModifiedStat: stats.Health,
		Modifier: func(stamina float64, health float64) float64 {
			return health + stamina*10
		},
	})
	character.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Agility,
		ModifiedStat: stats.Armor,
		Modifier: func(agility float64, armor float64) float64 {
			return armor + agility*2
		},
	})
}

// Empty implementation so its optional for Agents.
func (character *Character) ApplyGearBonuses() {}

// Returns a partially-filled PlayerStats proto for use in the CharacterStats api call.
func (character *Character) applyAllEffects(agent Agent, raidBuffs proto.RaidBuffs, partyBuffs proto.PartyBuffs, individualBuffs proto.IndividualBuffs) *proto.PlayerStats {
	playerStats := &proto.PlayerStats{}

	applyRaceEffects(agent)
	playerStats.BaseStats = character.SortAndApplyStatDependencies(character.stats).ToFloatArray()

	character.AddStats(character.Equip.Stats())
	character.applyItemEffects(agent)
	character.applyItemSetBonusEffects(agent)
	agent.ApplyGearBonuses()
	playerStats.GearStats = character.SortAndApplyStatDependencies(character.stats).ToFloatArray()

	agent.ApplyTalents()
	playerStats.TalentsStats = character.SortAndApplyStatDependencies(character.stats).ToFloatArray()

	applyBuffEffects(agent, raidBuffs, partyBuffs, individualBuffs)
	playerStats.BuffsStats = character.SortAndApplyStatDependencies(character.stats).ToFloatArray()

	applyConsumeEffects(agent, raidBuffs, partyBuffs)
	playerStats.ConsumesStats = character.SortAndApplyStatDependencies(character.stats).ToFloatArray()

	for _, petAgent := range character.Pets {
		applyPetBuffEffects(petAgent, raidBuffs, partyBuffs, individualBuffs)
	}

	return playerStats
}

// Apply effects from all equipped items.
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
		if applyEnchantEffect, ok := itemEffects[eq.Enchant.ID]; ok {
			applyEnchantEffect(agent)
		}

		if applyWeaponEffect, ok := weaponEffects[eq.Enchant.ID]; ok {
			applyWeaponEffect(agent, proto.ItemSlot(slot))
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

func (character *Character) AddStatsDynamic(sim *Simulation, stat stats.Stats) {
	character.Unit.AddStatsDynamic(sim, stat)

	if len(character.Pets) > 0 {
		for _, petAgent := range character.Pets {
			petAgent.GetPet().addOwnerStats(sim, stat)
		}
	}
}
func (character *Character) AddStatDynamic(sim *Simulation, stat stats.Stat, amount float64) {
	character.Unit.AddStatDynamic(sim, stat, amount)

	if len(character.Pets) > 0 {
		for _, petAgent := range character.Pets {
			petAgent.GetPet().addOwnerStat(sim, stat, amount)
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
	if character.HasMetaGemEquipped(34220) || character.HasMetaGemEquipped(32409) { // CSD and RED
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
func (character *Character) DefaultSpellCritMultiplier() float64 {
	return character.SpellCritMultiplier(1, 0)
}
func (character *Character) DefaultMeleeCritMultiplier() float64 {
	return character.MeleeCritMultiplier(1, 0)
}

func (character *Character) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
}
func (character *Character) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
	if character.Race == proto.Race_RaceDraenei {
		class := character.Class
		if class == proto.Class_ClassHunter ||
			class == proto.Class_ClassPaladin ||
			class == proto.Class_ClassWarrior {
			partyBuffs.DraeneiRacialMelee = true
		} else if class == proto.Class_ClassMage ||
			class == proto.Class_ClassPriest ||
			class == proto.Class_ClassShaman {
			partyBuffs.DraeneiRacialCaster = true
		}
	}

	if character.Consumes.Drums > 0 {
		partyBuffs.Drums = character.Consumes.Drums
	}

	if character.Equip[items.ItemSlotMainHand].ID == ItemIDAtieshMage {
		partyBuffs.AtieshMage += 1
	}
	if character.Equip[items.ItemSlotMainHand].ID == ItemIDAtieshWarlock {
		partyBuffs.AtieshWarlock += 1
	}

	if character.Equip[items.ItemSlotNeck].ID == ItemIDBraidedEterniumChain {
		partyBuffs.BraidedEterniumChain = true
	}
	if character.Equip[items.ItemSlotNeck].ID == ItemIDChainOfTheTwilightOwl {
		partyBuffs.ChainOfTheTwilightOwl = true
	}
	if character.Equip[items.ItemSlotNeck].ID == ItemIDEyeOfTheNight {
		partyBuffs.EyeOfTheNight = true
	}
	if character.Equip[items.ItemSlotNeck].ID == ItemIDJadePendantOfBlasting {
		partyBuffs.JadePendantOfBlasting = true
	}
}

func (character *Character) initialize(agent Agent) {
	character.majorCooldownManager.initialize(character)

	character.gcdAction = &PendingAction{
		Priority: ActionPriorityGCD,
		OnAction: func(sim *Simulation) {
			character.TryUseCooldowns(sim)
			if character.GCD.IsReady(sim) {
				agent.OnGCDReady(sim)
			}
		},
	}
}

func (character *Character) Finalize(playerStats *proto.PlayerStats) {
	if character.Env.IsFinalized() {
		return
	}

	character.StatDependencyManager.Finalize()
	character.stats = character.ApplyStatDependencies(character.stats)

	character.PseudoStats.ParryHaste = character.PseudoStats.CanParry

	character.Unit.finalize()

	character.majorCooldownManager.finalize(character)

	if playerStats != nil {
		playerStats.FinalStats = character.GetStats().ToFloatArray()
		playerStats.Sets = character.GetActiveSetBonusNames()
		playerStats.Cooldowns = character.GetMajorCooldownIDs()
	}
}

func (character *Character) init(sim *Simulation, agent Agent) {
	character.Unit.init(sim)
}

func (character *Character) reset(sim *Simulation, agent Agent) {
	character.ExpectedBonusMana = 0
	character.majorCooldownManager.reset(sim)
	character.Unit.reset(sim, agent)

	if character.Type == PlayerUnit {
		character.SetGCDTimer(sim, 0)
	}

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
			petAgent.GetPet().advance(sim, elapsedTime)
		}
	}
}

func (character *Character) HasTrinketEquipped(itemID int32) bool {
	return character.Equip[items.ItemSlotTrinket1].ID == itemID ||
		character.Equip[items.ItemSlotTrinket2].ID == itemID
}

func (character *Character) HasRingEquipped(itemID int32) bool {
	return character.Equip[items.ItemSlotFinger1].ID == itemID ||
		character.Equip[items.ItemSlotFinger2].ID == itemID
}

func (character *Character) HasMetaGemEquipped(gemID int32) bool {
	for _, gem := range character.Equip[items.ItemSlotHead].Gems {
		if gem.ID == gemID {
			return true
		}
	}
	return false
}

// Returns the MH weapon if one is equipped, and null otherwise.
func (character *Character) GetMHWeapon() *items.Item {
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
func (character *Character) GetOHWeapon() *items.Item {
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
func (character *Character) GetRangedWeapon() *items.Item {
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

func (character *Character) GetMetricsProto(numIterations int32) *proto.UnitMetrics {
	metrics := character.Metrics.ToProto(numIterations)
	metrics.Name = character.Name
	metrics.Auras = character.auraTracker.GetMetricsProto(numIterations)

	metrics.Pets = []*proto.UnitMetrics{}
	for _, petAgent := range character.Pets {
		metrics.Pets = append(metrics.Pets, petAgent.GetPet().GetMetricsProto(numIterations))
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
