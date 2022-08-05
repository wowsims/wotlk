package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var FrostStrikeActionID = core.ActionID{SpellID: 55268}
var FrostStrikeMHOutcome = core.OutcomeMiss
var FrostStrikeOHOutcome = core.OutcomeMiss

func (dk *Deathknight) newFrostStrikeHitSpell(isMH bool) *RuneSpell {
	baseDamage := 250.0 + dk.sigilOfTheVengefulHeartFrostStrike()
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, baseDamage, 0.55, true)
	if !isMH {
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, true, baseDamage, 0.55*dk.nervesOfColdSteelBonus(), true)
	}

	effect := core.SpellEffect{
		BonusCritRating:  (dk.annihilationCritBonus() + dk.darkrunedBattlegearCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: dk.bloodOfTheNorthCoeff(),
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) *
					dk.glacielRotBonus(hitEffect.Target) *
					dk.RoRTSBonus(hitEffect.Target) *
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

	return dk.RegisterSpell(nil, core.SpellConfig{
		ActionID:     FrostStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool:  core.SpellSchoolFrost,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (dk *Deathknight) registerFrostStrikeSpell() {
	baseCost := float64(core.NewRuneCost(
		core.Ternary(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfFrostStrike), uint8(32), 40), 0, 0, 0, 0,
	))

	dk.FrostStrikeMhHit = dk.newFrostStrikeHitSpell(true)
	dk.FrostStrikeOhHit = dk.newFrostStrikeHitSpell(false)

	dk.FrostStrike = dk.RegisterSpell(nil, core.SpellConfig{
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
				dk.LastOutcome = spellEffect.Outcome

				// Check for KM after both hits have passed
				if dk.LastOutcome.Matches(core.OutcomeLanded) {
					if dk.KillingMachineAura.IsActive() {
						dk.KillingMachineAura.Deactivate(sim)
					}
				}
			},
		}),
	})
}

func (dk *Deathknight) CanFrostStrike(sim *core.Simulation) bool {
	runeCost := core.RuneCost(dk.FrostStrike.BaseCost)
	return dk.CastCostPossible(sim, float64(runeCost.RunicPower()), 0, 0, 0) && dk.FrostStrike.IsReady(sim)
}

func (dk *Deathknight) CastFrostStrike(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanFrostStrike(sim) {
		return dk.FrostStrike.Cast(sim, target)
	}
	return false
}
