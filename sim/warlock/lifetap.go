package warlock

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerLifeTapSpell() {
	actionID := core.ActionID{SpellID: 57946}
	impLifetap := 1.0 + 0.1*float64(warlock.Talents.ImprovedLifeTap)
	manaMetrics := warlock.NewManaMetrics(actionID)

	var petManaMetrics *core.ResourceMetrics
	if warlock.Talents.ManaFeed && warlock.Pet != nil {
		petManaMetrics = warlock.Pet.NewManaMetrics(actionID)
	}

	warlock.LifeTap = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			restore := (2000.0 + 0.5*warlock.GetStat(stats.SpellPower)) * impLifetap
			warlock.AddMana(sim, restore, manaMetrics)

			if warlock.Talents.ManaFeed && warlock.Pet != nil {
				warlock.Pet.AddMana(sim, restore, petManaMetrics)
			}
			if warlock.GlyphOfLifeTapAura != nil {
				warlock.GlyphOfLifeTapAura.Activate(sim)
			}
		},
	})
}
