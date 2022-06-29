package warrior

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

var SunderArmorActionID = core.ActionID{SpellID: 25225}

func (warrior *Warrior) newSunderArmorSpell(isDevastateEffect bool) *core.Spell {
	cost := 15.0 - float64(warrior.Talents.ImprovedSunderArmor) - float64(warrior.Talents.FocusedRage)
	refundAmount := cost * 0.8
	warrior.SunderArmorAura = core.SunderArmorAura(warrior.CurrentTarget, 0)
	warrior.ExposeArmorAura = core.ExposeArmorAura(warrior.CurrentTarget, 2)

	config := core.SpellConfig{
		ActionID:    SunderArmorActionID,
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
	}
	if isDevastateEffect {
		config.ResourceType = 0
		config.BaseCost = 0
		config.Cast.DefaultCast.Cost = 0
		config.Cast.DefaultCast.GCD = 0
	}

	effect := core.SpellEffect{
		ProcMask: core.ProcMaskMeleeMHSpecial,

		ThreatMultiplier: 1,
		FlatThreatBonus:  301.5,

		OutcomeApplier: warrior.OutcomeFuncMeleeSpecialHit(),

		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				warrior.SunderArmorAura.Activate(sim)
				if warrior.SunderArmorAura.IsActive() {
					warrior.SunderArmorAura.AddStack(sim)
				}
			} else {
				warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
			}
		},
	}
	if isDevastateEffect {
		effect.OutcomeApplier = warrior.OutcomeFuncAlwaysHit()
		effect.OnInit = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if warrior.SunderArmorAura.GetStacks() == 5 {
				spellEffect.ThreatMultiplier = 0
				spellEffect.FlatThreatBonus = 0
			}
		}
	}

	config.ApplyEffects = core.ApplyEffectFuncDirectDamage(effect)
	return warrior.RegisterSpell(config)
}

func (warrior *Warrior) CanSunderArmor(sim *core.Simulation) bool {
	return warrior.CurrentRage() >= warrior.SunderArmor.DefaultCast.Cost && !warrior.ExposeArmorAura.IsActive()
}
