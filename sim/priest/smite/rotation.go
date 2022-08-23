package smite

import (
	"github.com/wowsims/wotlk/sim/core"
)

// TODO: probably do something different instead of making it global?
const (
	mbidx int = iota
	swdidx
	vtidx
	swpidx
)

func (spriest *SmitePriest) OnGCDReady(sim *core.Simulation) {
	spriest.tryUseGCD(sim)
}

func (spriest *SmitePriest) tryUseGCD(sim *core.Simulation) {
	spell := spriest.chooseSpell(sim)

	if success := spell.Cast(sim, spriest.CurrentTarget); !success {
		spriest.WaitForMana(sim, spell.CurCast.Cost)
	}
}

func (spriest *SmitePriest) chooseSpell(sim *core.Simulation) *core.Spell {
	if spriest.holyFireDotWillBeUp(sim) {
		if spriest.InnerFocus != nil && spriest.InnerFocus.IsReady(sim) {
			spriest.InnerFocus.Cast(sim, nil)
		}

		// Make sure we spam smite while dot is active.
		return spriest.Smite
	} else if !spriest.ShadowWordPainDot.IsActive() {
		return spriest.ShadowWordPain
	} else if spriest.rotation.UseDevouringPlague && !spriest.DevouringPlagueDot.IsActive() {
		return spriest.DevouringPlague
	} else if spriest.Penance != nil && spriest.Penance.IsReady(sim) {
		return spriest.Penance
	} else if spriest.rotation.UseShadowWordDeath && spriest.ShadowWordDeath.IsReady(sim) {
		return spriest.ShadowWordDeath
	} else if spriest.rotation.UseMindBlast && spriest.MindBlast.IsReady(sim) {
		return spriest.MindBlast
	} else if spriest.HolyFire.IsReady(sim) {
		return spriest.HolyFire
	} else {
		return spriest.Smite
	}
}

// Returns whether a Smite cast starting now would complete while Holy Fire is active.
func (spriest *SmitePriest) holyFireDotWillBeUp(sim *core.Simulation) bool {
	if !spriest.HolyFireDot.IsActive() {
		return false
	}

	smiteCastTime := spriest.ApplyCastSpeedForSpell(spriest.Smite.DefaultCast.CastTime, spriest.Smite)
	return smiteCastTime <= spriest.HolyFireDot.RemainingDuration(sim)
}
