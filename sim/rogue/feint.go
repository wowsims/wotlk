package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (rogue *Rogue) registerFeintSpell() {
	rogue.Feint = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48659},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		EnergyCost: core.EnergyCostOptions{
			Cost: core.TernaryFloat64(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfFeint), 0, 20),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * 10,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 0,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
		},
	})
	// Feint
	if rogue.Rotation.UseFeint {
		rogue.AddMajorCooldown(core.MajorCooldown{
			Spell:    rogue.Feint,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeDPS,
			//don't feint if you're gonna waste energy by using the gcd
			ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
				thresh := 55.0 //55 simmed best with standard settings for now 3/12/2023, will refine with the rotational refinements. 55 was definitely best for combat, didn't make a difference for muti
				return rogue.CurrentEnergy() <= thresh
			},
		})
	}
}
