package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (rogue *Rogue) registerPreparationCD() {
	if !rogue.Talents.Preparation {
		return
	}

	rogue.Preparation = rogue.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 14185},
		Flags:    core.SpellFlagAPL,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Minute*8 - time.Second*time.Duration(90*rogue.Talents.FilthyTricks),
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			// Spells affected by Preparation are: Cold Blood, Shadowstep, Vanish (Overkill/Master of Subtlety), Evasion, Sprint
			// If Glyph of Preparation is applied, Blade Flurry, Dismantle, and Kick are also affected
			var affectedSpells = []*core.Spell{rogue.ColdBlood, rogue.Shadowstep, rogue.Vanish}
			if rogue.GetCharacter().HasGlyph(int32(proto.RogueMajorGlyph_GlyphOfPreparation)) {
				affectedSpells = append(affectedSpells, rogue.BladeFlurry)
			}
			// Reset Cooldown on affected spells
			for _, affectedSpell := range affectedSpells {
				if affectedSpell != nil {
					affectedSpell.CD.Reset()
				}
			}
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    rogue.Preparation,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return !rogue.Vanish.CD.IsReady(sim)
		},
	})
}
