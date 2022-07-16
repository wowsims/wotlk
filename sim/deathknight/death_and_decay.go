package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (deathKnight *DeathKnight) registerDeathAndDecaySpell() {
	var actionID = core.ActionID{SpellID: 49938}
	deathKnight.DeathAndDecayDot = core.NewDot(core.Dot{
		Aura: deathKnight.RegisterAura(core.Aura{
			Label:    "Death and Decay",
			ActionID: actionID,
		}),
		NumberOfTicks: 10,
		TickLength:    time.Second * 1,
		TickEffects: core.TickFuncAOESnapshot(deathKnight.Env, core.SpellEffect{
			ProcMask:        core.ProcMaskEmpty,
			BonusSpellPower: 0.0,

			DamageMultiplier: core.TernaryFloat64(deathKnight.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfDeathAndDecay), 1.2, 1.0),
			ThreatMultiplier: 1,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return (62.0 + deathKnight.applyImpurity(hitEffect, spell.Unit)*0.0475) *
						deathKnight.rageOfRivendareBonus() *
						deathKnight.tundraStalkerBonus()
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: deathKnight.OutcomeFuncMagicHitAndCrit(deathKnight.spellCritMultiplier(false)),
			IsPeriodic:     false,
		}),
	})

	deathKnight.DeathAndDecay = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    deathKnight.NewTimer(),
				Duration: time.Second*30 - time.Second*5*time.Duration(deathKnight.Talents.Morbidity),
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dkSpellCost := deathKnight.DetermineOptimalCost(sim, 1, 1, 1)
			deathKnight.Spend(sim, spell, dkSpellCost)

			amountOfRunicPower := 15.0
			deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())

			deathKnight.DeathAndDecayDot.Apply(sim)
		},
	})

	deathKnight.DeathAndDecayDot.Spell = deathKnight.DeathAndDecay
}

func (deathKnight *DeathKnight) CanDeathAndDecay(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 0.0, 1, 1, 1) && deathKnight.DeathAndDecay.IsReady(sim)
}
