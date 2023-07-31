package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var TalentTreeSizes = [3]int{28, 27, 26}

type Warlock struct {
	core.Character
	Talents  *proto.WarlockTalents
	Options  *proto.Warlock_Options
	Rotation *proto.Warlock_Rotation

	Pet *WarlockPet

	ShadowBolt         *core.Spell
	Incinerate         *core.Spell
	Immolate           *core.Spell
	UnstableAffliction *core.Spell
	Corruption         *core.Spell
	Haunt              *core.Spell
	LifeTap            *core.Spell
	DarkPact           *core.Spell
	ChaosBolt          *core.Spell
	SoulFire           *core.Spell
	Conflagrate        *core.Spell
	DrainSoul          *core.Spell
	Shadowburn         *core.Spell
	SearingPain        *core.Spell

	CurseOfElements      *core.Spell
	CurseOfElementsAuras core.AuraArray
	CurseOfWeakness      *core.Spell
	CurseOfWeaknessAuras core.AuraArray
	CurseOfTongues       *core.Spell
	CurseOfTonguesAuras  core.AuraArray
	CurseOfAgony         *core.Spell
	CurseOfDoom          *core.Spell
	Seed                 *core.Spell
	SeedDamageTracker    []float64

	ShadowEmbraceAuras     core.AuraArray
	NightfallProcAura      *core.Aura
	EradicationAura        *core.Aura
	DemonicEmpowerment     *core.Spell
	DemonicEmpowermentAura *core.Aura
	DemonicPactAura        *core.Aura
	DemonicSoulAura        *core.Aura
	Metamorphosis          *core.Spell
	MetamorphosisAura      *core.Aura
	ImmolationAura         *core.Spell
	HauntDebuffAuras       core.AuraArray
	MoltenCoreAura         *core.Aura
	DecimationAura         *core.Aura
	PyroclasmAura          *core.Aura
	BackdraftAura          *core.Aura
	EmpoweredImpAura       *core.Aura
	GlyphOfLifeTapAura     *core.Aura
	SpiritsoftheDamnedAura *core.Aura

	Infernal *InfernalPet
	Inferno  *core.Spell

	// The sum total of demonic pact spell power * seconds.
	DPSPAggregate float64
	PreviousTime  time.Duration

	petStmBonusSP float64
	acl           []ActionCondition

	// contains for each target the time the last shadowbolt was casted onto them
	corrRefreshList []time.Duration
}

type ACLaction int

const (
	ACLCast ACLaction = iota
	ACLNext
	ACLRecast
)

type ActionCondition struct {
	Spell     *core.Spell
	Condition func(*core.Simulation) (ACLaction, *core.Unit)
}

func (warlock *Warlock) GetCharacter() *core.Character {
	return &warlock.Character
}

func (warlock *Warlock) GetWarlock() *Warlock {
	return warlock
}

func (warlock *Warlock) GrandSpellstoneBonus() float64 {
	return core.TernaryFloat64(warlock.Options.WeaponImbue == proto.Warlock_Options_GrandSpellstone, 0.01, 0)
}
func (warlock *Warlock) GrandFirestoneBonus() float64 {
	return core.TernaryFloat64(warlock.Options.WeaponImbue == proto.Warlock_Options_GrandFirestone, 0.01, 0)
}

func (warlock *Warlock) Initialize() {
	warlock.registerIncinerateSpell()
	warlock.registerShadowBoltSpell()
	warlock.registerImmolateSpell()
	warlock.registerCorruptionSpell()
	warlock.registerCurseOfElementsSpell()
	warlock.registerCurseOfWeaknessSpell()
	warlock.registerCurseOfTonguesSpell()
	warlock.registerCurseOfAgonySpell()
	warlock.registerCurseOfDoomSpell()
	warlock.registerLifeTapSpell()
	warlock.registerSeedSpell()
	warlock.registerSoulFireSpell()
	warlock.registerUnstableAfflictionSpell()
	warlock.registerDrainSoulSpell()
	warlock.registerConflagrateSpell()
	warlock.registerHauntSpell()
	warlock.registerChaosBoltSpell()
	warlock.registerDemonicEmpowermentSpell()
	warlock.registerMetamorphosisSpell()
	warlock.registerDarkPactSpell()
	warlock.registerShadowBurnSpell()
	warlock.registerSearingPainSpell()
	warlock.registerInfernoSpell()
	warlock.registerBlackBook()

	warlock.defineRotation()

	precastSpell := warlock.ShadowBolt
	if warlock.Rotation.Type == proto.Warlock_Rotation_Destruction {
		precastSpell = warlock.SoulFire
	}
	// Do this post-finalize so cast speed is updated with new stats
	warlock.Env.RegisterPostFinalizeEffect(func() {
		// if itemswap is enabled, correct for any possible haste changes
		var correction stats.Stats
		if warlock.ItemSwap.IsEnabled() {
			correction = warlock.ItemSwap.CalcStatChanges([]proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand,
				proto.ItemSlot_ItemSlotOffHand, proto.ItemSlot_ItemSlotRanged})

			warlock.AddStats(correction)
			warlock.MultiplyCastSpeed(1.0)
		}

		if warlock.Options.Summon != proto.Warlock_Options_NoSummon && warlock.Talents.DemonicKnowledge > 0 {
			warlock.RegisterPrepullAction(-999*time.Second, func(sim *core.Simulation) {
				// TODO: investigate a better way of handling this like a "reverse inheritance" for pets.
				// TODO: this will break if we ever get stamina/intellect from procs, but there aren't
				// many such effects and none that we care about
				bonus := (warlock.Pet.GetStat(stats.Stamina) + warlock.Pet.GetStat(stats.Intellect)) *
					(0.04 * float64(warlock.Talents.DemonicKnowledge))
				if bonus != warlock.petStmBonusSP {
					warlock.AddStatDynamic(sim, stats.SpellPower, bonus-warlock.petStmBonusSP)
					warlock.petStmBonusSP = bonus
				}
			})
		}

		if warlock.IsUsingAPL {
			return
		}

		precastSpellAt := -warlock.ApplyCastSpeedForSpell(precastSpell.DefaultCast.CastTime, precastSpell)

		warlock.RegisterPrepullAction(precastSpellAt, func(sim *core.Simulation) {
			precastSpell.Cast(sim, warlock.CurrentTarget)
		})
		if warlock.GlyphOfLifeTapAura != nil || warlock.SpiritsoftheDamnedAura != nil {
			warlock.RegisterPrepullAction(precastSpellAt-warlock.SpellGCD(), func(sim *core.Simulation) {
				warlock.LifeTap.Cast(sim, nil)
			})
		}
		if warlock.ItemSwap.IsEnabled() {
			warlock.AddStats(correction.Multiply(-1))
			warlock.MultiplyCastSpeed(1.0)
		}
	})
}

func (warlock *Warlock) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.BloodPact = core.MaxTristate(raidBuffs.BloodPact, core.MakeTristateValue(
		warlock.Options.Summon == proto.Warlock_Options_Imp,
		warlock.Talents.ImprovedImp == 2,
	))

	raidBuffs.FelIntelligence = core.MaxTristate(raidBuffs.FelIntelligence, core.MakeTristateValue(
		warlock.Options.Summon == proto.Warlock_Options_Felhunter,
		warlock.Talents.ImprovedFelhunter == 2,
	))
}

func (warlock *Warlock) Reset(sim *core.Simulation) {
	if sim.CurrentTime == 0 {
		warlock.petStmBonusSP = 0
	}

	warlock.ItemSwap.SwapItems(sim, []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand,
		proto.ItemSlot_ItemSlotOffHand, proto.ItemSlot_ItemSlotRanged}, false)
	warlock.corrRefreshList = make([]time.Duration, len(warlock.Env.Encounter.TargetUnits))
	warlock.setupCooldowns(sim)
}

func NewWarlock(character core.Character, options *proto.Player) *Warlock {
	warlockOptions := options.GetWarlock()

	warlock := &Warlock{
		Character: character,
		Talents:   &proto.WarlockTalents{},
		Options:   warlockOptions.Options,
		Rotation:  warlockOptions.Rotation,
	}
	core.FillTalentsProto(warlock.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)
	warlock.EnableManaBar()

	warlock.AddStatDependency(stats.Strength, stats.AttackPower, 1)

	if warlock.Options.Armor == proto.Warlock_Options_FelArmor {
		demonicAegisMultiplier := 1 + float64(warlock.Talents.DemonicAegis)*0.1
		amount := 180.0 * demonicAegisMultiplier
		warlock.AddStat(stats.SpellPower, amount)
		warlock.AddStatDependency(stats.Spirit, stats.SpellPower, 0.3*demonicAegisMultiplier)
	}

	if warlock.Options.Summon != proto.Warlock_Options_NoSummon {
		warlock.Pet = warlock.NewWarlockPet()
	}

	warlock.Infernal = warlock.NewInfernal()

	if warlock.Rotation.Type == proto.Warlock_Rotation_Affliction && warlock.Rotation.EnableWeaponSwap {
		warlock.EnableItemSwap(warlock.Rotation.WeaponSwap, 1, 1, 1)
	}

	warlock.applyWeaponImbue()

	return warlock
}

func RegisterWarlock() {
	core.RegisterAgentFactory(
		proto.Player_Warlock{},
		proto.Spec_SpecWarlock,
		func(character core.Character, options *proto.Player) core.Agent {
			return NewWarlock(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Warlock)
			if !ok {
				panic("Invalid spec value for Warlock!")
			}
			player.Spec = playerSpec
		},
	)
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceBloodElf, Class: proto.Class_ClassWarlock}] = stats.Stats{
		stats.Health:    7164,
		stats.Strength:  56,
		stats.Agility:   69,
		stats.Stamina:   89,
		stats.Intellect: 162,
		stats.Spirit:    164,
		stats.Mana:      3856,
		stats.SpellCrit: 1.697 * core.CritRatingPerCritChance,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceOrc, Class: proto.Class_ClassWarlock}] = stats.Stats{
		stats.Health:    7164,
		stats.Strength:  62,
		stats.Agility:   64,
		stats.Stamina:   90,
		stats.Intellect: 156,
		stats.Spirit:    168,
		stats.Mana:      3856,
		stats.SpellCrit: 1.697 * core.CritRatingPerCritChance,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassWarlock}] = stats.Stats{
		stats.Health:    7164,
		stats.Strength:  58,
		stats.Agility:   65,
		stats.Stamina:   89,
		stats.Intellect: 157,
		stats.Spirit:    171,
		stats.Mana:      3856,
		stats.SpellCrit: 1.697 * core.CritRatingPerCritChance,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassWarlock}] = stats.Stats{
		stats.Health:    7164,
		stats.Strength:  59,
		stats.Agility:   67,
		stats.Stamina:   89,
		stats.Intellect: 159,
		stats.Spirit:    166, // racial makes this 170
		stats.Mana:      3856,
		stats.SpellCrit: 1.697 * core.CritRatingPerCritChance,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceGnome, Class: proto.Class_ClassWarlock}] = stats.Stats{
		stats.Health:    7164,
		stats.Strength:  54,
		stats.Agility:   69,
		stats.Stamina:   89,
		stats.Intellect: 162, // racial makes this 170
		stats.Spirit:    166,
		stats.Mana:      3856,
		stats.SpellCrit: 1.697 * core.CritRatingPerCritChance,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.CritRatingPerCritChance,
	}
}

// Agent is a generic way to access underlying warlock on any of the agents.
type WarlockAgent interface {
	GetWarlock() *Warlock
}

func (warlock *Warlock) HasMajorGlyph(glyph proto.WarlockMajorGlyph) bool {
	return warlock.HasGlyph(int32(glyph))
}

func (warlock *Warlock) HasMinorGlyph(glyph proto.WarlockMinorGlyph) bool {
	return warlock.HasGlyph(int32(glyph))
}
