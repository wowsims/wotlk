package core

import (
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Could be in constants.go, but they won't be used anywhere else
// C values are divided by 100 so that we are working with 1% = 0.01
const Diminish_k_Druid = 0.972
const Diminish_k_Nondruid = 0.956
const Diminish_Cd_Druid = 116.890707 / 100
const Diminish_Cd_Nondruid = 88.129021 / 100
const Diminish_Cp = 47.003525 / 100
const Diminish_Cm = 16.000 / 100
const Diminish_kCd_Druid = (Diminish_k_Druid * Diminish_Cd_Druid)
const Diminish_kCd_Nondruid = (Diminish_k_Nondruid * Diminish_Cd_Nondruid)
const Diminish_kCp = (Diminish_k_Nondruid * Diminish_Cp)
const Diminish_kCm_Druid = (Diminish_k_Druid * Diminish_Cm)
const Diminish_kCm_Nondruid = (Diminish_k_Nondruid * Diminish_Cm)

// Diminishing Returns for tank avoidance
// Non-diminishing sources are added separately in spell outcome funcs

func (unit *Unit) GetDiminishedDodgeChance() float64 {

	// undiminished Dodge % = D
	// diminished Dodge % = (D * Cd)/((k*Cd) + D)

	dodgeChance :=
		unit.stats[stats.Dodge]/DodgeRatingPerDodgeChance/100 +
			unit.stats[stats.Defense]*DefenseRatingToChanceReduction

	if unit.PseudoStats.CanParry {
		return (dodgeChance * Diminish_Cd_Nondruid) / (Diminish_kCd_Nondruid + dodgeChance)
	} else {
		return (dodgeChance * Diminish_Cd_Druid) / (Diminish_kCd_Druid + dodgeChance)
	}
}

func (unit *Unit) GetDiminishedParryChance() float64 {

	// undiminished Parry % = P
	// diminished Parry % = (P * Cp)/((k*Cp) + P)

	parryChance :=
		unit.stats[stats.Parry]/ParryRatingPerParryChance/100 +
			unit.stats[stats.Defense]*DefenseRatingToChanceReduction

	return (parryChance * Diminish_Cp) / (Diminish_kCp + parryChance)

}

func (unit *Unit) GetDiminishedMissChance() float64 {

	// undiminished Miss % = M
	// diminished Miss % = (M * Cm)/((k*Cm) + M)

	missChance := unit.stats[stats.Defense] * DefenseRatingToChanceReduction

	if unit.PseudoStats.CanParry {
		return (missChance * Diminish_Cm) / (Diminish_kCm_Nondruid + missChance)
	} else {
		return (missChance * Diminish_Cm) / (Diminish_kCm_Druid + missChance)
	}
}
