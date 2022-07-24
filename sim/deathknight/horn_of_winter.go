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
				if aura.IsActive() {
					aura.Deactivate(sim)
					aura.Activate(sim)
				}
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

	dk.HornOfWinter = dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: 20 * time.Second,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if dk.Inputs.RefreshHornOfWinter {
				dk.HornOfWinterAura.Activate(sim)
				dk.HornOfWinterAura.Prioritize()
			}

			amountOfRunicPower := 10.0
			dk.AddRunicPower(sim, amountOfRunicPower, dk.HornOfWinter.RunicPowerMetrics())
		},
	})
}

func (dk *Deathknight) CanHornOfWinter(sim *core.Simulation) bool {
	return dk.HornOfWinter.IsReady(sim)
}

func (dk *Deathknight) ShouldHornOfWinter(sim *core.Simulation) bool {
	return dk.Inputs.RefreshHornOfWinter && dk.HornOfWinter.IsReady(sim) && !dk.HornOfWinterAura.IsActive()
}

func (dk *Deathknight) CastHornOfWinter(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanHornOfWinter(sim) {
		dk.HornOfWinter.Cast(sim, target)
		return true
	}
	return false
}
