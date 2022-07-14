package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (deathKnight *DeathKnight) registerDeathCoilSpell() {
	deathKnight.DeathCoil = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49895},
		SpellSchool: core.SpellSchoolShadow,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: 0.0,
			DamageMultiplier:     1.0,
			ThreatMultiplier:     1.0,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return (443.0 + hitEffect.MeleeAttackPower(spell.Unit)*0.15) *
						(1.0 + 0.05*float64(deathKnight.Talents.Morbidity))
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: deathKnight.OutcomeFuncMagicHitAndCrit(deathKnight.spellCritMultiplier(false)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				deathKnight.SpendRunicPower(sim, 40.0, spell.RunicPowerMetrics())
				if spellEffect.Landed() && deathKnight.Talents.UnholyBlight {
					deathKnight.UnholyBlight.Apply(sim)
				}
			},
		}),
	})
}

func (deathKnight *DeathKnight) CanDeathCoil(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 40.0, 0, 0, 0) && deathKnight.DeathCoil.IsReady(sim)
}
