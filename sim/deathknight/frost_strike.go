package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var FrostStrikeActionID = core.ActionID{SpellID: 55268}
var FrostStrikeMHOutcome = core.OutcomeMiss
var FrostStrikeOHOutcome = core.OutcomeMiss

func (dk *Deathknight) newFrostStrikeHitSpell(isMH bool) *core.Spell {
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 250.0, 0.55, true)
	if !isMH {
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, true, 250.0, 0.55*dk.nervesOfColdSteelBonus(), true)
	}

	effect := core.SpellEffect{
		BonusCritRating:  (dk.annihilationCritBonus() + dk.darkrunedBattlegearCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: dk.bloodOfTheNorthCoeff(),
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) *
					dk.glacielRotBonus(hitEffect.Target) *
					dk.rageOfRivendareBonus(hitEffect.Target) *
					dk.tundraStalkerBonus(hitEffect.Target) *
					dk.mercilessCombatBonus(sim)
			},
			TargetSpellCoefficient: 1,
		},

		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if isMH {
				FrostStrikeMHOutcome = spellEffect.Outcome
			} else {
				FrostStrikeOHOutcome = spellEffect.Outcome
			}
		},
	}

	dk.threatOfThassarianProcMasks(isMH, &effect, true, false, dk.killingMachineOutcomeMod)

	return dk.RegisterSpell(core.SpellConfig{
		ActionID:     FrostStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool:  core.SpellSchoolFrost,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (dk *Deathknight) registerFrostStrikeSpell() {
	baseCost := 40.0
	if dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfFrostStrike) {
		baseCost -= 8.0
	}

	dk.FrostStrikeMhHit = dk.newFrostStrikeHitSpell(true)
	dk.FrostStrikeOhHit = dk.newFrostStrikeHitSpell(false)

	dk.FrostStrike = dk.RegisterSpell(core.SpellConfig{
		ActionID:    FrostStrikeActionID.WithTag(3),
		SpellSchool: core.SpellSchoolFrost,
		Flags:       core.SpellFlagNoMetrics | core.SpellFlagNoLogs,

		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: baseCost,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			ThreatMultiplier: 1,

			OutcomeApplier: dk.OutcomeFuncAlwaysHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				dk.threatOfThassarianProc(sim, spellEffect, dk.FrostStrikeMhHit, dk.FrostStrikeOhHit)
				dk.threatOfThassarianAdjustMetrics(sim, spell, spellEffect, FrostStrikeMHOutcome)
				dk.LastCastOutcome = FrostStrikeMHOutcome

				// Check for KM after both hits have passed
				if dk.LastCastOutcome.Matches(core.OutcomeLanded) {
					if dk.KillingMachineAura.IsActive() {
						dk.KillingMachineAura.Deactivate(sim)
					}
				}
			},
		}),
	})
}

func (dk *Deathknight) CanFrostStrike(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 40.0, 0, 0, 0) && dk.FrostStrike.IsReady(sim)
}

func (dk *Deathknight) CastFrostStrike(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanFrostStrike(sim) {
		dk.FrostStrike.Cast(sim, target)
		return true
	}
	return false
}
