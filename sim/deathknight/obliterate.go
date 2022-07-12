package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) registerObliterateSpell() {
	baseCost := 15.0
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 467.0, 0.8, true)

	deathKnight.Obliterate = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 51425},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:             core.ProcMaskMeleeMHSpecial,
			BonusSpellCritRating: 5.0*float64(deathKnight.Talents.Rime)*core.CritRatingPerCritChance + 3.0*float64(deathKnight.Talents.Subversion)*core.CritRatingPerCritChance,
			DamageMultiplier:     1,
			ThreatMultiplier:     1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return weaponBaseDamage(sim, hitEffect, spell) *
						(1.0 +
							core.TernaryFloat64(deathKnight.FrostFeverDisease.IsActive(), 0.125, 0.0) +
							core.TernaryFloat64(deathKnight.BloodPlagueDisease.IsActive(), 0.125, 0.0) +
							core.TernaryFloat64(sim.IsExecutePhase35() && deathKnight.Talents.MercilessCombat > 0, 0.06*float64(deathKnight.Talents.MercilessCombat), 0.0))
				},
				TargetSpellCoefficient: 1,
			},

			OutcomeApplier: deathKnight.OutcomeFuncMeleeSpecialHitAndCrit(deathKnight.DefaultMeleeCritMultiplier()),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					deathKnight.SpendUnholyRune(sim, spell.UnholyRuneMetrics())
					deathKnight.SpendFrostRune(sim, spell.FrostRuneMetrics())
					amountOfRunicPower := 10.0 + 2.5*float64(deathKnight.Talents.ChillOfTheGrave)
					deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
				}
			},
		}),
	})
}

func (deathKnight *DeathKnight) CanObliterate(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 15.0, 0, 1, 1, 0) && deathKnight.Obliterate.IsReady(sim)
}
