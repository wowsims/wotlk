package warrior

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

// TODO: GlyphOfSunderArmor will require refactoring this a bit

func (warrior *Warrior) newSunderArmorSpell(isDevastateEffect bool) *core.Spell {
	warrior.SunderArmorAuras = warrior.NewEnemyAuraArray(core.SunderArmorAura)

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
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.CanApplySunderAura(target)
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  360,

		RelatedAuras: []core.AuraArray{warrior.SunderArmorAuras},
	}

	extraStack := isDevastateEffect && warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfDevastate)
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
				if extraStack {
					aura.AddStack(sim)
				}
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
