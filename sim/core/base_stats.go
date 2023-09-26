package core

import (
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

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

// These are also scattered in various dbc/casc files,
// `octbasempbyclass.txt`, `combatratings.txt`, `chancetospellcritbase.txt`, etc.

var RaceOffsets = map[proto.Race]stats.Stats{
	proto.Race_RaceUnknown: stats.Stats{},
	proto.Race_RaceHuman:   stats.Stats{},
	proto.Race_RaceOrc: {
		stats.Agility:   -3,
		stats.Strength:  3,
		stats.Intellect: -3,
		stats.Spirit:    2,
		stats.Stamina:   1,
	},
	proto.Race_RaceDwarf: {
		stats.Agility:   -4,
		stats.Strength:  5,
		stats.Intellect: -1,
		stats.Spirit:    -1,
		stats.Stamina:   1,
	},
	proto.Race_RaceNightElf: {
		stats.Agility:   4,
		stats.Strength:  -4,
		stats.Intellect: 0,
		stats.Spirit:    0,
		stats.Stamina:   0,
	},
	proto.Race_RaceUndead: {
		stats.Agility:   -2,
		stats.Strength:  -1,
		stats.Intellect: -2,
		stats.Spirit:    5,
		stats.Stamina:   0,
	},
	proto.Race_RaceTauren: {
		stats.Agility:   -4,
		stats.Strength:  5,
		stats.Intellect: -4,
		stats.Spirit:    2,
		stats.Stamina:   1,
	},
	proto.Race_RaceGnome: {
		stats.Agility:   2,
		stats.Strength:  -5,
		stats.Intellect: 3,
		stats.Spirit:    0,
		stats.Stamina:   0,
	},
	proto.Race_RaceTroll: {
		stats.Agility:   2,
		stats.Strength:  1,
		stats.Intellect: -4,
		stats.Spirit:    1,
		stats.Stamina:   0,
	},
	proto.Race_RaceBloodElf: {
		stats.Agility:   2,
		stats.Strength:  -3,
		stats.Intellect: 3,
		stats.Spirit:    -2,
		stats.Stamina:   0,
	},
	proto.Race_RaceDraenei: {
		stats.Agility:   -3,
		stats.Strength:  1,
		stats.Intellect: 0,
		stats.Spirit:    2,
		stats.Stamina:   0,
	},
}

var ClassBaseStats = map[proto.Class]stats.Stats{
	proto.Class_ClassUnknown: {},
	proto.Class_ClassWarrior: {
		stats.Health:      8121,
		stats.Agility:     113,
		stats.Strength:    174,
		stats.Intellect:   36,
		stats.Spirit:      59,
		stats.Stamina:     159,
		stats.AttackPower: float64(CharacterLevel)*3.0 - 20,
	},
	proto.Class_ClassPaladin: {
		stats.Health:      6934,
		stats.Agility:     90,
		stats.Strength:    151,
		stats.Intellect:   98,
		stats.Spirit:      105,
		stats.Stamina:     143,
		stats.AttackPower: float64(CharacterLevel)*3.0 - 20,
	},
	proto.Class_ClassHunter: {
		stats.Health:            7324,
		stats.Agility:           181,
		stats.Strength:          74,
		stats.Intellect:         90,
		stats.Spirit:            97,
		stats.Stamina:           128,
		stats.AttackPower:       float64(CharacterLevel)*2.0 - 20,
		stats.RangedAttackPower: float64(CharacterLevel)*2.0 - 10,
	},
	proto.Class_ClassRogue: {
		stats.Health:      7604,
		stats.Agility:     189,
		stats.Strength:    113,
		stats.Intellect:   43,
		stats.Spirit:      67,
		stats.Stamina:     105,
		stats.AttackPower: float64(CharacterLevel)*2.0 - 20,
	},
	proto.Class_ClassPriest: {
		stats.Health:    6960,
		stats.Agility:   51,
		stats.Strength:  43,
		stats.Intellect: 174,
		stats.Spirit:    181,
		stats.Stamina:   67,
	},
	proto.Class_ClassDeathknight: {
		stats.Health:      8121,
		stats.Agility:     112,
		stats.Strength:    175,
		stats.Intellect:   35,
		stats.Spirit:      59,
		stats.Stamina:     160,
		stats.AttackPower: float64(CharacterLevel)*3.0 - 20,
	},
	proto.Class_ClassShaman: {
		stats.Health:      6960,
		stats.Agility:     74,
		stats.Strength:    120,
		stats.Intellect:   128,
		stats.Spirit:      143,
		stats.Stamina:     136,
		stats.AttackPower: float64(CharacterLevel)*2.0 - 20,
	},
	proto.Class_ClassMage: {
		stats.Health:    6963,
		stats.Agility:   43,
		stats.Strength:  36,
		stats.Intellect: 181,
		stats.Spirit:    174,
		stats.Stamina:   59,
	},
	proto.Class_ClassWarlock: {
		stats.Health:      7164,
		stats.Agility:     67,
		stats.Strength:    59,
		stats.Intellect:   159,
		stats.Spirit:      166,
		stats.Stamina:     89,
		stats.AttackPower: -10,
	},
	proto.Class_ClassDruid: {
		stats.Health:      7417,
		stats.Agility:     82,
		stats.Strength:    89,
		stats.Intellect:   143,
		stats.Spirit:      159,
		stats.Stamina:     98,
		stats.AttackPower: -20,
	},
}

func AddBaseStatsCombo(r proto.Race, c proto.Class) {
	BaseStats[BaseStatsKey{Race: r, Class: c}] = ClassBaseStats[c].Add(RaceOffsets[r]).Add(ExtraClassBaseStats[c])
}

func init() {
	AddBaseStatsCombo(proto.Race_RaceTauren, proto.Class_ClassDruid)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassDruid)

	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassDeathknight)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassDeathknight)
	AddBaseStatsCombo(proto.Race_RaceGnome, proto.Class_ClassDeathknight)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassDeathknight)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassDeathknight)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassDeathknight)
	AddBaseStatsCombo(proto.Race_RaceTauren, proto.Class_ClassDeathknight)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassDeathknight)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassDeathknight)
	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassDeathknight)

	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceTauren, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassHunter)

	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceGnome, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassMage)

	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassPaladin)
	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassPaladin)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassPaladin)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassPaladin)

	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassPriest)

	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceGnome, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassRogue)

	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassShaman)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassShaman)
	AddBaseStatsCombo(proto.Race_RaceTauren, proto.Class_ClassShaman)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassShaman)

	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassWarlock)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassWarlock)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassWarlock)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassWarlock)
	AddBaseStatsCombo(proto.Race_RaceGnome, proto.Class_ClassWarlock)

	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceGnome, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceTauren, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassWarrior)
}
