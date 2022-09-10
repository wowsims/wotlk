package druid

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerLacerateSpell() {
	actionID := core.ActionID{SpellID: 48568}

	cost := 15.0 - float64(druid.Talents.ShreddingAttacks)
	refundAmount := cost * 0.8

	tickDamage := 320.0 / 5
	initialDamage := 88.0
	if druid.Equip[items.ItemSlotRanged].ID == 27744 { // Idol of Ursoc
		tickDamage += 8
		initialDamage += 8
	}

	lbdm := core.TernaryFloat64(druid.HasSetBonus(ItemSetLasherweaveBattlegear, 2), 1.2, 1.0)
	dwdm := core.TernaryFloat64(druid.HasSetBonus(ItemSetDreamwalkerBattlegear, 2), 1.05, 1.0)
	t9bonus := core.TernaryFloat64(druid.HasT9FeralSetBonus(2), 1.05, 1.0)

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
			DamageMultiplier: lbdm * dwdm,
			ThreatMultiplier: 0.5,
			FlatThreatBonus:  267,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					damage := initialDamage + 0.01*hitEffect.MeleeAttackPower(spell.Unit)
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
						druid.LacerateDot.TakeSnapshot(sim, true)
					} else {
						druid.LacerateDot.Activate(sim)
						druid.LacerateDot.SetStacks(sim, 1)
						druid.LacerateDot.TakeSnapshot(sim, true)
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
		TickEffects: core.TickFuncSnapshot(druid.CurrentTarget, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: lbdm * t9bonus,
			ThreatMultiplier: 0.5,
			IsPeriodic:       true,
			BaseDamage: core.MultiplyByStacks(core.BaseDamageConfig{
				Calculator:             core.BaseDamageFuncMelee(tickDamage, tickDamage, 0.01),
				TargetSpellCoefficient: 0,
			}, dotAura),
			OutcomeApplier: druid.PrimalGoreOutcomeFuncTick(),
		}),
	})
}

func (druid *Druid) CanLacerate(sim *core.Simulation) bool {
	return druid.CurrentRage() >= druid.Lacerate.DefaultCast.Cost
}
