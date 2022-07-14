package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

// TODO: make this an AoE spell, idk how to so for now its single target
func (deathKnight *DeathKnight) registerHowlingBlastSpell() {
	if !deathKnight.Talents.HowlingBlast {
		return
	}
	//target := deathKnight.CurrentTarget

	glacierRotCoeff := 0.0
	if deathKnight.Talents.GlacierRot == 1 {
		glacierRotCoeff = 0.07
	} else if deathKnight.Talents.GlacierRot == 2 {
		glacierRotCoeff = 0.13
	} else if deathKnight.Talents.GlacierRot == 3 {
		glacierRotCoeff = 0.20
	}

	guileOfGorefiend := deathKnight.Talents.GuileOfGorefiend > 0

	deathKnight.HowlingBlast = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 51411},
		SpellSchool: core.SpellSchoolFrost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    deathKnight.NewTimer(),
				Duration: 8.0 * time.Second,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: 0.0,
			DamageMultiplier:     1.0,
			ThreatMultiplier:     1.0,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					roll := (562.0-518.0)*sim.RandomFloat("Howling Blast") + 518.0
					return (roll + hitEffect.MeleeAttackPower(spell.Unit)*0.1) *
						(1.0 +
							core.TernaryFloat64(deathKnight.DiseasesAreActive() && deathKnight.Talents.GlacierRot > 0, glacierRotCoeff, 0.0) +
							core.TernaryFloat64(deathKnight.DiseasesAreActive(), 0.05*float64(deathKnight.Talents.TundraStalker), 0.0) +
							core.TernaryFloat64(sim.IsExecutePhase35() && deathKnight.Talents.MercilessCombat > 0, 0.06*float64(deathKnight.Talents.MercilessCombat), 0.0))
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: deathKnight.killingMachineOutcomeMod(deathKnight.OutcomeFuncMagicHitAndCrit(deathKnight.spellCritMultiplier(guileOfGorefiend))),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					if !deathKnight.HowlingBlastCostless {
						dkSpellCost := deathKnight.DetermineOptimalCost(sim, 0, 1, 1)
						deathKnight.Spend(sim, spell, dkSpellCost)
					} else {
						deathKnight.HowlingBlastCostless = false
					}

					amountOfRunicPower := 15.0 + 2.5*float64(deathKnight.Talents.ChillOfTheGrave)
					deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
				}
			},
		}),
	})
}

func (deathKnight *DeathKnight) CanHowlingBlast(sim *core.Simulation) bool {
	if deathKnight.HowlingBlastCostless {
		return deathKnight.HowlingBlast.IsReady(sim)
	}
	return deathKnight.CastCostPossible(sim, 0.0, 0, 1, 1) && deathKnight.HowlingBlast.IsReady(sim)
}
