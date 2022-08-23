package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var ExposeArmorActionID = core.ActionID{SpellID: 8647}

func (rogue *Rogue) makeExposeArmor(comboPoints int32) *core.Spell {
	baseCost := 25.0 - float64(rogue.Talents.ImprovedExposeArmor)*5
	refundAmount := 0.4 * float64(rogue.Talents.QuickRecovery)
	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:     ExposeArmorActionID.WithTag(comboPoints),
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics | rogue.finisherFlags(),
		ResourceType: stats.Energy,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			ModifyCast:  rogue.CastModifier,
			IgnoreHaste: true,
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			ThreatMultiplier: 1,
			OutcomeApplier:   rogue.OutcomeFuncMeleeSpecialHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					rogue.ExposeArmorAura.Duration = rogue.exposeArmorDurations[comboPoints]
					rogue.ExposeArmorAura.Activate(sim)
					rogue.ApplyFinisher(sim, spell)
				} else {
					if refundAmount > 0 {
						rogue.AddEnergy(sim, spell.CurCast.Cost*refundAmount, rogue.QuickRecoveryMetrics)
					}
				}
			},
		}),
	})
}

func (rogue *Rogue) registerExposeArmorSpell() {
	rogue.ExposeArmorAura = core.ExposeArmorAura(rogue.CurrentTarget, rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfExposeArmor))
	onExpire := rogue.ExposeArmorAura.OnExpire
	if rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Once {
		rogue.ExposeArmorAura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
			onExpire(aura, sim)
			if rogue.initialArmorDebuffAura != nil {
				rogue.initialArmorDebuffAura.Activate(sim)
				if rogue.initialArmorDebuffAura.ActionID.SameActionIgnoreTag(core.SunderArmorActionID) {
					core.StartPeriodicAction(sim, core.SunderArmorPeriodicActionOptions(rogue.initialArmorDebuffAura))
				} else if rogue.initialArmorDebuffAura.ActionID.SameActionIgnoreTag(ExposeArmorActionID) {
					core.StartPeriodicAction(sim, core.ExposeArmorPeriodicActonOptions(rogue.initialArmorDebuffAura))
				} else if rogue.initialArmorDebuffAura.ActionID.SameActionIgnoreTag(core.AcidSpitActionID) {
					core.StartPeriodicAction(sim, core.AcidSpitPeriodicActionOptions(rogue.initialArmorDebuffAura))
				}
			}
		}
	}
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
