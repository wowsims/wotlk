package paladin

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (paladin *Paladin) registerSpiritualAttunement() {
	if paladin.Talents.SpiritualAttunement == 0 {
		return
	}

	paladin.SpiritualAttunementMetrics = paladin.NewManaMetrics(core.ActionID{SpellID: 33776})

	paladin.RegisterAura(core.Aura{
		Label:    "Spiritual Attunement",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Damage > 0 {
				paladin.AddMana(sim, spellEffect.Damage, paladin.SpiritualAttunementMetrics, false)
			}
		},
	})
}
