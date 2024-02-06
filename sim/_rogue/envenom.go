package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

// TODO: Link to Rune of Envenom
// TODO: Level based damage scaling
func (rogue *Rogue) registerEnvenom() {
	rogue.EnvenomAura = rogue.RegisterAura(core.Aura{
		Label:    "Envenom",
		ActionID: core.ActionID{SpellID: 57993},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.instantPoisonProcChanceBonus(0.75)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.instantPoisonProcChanceBonus(0.0)
		},
	})

	rogue.Envenom = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 57993},
		SpellSchool:  core.SpellSchoolNature,
		ProcMask:     core.ProcMaskMeleeMHSpecial, // not core.ProcMaskSpellDamage
		Flags:        core.SpellFlagMeleeMetrics | rogue.finisherFlags() | SpellFlagColdBlooded | core.SpellFlagAPL,
		MetricSplits: 6,

		EnergyCost: core.EnergyCostOptions{
			Cost:   35,
			Refund: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(spell.Unit.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0 && rogue.DeadlyPoison.Dot(target).IsActive()
		},

		DamageMultiplier: 1 +
			[]float64{0.0, 0.04, 0.08, 0.12, 0.16, 0.2}[rogue.Talents.VilePoisons],
		CritMultiplier:   rogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			comboPoints := rogue.ComboPoints()
			// - the aura is active even if the attack fails to land
			// - the aura is applied before the hit effect
			// See: https://github.com/where-fore/rogue-wotlk/issues/32
			rogue.EnvenomAura.Duration = time.Second * time.Duration(1+comboPoints)
			rogue.EnvenomAura.Activate(sim)

			dp := rogue.DeadlyPoison.Dot(target)
			// - 215 base is scaled by consumed doses (<= comboPoints)
			// - apRatio is independent of consumed doses (== comboPoints)
			consumed := min(dp.GetStacks(), comboPoints)
			baseDamage := 215*float64(consumed) + 0.09*float64(comboPoints)*spell.MeleeAttackPower()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				rogue.ApplyFinisher(sim, spell)
			} else {
				spell.IssueRefund(sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}

func (rogue *Rogue) EnvenomDuration(comboPoints int32) time.Duration {
	return time.Second * (1 + time.Duration(comboPoints))
}
