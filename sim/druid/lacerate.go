package druid

import (
	"strconv"
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (druid *Druid) registerLacerateSpell() {
	actionID := core.ActionID{SpellID: 33745}

	cost := 15.0 - float64(druid.Talents.ShreddingAttacks)
	refundAmount := cost * 0.8

	tickDamage := 155.0 / 5
	if ItemSetNordrassilHarness.CharacterHasSetBonus(&druid.Character, 4) {
		tickDamage += 15
	}
	if druid.Equip[items.ItemSlotRanged].ID == 27744 { // Idol of Ursoc
		tickDamage += 8
	}

	mangleAura := core.MangleAura(druid.CurrentTarget)

	druid.Lacerate = druid.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			ModifyCast:  druid.ApplyClearcasting,
			IgnoreHaste: true,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			DamageMultiplier: 1,
			ThreatMultiplier: 0.5,
			FlatThreatBonus:  267,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					damage := tickDamage + 0.01*hitEffect.MeleeAttackPower(spell.Unit)
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
					if druid.LacerateDot.IsActive() {
						druid.LacerateDot.Refresh(sim)
						druid.LacerateDot.AddStack(sim)
						druid.LacerateDot.TakeSnapshot(sim)
					} else {
						druid.LacerateDot.Apply(sim)
						druid.LacerateDot.SetStacks(sim, 1)
					}
				} else {
					druid.AddRage(sim, refundAmount, druid.RageRefundMetrics)
				}
			},
		}),
	})

	dotAura := druid.CurrentTarget.RegisterAura(core.Aura{
		Label:     "Lacerate-" + strconv.Itoa(int(druid.Index)),
		ActionID:  actionID,
		MaxStacks: 5,
		Duration:  time.Second * 15,
	})
	druid.LacerateDot = core.NewDot(core.Dot{
		Spell:         druid.Lacerate,
		Aura:          dotAura,
		NumberOfTicks: 5,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 0.5,
			IsPeriodic:       true,
			BaseDamage: core.MultiplyByStacks(core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return tickDamage + 0.01*hitEffect.MeleeAttackPower(spell.Unit)
				},
				TargetSpellCoefficient: 0,
			}, dotAura),
			OutcomeApplier: druid.OutcomeFuncTick(),
		})),
	})
}

func (druid *Druid) CanLacerate(sim *core.Simulation) bool {
	return druid.CurrentRage() >= druid.Lacerate.DefaultCast.Cost
}
