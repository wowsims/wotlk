package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (rogue *Rogue) registerFeintSpell() {
	cost := 20.0
	castModifier := rogue.CastModifier
	if rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfFeint) {
		castModifier = func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
			cast.Cost = 0
		}
	}
	rogue.Feint = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48659},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics,
		ResourceType: stats.Energy,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  time.Second,
			},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * 10,
			},
			IgnoreHaste: true,
			ModifyCast:  castModifier,
		},

		DamageMultiplier: 0.0,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCrit)
		},
	})
	// Feint
	if rogue.Rotation.UseFeint {
		rogue.AddMajorCooldown(core.MajorCooldown{
			Spell:    rogue.Feint,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeDPS,
		})
	}
}
