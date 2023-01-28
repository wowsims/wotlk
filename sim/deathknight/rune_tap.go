package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
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

	dk.RuneTap = dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost: 1,
		},
		Cast: core.CastConfig{
			// TODO: Does not invoke GCD?
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			maxHealth := dk.MaxHealth()
			dk.GainHealth(sim, (maxHealth*(0.1+glyphHealBonus))*(1.0+healthGainMult)*(1.0+core.TernaryFloat64(dk.VampiricBloodAura.IsActive(), 0.35, 0.0)), healthMetrics)
		},
	})

	if !dk.Inputs.IsDps {
		// dk.AddMajorCooldown(core.MajorCooldown{
		// 	Spell:    dk.RuneTap,
		// 	Type:     core.CooldownTypeSurvival,
		// 	Priority: core.CooldownPriorityDefault,
		// 	CanActivate: func(sim *core.Simulation, character *core.Character) bool {
		// 		success := dk.RuneTap.CanCast(sim, nil)
		// 		if !success && dk.BloodTap.IsReady(sim) {
		// 			dk.BloodTap.Cast(sim, nil)
		// 			success = dk.RuneTap.CanCast(sim, nil)
		// 		}
		// 		return success
		// 	},
		// })
	}
}
