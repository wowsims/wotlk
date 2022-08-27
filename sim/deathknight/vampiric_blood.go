package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerVampiricBloodSpell() {
	if !dk.Talents.VampiricBlood {
		return
	}

	actionID := core.ActionID{SpellID: 55233}
	healthMetrics := dk.NewHealthMetrics(actionID)

	cdTimer := dk.NewTimer()
	cd := time.Minute * 1

	var bonusHealth float64
	dk.VampiricBloodAura = dk.RegisterAura(core.Aura{
		Label:    "Vampiric Blood",
		ActionID: actionID,
		Duration: time.Second*10 + core.TernaryDuration(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfVampiricBlood), 5*time.Second, 0),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			bonusHealth = dk.MaxHealth() * 0.15
			dk.AddStatsDynamic(sim, stats.Stats{stats.Health: bonusHealth})
			dk.GainHealth(sim, bonusHealth, healthMetrics)
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.AddStatsDynamic(sim, stats.Stats{stats.Health: -bonusHealth})
		},
	})

	baseCost := float64(core.NewRuneCost(10, 1, 0, 0, 0))
	rs := &RuneSpell{}
	dk.VampiricBlood = dk.RegisterSpell(rs, core.SpellConfig{
		ActionID:     actionID,
		Flags:        core.SpellFlagNoOnCastComplete,
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				// TODO: does not invoke the GCD?
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.VampiricBloodAura.Activate(sim)
			rs.DoCost(sim)
		},
	}, func(sim *core.Simulation) bool {
		return dk.CastCostPossible(sim, 0, 1, 0, 0) && dk.VampiricBlood.IsReady(sim)
	}, nil)

	if !dk.Inputs.IsDps {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell:    dk.VampiricBlood.Spell,
			Type:     core.CooldownTypeSurvival,
			Priority: core.CooldownPriorityLow,
		})
	}
}
