package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (rogue *Rogue) makeEnvenom(comboPoints int32) *core.Spell {
	refundAmount := 0.4 * float64(rogue.Talents.QuickRecovery)
	baseDamage := 215 + core.TernaryFloat64(rogue.HasSetBonus(ItemSetDeathmantle, 2), 40, 0)
	apRatio := 0.09
	chanceToRetainStacks := rogue.Talents.MasterPoisoner / 3.0

	cost := 35.0
	if rogue.HasSetBonus(ItemSetAssassination, 4) {
		cost -= 10
	}

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 57993, Tag: comboPoints},
		SpellSchool: core.SpellSchoolNature,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreResists | rogue.finisherFlags(),

		ResourceType: stats.Energy,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  time.Second,
			},
			ModifyCast:  rogue.applyDeathmantle,
			IgnoreHaste: true,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			DamageMultiplier: 1 + []float64{0.0, 0.07, 0.14, 0.2}[rogue.Talents.VilePoisons],
			ThreatMultiplier: 1,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					target := hitEffect.Target
					deadlyPoisonStacks := rogue.DeadlyPoisonDots[target.Index].GetStacks()
					doses := float64(core.MinInt32(deadlyPoisonStacks, comboPoints))
					return baseDamage*doses + apRatio*doses*hitEffect.MeleeAttackPower(spell.Unit)
				},
				TargetSpellCoefficient: 0,
			},
			OutcomeApplier: rogue.OutcomeFuncMeleeSpecialHitAndCrit(rogue.MeleeCritMultiplier(true, false)),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				target := spellEffect.Target
				if spellEffect.Landed() {
					rogue.ApplyFinisher(sim, spell)
					rogue.ApplyCutToTheChase(sim)
					deadlyPoisonStacks := rogue.DeadlyPoisonDots[target.Index].GetStacks()
					doses := core.MinInt32(deadlyPoisonStacks, comboPoints)
					if chanceToRetainStacks < 1 && sim.RandomFloat("Master Poisoner") > float64(chanceToRetainStacks) {
						rogue.DeadlyPoisonDots[target.Index].Cancel(sim)
					}
					rogue.EnvenomAura.Duration = time.Second * time.Duration(1+doses)
					rogue.EnvenomAura.Activate(sim)
				} else {
					if refundAmount > 0 {
						rogue.AddEnergy(sim, spell.CurCast.Cost*refundAmount, rogue.QuickRecoveryMetrics)
					}
				}
			},
		}),
	})
}

func (rogue *Rogue) registerEnvenom() {
	rogue.EnvenomAura = rogue.RegisterAura(core.Aura{
		Label:    "Envenom",
		ActionID: core.ActionID{SpellID: 57993},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.DeadlyPoisonProcChanceBonus += 0.15
			rogue.InstantPoisonProcChanceBonus += 0.75
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.DeadlyPoisonProcChanceBonus -= 0.15
			rogue.InstantPoisonProcChanceBonus -= 0.75
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
