package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (mage *Mage) OnGCDReady(sim *core.Simulation) {
	mage.tryUseGCD(sim)
}

func (mage *Mage) tryUseGCD(sim *core.Simulation) {
	spell := mage.chooseSpell(sim)
	if success := spell.Cast(sim, mage.CurrentTarget); !success {
		mage.WaitForMana(sim, spell.CurCast.Cost)
	}
}

func (mage *Mage) chooseSpell(sim *core.Simulation) *core.Spell {
	if mage.Rotation.MaintainImprovedScorch && (!mage.ScorchAura.IsActive() || mage.ScorchAura.RemainingDuration(sim) < time.Millisecond*4000) {
		return mage.Scorch
	}

	if mage.Rotation.Type == proto.Mage_Rotation_Arcane {
		spell := mage.doArcaneRotation(sim)
		if spell == mage.ArcaneBlast {
			mage.arcaneBlastStreak++
		}
		return spell
	} else if mage.Rotation.Type == proto.Mage_Rotation_Fire {
		return mage.doFireRotation(sim)
	} else if mage.Rotation.Type == proto.Mage_Rotation_Frost {
		return mage.doFrostRotation(sim)
	} else {
		return mage.doAoeRotation(sim)
	}
}

func (mage *Mage) doArcaneRotation(sim *core.Simulation) *core.Spell {
	// AB until the end.
	if mage.canBlast(sim) {
		return mage.ArcaneBlast
	}

	// Extra ABs before first AP.
	if sim.CurrentTime < time.Second*10 && !mage.ArcanePowerAura.IsActive() && mage.arcanePowerMCD != nil && mage.arcanePowerMCD.TimeToNextCast(sim) < time.Second*5 {
		return mage.ArcaneBlast
	}

	// Extra ABs during first AP.
	if sim.CurrentTime < time.Second*60 && mage.ArcanePowerAura.IsActive() && mage.arcaneBlastStreak < mage.Rotation.ExtraBlastsDuringFirstAp+4 {
		return mage.ArcaneBlast
	}

	abStacks := mage.ArcaneBlastAura.GetStacks()
	hasMissileBarrage := mage.MissileBarrageAura.IsActive() && mage.MissileBarrageAura.TimeActive(sim) > mage.ReactionTime

	// AM if we have MB and below n AB stacks.
	if hasMissileBarrage && abStacks < mage.Rotation.MissileBarrageBelowArcaneBlastStacks {
		return mage.ArcaneMissiles
	}

	// AM if we have MB and below mana %.
	manaPercent := mage.CurrentManaPercent()
	if hasMissileBarrage && manaPercent < mage.Rotation.MissileBarrageBelowManaPercent {
		return mage.ArcaneMissiles
	}

	// AM if we don't have barrage and over mana %.
	if !hasMissileBarrage && manaPercent > mage.Rotation.BlastWithoutMissileBarrageAboveManaPercent {
		return mage.ArcaneBlast
	}

	// If we've reached max desired stacks, use AM / ABarr. Otherwise blast.
	maxAbStacks := int32(4)
	if manaPercent < mage.Rotation.Only_3ArcaneBlastStacksBelowManaPercent {
		maxAbStacks = 3
	}
	if abStacks < maxAbStacks {
		return mage.ArcaneBlast
	} else if mage.ArcaneBarrage != nil && mage.ArcaneBarrage.IsReady(sim) {
		return mage.ArcaneBarrage
	} else {
		return mage.ArcaneMissiles
	}
}

func (mage *Mage) canBlast(sim *core.Simulation) bool {
	// Save computation by assuming we can't blast for 30+ seconds.
	remainingDur := sim.GetRemainingDuration()
	if remainingDur > time.Second*30 {
		return false
	}

	castTime := mage.ApplyCastSpeed(ArcaneBlastBaseCastTime)
	manaCost := mage.ArcaneBlast.DefaultCast.Cost

	stacks := float64(mage.ArcaneBlastAura.GetStacks())
	curMana := mage.CurrentMana()
	for curTime := time.Duration(0); curTime <= remainingDur; curTime += castTime {
		if stacks < 4 {
			stacks++
		}
		curMana -= manaCost * 1.75 * stacks
		if curMana < 0 {
			return false
		}
	}
	return true
}

func (mage *Mage) doFireRotation(sim *core.Simulation) *core.Spell {
	noBomb := mage.LivingBomb != nil && !mage.LivingBombDot.IsActive()
	if noBomb && !mage.heatingUp {
		return mage.LivingBomb
	}

	hasHotStreak := mage.HotStreakAura.IsActive() && mage.HotStreakAura.TimeActive(sim) > mage.ReactionTime
	if hasHotStreak && mage.Pyroblast != nil {
		return mage.Pyroblast
	}

	if noBomb {
		return mage.LivingBomb
	}

	if mage.Rotation.PrimaryFireSpell == proto.Mage_Rotation_Fireball {
		return mage.Fireball
	} else if mage.Rotation.PrimaryFireSpell == proto.Mage_Rotation_FrostfireBolt {
		return mage.FrostfireBolt
	} else {
		return mage.Scorch
	}
}

func (mage *Mage) doFrostRotation(sim *core.Simulation) *core.Spell {
	hasBrainFreeze := mage.BrainFreezeAura.IsActive() && mage.BrainFreezeAura.TimeActive(sim) > mage.ReactionTime
	if mage.FingersOfFrostAura.IsActive() {
		if mage.DeepFreeze != nil && mage.DeepFreeze.IsReady(sim) {
			return mage.DeepFreeze
		} else if hasBrainFreeze {
			return mage.FrostfireBolt
		} else if mage.Rotation.UseIceLance {
			return mage.IceLance
		}
	}

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
