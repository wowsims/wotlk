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
	cd := time.Minute*1 - dk.thassariansPlateCooldownReduction(dk.VampiricBlood)

	var bonusHealth float64
	dk.VampiricBloodAura = dk.RegisterAura(core.Aura{
		Label:    "Vampiric Blood",
		ActionID: actionID,
		Duration: time.Second*10 + core.TernaryDuration(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfVampiricBlood), 5*time.Second, 0),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			bonusHealth = dk.MaxHealth() * 0.15
			dk.AddStatsDynamic(sim, stats.Stats{stats.Health: bonusHealth})
			dk.GainHealth(sim, bonusHealth, healthMetrics)
			dk.PseudoStats.HealingTakenMultiplier *= 1.35
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.AddStatsDynamic(sim, stats.Stats{stats.Health: -bonusHealth})
			dk.PseudoStats.HealingTakenMultiplier /= 1.35
		},
	})

	dk.VampiricBlood = dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost:  1,
			RunicPowerGain: 10,
		},
		Cast: core.CastConfig{
			// TODO: does not invoke the GCD?
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.VampiricBloodAura.Activate(sim)
		},
	})

	if !dk.Inputs.IsDps {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell: dk.VampiricBlood,
			Type:  core.CooldownTypeSurvival,
		})
	}
}
