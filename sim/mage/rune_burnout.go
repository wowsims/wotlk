package mage

import (
	"github.com/wowsims/classic/sod/sim/core"
	"github.com/wowsims/classic/sod/sim/core/stats"
)

func (mage *Mage) registerRuneBurnout() {
	if !mage.HasRuneById(MageRuneChestBurnout) {
		return
	}

	actionID := core.ActionID{SpellID: 412286}
	metric := mage.NewManaMetrics(actionID)

	mage.RegisterAura(core.Aura{
		Label:    "Burnout",
		ActionID: core.ActionID{SpellID: 412286},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.SpellCrit, 15*core.SpellCritRatingPerCritChance)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.SpellCrit, -(15 * core.SpellCritRatingPerCritChance))
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(SpellFlagMage) && !result.DidCrit() {
				return
			}
			aura.Unit.SpendMana(sim, aura.Unit.BaseMana*0.01, metric)
		},
	})
}
