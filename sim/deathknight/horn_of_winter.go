package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *Deathknight) registerHornOfWinterSpell() {
	actionID := core.ActionID{SpellID: 57623}
	duration := time.Minute * time.Duration((2.0 + core.TernaryFloat64(deathKnight.HasMinorGlyph(proto.DeathknightMinorGlyph_GlyphOfHornOfWinter), 1.0, 0.0)))

	bonusStats := stats.Stats{stats.Strength: 155.0, stats.Agility: 155.0}
	negativeStats := bonusStats.Multiply(-1)

	deathKnight.HornOfWinterAura = deathKnight.RegisterAura(core.Aura{
		Label:    "Horn of Winter",
		ActionID: actionID,
		Duration: duration,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			if deathKnight.Options.PrecastHornOfWinter && deathKnight.RefreshHornOfWinter {
				if aura.IsActive() {
					aura.Deactivate(sim)
					aura.Activate(sim)
				}
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if !deathKnight.OtherRelevantStrAgiActive {
				deathKnight.HornOfWinterAura.Unit.AddStatsDynamic(sim, bonusStats)
			}
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if !deathKnight.OtherRelevantStrAgiActive {
				deathKnight.HornOfWinterAura.Unit.AddStatsDynamic(sim, negativeStats)
			}
		},
	})

	deathKnight.HornOfWinter = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = deathKnight.getModifiedGCD()
			},
			CD: core.Cooldown{
				Timer:    deathKnight.NewTimer(),
				Duration: 20 * time.Second,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if deathKnight.RefreshHornOfWinter {
				deathKnight.HornOfWinterAura.Activate(sim)
				deathKnight.HornOfWinterAura.Prioritize()
			}

			amountOfRunicPower := 10.0
			deathKnight.AddRunicPower(sim, amountOfRunicPower, deathKnight.HornOfWinter.RunicPowerMetrics())
		},
	})
}

func (deathKnight *Deathknight) CanHornOfWinter(sim *core.Simulation) bool {
	return deathKnight.HornOfWinter.IsReady(sim)
}

func (deathKnight *Deathknight) ShouldHornOfWinter(sim *core.Simulation) bool {
	return deathKnight.RefreshHornOfWinter && deathKnight.HornOfWinter.IsReady(sim) && !deathKnight.HornOfWinterAura.IsActive()
}

func (deathKnight *Deathknight) CastHornOfWinter(sim *core.Simulation, target *core.Unit) bool {
	if deathKnight.CanHornOfWinter(sim) {
		deathKnight.HornOfWinter.Cast(sim, target)
		return true
	}
	return false
}
