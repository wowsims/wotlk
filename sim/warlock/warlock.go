package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type Warlock struct {
	core.Character
	Talents  *proto.WarlockTalents
	Options  *proto.Warlock_Options
	Rotation *proto.Warlock_Rotation

	procTrackers []*ProcTracker
	majorCds     []*core.MajorCooldown

	Pet *WarlockPet

	ShadowBolt            *core.Spell
	Incinerate            *core.Spell
	Immolate              *core.Spell
	ImmolateDot           *core.Dot
	UnstableAffliction    *core.Spell
	UnstableAfflictionDot *core.Dot
	Corruption            *core.Spell
	CorruptionDot         *core.Dot
	Haunt                 *core.Spell
	LifeTap               *core.Spell
	DarkPact              *core.Spell
	ChaosBolt             *core.Spell
	SoulFire              *core.Spell
	Conflagrate           *core.Spell
	ConflagrateDot        *core.Dot
	DrainSoul             *core.Spell
	DrainSoulDot          *core.Dot
	Shadowburn            *core.Spell

	CurseOfElements     *core.Spell
	CurseOfElementsAura *core.Aura
	CurseOfWeakness     *core.Spell
	CurseOfWeaknessAura *core.Aura
	CurseOfTongues      *core.Spell
	CurseOfTonguesAura  *core.Aura
	CurseOfAgony        *core.Spell
	CurseOfAgonyDot     *core.Dot
	CurseOfDoom         *core.Spell
	CurseOfDoomDot      *core.Dot
	Seeds               []*core.Spell
	SeedDots            []*core.Dot

	NightfallProcAura          *core.Aura
	EradicationAura            *core.Aura
	DemonicEmpowerment         *core.Spell
	DemonicEmpowermentAura     *core.Aura
	DemonicPactAura            *core.Aura
	DemonicSoulAura            *core.Aura
	Metamorphosis              *core.Spell
	MetamorphosisAura          *core.Aura
	ImmolationAura             *core.Spell
	ImmolationAuraDot          *core.Dot
	HauntDebuffAura            *core.Aura
	MoltenCoreAura             *core.Aura
	DecimationAura             *core.Aura
	PyroclasmAura              *core.Aura
	BackdraftAura              *core.Aura
	EmpoweredImpAura           *core.Aura
	GlyphOfLifeTapAura         *core.Aura
	FakeSpiritsoftheDamnedAura *core.Aura

	// Rotation related memory
	CorruptionRolloverPower float64
	DrainSoulRolloverPower  float64
	// The sum total of demonic pact spell power * seconds.
	DPSPAggregate  float64
	PreviousTime   time.Duration
	SpellsRotation []SpellRotation

	petStmBonusSP float64

	// set bonus cache
	T7TwoSetBonus   bool
	T7FourSetBonus  bool
	T8TwoSetBonus   bool
	T8FourSetBonus  bool
	T9TwoSetBonus   bool
	T9FourSetBonus  bool
	T10TwoSetBonus  bool
	T10FourSetBonus bool
}

type SpellRotation struct {
	Spell    *core.Spell
	CastIn   CastReadyness
	Priority int
}

type CastReadyness func(*core.Simulation) time.Duration

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

	warlock.registerFakeSpiritsoftheDamnedAura()
	warlock.registerSetBonuses()
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
	if warlock.Talents.Metamorphosis {
		warlock.registerMetamorphosisSpell()
		warlock.registerImmolationAuraSpell()
	}
	warlock.registerDarkPactSpell()
	warlock.registerShadowBurnSpell()

	warlock.defineRotation()
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

func (warlock *Warlock) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (warlock *Warlock) Reset(sim *core.Simulation) {
	if sim.CurrentTime == 0 {
		warlock.petStmBonusSP = 0
	}
}

func NewWarlock(character core.Character, options *proto.Player) *Warlock {
	warlockOptions := options.GetWarlock()

	warlock := &Warlock{
		Character: character,
		Talents:   warlockOptions.Talents,
		Options:   warlockOptions.Options,
		Rotation:  warlockOptions.Rotation,
		// manaTracker:           common.NewManaSpendingRateTracker(),
	}
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
