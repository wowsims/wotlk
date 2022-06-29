package druid

import (
	"strconv"
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (druid *Druid) registerRakeSpell() {
	actionID := core.ActionID{SpellID: 27003}

	cost := 40.0 - float64(druid.Talents.Ferocity)
	refundAmount := cost * 0.8

	mangleAura := core.MangleAura(druid.CurrentTarget)

	druid.Rake = druid.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Energy,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  time.Second,
			},
			ModifyCast:  druid.ApplyClearcasting,
			IgnoreHaste: true,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					damage := 78 + 0.01*hitEffect.MeleeAttackPower(spell.Unit)
					if mangleAura.IsActive() {
						return damage * 1.3
					} else {
						return damage
					}
				},
				TargetSpellCoefficient: 0,
			},
			OutcomeApplier: druid.OutcomeFuncMeleeSpecialHitAndCrit(druid.MeleeCritMultiplier()),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					if druid.RakeDot.IsActive() {
						druid.RakeDot.Refresh(sim)
						druid.RakeDot.TakeSnapshot(sim)
					} else {
						druid.RakeDot.Apply(sim)
					}
				} else {
					druid.AddEnergy(sim, refundAmount, druid.EnergyRefundMetrics)
				}
			},
		}),
	})

	dotAura := druid.CurrentTarget.RegisterAura(core.Aura{
		Label:    "Rake-" + strconv.Itoa(int(druid.Index)),
		ActionID: actionID,
		Duration: time.Second * 9,
	})
	druid.RakeDot = core.NewDot(core.Dot{
		Spell:         druid.Rake,
		Aura:          dotAura,
		NumberOfTicks: 3,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			IsPeriodic:       true,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return 36 + 0.02*hitEffect.MeleeAttackPower(spell.Unit)
				},
				TargetSpellCoefficient: 0,
			},
			OutcomeApplier: druid.OutcomeFuncTick(),
		})),
	})
}

func (druid *Druid) CanRake(sim *core.Simulation) bool {
	return druid.CurrentEnergy() >= druid.Rake.DefaultCast.Cost
}
