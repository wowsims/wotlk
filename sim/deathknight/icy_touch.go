package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) killingMachineOutcomeMod(outcomeApplier core.OutcomeApplier) core.OutcomeApplier {
	return func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect, attackTable *core.AttackTable) {
		if deathKnight.KillingMachineAura.IsActive() {
			deathKnight.AddStatDynamic(sim, stats.SpellCrit, 100*core.CritRatingPerCritChance)
			outcomeApplier(sim, spell, spellEffect, attackTable)
			deathKnight.AddStatDynamic(sim, stats.SpellCrit, -100*core.CritRatingPerCritChance)
		} else {
			outcomeApplier(sim, spell, spellEffect, attackTable)
		}
	}
}

func (deathKnight *DeathKnight) registerIcyTouchSpell() {
	target := deathKnight.CurrentTarget

	itAura := core.IcyTouchAura(target, deathKnight.Talents.ImprovedIcyTouch)
	deathKnight.IcyTouchAura = itAura

	glacierRotCoeff := 0.0
	if deathKnight.Talents.GlacierRot == 1 {
		glacierRotCoeff = 0.07
	} else if deathKnight.Talents.GlacierRot == 2 {
		glacierRotCoeff = 0.13
	} else if deathKnight.Talents.GlacierRot == 3 {
		glacierRotCoeff = 0.20
	}

	deathKnight.IcyTouch = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 59131},
		SpellSchool: core.SpellSchoolFrost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: 5.0 * float64(deathKnight.Talents.Rime) * core.CritRatingPerCritChance,
			DamageMultiplier:     1.0,
			ThreatMultiplier:     7.0,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					roll := (245.0-227.0)*sim.RandomFloat("Icy Touch") + 227.0
					return (roll + hitEffect.MeleeAttackPower(spell.Unit)*0.1) *
						(1.0 +
							0.05*float64(deathKnight.Talents.ImprovedIcyTouch) +
							core.TernaryFloat64(deathKnight.DiseasesAreActive() && deathKnight.Talents.GlacierRot > 0, glacierRotCoeff, 0.0) +
							core.TernaryFloat64(sim.IsExecutePhase35() && deathKnight.Talents.MercilessCombat > 0, 0.06*float64(deathKnight.Talents.MercilessCombat), 0.0))
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: deathKnight.killingMachineOutcomeMod(deathKnight.OutcomeFuncMagicHitAndCrit(deathKnight.spellCritMultiplier(false))),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					dkSpellCost := deathKnight.DetermineOptimalCost(sim, 0, 1, 0)
					deathKnight.Spend(sim, spell, dkSpellCost)

					deathKnight.FrostFeverDisease.Apply(sim)

					// TODO: Temporary application of ebon plague until dot auras
					// properly run their events to control ebon plague
					deathKnight.checkForEbonPlague(sim)

					amountOfRunicPower := 10.0 + 2.5*float64(deathKnight.Talents.ChillOfTheGrave)
					deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())

					deathKnight.IcyTouchAura.Activate(sim)

					if deathKnight.IcyTouchAura.IsActive() {
						deathKnight.IcyTalonsAura.Activate(sim)
					}
				}
			},
		}),
	})
}

func (deathKnight *DeathKnight) CanIcyTouch(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 0.0, 0, 1, 0) && deathKnight.IcyTouch.IsReady(sim)
}
