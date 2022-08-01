package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (mage *Mage) OnGCDReady(sim *core.Simulation) {
	mage.tryUseGCD(sim)

	if mage.GCD.IsReady(sim) && (!mage.IsWaiting() && !mage.IsWaitingForMana()) {
		panic("failed to use our gcd")
	}
}

func (mage *Mage) tryUseGCD(sim *core.Simulation) {
	var spell *core.Spell
	if mage.RotationType == proto.Mage_Rotation_Arcane {
		spell = mage.doArcaneRotation(sim)
	} else if mage.RotationType == proto.Mage_Rotation_Fire {
		spell = mage.doFireRotation(sim)
	} else {
		spell = mage.doFrostRotation(sim)
	}

	if success := spell.Cast(sim, mage.CurrentTarget); success {
		return
	} else {
		mage.Metrics.MarkOOM(&mage.Unit, sim.CurrentTime)
		mage.WaitForMana(sim, spell.CurCast.Cost)
	}
}

func (mage *Mage) doArcaneRotation(sim *core.Simulation) *core.Spell {
	if mage.UseAoeRotation {
		return mage.doAoeRotation(sim)
	}

	//Going oom should be visible to person simming, but standard rotation should not oom with reasonable evocate/gem usage
	numStacks := mage.ArcaneBlastAura.GetStacks()
	if mage.ArcaneRotation.MinBlastBeforeMissiles > numStacks || !mage.MissileBarrageAura.IsActive() {
		return mage.ArcaneBlast
	} else {
		return mage.ArcaneMissiles
	}
}

func (mage *Mage) doFireRotation(sim *core.Simulation) *core.Spell {
	if mage.FireRotation.MaintainImprovedScorch && mage.ScorchAura != nil && (mage.ScorchAura.GetStacks() < 5 || mage.ScorchAura.RemainingDuration(sim) < time.Millisecond*5500) {
		return mage.Scorch
	}

	if mage.UseAoeRotation {
		return mage.doAoeRotation(sim)
	}

	if mage.FireRotation.WeaveFireBlast && mage.FireBlast.IsReady(sim) {
		return mage.FireBlast
	}

	if mage.FireRotation.PrimarySpell == proto.Mage_Rotation_FireRotation_Fireball {
		return mage.Fireball
	} else {
		return mage.Scorch
	}
}

func (mage *Mage) doFrostRotation(sim *core.Simulation) *core.Spell {
	if mage.UseAoeRotation {
		return mage.doAoeRotation(sim)
	}

	return mage.Frostbolt
}

func (mage *Mage) doAoeRotation(sim *core.Simulation) *core.Spell {
	if mage.AoeRotation.Rotation == proto.Mage_Rotation_AoeRotation_ArcaneExplosion {
		return mage.ArcaneExplosion
	} else if mage.AoeRotation.Rotation == proto.Mage_Rotation_AoeRotation_Flamestrike {
		return mage.Flamestrike
	} else {
		return mage.Blizzard
	}
}
