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
		stats.Spirit:    3,
		stats.Stamina:   2,
	},
	proto.Race_RaceDwarf: {
		stats.Agility:   -4,
		stats.Strength:  2,
		stats.Intellect: -1,
		stats.Spirit:    -1,
		stats.Stamina:   3,
	},
	proto.Race_RaceNightElf: {
		stats.Agility:   5,
		stats.Strength:  -3,
		stats.Intellect: 0,
		stats.Spirit:    0,
		stats.Stamina:   -1,
	},
	proto.Race_RaceUndead: {
		stats.Agility:   -2,
		stats.Strength:  -1,
		stats.Intellect: -2,
		stats.Spirit:    5,
		stats.Stamina:   1,
	},
	proto.Race_RaceTauren: {
		stats.Agility:   -5,
		stats.Strength:  5,
		stats.Intellect: -5,
		stats.Spirit:    2,
		stats.Stamina:   2,
	},
	proto.Race_RaceGnome: {
		stats.Agility:   3,
		stats.Strength:  -5,
		stats.Intellect: 3,
		stats.Spirit:    0,
		stats.Stamina:   -1,
	},
	proto.Race_RaceTroll: {
		stats.Agility:   2,
		stats.Strength:  1,
		stats.Intellect: -4,
		stats.Spirit:    1,
		stats.Stamina:   1,
	},
}

var ClassBaseCrit = map[proto.Class]stats.Stats{
	proto.Class_ClassUnknown: {},
	proto.Class_ClassWarrior: {
		stats.SpellCrit: 0.0000 * CritRatingPerCritChance,
		stats.MeleeCrit: 0.0000 * CritRatingPerCritChance,
	},
	proto.Class_ClassPaladin: {
		stats.SpellCrit: 3.5000 * CritRatingPerCritChance,
		stats.MeleeCrit: 0.7000 * CritRatingPerCritChance,
	},
	proto.Class_ClassHunter: {
		stats.SpellCrit: 3.6000 * CritRatingPerCritChance,
		stats.MeleeCrit: 0.0000 * CritRatingPerCritChance,
	},
	proto.Class_ClassRogue: {
		stats.SpellCrit: 0.0000 * CritRatingPerCritChance,
		stats.MeleeCrit: 0.0000 * CritRatingPerCritChance,
	},
	proto.Class_ClassPriest: {
		stats.SpellCrit: 0.8000 * CritRatingPerCritChance,
		stats.MeleeCrit: 3.0000 * CritRatingPerCritChance,
	},
	proto.Class_ClassShaman: {
		stats.SpellCrit: 2.3000 * CritRatingPerCritChance,
		stats.MeleeCrit: 1.7000 * CritRatingPerCritChance,
	},
	proto.Class_ClassMage: {
		stats.SpellCrit: 0.2000 * CritRatingPerCritChance,
		stats.MeleeCrit: 3.2000 * CritRatingPerCritChance,
	},
	proto.Class_ClassWarlock: {
		stats.SpellCrit: 1.7000 * CritRatingPerCritChance,
		stats.MeleeCrit: 2.0000 * CritRatingPerCritChance,
	},
	proto.Class_ClassDruid: {
		stats.SpellCrit: 1.8000 * CritRatingPerCritChance,
		stats.MeleeCrit: 0.9000 * CritRatingPerCritChance,
	},
}

// Melee/Ranged crit agi scaling
// TODO: Level 40 and 50 values!
var CritPerAgiAtLevel = map[proto.Class]map[int]float64{
	proto.Class_ClassUnknown: {25: 0.0, 45: 0.0, 50: 0.0, 60: 0.0},
	proto.Class_ClassWarrior: {25: 0.1111, 40: 0.0755, 50: 0.0604, 60: 0.0500},
	proto.Class_ClassPaladin: {25: 0.1075, 40: 0.0753, 50: 0.0618, 60: 0.0506},
	proto.Class_ClassHunter:  {25: 0.0515, 40: 0.0312, 50: 0.0239, 60: 0.0189},
	proto.Class_ClassRogue:   {25: 0.0952, 40: 0.0572, 50: 0.0440, 60: 0.0345},
	proto.Class_ClassPriest:  {25: 0.0769, 40: 0.0640, 50: 0.0565, 60: 0.0500},
	proto.Class_ClassShaman:  {25: 0.0971, 40: 0.0722, 50: 0.0602, 60: 0.0508},
	proto.Class_ClassMage:    {25: 0.0720, 40: 0.0623, 50: 0.0566, 60: 0.0514},
	proto.Class_ClassWarlock: {25: 0.0909, 40: 0.0639, 50: 0.0551, 60: 0.0500},
	proto.Class_ClassDruid:   {25: 0.1025, 40: 0.0730, 50: 0.0599, 60: 0.0500},
}

// Spell crit int scaling
// TODO: Level 40 and 50 values!
var CritPerIntAtLevel = map[proto.Class]map[int]float64{
	proto.Class_ClassUnknown: {25: 0.0, 45: 0.0, 50: 0.0, 60: 0.0},
	proto.Class_ClassWarrior: {25: 0.0, 40: 0.0, 50: 0.0, 60: 0.0},
	proto.Class_ClassPaladin: {25: 0.0357, 40: 0.0250, 50: 0.0203, 60: 0.0167},
	proto.Class_ClassHunter:  {25: 0.0350, 40: 0.0246, 50: 0.0200, 60: 0.0165},
	proto.Class_ClassRogue:   {25: 0.0, 40: 0.0, 50: 0.0, 60: 0.0},
	proto.Class_ClassPriest:  {25: 0.0457, 40: 0.0277, 50: 0.0212, 60: 0.0168},
	proto.Class_ClassShaman:  {25: 0.0422, 40: 0.0269, 50: 0.0210, 60: 0.0169},
	proto.Class_ClassMage:    {25: 0.0475, 40: 0.0283, 50: 0.0214, 60: 0.0168},
	proto.Class_ClassWarlock: {25: 0.0429, 40: 0.0267, 50: 0.0207, 60: 0.0165},
	proto.Class_ClassDruid:   {25: 0.0427, 40: 0.0268, 50: 0.0208, 60: 0.0167},
}

var ClassBaseStats = map[proto.Class]map[int]stats.Stats{
	proto.Class_ClassUnknown: {},
	proto.Class_ClassWarrior: {
		25: {
			stats.Health:      274,
			stats.Mana:        0,
			stats.Agility:     39,
			stats.Strength:    54,
			stats.Intellect:   23,
			stats.Spirit:      28,
			stats.Stamina:     50,
			stats.AttackPower: 25*3 - 20,
		},
		40: {
			stats.Health:      649,
			stats.Mana:        0,
			stats.Agility:     54,
			stats.Strength:    79,
			stats.Intellect:   26,
			stats.Spirit:      34,
			stats.Stamina:     72,
			stats.AttackPower: 40*3 - 20,
		},
		50: {
			stats.Health:      1079,
			stats.Mana:        0,
			stats.Agility:     66,
			stats.Strength:    98,
			stats.Intellect:   28,
			stats.Spirit:      39,
			stats.Stamina:     90,
			stats.AttackPower: 50*3 - 20,
		},
		60: {
			stats.Health:      1689,
			stats.Mana:        0,
			stats.Agility:     80,
			stats.Strength:    120,
			stats.Intellect:   30,
			stats.Spirit:      45,
			stats.Stamina:     110,
			stats.AttackPower: 60*3 - 20,
		},
	},
	proto.Class_ClassPaladin: {
		25: {
			stats.Health:      266,
			stats.Mana:        552,
			stats.Agility:     34,
			stats.Strength:    48,
			stats.Intellect:   36,
			stats.Spirit:      38,
			stats.Stamina:     47,
			stats.AttackPower: 25*3 - 20,
		},
		40: {
			stats.Health:      621,
			stats.Mana:        987,
			stats.Agility:     46,
			stats.Strength:    70,
			stats.Intellect:   49,
			stats.Spirit:      52,
			stats.Stamina:     67,
			stats.AttackPower: 40*3 - 20,
		},
		50: {
			stats.Health:      966,
			stats.Mana:        1257,
			stats.Agility:     55,
			stats.Strength:    86,
			stats.Intellect:   59,
			stats.Spirit:      63,
			stats.Stamina:     82,
			stats.AttackPower: 50*3 - 20,
		},
		60: {
			stats.Health:      1381,
			stats.Mana:        1512,
			stats.Agility:     65,
			stats.Strength:    105,
			stats.Intellect:   70,
			stats.Spirit:      75,
			stats.Stamina:     100,
			stats.AttackPower: 60*3 - 20,
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
			stats.Health:      318,
			stats.Mana:        0,
			stats.Agility:     57,
			stats.Strength:    40,
			stats.Intellect:   25,
			stats.Spirit:      30,
			stats.Stamina:     38,
			stats.AttackPower: 25*2 - 20,
		},
		40: {
			stats.Health:      703,
			stats.Mana:        0,
			stats.Agility:     84,
			stats.Strength:    55,
			stats.Intellect:   29,
			stats.Spirit:      37,
			stats.Stamina:     52,
			stats.AttackPower: 40*2 - 20,
		},
		50: {
			stats.Health:      1068,
			stats.Mana:        0,
			stats.Agility:     106,
			stats.Strength:    67,
			stats.Intellect:   32,
			stats.Spirit:      43,
			stats.Stamina:     63,
			stats.AttackPower: 50*2 - 20,
		},
		60: {
			stats.Health:      1523,
			stats.Mana:        0,
			stats.Agility:     130,
			stats.Strength:    80,
			stats.Intellect:   35,
			stats.Spirit:      50,
			stats.Stamina:     75,
			stats.AttackPower: 60*2 - 20,
		},
	},
	proto.Class_ClassPriest: {
		25: {
			stats.Health:      302,
			stats.Mana:        497,
			stats.Agility:     26,
			stats.Strength:    25,
			stats.Intellect:   53,
			stats.Spirit:      55,
			stats.Stamina:     30,
			stats.AttackPower: -10,
		},
		40: {
			stats.Health:      637,
			stats.Mana:        911,
			stats.Agility:     31,
			stats.Strength:    29,
			stats.Intellect:   78,
			stats.Spirit:      81,
			stats.Stamina:     37,
			stats.AttackPower: -10,
		},
		50: {
			stats.Health:      972,
			stats.Mana:        1166,
			stats.Agility:     35,
			stats.Strength:    32,
			stats.Intellect:   98,
			stats.Spirit:      102,
			stats.Stamina:     43,
			stats.AttackPower: -10,
		},
		60: {
			stats.Health:      1397,
			stats.Mana:        1376,
			stats.Agility:     40,
			stats.Strength:    35,
			stats.Intellect:   120,
			stats.Spirit:      125,
			stats.Stamina:     50,
			stats.AttackPower: -10,
		},
	},
	proto.Class_ClassShaman: {
		25: {
			stats.Health:      257,
			stats.Mana:        505,
			stats.Agility:     31,
			stats.Strength:    41,
			stats.Intellect:   43,
			stats.Spirit:      47,
			stats.Stamina:     45,
			stats.AttackPower: 25*2 - 20,
		},
		40: {
			stats.Health:      610,
			stats.Mana:        975,
			stats.Agility:     40,
			stats.Strength:    58,
			stats.Intellect:   61,
			stats.Spirit:      67,
			stats.Stamina:     63,
			stats.AttackPower: 40*2 - 20,
		},
		50: {
			stats.Health:      947,
			stats.Mana:        1255,
			stats.Agility:     47,
			stats.Strength:    70,
			stats.Intellect:   74,
			stats.Spirit:      82,
			stats.Stamina:     78,
			stats.AttackPower: 50*2 - 20,
		},
		60: {
			stats.Health:      1423,
			stats.Mana:        1520,
			stats.Agility:     55,
			stats.Strength:    85,
			stats.Intellect:   90,
			stats.Spirit:      100,
			stats.Stamina:     95,
			stats.AttackPower: 60*2 - 20,
		},
	},
	proto.Class_ClassMage: {
		25: {
			stats.Health:      315,
			stats.Mana:        481,
			stats.Agility:     25,
			stats.Strength:    23,
			stats.Intellect:   55,
			stats.Spirit:      53,
			stats.Stamina:     28,
			stats.AttackPower: -10,
		},
		40: {
			stats.Health:      630,
			stats.Mana:        853,
			stats.Agility:     29,
			stats.Strength:    26,
			stats.Intellect:   81,
			stats.Spirit:      78,
			stats.Stamina:     34,
			stats.AttackPower: -10,
		},
		50: {
			stats.Health:      955,
			stats.Mana:        1048,
			stats.Agility:     32,
			stats.Strength:    28,
			stats.Intellect:   102,
			stats.Spirit:      98,
			stats.Stamina:     39,
			stats.AttackPower: -10,
		},
		60: {
			stats.Health:      1370,
			stats.Mana:        1213,
			stats.Agility:     35,
			stats.Strength:    30,
			stats.Intellect:   125,
			stats.Spirit:      120,
			stats.Stamina:     45,
			stats.AttackPower: -10,
		},
	},
	proto.Class_ClassWarlock: {
		25: {
			stats.Health:      279,
			stats.Mana:        498,
			stats.Agility:     30,
			stats.Strength:    28,
			stats.Intellect:   50,
			stats.Spirit:      52,
			stats.Stamina:     35,
			stats.AttackPower: -10,
		},
		40: {
			stats.Health:      634,
			stats.Mana:        923,
			stats.Agility:     37,
			stats.Strength:    34,
			stats.Intellect:   72,
			stats.Spirit:      75,
			stats.Stamina:     46,
			stats.AttackPower: -10,
		},
		50: {
			stats.Health:      979,
			stats.Mana:        1163,
			stats.Agility:     43,
			stats.Strength:    39,
			stats.Intellect:   90,
			stats.Spirit:      94,
			stats.Stamina:     55,
			stats.AttackPower: -10,
		},
		60: {
			stats.Health:      1414,
			stats.Mana:        1373,
			stats.Agility:     50,
			stats.Strength:    45,
			stats.Intellect:   110,
			stats.Spirit:      115,
			stats.Stamina:     65,
			stats.AttackPower: -10,
		},
	},
	proto.Class_ClassDruid: {
		25: {
			stats.Health:      318,
			stats.Mana:        479,
			stats.Agility:     33,
			stats.Strength:    35,
			stats.Intellect:   47,
			stats.Spirit:      50,
			stats.Stamina:     36,
			stats.AttackPower: -20,
		},
		40: {
			stats.Health:      683,
			stats.Mana:        854,
			stats.Agility:     43,
			stats.Strength:    46,
			stats.Intellect:   67,
			stats.Spirit:      72,
			stats.Stamina:     49,
			stats.AttackPower: -20,
		},
		50: {
			stats.Health:      1038,
			stats.Mana:        1064,
			stats.Agility:     51,
			stats.Strength:    55,
			stats.Intellect:   82,
			stats.Spirit:      90,
			stats.Stamina:     59,
			stats.AttackPower: -20,
		},
		60: {
			stats.Health:      1483,
			stats.Mana:        1244,
			stats.Agility:     60,
			stats.Strength:    65,
			stats.Intellect:   100,
			stats.Spirit:      110,
			stats.Stamina:     70,
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

	return starting.Add(RaceOffsets[r]).Add(ClassBaseCrit[c])
}
