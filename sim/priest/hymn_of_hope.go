package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

// TODO: This currently only affects the caster, not other raid members.
func (priest *Priest) RegisterHymnOfHopeCD() {
	actionID := core.ActionID{SpellID: 64901}
	manaMetrics := priest.NewManaMetrics(actionID)

	numTicks := 4 + core.TernaryInt32(priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfHymnOfHope), 1, 0)

	hymnOfHopeSpell := priest.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagHelpful,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Minute * 6,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   spell.Unit.ApplyCastSpeedForSpell(time.Second*2, spell),
				NumTicks: int(numTicks),
				OnAction: func(sim *core.Simulation) {
					// This is 3%, but it increases the target's max mana by 20% for the duration
					// so just simplify to 3 * 1.2 = 3.6%.
					priest.AddMana(sim, priest.MaxMana()*0.036, manaMetrics)
				},
			})
		},
	})

	priest.AddMajorCooldown(core.MajorCooldown{
		Spell: hymnOfHopeSpell,
		Type:  core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return character.CurrentManaPercent() < 0.1
		},
	})
}
