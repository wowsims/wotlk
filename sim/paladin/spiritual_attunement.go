package paladin

import (
	"github.com/wowsims/tbc/sim/core"
)

func (paladin *Paladin) registerSpiritualAttunement() {
	coeff := 0.1 * core.TernaryFloat64(ItemSetLightbringerArmor.CharacterHasSetBonus(&paladin.Character, 2), 1.1, 1)
	paladin.SpiritualAttunementMetrics = paladin.NewManaMetrics(core.ActionID{SpellID: 33776})

	paladin.RegisterAura(core.Aura{
		Label:    "Spiritual Attunement",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Damage > 0 {
				paladin.AddMana(sim, spellEffect.Damage*coeff, paladin.SpiritualAttunementMetrics, false)
			}
		},
	})
}
