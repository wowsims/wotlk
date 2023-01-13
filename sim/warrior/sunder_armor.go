package warrior

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

// TODO: GlyphOfSunderArmor will require refactoring this a bit

func (warrior *Warrior) newSunderArmorSpell(isDevastateEffect bool) *core.Spell {
	warrior.SunderArmorAura = core.SunderArmorAura(warrior.CurrentTarget)

	config := core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47467},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		RageCost: core.RageCostOptions{
			Cost:   15 - float64(warrior.Talents.FocusedRage) - float64(warrior.Talents.Puncture),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  360,
	}

	extraStack := isDevastateEffect && warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfDevastate)
	if isDevastateEffect {
		config.RageCost = core.RageCostOptions{}
		config.Cast.DefaultCast.GCD = 0

		// In wrath sunder from devastate generates no threat
		config.ThreatMultiplier = 0
		config.FlatThreatBonus = 0
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
			warrior.SunderArmorAura.Activate(sim)
			if warrior.SunderArmorAura.IsActive() {
				warrior.SunderArmorAura.AddStack(sim)
				if extraStack {
					warrior.SunderArmorAura.AddStack(sim)
				}
			}
		} else {
			spell.IssueRefund(sim)
		}

		spell.DealOutcome(sim, result)
	}
	return warrior.RegisterSpell(config)
}

func (warrior *Warrior) CanSunderArmor(sim *core.Simulation) bool {
	return warrior.CurrentRage() >= warrior.SunderArmor.DefaultCast.Cost && warrior.CanApplySunderAura()
}
func (warrior *Warrior) CanApplySunderAura() bool {
	return warrior.SunderArmorAura.IsActive() || !warrior.SunderArmorAura.ExclusiveEffects[0].Category.AnyActive()
}
