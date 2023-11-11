package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"

	"github.com/wowsims/wotlk/sim/core"
)

func (rogue *Rogue) registerShadowDanceCD() {
	if !rogue.Talents.ShadowDance {
		return
	}

	duration := time.Second * 6
	if rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfShadowDance) {
		duration = time.Second * 8
	}

	actionID := core.ActionID{SpellID: 51713}

	rogue.ShadowDanceAura = rogue.RegisterAura(core.Aura{
		Label:    "Shadow Dance",
		ActionID: actionID,
		Duration: duration,
		// Can now cast opening abilities outside of stealth
		// Covered in rogue.go by IsStealthed()
	})

	rogue.ShadowDance = rogue.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Minute,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			rogue.ShadowDanceAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    rogue.ShadowDance,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			return rogue.GCD.IsReady(s) && rogue.ComboPoints() <= 2 && rogue.CurrentEnergy() >= 60
		},
	})
}
