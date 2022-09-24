package paladin

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (paladin *Paladin) registerSpiritualAttunement() {
	if paladin.Talents.SpiritualAttunement == 0 {
		return
	}

	// No longer baseline in WotLK, affected by talent points and glyphs. Ignoring the old set bonus here.
	SpiritualAttunementScalar := (0.05*float64(paladin.Talents.SpiritualAttunement) + core.TernaryFloat64(paladin.HasMajorGlyph(41096), 0.02, 0))

	paladin.SpiritualAttunementMetrics = paladin.NewManaMetrics(core.ActionID{SpellID: 33776})

	paladin.RegisterAura(core.Aura{
		Label:    "Spiritual Attunement",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			// We assume we were instantly healed for the damage.
			if spellEffect.Damage > 0 {
				paladin.AddMana(sim, spellEffect.Damage*SpiritualAttunementScalar, paladin.SpiritualAttunementMetrics, false)
			}
		},
		OnPeriodicDamageTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			// We assume we were instantly healed for the damage.
			if spellEffect.Damage > 0 {
				paladin.AddMana(sim, spellEffect.Damage*SpiritualAttunementScalar, paladin.SpiritualAttunementMetrics, false)
			}
		},
	})
}
