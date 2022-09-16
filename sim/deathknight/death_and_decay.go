package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) OutcomeDeathAndDecaySpecial() core.OutcomeApplier {
	return func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect, attackTable *core.AttackTable) {
		if spellEffect.MagicHitCheck(sim, spell, attackTable) {
			if sim.RandomFloat("Fixed Crit Roll") < dk.dndCritSnapshot {
				spellEffect.Outcome = core.OutcomeCrit
				spell.SpellMetrics[spellEffect.Target.UnitIndex].Crits++
				spellEffect.Damage *= dk.DefaultMeleeCritMultiplier()
			} else {
				spellEffect.Outcome = core.OutcomeHit
				spell.SpellMetrics[spellEffect.Target.UnitIndex].Hits++
			}
		} else {
			spellEffect.Outcome = core.OutcomeMiss
			spell.SpellMetrics[spellEffect.Target.UnitIndex].Misses++
			spellEffect.Damage = 0
		}
	}
}

func (dk *Deathknight) registerDeathAndDecaySpell() {

	actionID := core.ActionID{SpellID: 49938}
	glyphBonus := core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDeathAndDecay), 1.2, 1.0)

	doSnapshot := false
	baseCost := float64(core.NewRuneCost(15, 1, 1, 1, 0))
	dk.DeathAndDecay = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolShadow,
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
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Second*30 - time.Second*5*time.Duration(dk.Talents.Morbidity),
			},
		},

		ThreatMultiplier: 1.9,

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			doSnapshot = true
			dk.dndApSnapshot = 0.0
			dk.dndCritSnapshot = 0.0
			dk.DeathAndDecayDot.Apply(sim)
			dk.DeathAndDecayDot.TickOnce()
		},
	}, func(sim *core.Simulation) bool {
		return dk.CastCostPossible(sim, 0.0, 1, 1, 1) && dk.DeathAndDecay.IsReady(sim)
	}, nil)

	dk.DeathAndDecayDot = core.NewDot(core.Dot{
		Aura: dk.RegisterAura(core.Aura{
			Label:    "Death and Decay",
			ActionID: actionID,
		}),
		NumberOfTicks: 10,
		TickLength:    time.Second * 1,
		TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncAOEDamage(dk.Env, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: glyphBonus * dk.scourgelordsPlateDamageBonus(),
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					if doSnapshot {
						dk.dndCritSnapshot = hitEffect.SpellCritChance(spell.Unit, spell)
						dk.dndApSnapshot = 62.0 + dk.getImpurityBonus(hitEffect, spell.Unit)*0.0475
						doSnapshot = false
					}
					return dk.dndApSnapshot * dk.RoRTSBonus(hitEffect.Target)
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: dk.OutcomeDeathAndDecaySpecial(),
			IsPeriodic:     false,
		})),
	})
	dk.DeathAndDecayDot.Spell = dk.DeathAndDecay.Spell
}
