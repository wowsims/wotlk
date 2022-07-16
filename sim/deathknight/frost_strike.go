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

	bloodOfTheNorthCoeff := 1.0
	if deathKnight.Talents.BloodOfTheNorth == 1 {
		bloodOfTheNorthCoeff = 1.03
	} else if deathKnight.Talents.BloodOfTheNorth == 2 {
		bloodOfTheNorthCoeff = 1.06
	} else if deathKnight.Talents.BloodOfTheNorth == 3 {
		bloodOfTheNorthCoeff = 1.1
	}

	glacierRotCoeff := 1.0
	if deathKnight.Talents.GlacierRot == 1 {
		glacierRotCoeff = 1.07
	} else if deathKnight.Talents.GlacierRot == 2 {
		glacierRotCoeff = 1.13
	} else if deathKnight.Talents.GlacierRot == 3 {
		glacierRotCoeff = 1.20
	}

	guileOfGorefiend := deathKnight.Talents.GuileOfGorefiend > 0

	effect := core.SpellEffect{
		BonusCritRating:  (1.0 * float64(deathKnight.Talents.Annihilation)) * core.CritRatingPerCritChance,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) *
					bloodOfTheNorthCoeff *
					core.TernaryFloat64(deathKnight.DiseasesAreActive() && deathKnight.Talents.GlacierRot > 0, glacierRotCoeff, 1.0) *
					core.TernaryFloat64(deathKnight.DiseasesAreActive(), 1.0+0.05*float64(deathKnight.Talents.TundraStalker), 1.0) *
					core.TernaryFloat64(deathKnight.BloodPlagueDisease.IsActive(), 1.0+0.02*float64(deathKnight.Talents.RageOfRivendare), 1.0) *
					core.TernaryFloat64(sim.IsExecutePhase35() && deathKnight.Talents.MercilessCombat > 0, 1.0+0.06*float64(deathKnight.Talents.MercilessCombat), 1.0)
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
		effect.OutcomeApplier = deathKnight.killingMachineOutcomeMod(deathKnight.OutcomeFuncMeleeSpecialHitAndCrit(deathKnight.critMultiplier(guileOfGorefiend)))
	} else {
		effect.ProcMask = core.ProcMaskMeleeOHSpecial
		effect.OutcomeApplier = deathKnight.killingMachineOutcomeMod(deathKnight.OutcomeFuncMeleeSpecialNoBlockDodgeParry(deathKnight.critMultiplier(guileOfGorefiend)))
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

	mhHitSpell := deathKnight.newFrostStrikeHitSpell(true)
	ohHitSpell := deathKnight.newFrostStrikeHitSpell(false)

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
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			ThreatMultiplier: 1,

			OutcomeApplier: deathKnight.OutcomeFuncAlwaysHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				deathKnight.threatOfThassarianProc(sim, spellEffect, mhHitSpell, ohHitSpell)
				deathKnight.threatOfThassarianAdjustMetrics(sim, spell, spellEffect, FrostStrikeMHOutcome)
			},
		}),
	})
}

func (deathKnight *DeathKnight) CanFrostStrike(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 40.0, 0, 0, 0) && deathKnight.FrostStrike.IsReady(sim)
}
