package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (dk *Deathknight) GetRuneTapHealing() float64 {
	maxHealth := dk.MaxHealth()
	return maxHealth * dk.bonusCoeffs.runeTapHealing * core.TernaryFloat64(dk.VampiricBloodAura.IsActive(), 1.35, 1.0)
}

func (dk *Deathknight) registerRuneTapSpell() {
	if !dk.Talents.RuneTap {
		return
	}

	actionID := core.ActionID{SpellID: 48982}
	cdTimer := dk.NewTimer()
	cd := time.Minute * 1
	healthMetrics := dk.NewHealthMetrics(actionID)

	dk.bonusCoeffs.runeTapHealing = []float64{1.0, 1.33, 1.66, 2.0}[dk.Talents.ImprovedRuneTap] *
		core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfRuneTap), 0.11, 0.10)

	dk.RuneTap = dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost: 1,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.GainHealth(sim, dk.GetRuneTapHealing(), healthMetrics)
		},
	})
}
