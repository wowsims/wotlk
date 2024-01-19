package core

import (
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

type BaseStatsKey struct {
	Race  proto.Race
	Class proto.Class
	Level int
}

var BaseStats = map[BaseStatsKey]stats.Stats{}

// To calculate base stats, get a naked toon of desired level of the race/class you want, ideally without any talents to mess up base stats.
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
}

// TODO: Classic base stats
var ClassBaseStats = map[proto.Class]map[int]stats.Stats{
	proto.Class_ClassUnknown: {},
	proto.Class_ClassWarrior: {
		25: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     44,
			stats.Strength:    69,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 25*3 - 20,
		},
		40: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     59,
			stats.Strength:    94,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 40*3 - 20,
		},
		50: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     71,
			stats.Strength:    113,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 50*3 - 20,
		},
		60: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     85,
			stats.Strength:    135,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 60*3 - 20,
		},
	},
	proto.Class_ClassPaladin: {
		25: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		40: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		50: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		60: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
	},
	proto.Class_ClassHunter: {
		25: {
			stats.Health:            292,
			stats.Mana:              611,
			stats.Agility:           55,
			stats.Strength:          31,
			stats.Intellect:         34,
			stats.Spirit:            37,
			stats.Stamina:           43,
			stats.AttackPower:       25*2 - 20,
			stats.RangedAttackPower: 25*2 - 20,
		},
		40: {
			stats.Health:            667,
			stats.Mana:              1105,
			stats.Agility:           81,
			stats.Strength:          40,
			stats.Intellect:         46,
			stats.Spirit:            49,
			stats.Stamina:           61,
			stats.AttackPower:       40*2 - 20,
			stats.RangedAttackPower: 40*2 - 20,
		},
		50: {
			stats.Health:            1047,
			stats.Mana:              1420,
			stats.Agility:           102,
			stats.Strength:          47,
			stats.Intellect:         55,
			stats.Spirit:            59,
			stats.Stamina:           74,
			stats.AttackPower:       50*2 - 20,
			stats.RangedAttackPower: 50*2 - 20,
		},
		60: {
			stats.Health:            1467,
			stats.Mana:              1720,
			stats.Agility:           125,
			stats.Strength:          55,
			stats.Intellect:         65,
			stats.Spirit:            70,
			stats.Stamina:           90,
			stats.AttackPower:       60*2 - 20,
			stats.RangedAttackPower: 60*2 - 20,
		},
	},
	proto.Class_ClassRogue: {
		25: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		40: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		50: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		60: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
	},
	proto.Class_ClassPriest: {
		25: {
			stats.Health:    222,
			stats.Mana:      217,
			stats.Agility:   26,
			stats.Strength:  22,
			stats.Intellect: 53,
			stats.Spirit:    55,
			stats.Stamina:   44,
		},
		40: {
			stats.Health:    457,
			stats.Mana:      631,
			stats.Agility:   0,
			stats.Strength:  26,
			stats.Intellect: 78,
			stats.Spirit:    81,
			stats.Stamina:   39,
		},
		50: {
			stats.Health:    792,
			stats.Mana:      886,
			stats.Agility:   35,
			stats.Strength:  29,
			stats.Intellect: 98,
			stats.Spirit:    102,
			stats.Stamina:   45,
		},
		60: {
			stats.Health:    1217,
			stats.Mana:      1096,
			stats.Agility:   40,
			stats.Strength:  32,
			stats.Intellect: 120,
			stats.Spirit:    125,
			stats.Stamina:   52,
		},
	},
	proto.Class_ClassShaman: {
		25: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		40: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		50: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		60: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
	},
	proto.Class_ClassMage: {
		25: {
			stats.Health:      135,
			stats.Mana:        201,
			stats.Agility:     26,
			stats.Strength:    23,
			stats.Intellect:   61,
			stats.Spirit:      53,
			stats.Stamina:     27,
			stats.AttackPower: 26,
		},
		40: {
			stats.Health:      450,
			stats.Mana:        573,
			stats.Agility:     30,
			stats.Strength:    26,
			stats.Intellect:   89,
			stats.Spirit:      78,
			stats.Stamina:     33,
			stats.AttackPower: 32,
		},
		50: {
			stats.Health:      775,
			stats.Mana:        768,
			stats.Agility:     33,
			stats.Strength:    28,
			stats.Intellect:   112,
			stats.Spirit:      98,
			stats.Stamina:     38,
			stats.AttackPower: 36,
		},
		60: {
			stats.Health:      1190,
			stats.Mana:        2085,
			stats.Agility:     36,
			stats.Strength:    30,
			stats.Intellect:   136,
			stats.Spirit:      120,
			stats.Stamina:     44,
			stats.AttackPower: 40,
		},
	},
	proto.Class_ClassWarlock: {
		25: {
			stats.Health:      99,
			stats.Mana:        498,
			stats.Strength:    28,
			stats.Agility:     30,
			stats.Stamina:     35,
			stats.Intellect:   50,
			stats.Spirit:      52,
			stats.AttackPower: 36,
		},
		40: {
			stats.Health:      454,
			stats.Mana:        643,
			stats.Agility:     38,
			stats.Strength:    34,
			stats.Intellect:   79,
			stats.Spirit:      75,
			stats.Stamina:     45,
			stats.AttackPower: 48,
		},
		50: {
			stats.Health:      799,
			stats.Mana:        883,
			stats.Agility:     44,
			stats.Strength:    39,
			stats.Intellect:   99,
			stats.Spirit:      93,
			stats.Stamina:     54,
			stats.AttackPower: 58,
		},
		60: {
			stats.Health:      1234,
			stats.Mana:        1093,
			stats.Agility:     51,
			stats.Strength:    45,
			stats.Intellect:   121,
			stats.Spirit:      115,
			stats.Stamina:     64,
			stats.AttackPower: 70,
		},
	},
	proto.Class_ClassDruid: {
		25: {
			stats.Health:      138,
			stats.Mana:        479,
			stats.Agility:     34,
			stats.Strength:    36,
			stats.Intellect:   47,
			stats.Spirit:      50,
			stats.Stamina:     35,
			stats.AttackPower: -20,
		},
		40: {
			stats.Health:      503,
			stats.Mana:        1005,
			stats.Agility:     44,
			stats.Strength:    47,
			stats.Intellect:   67,
			stats.Spirit:      72,
			stats.Stamina:     48,
			stats.AttackPower: -20,
		},
		50: {
			stats.Health:      858,
			stats.Mana:        784,
			stats.Agility:     52,
			stats.Strength:    56,
			stats.Intellect:   82,
			stats.Spirit:      90,
			stats.Stamina:     58,
			stats.AttackPower: -20,
		},
		60: {
			stats.Health:      1303,
			stats.Mana:        964,
			stats.Agility:     61,
			stats.Strength:    66,
			stats.Intellect:   100,
			stats.Spirit:      110,
			stats.Stamina:     69,
			stats.AttackPower: -20,
		},
	},
}

// Retrieves base stats, with race offsets, and crit rating adjustments per level
func getBaseStatsCombo(r proto.Race, c proto.Class, level int) stats.Stats {
	if level == 0 {
		level = 60
	}

	starting := ClassBaseStats[c][level]

	return starting.Add(RaceOffsets[r]).Add(ExtraClassBaseStats[c][level])
}
