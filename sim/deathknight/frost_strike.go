package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var FrostStrikeActionID = core.ActionID{SpellID: 55268}
var FrostStrikeMHOutcome = core.OutcomeHit
var FrostStrikeOHOutcome = core.OutcomeHit

func (deathKnight *DeathKnight) newFrostStrikeHitSpell(isMH bool) *core.Spell {
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, false, 138.0, 0.55, true)
	if !isMH {
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 138.0, 0.55, true)
	}

	effect := core.SpellEffect{
		BonusCritRating:  (1.0 * float64(deathKnight.Talents.Annihilation)) * core.CritRatingPerCritChance,
		DamageMultiplier: deathKnight.bloodOfTheNorthCoeff(),
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) *
					deathKnight.glacielRotBonus() *
					deathKnight.rageOfRivendareBonus() *
					deathKnight.tundraStalkerBonus() *
					deathKnight.mercilessCombatBonus(sim)
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

	if isMH {
		effect.ProcMask = core.ProcMaskMeleeMHSpecial
		effect.OutcomeApplier = deathKnight.killingMachineOutcomeMod(deathKnight.OutcomeFuncMeleeSpecialHitAndCrit(deathKnight.critMultiplier()))
	} else {
		effect.ProcMask = core.ProcMaskMeleeOHSpecial
		effect.OutcomeApplier = deathKnight.killingMachineOutcomeMod(deathKnight.OutcomeFuncMeleeSpecialNoBlockDodgeParry(deathKnight.critMultiplier()))
	}

	return deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:     FrostStrikeActionID,
		SpellSchool:  core.SpellSchoolFrost,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (deathKnight *DeathKnight) registerFrostStrikeSpell() {
	baseCost := 40.0
	if deathKnight.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfFrostStrike) {
		baseCost -= 8.0
	}

	deathKnight.FrostStrikeMhHit = deathKnight.newFrostStrikeHitSpell(true)
	deathKnight.FrostStrikeOhHit = deathKnight.newFrostStrikeHitSpell(false)

	deathKnight.FrostStrike = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    FrostStrikeActionID,
		SpellSchool: core.SpellSchoolFrost,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: baseCost,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = deathKnight.getModifiedGCD()
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			ThreatMultiplier: 1,

			OutcomeApplier: deathKnight.OutcomeFuncAlwaysHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				deathKnight.threatOfThassarianProc(sim, spellEffect, deathKnight.FrostStrikeMhHit, deathKnight.FrostStrikeOhHit)
				deathKnight.threatOfThassarianAdjustMetrics(sim, spell, spellEffect, FrostStrikeMHOutcome)
			},
		}),
	})
}

func (deathKnight *DeathKnight) CanFrostStrike(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 40.0, 0, 0, 0) && deathKnight.FrostStrike.IsReady(sim)
}
