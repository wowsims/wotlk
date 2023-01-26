package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (rogue *Rogue) makeEnvenom(comboPoints int32) *core.Spell {
	apRatio := 0.09 * float64(comboPoints)
	chanceToRetainStacks := []float64{0, 0.33, 0.66, 1}[rogue.Talents.MasterPoisoner]

	// TODO Envenom can only be cast if the target is afflicted by Deadly Poison
	//  The current rotation code doesn't handle cast failures gracefully, so this is hard to
	//  work around at the moment
	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 57993, Tag: comboPoints},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskMeleeMHSpecial, // not core.ProcMaskSpellDamage
		Flags:       core.SpellFlagMeleeMetrics | rogue.finisherFlags(),

		EnergyCost: core.EnergyCostOptions{
			Cost:          35,
			Refund:        0.4 * float64(rogue.Talents.QuickRecovery),
			RefundMetrics: rogue.QuickRecoveryMetrics,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1 +
			0.02*float64(rogue.Talents.FindWeakness) +
			[]float64{0.0, 0.07, 0.14, 0.2}[rogue.Talents.VilePoisons],
		CritMultiplier:   rogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// - the aura is active even if the attack fails to land
			// - the aura is applied before the hit effect
			// See: https://github.com/where-fore/rogue-wotlk/issues/32
			rogue.EnvenomAura.Duration = time.Second * time.Duration(1+comboPoints)
			rogue.EnvenomAura.Activate(sim)

			dp := rogue.DeadlyPoison.Dot(target)
			// - 215 base is scaled by consumed doses (<= comboPoints)
			// - apRatio is independent of consumed doses (== comboPoints)
			consumed := core.MinInt32(dp.GetStacks(), comboPoints)
			baseDamage := 215*float64(consumed) + apRatio*spell.MeleeAttackPower()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				rogue.ApplyFinisher(sim, spell)
				rogue.ApplyCutToTheChase(sim)
				if !sim.Proc(chanceToRetainStacks, "Master Poisoner") {
					if newStacks := dp.GetStacks() - comboPoints; newStacks > 0 {
						dp.SetStacks(sim, newStacks)
					} else {
						dp.Cancel(sim)
					}
				}
			} else {
				spell.IssueRefund(sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}

func (rogue *Rogue) registerEnvenom() {
	rogue.EnvenomAura = rogue.RegisterAura(core.Aura{
		Label:    "Envenom",
		ActionID: core.ActionID{SpellID: 57993},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.deadlyPoisonProcChanceBonus += 0.15
			rogue.UpdateInstantPoisonPPM(0.75)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.deadlyPoisonProcChanceBonus -= 0.15
			rogue.UpdateInstantPoisonPPM(0.0)
		},
	})
	rogue.Envenom = [6]*core.Spell{
		nil,
		rogue.makeEnvenom(1),
		rogue.makeEnvenom(2),
		rogue.makeEnvenom(3),
		rogue.makeEnvenom(4),
		rogue.makeEnvenom(5),
	}
}
