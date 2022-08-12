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
	if mage.Rotation.Type == proto.Mage_Rotation_Arcane {
		spell = mage.doArcaneRotation(sim)
	} else if mage.Rotation.Type == proto.Mage_Rotation_Fire {
		spell = mage.doFireRotation(sim)
	} else if mage.Rotation.Type == proto.Mage_Rotation_Frost {
		spell = mage.doFrostRotation(sim)
	} else {
		spell = mage.doAoeRotation(sim)
	}

	if success := spell.Cast(sim, mage.CurrentTarget); success {
		return
	} else {
		mage.Metrics.MarkOOM(&mage.Unit, sim.CurrentTime)
		mage.WaitForMana(sim, spell.CurCast.Cost)
	}
}

// 4 ABs used < x always fish for AM
// 4 ABs used > y always cast AM as soon as barrage procs
func (mage *Mage) doArcaneRotation(sim *core.Simulation) *core.Spell {
	numStacks := mage.ArcaneBlastAura.GetStacks()

	if sim.GetRemainingDuration() < 12*time.Second {
		mage.DisableMajorCooldown(core.ActionID{SpellID: EvocationId})
	}

	burstDuration := time.Duration(mage.Character.CurrentManaPercent()*40) * time.Second
	if sim.GetRemainingDuration() < burstDuration {
		mage.DisableMajorCooldown(core.ActionID{SpellID: EvocationId})
		if mage.Character.CurrentMana() < mage.ArcaneBlast.CurCast.Cost {
			return mage.ArcaneMissiles
		} else {
			return mage.ArcaneBlast
		}
	}

	if mage.Rotation.MinBlastBeforeMissiles > numStacks {
		if mage.isMissilesBarrageVisible && mage.Rotation.Num_4StackBlastsToEarlyMissiles < mage.num4CostAB {
			return mage.ArcaneMissiles
		} else {
			return mage.ArcaneBlast
		}
	} else {
		if mage.extraABsAP > 0 && mage.GetAura("Arcane Power").IsActive() {
			mage.extraABsAP--
			return mage.ArcaneBlast
		}

		if mage.isMissilesBarrageVisible || mage.Rotation.Num_4StackBlastsToMissilesGamble < mage.num4CostAB {
			return mage.ArcaneMissiles
		} else {
			return mage.ArcaneBlast
		}
	}
}

func (mage *Mage) doFireRotation(sim *core.Simulation) *core.Spell {
	if mage.Rotation.MaintainImprovedScorch && mage.ScorchAura != nil && (!mage.ScorchAura.IsActive() || mage.ScorchAura.RemainingDuration(sim) < time.Millisecond*4000) {
		return mage.Scorch
	}

	if mage.HotStreakAura.IsActive() {
		return mage.Pyroblast
	}

	if !mage.LivingBombNotActive.Empty() {
		return mage.LivingBomb
	}

	if mage.Rotation.PrimaryFireSpell == proto.Mage_Rotation_Fireball {
		return mage.Fireball
	} else {
		return mage.FrostfireBolt
	}
}

func (mage *Mage) doFrostRotation(sim *core.Simulation) *core.Spell {
	return mage.Frostbolt
}

func (mage *Mage) doAoeRotation(sim *core.Simulation) *core.Spell {
	if mage.Rotation.Aoe == proto.Mage_Rotation_ArcaneExplosion {
		return mage.ArcaneExplosion
	} else if mage.Rotation.Aoe == proto.Mage_Rotation_Flamestrike {
		return mage.Flamestrike
	} else {
		return mage.Blizzard
	}
}
