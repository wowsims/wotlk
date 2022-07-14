package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) registerDeathCoilSpell() {
	deathKnight.DeathCoil = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49895},
		SpellSchool: core.SpellSchoolShadow,

		ResourceType: stats.RunicPower,
		BaseCost:     40.0,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: 40.0,
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
						(1.0 +
							0.05*float64(deathKnight.Talents.Morbidity) +
							core.TernaryFloat64(deathKnight.BloodPlagueDisease.IsActive(), 0.02*float64(deathKnight.Talents.RageOfRivendare), 0.0))
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: deathKnight.OutcomeFuncMagicHitAndCrit(deathKnight.spellCritMultiplier(false)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
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
