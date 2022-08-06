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

func (mage *Mage) doArcaneRotation(sim *core.Simulation) *core.Spell {
	//Going oom should be visible to person simming, but standard rotation should not oom with reasonable evocate/gem usage
	numStacks := mage.ArcaneBlastAura.GetStacks()
	if mage.Rotation.MinBlastBeforeMissiles > numStacks || !mage.MissileBarrageAura.IsActive() {
		return mage.ArcaneBlast
	} else {
		return mage.ArcaneMissiles
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
