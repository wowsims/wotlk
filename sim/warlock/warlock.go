package warlock

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type Warlock struct {
	core.Character
	Talents  proto.WarlockTalents
	Options  proto.Warlock_Options
	Rotation proto.Warlock_Rotation

	Shadowbolt     *core.Spell
	Incinerate     *core.Spell
	Immolate       *core.Spell
	ImmolateDot    *core.Dot
	UnstableAff    *core.Spell
	UnstableAffDot *core.Dot
	Corruption     *core.Spell
	CorruptionDot  *core.Dot
	/*	Haunt		   *core.Spell
		Haunt		   *core.Aura

		DemonicEmpowerment		   *core.Aura
	*/
	LifeTap *core.Spell

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

	Seeds    []*core.Spell
	SeedDots []*core.Dot

	NightfallProcAura *core.Aura
	ShadowEmbraceAura *core.Aura

	Pet *WarlockPet

	DoingRegen bool
}

func (warlock *Warlock) GetCharacter() *core.Character {
	return &warlock.Character
}

func (warlock *Warlock) GetWarlock() *Warlock {
	return warlock
}

func (warlock *Warlock) Initialize() {
	warlock.registerIncinerateSpell()
	warlock.registerShadowboltSpell()
	warlock.registerImmolateSpell()
	warlock.registerCorruptionSpell()
	warlock.registerCurseOfElementsSpell()
	warlock.registerCurseOfWeaknessSpell()
	warlock.registerCurseOfTonguesSpell()
	warlock.registerCurseOfAgonySpell()
	warlock.registerCurseOfDoomSpell()
	warlock.registerLifeTapSpell()
	if warlock.Talents.UnstableAffliction {
		warlock.registerUnstableAffSpell()
	}
	warlock.registerSeedSpell()
}

func (warlock *Warlock) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.BloodPact = core.MaxTristate(raidBuffs.BloodPact, core.MakeTristateValue(
		warlock.Options.Summon == proto.Warlock_Options_Imp,
		warlock.Talents.ImprovedImp == 2))
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

	warlock.Character.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Strength,
		ModifiedStat: stats.AttackPower,
		Modifier: func(strength float64, attackPower float64) float64 {
			return attackPower + strength*2
		},
	})

	if warlock.Options.Armor == proto.Warlock_Options_FelArmor {
		amount := 100.0
		amount *= 1 + float64(warlock.Talents.DemonicAegis)*0.1
		warlock.AddStat(stats.SpellPower, amount)
	}

	/*	if warlock.Talents.DemonicSacrifice && warlock.Options.SacrificeSummon {
		switch warlock.Options.Summon {
		case proto.Warlock_Options_Succubus:
			warlock.PseudoStats.ShadowDamageDealtMultiplier *= 1.15
		case proto.Warlock_Options_Imp:
			warlock.PseudoStats.FireDamageDealtMultiplier *= 1.15
		case proto.Warlock_Options_Felgaurd:
			warlock.PseudoStats.ShadowDamageDealtMultiplier *= 1.10
		}
	} else*/
	if warlock.Options.Summon != proto.Warlock_Options_NoSummon {
		warlock.Pet = warlock.NewWarlockPet()
	}

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
		stats.Health:      7906,
		stats.Strength:    56,
		stats.Agility:     69,
		stats.Stamina:     95,
		stats.Intellect:   163,
		stats.Spirit:      165,
		stats.Mana:        6021,
		stats.SpellCrit:   1.697 * core.CritRatingPerCritChance,
		stats.AttackPower: 102,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceOrc, Class: proto.Class_ClassWarlock}] = stats.Stats{
		stats.Health:      7946,
		stats.Strength:    62,
		stats.Agility:     64,
		stats.Stamina:     99,
		stats.Intellect:   156,
		stats.Spirit:      169,
		stats.Mana:        5916,
		stats.SpellCrit:   1.697 * core.CritRatingPerCritChance,
		stats.AttackPower: 114,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassWarlock}] = stats.Stats{
		stats.Health:      7946,
		stats.Strength:    58,
		stats.Agility:     65,
		stats.Stamina:     99,
		stats.Intellect:   157,
		stats.Spirit:      171,
		stats.Mana:        5931,
		stats.SpellCrit:   1.697 * core.CritRatingPerCritChance,
		stats.AttackPower: 106,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassWarlock}] = stats.Stats{
		stats.Health:      7926,
		stats.Strength:    59,
		stats.Agility:     67,
		stats.Stamina:     97,
		stats.Intellect:   159,
		stats.Spirit:      178,
		stats.Mana:        5961,
		stats.SpellCrit:   1.697 * core.CritRatingPerCritChance,
		stats.AttackPower: 108,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceGnome, Class: proto.Class_ClassWarlock}] = stats.Stats{
		stats.Health:      7916,
		stats.Strength:    54,
		stats.Agility:     70,
		stats.Stamina:     96,
		stats.Intellect:   178,
		stats.Spirit:      166,
		stats.Mana:        6246,
		stats.SpellCrit:   1.697 * core.CritRatingPerCritChance,
		stats.AttackPower: 98,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.CritRatingPerCritChance,
	}
}

// Agent is a generic way to access underlying warlock on any of the agents.
type WarlockAgent interface {
	GetWarlock() *Warlock
}
