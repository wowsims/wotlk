package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerHornOfWinterSpell() {
	actionID := core.ActionID{SpellID: 57623}
	duration := time.Minute * time.Duration((2.0 + core.TernaryFloat64(dk.HasMinorGlyph(proto.DeathknightMinorGlyph_GlyphOfHornOfWinter), 1.0, 0.0)))

	bonusStats := stats.Stats{stats.Strength: 155.0, stats.Agility: 155.0}
	negativeStats := bonusStats.Multiply(-1)

	dk.HornOfWinterAura = dk.RegisterAura(core.Aura{
		Label:    "Horn of Winter",
		ActionID: actionID,
		Duration: duration,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			if dk.Inputs.PrecastHornOfWinter && dk.Inputs.RefreshHornOfWinter {
				aura.Activate(sim)
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if !dk.OtherRelevantStrAgiActive {
				dk.HornOfWinterAura.Unit.AddStatsDynamic(sim, bonusStats)
			}
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if !dk.OtherRelevantStrAgiActive {
				dk.HornOfWinterAura.Unit.AddStatsDynamic(sim, negativeStats)
			}
		},
	})

	dk.HornOfWinter = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: 20 * time.Second,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if dk.Inputs.RefreshHornOfWinter {
				dk.HornOfWinterAura.Activate(sim)
			}
			dk.AddRunicPower(sim, 10, dk.HornOfWinter.RunicPowerMetrics())
		},
	})
}

func (dk *Deathknight) ShouldHornOfWinter(sim *core.Simulation) bool {
	return dk.Inputs.RefreshHornOfWinter && dk.HornOfWinter.IsReady(sim) && !dk.HornOfWinterAura.IsActive()
}
