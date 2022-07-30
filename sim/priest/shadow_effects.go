package priest

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (priest *Priest) ApplyMisery(sim *core.Simulation, target *core.Unit) {
	if priest.MiseryAura != nil {
		priest.MiseryAura.Activate(sim)
	}
}

//func (priest *Priest) ApplyShadowWeaving(sim *core.Simulation, target *core.Unit) {
//if priest.ShadowWeavingAura == nil {
//	return
//}

//if priest.Talents.ShadowWeaving < 5 && sim.RandomFloat("Shadow Weaving") > 0.2*float64(priest.Talents.ShadowWeaving) {
//	return
//}

//	priest.ShadowWeavingAura.Activate(sim)
//if priest.ShadowWeavingAura.IsActive() {
//	priest.ShadowWeavingAura.AddStack(sim)
//}
//}

func (priest *Priest) ApplyShadowOnHitEffects() {
	// This is a combined aura for all priest major on hit effects.
	//  Shadow Weaving, Vampiric Touch, and Misery
	priest.RegisterAura(core.Aura{
		Label:    "Priest Shadow Effects",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		//OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) { // Needs to be replaced by replensishment (when VT is active MB now gives mana return)
		//	priest.applySWforMF(sim, spellEffect.Damage)
		//	},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}
			//priest.ApplyShadowWeaving(sim, spellEffect.Target)
			// priest.ApplyVampiricTouchManaReturn(sim, spellEffect.Damage) // Needs to be replaced by replensishment (when VT is active MB now gives mana return)

			if spell == priest.ShadowWordPain || spell == priest.VampiricTouch || spell.ActionID.SpellID == priest.MindFlay[1].ActionID.SpellID {
				priest.ApplyMisery(sim, spellEffect.Target)
			}
		},
	})
}
