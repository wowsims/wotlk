package warrior

import (
	"github.com/wowsims/sod/sim/core"
)

func (warrior *Warrior) newSunderArmorSpell(isDevastateEffect bool) *core.Spell {
	warrior.SunderArmorAuras = warrior.NewEnemyAuraArray(core.SunderArmorAura)
	spellID := map[int32]int32{
		25: 7405,
		40: 8380,
		50: 11596,
		60: 11597,
	}[warrior.Level]

	config := core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		RageCost: core.RageCostOptions{
			Cost:   15,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.CanApplySunderAura(target)
		},

		ThreatMultiplier: 1,
		// TODO Warrior: set threat according to spell's level
		FlatThreatBonus: 360,

		RelatedAuras: []core.AuraArray{warrior.SunderArmorAuras},
	}

	if isDevastateEffect {
		config.RageCost = core.RageCostOptions{}
		config.Cast.DefaultCast.GCD = 0
		config.ExtraCastCondition = nil

		// In wrath sunder from devastate generates no threat
		config.ThreatMultiplier = 0
		config.FlatThreatBonus = 0
	} else {
		config.Flags |= core.SpellFlagAPL
	}

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		var result *core.SpellResult
		if isDevastateEffect {
			result = spell.CalcOutcome(sim, target, spell.OutcomeAlwaysHit)
		} else {
			result = spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			result.Threat = spell.ThreatFromDamage(result.Outcome, 0.05*spell.MeleeAttackPower())
		}

		if result.Landed() {
			aura := warrior.SunderArmorAuras.Get(target)
			aura.Activate(sim)
			if aura.IsActive() {
				aura.AddStack(sim)
			}
		} else {
			spell.IssueRefund(sim)
		}

		spell.DealOutcome(sim, result)
	}
	return warrior.RegisterSpell(config)
}

func (warrior *Warrior) CanApplySunderAura(target *core.Unit) bool {
	return warrior.SunderArmorAuras.Get(target).IsActive() || !warrior.SunderArmorAuras.Get(target).ExclusiveEffects[0].Category.AnyActive()
}
