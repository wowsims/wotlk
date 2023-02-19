package warlock

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerLifeTapSpell() {
	actionID := core.ActionID{SpellID: 57946}
	baseRestore := 2000.0 * (1.0 + 0.1*float64(warlock.Talents.ImprovedLifeTap))
	manaMetrics := warlock.NewManaMetrics(actionID)

	petRestore := core.TernaryFloat64(warlock.Talents.ManaFeed, 1, 0)
	var petManaMetrics []*core.ResourceMetrics
	if warlock.Talents.ManaFeed {
		for _, pet := range warlock.Pets {
			petManaMetrics = append(petManaMetrics, pet.GetPet().NewManaMetrics(actionID))
		}
	}

	warlock.LifeTap = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Life tap adds 0.5*sp to mana restore
			restore := baseRestore + 0.5*warlock.GetStat(stats.SpellPower)
			warlock.AddMana(sim, restore, manaMetrics)

			if warlock.Talents.ManaFeed {
				for i, pet := range warlock.Pets {
					pet.GetPet().AddMana(sim, restore*petRestore, petManaMetrics[i])
				}
			}
			if warlock.GlyphOfLifeTapAura != nil {
				warlock.GlyphOfLifeTapAura.Activate(sim)
			}
		},
	})
}
