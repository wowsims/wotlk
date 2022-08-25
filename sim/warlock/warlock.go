package warlock

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
	"time"
)

type Warlock struct {
	core.Character
	Talents  proto.WarlockTalents
	Options  proto.Warlock_Options
	Rotation proto.Warlock_Rotation

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
	DrainSoulChannelling  *core.Spell
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

	NightfallProcAura      *core.Aura
	EradicationAura        *core.Aura
	DemonicEmpowerment     *core.Spell
	DemonicEmpowermentAura *core.Aura
	Metamorphosis          *core.Spell
	MetamorphosisAura      *core.Aura
	ImmolationAura         *core.Spell
	ImmolationAuraDot      *core.Dot
	MoltenCoreAura         *core.Aura
	DecimationAura         *core.Aura
	PyroclasmAura          *core.Aura
	BackdraftAura          *core.Aura
	EmpoweredImpAura       *core.Aura
	GlyphOfLifeTapAura     *core.Aura

	// Rotation related memory
	CorruptionRolloverPower float64
	DPSPAverage            float64
	PreviousTime           time.Duration
	SpellsRotation         []SpellRotation
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

	if warlock.Talents.DemonicEmpowerment {
		warlock.registerDemonicEmpowermentSpell()
	}
	if warlock.Talents.Metamorphosis {
		warlock.registerMetamorphosisSpell()
		warlock.registerImmolationAuraSpell()
	}
	if warlock.Talents.DarkPact {
		warlock.registerDarkPactSpell()
	}
	if warlock.Talents.Shadowburn {
		warlock.registerShadowBurnSpell()
	}

	warlock.defineRotation()
}

func (warlock *Warlock) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.BloodPact = core.MaxTristate(raidBuffs.BloodPact, core.MakeTristateValue(
		warlock.Options.Summon == proto.Warlock_Options_Imp,
		warlock.Talents.ImprovedImp == 2,
	))

	if warlock.Talents.DemonicPact > 0 {
		raidBuffs.DemonicPact = int32(float64(stats.SpellPower) * 0.02 * float64(warlock.Talents.DemonicPact) * 1.111)
		// * 1.1 because the buff gets 10% better after the first refresh and so on every 20s
	}
}

func (warlock *Warlock) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (warlock *Warlock) Reset(sim *core.Simulation) {

}

func NewWarlock(character core.Character, options proto.Player) *Warlock {
	warlockOptions := options.GetWarlock()

	warlock := &Warlock{
		Character: character,
		Talents:   *warlockOptions.Talents,
		Options:   *warlockOptions.Options,
		Rotation:  *warlockOptions.Rotation,
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
		func(character core.Character, options proto.Player) core.Agent {
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
