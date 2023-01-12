package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (rogue *Rogue) makeExposeArmor(comboPoints int32) *core.Spell {
	actionID := core.ActionID{SpellID: 8647}

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID.WithTag(comboPoints),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | rogue.finisherFlags(),

		EnergyCost: core.EnergyCostOptions{
			Cost:          25.0 - 5*float64(rogue.Talents.ImprovedExposeArmor),
			Refund:        0.4 * float64(rogue.Talents.QuickRecovery),
			RefundMetrics: rogue.QuickRecoveryMetrics,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				rogue.ExposeArmorAura.Duration = rogue.exposeArmorDurations[comboPoints]
				rogue.ExposeArmorAura.Activate(sim)
				rogue.ApplyFinisher(sim, spell)
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}

func (rogue *Rogue) registerExposeArmorSpell() {
	rogue.ExposeArmorAura = core.ExposeArmorAura(rogue.CurrentTarget, rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfExposeArmor))
	durationBonus := core.TernaryDuration(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfExposeArmor), time.Second*12, 0)
	rogue.exposeArmorDurations = [6]time.Duration{
		0,
		time.Second*6 + durationBonus,
		time.Second*12 + durationBonus,
		time.Second*18 + durationBonus,
		time.Second*24 + durationBonus,
		time.Second*30 + durationBonus,
	}
	rogue.ExposeArmor = [6]*core.Spell{
		nil,
		rogue.makeExposeArmor(1),
		rogue.makeExposeArmor(2),
		rogue.makeExposeArmor(3),
		rogue.makeExposeArmor(4),
		rogue.makeExposeArmor(5),
	}
}
