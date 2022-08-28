package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerRuneTapSpell() {
	if !dk.Talents.RuneTap {
		return
	}

	actionID := core.ActionID{SpellID: 48982}
	cdTimer := dk.NewTimer()
	cd := time.Minute * 1
	healthMetrics := dk.NewHealthMetrics(actionID)

	healthGainMult := 0.0
	if dk.Talents.ImprovedRuneTap == 1 {
		healthGainMult = 0.33
	} else if dk.Talents.ImprovedRuneTap == 2 {
		healthGainMult = 0.66
	} else if dk.Talents.ImprovedRuneTap == 3 {
		healthGainMult = 1.0
	}

	glyphHealBonus := core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfRuneTap), 0.01, 0.0)

	baseCost := float64(core.NewRuneCost(0, 1, 0, 0, 0))
	rs := &RuneSpell{}
	dk.RuneTap = dk.RegisterSpell(rs, core.SpellConfig{
		ActionID:     actionID,
		Flags:        core.SpellFlagNoOnCastComplete,
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				// TODO: Does not invoke GCD?
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			maxHealth := dk.MaxHealth()
			dk.GainHealth(sim, (maxHealth*(0.1+glyphHealBonus))*(1.0+healthGainMult)*(1.0+core.TernaryFloat64(dk.VampiricBloodAura.IsActive(), 0.35, 0.0)), healthMetrics)
			rs.DoCost(sim)
		},
	}, func(sim *core.Simulation) bool {
		return dk.CastCostPossible(sim, 0, 1, 0, 0) && dk.RuneTap.IsReady(sim)
	}, nil)

	if !dk.Inputs.IsDps {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell:    dk.RuneTap.Spell,
			Type:     core.CooldownTypeSurvival,
			Priority: core.CooldownPriorityLow,
		})
	}
}
