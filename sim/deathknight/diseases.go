package deathknight

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (deathKnight *DeathKnight) registerDiseaseDots() {
	deathKnight.registerFrostFever()
	deathKnight.registerBloodPlague()
}

func (deathKnight *DeathKnight) registerFrostFever() {
	actionID := core.ActionID{SpellID: 55095}
	target := deathKnight.CurrentTarget

	deathKnight.FrostFever = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFrost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: 0.0,
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			OutcomeApplier:   deathKnight.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					deathKnight.FrostFeverDisease.Apply(sim)
				}
			},
		}),
	})

	deathKnight.FrostFeverDisease = core.NewDot(core.Dot{
		Spell: deathKnight.FrostFever,
		Aura: target.RegisterAura(core.Aura{
			Label:    "FrostFever-" + strconv.Itoa(int(deathKnight.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 5 + int(deathKnight.Talents.Epidemic),
		TickLength:    time.Second * 3,

		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			IsPeriodic:       true,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return (127.0 + 80.0*0.32) + hitEffect.MeleeAttackPower(spell.Unit)*0.055
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: deathKnight.OutcomeFuncTick(),
		}),
	})
}

func (deathKnight *DeathKnight) registerBloodPlague() {
	actionID := core.ActionID{SpellID: 55078}
	target := deathKnight.CurrentTarget

	deathKnight.BloodPlague = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: 0.0,
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			OutcomeApplier:   deathKnight.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					deathKnight.BloodPlagueDisease.Apply(sim)
				}
			},
		}),
	})

	deathKnight.BloodPlagueDisease = core.NewDot(core.Dot{
		Spell: deathKnight.BloodPlague,
		Aura: target.RegisterAura(core.Aura{
			Label:    "BloodPlague-" + strconv.Itoa(int(deathKnight.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 5 + int(deathKnight.Talents.Epidemic),
		TickLength:    time.Second * 3,

		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			IsPeriodic:       true,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return (127.0 + 80.0*0.32) + hitEffect.MeleeAttackPower(spell.Unit)*0.055
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: deathKnight.OutcomeFuncTick(),
		}),
	})
}
