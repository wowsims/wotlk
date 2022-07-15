package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

func OutcomeEitherWeaponHitOrCrit(mhOutcome core.HitOutcome, ohOutcome core.HitOutcome) bool {
	return mhOutcome == core.OutcomeHit || mhOutcome == core.OutcomeCrit || ohOutcome == core.OutcomeHit || ohOutcome == core.OutcomeCrit
}

func ToTChance(deathKnight *DeathKnight) float64 {
	threatOfThassarianChance := 0.0
	if deathKnight.Talents.ThreatOfThassarian == 1 {
		threatOfThassarianChance = 0.30
	} else if deathKnight.Talents.ThreatOfThassarian == 2 {
		threatOfThassarianChance = 0.60
	} else if deathKnight.Talents.ThreatOfThassarian == 3 {
		threatOfThassarianChance = 1.0
	}
	return threatOfThassarianChance
}

func ToTWillCast(sim *core.Simulation, totChance float64) bool {
	ohWillCast := sim.RandomFloat("Threat of Thassarian") <= totChance
	return ohWillCast
}

func ToTAdjustMetrics(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect, mhOutcome core.HitOutcome) {
	spell.SpellMetrics[spellEffect.Target.TableIndex].Casts -= 1
	if mhOutcome == core.OutcomeHit {
		spell.SpellMetrics[spellEffect.Target.TableIndex].Hits -= 1
	} else if mhOutcome == core.OutcomeCrit {
		spell.SpellMetrics[spellEffect.Target.TableIndex].Hits -= 1
	} else {
		spell.SpellMetrics[spellEffect.Target.TableIndex].Hits -= 2
	}
}
