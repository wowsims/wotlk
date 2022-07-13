package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type DeathKnight struct {
	core.Character
	Talents  proto.DeathKnightTalents
	Options  proto.DeathKnight_Options
	Rotation proto.DeathKnight_Rotation

	Presence Presence

	IcyTouch     *core.Spell
	PlagueStrike *core.Spell
	Obliterate   *core.Spell
	BloodStrike  *core.Spell
	HowlingBlast *core.Spell
	//FrostStrike      *core.Spell
	//HornOfWinter     *core.Spell
	//UnbreakableArmor *core.Spell
	//ArmyOfTheDead    *core.Spell
	//RaiseDead        *core.Spell

	// "CDs"
	BloodTap *core.Spell

	// Diseases
	FrostFever         *core.Spell
	FrostFeverDisease  *core.Dot
	BloodPlague        *core.Spell
	BloodPlagueDisease *core.Dot

	// Talent Auras
	KillingMachineAura *core.Aura
	IcyTalonsAura      *core.Aura

	// Presences
	BloodPresenceAura  *core.Aura
	FrostPresenceAura  *core.Aura
	UnholyPresenceAura *core.Aura

	// Debuffs
	IcyTouchAura *core.Aura
}

func (deathKnight *DeathKnight) GetCharacter() *core.Character {
	return &deathKnight.Character
}

func (deathKnight *DeathKnight) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {

}

func (deathKnight *DeathKnight) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	if deathKnight.Talents.AbominationsMight > 0 {
		raidBuffs.AbominationsMight = true
	}

	if deathKnight.Talents.ImprovedIcyTalons {
		raidBuffs.IcyTalons = true
	}
}

func (deathKnight *DeathKnight) Initialize() {
	deathKnight.registerPresences()
	deathKnight.registerIcyTouchSpell()
	deathKnight.registerPlagueStrikeSpell()
	deathKnight.registerObliterateSpell()
	deathKnight.registerBloodStrikeSpell()
	deathKnight.registerBloodTapSpell()
	deathKnight.registerHowlingBlastSpell()
	deathKnight.registerDiseaseDots()
}

func (deathKnight *DeathKnight) Reset(sim *core.Simulation) {
	deathKnight.ResetRunicPowerBar(sim)
	deathKnight.BloodPresenceAura.Activate(sim)
	deathKnight.Presence = BloodPresence
}

func NewDeathKnight(character core.Character, options proto.Player) *DeathKnight {
	deathKnightOptions := options.GetDeathKnight()

	deathKnight := &DeathKnight{
		Character: character,
		Talents:   *deathKnightOptions.Talents,
		Options:   *deathKnightOptions.Options,
		Rotation:  *deathKnightOptions.Rotation,
	}

	maxRunicPower := 100.0
	if deathKnight.Talents.RunicPowerMastery == 1 {
		maxRunicPower = 115.0
	} else if deathKnight.Talents.RunicPowerMastery == 2 {
		maxRunicPower = 130.0
	}
	deathKnight.EnableRunicPowerBar(
		maxRunicPower,
		func(sim *core.Simulation) {},
		func(sim *core.Simulation) {
			if deathKnight.GCD.IsReady(sim) {
				deathKnight.tryUseGCD(sim)
			}
		},
		func(sim *core.Simulation) {
			if deathKnight.GCD.IsReady(sim) {
				deathKnight.tryUseGCD(sim)
			}
		},
		func(sim *core.Simulation) {
			if deathKnight.GCD.IsReady(sim) {
				deathKnight.tryUseGCD(sim)
			}
		},
		func(sim *core.Simulation) {
			if deathKnight.GCD.IsReady(sim) {
				deathKnight.tryUseGCD(sim)
			}
		},
	)

	deathKnight.EnableAutoAttacks(deathKnight, core.AutoAttackOptions{
		MainHand:       deathKnight.WeaponFromMainHand(deathKnight.DefaultMeleeCritMultiplier()),
		OffHand:        deathKnight.WeaponFromOffHand(deathKnight.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	deathKnight.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Agility,
		ModifiedStat: stats.MeleeCrit,
		Modifier: func(agility float64, meleecrit float64) float64 {
			return meleecrit + (agility/62.5)*core.CritRatingPerCritChance
		},
	})
	deathKnight.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Agility,
		ModifiedStat: stats.Dodge,
		Modifier: func(agility float64, dodge float64) float64 {
			return dodge + (agility/84.74576271)*core.DodgeRatingPerDodgeChance
		},
	})
	deathKnight.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Strength,
		ModifiedStat: stats.AttackPower,
		Modifier: func(strength float64, attackPower float64) float64 {
			return attackPower + strength*2
		},
	})

	return deathKnight
}

func RegisterDeathKnight() {
	core.RegisterAgentFactory(
		proto.Player_DeathKnight{},
		proto.Spec_SpecDeathKnight,
		func(character core.Character, options proto.Player) core.Agent {
			return NewDeathKnight(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_DeathKnight)
			if !ok {
				panic("Invalid spec value for DeathKnight!")
			}
			player.Spec = playerSpec
		},
	)
}

func (deathKnight *DeathKnight) DiseasesAreActive() bool {
	return deathKnight.FrostFeverDisease.IsActive() || deathKnight.BloodPlagueDisease.IsActive()
}

func (deathKnight *DeathKnight) secondaryCritModifier(applyGuile bool) float64 {
	secondaryModifier := 0.0
	if applyGuile {
		secondaryModifier += 0.15 * float64(deathKnight.Talents.GuileOfGorefiend)
	}
	return secondaryModifier
}
func (deathKnight *DeathKnight) critMultiplier(applyGuile bool) float64 {
	return deathKnight.MeleeCritMultiplier(1.0, deathKnight.secondaryCritModifier(applyGuile))
}
func (deathKnight *DeathKnight) spellCritMultiplier(applyGuile bool) float64 {
	return deathKnight.SpellCritMultiplier(1.0, deathKnight.secondaryCritModifier(applyGuile))
}
func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceGnome, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceNightElf, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceOrc, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTauren, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceBloodElf, Class: proto.Class_ClassDeathKnight}] = stats.Stats{
		stats.Health:      7941,
		stats.Strength:    180,
		stats.Agility:     112,
		stats.Stamina:     160,
		stats.Intellect:   35,
		stats.Spirit:      63,
		stats.AttackPower: 220,
		stats.MeleeCrit:   3.188 * core.CritRatingPerCritChance,
		stats.Dodge:       3.664 * core.DodgeRatingPerDodgeChance,
	}
}

// Agent is a generic way to access underlying warrior on any of the agents.
type DeathKnightAgent interface {
	GetDeathKnight() *DeathKnight
}
