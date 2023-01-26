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

	// Spells affected by Preparation are: Cold Blood, Shadowstep, Vanish (Overkill/Master of Subtlety), Evasion, Sprint
	// If Glyph of Preparation is applied, Blade Flurry, Dismantle, and Kick are also affected
	var affectedSpells = []*core.Spell{rogue.ColdBlood, rogue.Shadowstep, rogue.MasterOfSubtlety, rogue.Overkill}
	if rogue.GetCharacter().HasGlyph(int32(proto.RogueMajorGlyph_GlyphOfPreparation)) {
		affectedSpells = append(affectedSpells, rogue.BladeFlurry)
	}

	rogue.Preparation = rogue.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 14185},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second*8*60 - time.Second*time.Duration(90*rogue.Talents.FilthyTricks),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
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
			return !rogue.MasterOfSubtlety.CD.IsReady(sim)
		},
	})
}
