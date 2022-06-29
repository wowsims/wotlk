package warlock

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
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
	SiphonLife     *core.Spell
	SiphonLifeDot  *core.Dot

	LifeTap *core.Spell

	CurseOfElements         *core.Spell
	CurseOfElementsAura     *core.Aura
	CurseOfRecklessness     *core.Spell
	CurseOfRecklessnessAura *core.Aura
	CurseOfTongues          *core.Spell
	CurseOfTonguesAura      *core.Aura
	CurseOfAgony            *core.Spell
	CurseOfAgonyDot         *core.Dot
	CurseOfDoom             *core.Spell
	CurseOfDoomDot          *core.Dot

	AmplifyCurse     *core.Spell
	AmplifyCurseAura *core.Aura

	Seeds    []*core.Spell
	SeedDots []*core.Dot

	NightfallProcAura *core.Aura
	ImpShadowboltAura *core.Aura

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
	warlock.registerCurseOfRecklessnessSpell()
	warlock.registerCurseOfTonguesSpell()
	warlock.registerCurseOfAgonySpell()
	warlock.registerCurseOfDoomSpell()
	warlock.registerLifeTapSpell()
	if warlock.Talents.UnstableAffliction {
		warlock.registerUnstableAffSpell()
	}
	if warlock.Talents.SiphonLife {
		warlock.registerSiphonLifeSpell()
	}
	warlock.registerSeedSpell()
}

func (warlock *Warlock) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
}

func (warlock *Warlock) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
	partyBuffs.BloodPact = core.MaxTristate(partyBuffs.BloodPact, core.MakeTristateValue(
		warlock.Options.Summon == proto.Warlock_Options_Imp &&
			(!warlock.Talents.DemonicSacrifice || !warlock.Options.SacrificeSummon),
		warlock.Talents.ImprovedImp == 2))
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
		SourceStat:   stats.Intellect,
		ModifiedStat: stats.SpellCrit,
		Modifier: func(intellect float64, spellCrit float64) float64 {
			return spellCrit + (intellect/81.92)*core.SpellCritRatingPerCritChance
		},
	})

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

	if warlock.Talents.DemonicSacrifice && warlock.Options.SacrificeSummon {
		switch warlock.Options.Summon {
		case proto.Warlock_Options_Succubus:
			warlock.PseudoStats.ShadowDamageDealtMultiplier *= 1.15
		case proto.Warlock_Options_Imp:
			warlock.PseudoStats.FireDamageDealtMultiplier *= 1.15
		case proto.Warlock_Options_Felgaurd:
			warlock.PseudoStats.ShadowDamageDealtMultiplier *= 1.10
		}
	} else if warlock.Options.Summon != proto.Warlock_Options_NoSummon {
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
		stats.Health:      3310,
		stats.Strength:    48,
		stats.Agility:     60,
		stats.Stamina:     75,
		stats.Intellect:   137,
		stats.Spirit:      130,
		stats.Mana:        2335,
		stats.SpellCrit:   1.697 * core.SpellCritRatingPerCritChance,
		stats.AttackPower: 86,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.MeleeCritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceOrc, Class: proto.Class_ClassWarlock}] = stats.Stats{
		stats.Health:      3310,
		stats.Strength:    54,
		stats.Agility:     55,
		stats.Stamina:     78,
		stats.Intellect:   130,
		stats.Spirit:      134,
		stats.Mana:        2335,
		stats.SpellCrit:   1.697 * core.SpellCritRatingPerCritChance,
		stats.AttackPower: 86,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.MeleeCritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassWarlock}] = stats.Stats{
		stats.Health:      3310,
		stats.Strength:    50,
		stats.Agility:     56,
		stats.Stamina:     77,
		stats.Intellect:   131,
		stats.Spirit:      144,
		stats.Mana:        2335,
		stats.SpellCrit:   1.697 * core.SpellCritRatingPerCritChance,
		stats.AttackPower: 86,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.MeleeCritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassWarlock}] = stats.Stats{
		stats.Health:      3300,
		stats.Strength:    51,
		stats.Agility:     58,
		stats.Stamina:     76,
		stats.Intellect:   133,
		stats.Spirit:      144,
		stats.Mana:        2395,
		stats.SpellCrit:   1.697 * core.SpellCritRatingPerCritChance,
		stats.AttackPower: 86,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.MeleeCritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceGnome, Class: proto.Class_ClassWarlock}] = stats.Stats{
		stats.Health:      3310,
		stats.Strength:    46,
		stats.Agility:     61,
		stats.Stamina:     75,
		stats.Intellect:   143,
		stats.Spirit:      131,
		stats.Mana:        2335,
		stats.SpellCrit:   1.697 * core.SpellCritRatingPerCritChance,
		stats.AttackPower: 86,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.MeleeCritRatingPerCritChance,
	}
}

// Agent is a generic way to access underlying warlock on any of the agents.
type WarlockAgent interface {
	GetWarlock() *Warlock
}
