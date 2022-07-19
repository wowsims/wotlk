package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) registerHornOfWinterSpell() {
	actionID := core.ActionID{SpellID: 57623}
	duration := time.Minute * time.Duration((2.0 + core.TernaryFloat64(deathKnight.HasMinorGlyph(proto.DeathKnightMinorGlyph_GlyphOfHornOfWinter), 1.0, 0.0)))

	strengthBonus := 155.0
	agilityBonus := 155.0
	deathKnight.HornOfWinterAura = deathKnight.RegisterAura(core.Aura{
		Label:    "Horn of Winter",
		ActionID: actionID,
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if !deathKnight.OtherRelevantStrAgiActive {
				bonusStats := deathKnight.ApplyStatDependencies(stats.Stats{stats.Strength: strengthBonus, stats.Agility: agilityBonus})
				deathKnight.HornOfWinterAura.Unit.AddStatsDynamic(sim, bonusStats)
			}
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if !deathKnight.OtherRelevantStrAgiActive {
				bonusStats := deathKnight.ApplyStatDependencies(stats.Stats{stats.Strength: -strengthBonus, stats.Agility: -agilityBonus})
				deathKnight.HornOfWinterAura.Unit.AddStatsDynamic(sim, bonusStats)
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
			if deathKnight.Rotation.RefreshHornOfWinter {
				deathKnight.HornOfWinterAura.Activate(sim)
				deathKnight.HornOfWinterAura.Prioritize()
			}

			amountOfRunicPower := 10.0
			deathKnight.AddRunicPower(sim, amountOfRunicPower, deathKnight.HornOfWinter.RunicPowerMetrics())
		},
	})
}

func (deathKnight *DeathKnight) CanHornOfWinter(sim *core.Simulation) bool {
	return deathKnight.HornOfWinter.IsReady(sim)
}

func (deathKnight *DeathKnight) ShouldHornOfWinter(sim *core.Simulation) bool {
	return deathKnight.Rotation.RefreshHornOfWinter && deathKnight.HornOfWinter.IsReady(sim) && !deathKnight.HornOfWinterAura.IsActive()
}
