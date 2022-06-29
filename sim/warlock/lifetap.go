package warlock

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warlock *Warlock) registerLifeTapSpell() {
	actionID := core.ActionID{SpellID: 27222}
	baseRestore := 582.0 * (1.0 + 0.1*float64(warlock.Talents.ImprovedLifeTap))
	manaMetrics := warlock.NewManaMetrics(actionID)

	petRestore := 0.3333 * float64(warlock.Talents.ManaFeed)
	var petManaMetrics []*core.ResourceMetrics
	if warlock.Talents.ManaFeed > 0 {
		for _, pet := range warlock.Pets {
			petManaMetrics = append(petManaMetrics, pet.GetPet().NewManaMetrics(actionID))
		}
	}

	warlock.LifeTap = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			ThreatMultiplier: 1,
			FlatThreatBonus:  1,
			OutcomeApplier:   warlock.OutcomeFuncAlwaysHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				// Life tap adds 0.8*sp to mana restore
				restore := baseRestore + (warlock.GetStat(stats.SpellPower)+warlock.GetStat(stats.ShadowSpellPower))*0.8
				warlock.AddMana(sim, restore, manaMetrics, true)

				if warlock.Talents.ManaFeed > 0 {
					for i, pet := range warlock.Pets {
						pet.GetPet().AddMana(sim, restore*petRestore, petManaMetrics[i], true)
					}
				}
			},
		}),
	})
}
