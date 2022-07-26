package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (dk *Deathknight) registerDeathAndDecaySpell() {

	actionID := core.ActionID{SpellID: 49938}
	glyphBonus := core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDeathAndDecay), 1.2, 1.0)

	dk.DeathAndDecayDot = core.NewDot(core.Dot{
		Aura: dk.RegisterAura(core.Aura{
			Label:    "Death and Decay",
			ActionID: actionID,
		}),
		NumberOfTicks: 10,
		TickLength:    time.Second * 1,
		TickEffects: core.TickFuncAOESnapshot(dk.Env, core.SpellEffect{
			ProcMask:        core.ProcMaskEmpty,
			BonusSpellPower: 0.0,

			DamageMultiplier: glyphBonus * dk.scourgelordsPlateDamageBonus(),
			ThreatMultiplier: 1,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return (62.0 + dk.applyImpurity(hitEffect, spell.Unit)*0.0475) *
						dk.rageOfRivendareBonus(hitEffect.Target) *
						dk.tundraStalkerBonus(hitEffect.Target)
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: dk.OutcomeFuncMagicHitAndCrit(dk.spellCritMultiplier()),
			IsPeriodic:     false,
		}),
	})

	dk.DeathAndDecay = dk.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Second*30 - time.Second*5*time.Duration(dk.Talents.Morbidity),
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dk.DeathAndDecayDot.Apply(sim)
			dk.DeathAndDecayDot.TickOnce()

			dkSpellCost := dk.DetermineCost(sim, core.DKCastEnum_BFU)
			dk.Spend(sim, spell, dkSpellCost)

			amountOfRunicPower := 15.0
			dk.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
		},
	})

	dk.DeathAndDecayDot.Spell = dk.DeathAndDecay
}

func (dk *Deathknight) CanDeathAndDecay(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 0.0, 1, 1, 1) && dk.DeathAndDecay.IsReady(sim)
}

func (dk *Deathknight) CastDeathAndDecay(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanDeathAndDecay(sim) {
		dk.DeathAndDecay.Cast(sim, target)
		return true
	}
	return false
}
