package smite

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (spriest *SmitePriest) OnGCDReady(sim *core.Simulation) {
	spriest.tryUseGCD(sim)
}

func (spriest *SmitePriest) tryUseGCD(sim *core.Simulation) {
	var spell *core.Spell

	if spriest.rotation.MemeDream {
		spell = spriest.chooseSpellMemeDream(sim)
	} else {
		spell = spriest.chooseSpell(sim)
	}

	if spell == nil {
		// nil means wait for HF.
		spriest.WaitUntil(sim, spriest.HolyFire.ReadyAt())
		return
	}

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
	} else if spriest.rotation.UseDevouringPlague && !spriest.DevouringPlagueDot.IsActive() {
		return spriest.DevouringPlague
	} else if !spriest.ShadowWordPainDot.IsActive() {
		return spriest.ShadowWordPain
	} else if spriest.HolyFire.IsReady(sim) {
		return spriest.HolyFire
	} else if spriest.HolyFire.TimeToReady(sim) <= spriest.allowedHFDelay {
		return nil
	} else if spriest.ImprovedSpiritTap.IsActive() {
		return spriest.Smite
	} else if spriest.Penance != nil && spriest.Penance.IsReady(sim) {
		return spriest.Penance
	} else if spriest.rotation.UseShadowWordDeath && spriest.ShadowWordDeath.IsReady(sim) {
		return spriest.ShadowWordDeath
	} else if spriest.rotation.UseMindBlast && spriest.MindBlast.IsReady(sim) {
		return spriest.MindBlast
	} else if spriest.Talents.MindFlay {
		mfTickLength := spriest.MindFlayTickDuration()
		hfTimeToReady := spriest.HolyFire.TimeToReady(sim)
		numTicks := core.MinInt(3, int(hfTimeToReady/mfTickLength+1))
		return spriest.MindFlay[numTicks]
	} else {
		return spriest.Smite
	}
}

func (spriest *SmitePriest) chooseSpellMemeDream(sim *core.Simulation) *core.Spell {
	if spriest.rotation.UseDevouringPlague && !spriest.DevouringPlagueDot.IsActive() {
		return spriest.DevouringPlague
	} else if !spriest.ShadowWordPainDot.IsActive() {
		return spriest.ShadowWordPain
	} else if spriest.HolyFire.IsReady(sim) {
		return spriest.HolyFire
	} else if spriest.HolyFire.TimeToReady(sim) <= spriest.allowedHFDelay {
		return nil
	} else {
		if spriest.InnerFocus != nil && spriest.InnerFocus.IsReady(sim) {
			spriest.InnerFocus.Cast(sim, nil)
		}

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
