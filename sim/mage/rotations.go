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

//set options for total mana available. Can look at duration and current time to determine how many 4 cost ABs you can still use
//expected mana consumption is cost of 4ab am rotation + estimated number of 4 cost ABs over a fight length
//Ideally you'd after each cast update how many 4 cost ABs you still can use

//M is how much mana a fight will need on average. TM is total mana available over a fight. TM - M is average extra mana available
// TM - M / 4AB Cost is number of 4ABs that you can use on average or E. You could just decrement E every time you use a 4 cost AB
// if E goes negative then you could switch to a conserve rotation, but that's kind of worse than reassessing mana

// The true standard would be recalculate TM - M after every cast. TM is pretty simple just subtract mana you use from the cast
// dM is how much mana the rest of the fight will need. This should be pretty much constant as M is just a defined function.
// M = std rotation cost * duration + avg extra 4 cost ABs over duration. Depends heavily on how hard calculating the cost of a bad chain
// is over a given duration. Technically AP slightly increases your rotation cost as does haste. Haste adds duration to the fight
// but it's annoying as fuck to roll out that duration

// After every cast you could take the spellcast speed to the real cast duration. Then you could have a total duration increased by haste
// and subtract out the real cast duration every step.

// Could compare MPS vs MCS of steady state. MCS should be relatively stable over any fight length on average. It might be higher than steady state
// cost because of bad chains, but those will average out and if simmers could see typical 4 cost ABs available

// On average your MCS should be greater than MPS so it's just a matter of fight duration to hit oom. Thus you can guesstimate the remaining mana at different fight durations
// If you have the estimate of num 4ABs you can just subtract 1 each time you cast a 4AB. You could make settable conditions like
// if 4 AB remaining count > x always fish for AM
// if 4 AB remaining count < y always cast AM as soon as barrage procs
// clearcasting could add a little bit or just be ignored as it's a smaller gain

// f(M CDs, length) => num 4 cost ABs

// 4 ABs used < x always fish for AM
// 4 ABs used > y always cast AM as soon as barrage procs

// So longer fight durations just mean less 4 cost ABs. Can we just make a metric directly off fight durations like

// simpler option. Calculate total mana available give threshold options for Blind AM, ASAP Missile Barrage, and

// also smart evocate and smart last 10 second burn
func (mage *Mage) doArcaneRotation(sim *core.Simulation) *core.Spell {
	// sim.Duration
	//Going oom should be visible to person simming, but standard rotation should not oom with reasonable evocate/gem usage
	numStacks := mage.ArcaneBlastAura.GetStacks()

	if sim.GetRemainingDuration() < 10*time.Second {
		if mage.manaTracker.ProjectedRemainingMana(sim, &mage.Character) > mage.manaTracker.ProjectedManaCost(sim, &mage.Character) {
			if mage.Character.CurrentMana() < mage.ArcaneBlast.CurCast.Cost {
				return mage.ArcaneMissiles
			} else {
				return mage.ArcaneBlast
			}
		} else {
			return mage.ArcaneMissiles
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
