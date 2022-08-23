package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

// TODO: This currently only affects the caster, not other raid members.
func (priest *Priest) RegisterHymnOfHopeCD() {
	actionID := core.ActionID{SpellID: 64901}
	manaMetrics := priest.NewManaMetrics(actionID)

	numTicks := int32(4)

	channelTime := time.Duration(numTicks) * time.Second * 2
	manaPerTick := 0.0
	priest.Env.RegisterPostFinalizeEffect(func() {
		// This is 3%, but it increases the target's max mana by 20% for the duration
		// so just simplify to 3 * 1.2 = 3.6%.
		manaPerTick = priest.MaxMana() * 0.036
	})

	hymnOfHopeSpell := priest.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:         core.GCDDefault,
				ChannelTime: channelTime,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Minute * 6,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			period := spell.CurCast.ChannelTime / time.Duration(numTicks)
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   period,
				NumTicks: int(numTicks),
				OnAction: func(sim *core.Simulation) {
					priest.AddMana(sim, manaPerTick, manaMetrics, true)
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
