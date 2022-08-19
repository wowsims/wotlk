package druid

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerRipSpell() {
	actionID := core.ActionID{SpellID: 49800}
	baseCost := 30.0
	refundAmount := baseCost * (0.4 * float64(druid.Talents.PrimalPrecision))

	druid.Rip = druid.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreResists,

		ResourceType: stats.Energy,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			ModifyCast:  druid.ApplyClearcasting,
			IgnoreHaste: true,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			OutcomeApplier:   druid.OutcomeFuncMeleeSpecialHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					druid.RipDot.Apply(sim)
					druid.SpendComboPoints(sim, spell.ComboPointMetrics())
				} else if refundAmount > 0 {
					druid.AddEnergy(sim, refundAmount, druid.PrimalPrecisionRecoveryMetrics)
				}
			},
		}),
	})

	target := druid.CurrentTarget
	druid.RipDot = core.NewDot(core.Dot{
		Spell: druid.Rip,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Rip-" + strconv.Itoa(int(druid.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 6,
		TickLength:    time.Second * 2,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 + core.TernaryFloat64(druid.HasSetBonus(ItemSetThunderheartHarness, 4), 0.15, 0),
			ThreatMultiplier: 1,
			IsPeriodic:       true,
			BaseDamage: core.BuildBaseDamageConfig(func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				comboPoints := float64(druid.ComboPoints())
				attackPower := hitEffect.MeleeAttackPower(spell.Unit)

				bonusTickDamage := 0.0
				if druid.Equip[items.ItemSlotRanged].ID == 28372 { // Idol of Feral Shadows
					bonusTickDamage += 7 * float64(comboPoints)
				}

				return (36.0+93.0*comboPoints+0.01*comboPoints*attackPower)/6.0 + bonusTickDamage
			}, 0),
			OutcomeApplier: druid.PrimalGoreOutcomeFuncTick(),
		}),
	})
}
