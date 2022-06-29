package warrior

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warrior *Warrior) registerDevastateSpell() {
	cost := 15.0 - float64(warrior.Talents.ImprovedSunderArmor) - float64(warrior.Talents.FocusedRage)
	refundAmount := cost * 0.8

	normalBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 0, 0.5, true)

	warrior.Devastate = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 30022},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeMHSpecial,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			FlatThreatBonus:  100,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					// Bonus 35 damage / stack of sunder. Counts stacks AFTER cast but only if stacks > 0.
					sunderBonus := 0.0
					saStacks := warrior.SunderArmorAura.GetStacks()
					if saStacks != 0 {
						sunderBonus = 35 * float64(core.MinInt32(saStacks+1, 5))
					}

					return normalBaseDamage(sim, hitEffect, spell) + sunderBonus
				},
				TargetSpellCoefficient: 0,
			},
			OutcomeApplier: warrior.OutcomeFuncMeleeWeaponSpecialHitAndCrit(warrior.critMultiplier(true)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					if !warrior.ExposeArmorAura.IsActive() {
						warrior.SunderArmorDevastate.Cast(sim, spellEffect.Target)
					}
				} else {
					warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
				}
			},
		}),
	})
}

func (warrior *Warrior) CanDevastate(sim *core.Simulation) bool {
	return warrior.CurrentRage() >= warrior.Devastate.DefaultCast.Cost
}
