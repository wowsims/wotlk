package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (rogue *Rogue) registerExposeArmorSpell() {
	baseCost := 25.0
	refundAmount := 0.4 * float64(rogue.Talents.QuickRecovery)

	rogue.ExposeArmorAura = core.ExposeArmorAura(rogue.CurrentTarget, rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfExposeArmor))

	rogue.ExposeArmor = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 26866, Tag: 5},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics | rogue.finisherFlags(),

		ResourceType: stats.Energy,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			ModifyCast:  rogue.applyDeathmantle,
			IgnoreHaste: true,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			ThreatMultiplier: 1,
			OutcomeApplier:   rogue.OutcomeFuncMeleeSpecialHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
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
