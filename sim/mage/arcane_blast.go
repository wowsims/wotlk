package mage

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

const ArcaneBlastBaseManaCost = 195.0
const ArcaneBlastBaseCastTime = time.Millisecond * 2500

func (mage *Mage) newArcaneBlastSpell(numStacks int32) *core.Spell {
	mage.ArcaneBlastAura = mage.GetOrRegisterAura(core.Aura{
		Label:     "Arcane Blast",
		ActionID:  core.ActionID{SpellID: 36032},
		Duration:  time.Second * 8,
		MaxStacks: 3,
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			// Reset the mana cost on expiration.
			for i := int32(0); i < 4; i++ {
				mage.ArcaneBlast[i].CurCast.Cost = core.MaxFloat(0, mage.ArcaneBlast[i].CurCast.Cost-3.0*ArcaneBlastBaseManaCost*0.75)
			}
		},
	})

	actionID := core.ActionID{SpellID: 30451, Tag: numStacks + 1}

	return mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolArcane,
		Flags:       SpellFlagMage,

		ResourceType: stats.Mana,
		BaseCost:     ArcaneBlastBaseManaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: ArcaneBlastBaseManaCost * (1 + 0.75*float64(numStacks) + core.TernaryFloat64(mage.hasTristfal, 0.2, 0)),

				GCD:      core.GCDDefault,
				CastTime: ArcaneBlastBaseCastTime - time.Duration(numStacks)*time.Second/3,
			},
			OnCastComplete: func(sim *core.Simulation, _ *core.Spell) {
				mage.ArcaneBlastAura.Activate(sim)
				mage.ArcaneBlastAura.AddStack(sim)
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellHitRating:  float64(mage.Talents.ArcaneFocus) * 2 * core.SpellHitRatingPerHitChance,
			BonusSpellCritRating: float64(mage.Talents.ArcaneImpact) * 2 * core.SpellCritRatingPerCritChance,

			DamageMultiplier: mage.spellDamageMultiplier * core.TernaryFloat64(mage.hasTristfal, 1.2, 1),
			ThreatMultiplier: 1 - 0.2*float64(mage.Talents.ArcaneSubtlety),

			BaseDamage:     core.BaseDamageConfigMagic(668, 772, 2.5/3.5),
			OutcomeApplier: mage.OutcomeFuncMagicHitAndCrit(mage.SpellCritMultiplier(1, 0.25*float64(mage.Talents.SpellPower))),
		}),
	})
}

func (mage *Mage) ArcaneBlastCastTime(numStacks int32) time.Duration {
	castTime := mage.ArcaneBlast[numStacks].DefaultCast.CastTime
	castTime = mage.ApplyCastSpeed(castTime)
	return castTime
}

func (mage *Mage) ArcaneBlastManaCost(numStacks int32) float64 {
	cost := mage.ArcaneBlast[numStacks].DefaultCast.Cost
	cost = mage.ArcaneBlast[numStacks].ApplyCostModifiers(cost)
	return cost
}

// Whether Arcane Blast stacks will fall off before a new blast could finish casting.
func (mage *Mage) willDropArcaneBlastStacks(sim *core.Simulation, castTime time.Duration, numStacks int32) bool {
	remainingBuffTime := mage.ArcaneBlastAura.RemainingDuration(sim)
	return numStacks == 0 || remainingBuffTime < castTime
}

// Determines whether we can spam arcane blast for the remainder of the encounter.
func (mage *Mage) canBlast(sim *core.Simulation, curManaCost float64, curCastTime time.Duration, numStacks int32, willDropStacks bool) bool {
	numStacksAfterFirstCast := numStacks + 1
	if willDropStacks {
		numStacksAfterFirstCast = 1
	}

	remainingDuration := sim.GetRemainingDuration()
	projectedRemainingMana := mage.manaTracker.ProjectedRemainingMana(sim, mage.GetCharacter())

	extraManaCost := 0.0
	if mage.hasTristfal {
		extraManaCost = 39
	}

	// First cast, which is curArcaneBlast
	projectedRemainingMana -= curManaCost
	remainingDuration -= curCastTime
	if projectedRemainingMana < 0 {
		return false
	} else if remainingDuration < 0 {
		return true
	}

	// Second cast
	if numStacksAfterFirstCast == 1 {
		projectedRemainingMana -= ArcaneBlastBaseManaCost + (1.0 * ArcaneBlastBaseManaCost * 0.75) + extraManaCost
		remainingDuration -= mage.ApplyCastSpeed(ArcaneBlastBaseCastTime - (1 * time.Second / 3))
		if projectedRemainingMana < 0 {
			return false
		} else if remainingDuration < 0 {
			return true
		}
	}

	// Third cast
	if numStacksAfterFirstCast < 3 {
		projectedRemainingMana -= ArcaneBlastBaseManaCost + (2.0 * ArcaneBlastBaseManaCost * 0.75) + extraManaCost
		remainingDuration -= mage.ApplyCastSpeed(ArcaneBlastBaseCastTime - (2 * time.Second / 3))
		if projectedRemainingMana < 0 {
			return false
		} else if remainingDuration < 0 {
			return true
		}
	}

	// Everything after this will be full stack blasts.
	manaCost := ArcaneBlastBaseManaCost + (3.0 * ArcaneBlastBaseManaCost * 0.75) + extraManaCost
	castTime := mage.ApplyCastSpeed(ArcaneBlastBaseCastTime - (3 * time.Second / 3))
	numCasts := remainingDuration / castTime // time.Duration is an integer so we don't need to call math.Floor()
	totalManaCost := manaCost * float64(numCasts)

	clearcastProcChance := 0.02 * float64(mage.Talents.ArcaneConcentration)
	estimatedClearcastProcs := int(float64(numCasts) * clearcastProcChance)
	totalManaCost -= manaCost * float64(estimatedClearcastProcs)

	return totalManaCost < projectedRemainingMana
}
