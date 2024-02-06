package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (rogue *Rogue) registerSwordSpecialization(mask core.ProcMask) {
	// https://wotlk.wowhead.com/spell=13964/sword-specialization, proc mask = 20.
	var swordSpecSpell *core.Spell
	icd := core.Cooldown{
		Timer:    rogue.NewTimer(),
		Duration: time.Millisecond * 500,
	}
	procChance := 0.01 * float64(rogue.Talents.SwordSpecialization)

	rogue.RegisterAura(core.Aura{
		Label:    "Sword Specialization",
		Duration: core.NeverExpires,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			config := *rogue.AutoAttacks.MHConfig()
			config.ActionID = core.ActionID{SpellID: 13964}
			swordSpecSpell = rogue.GetOrRegisterSpell(config)
		},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}
			if !spell.ProcMask.Matches(mask) {
				return
			}
			if !icd.IsReady(sim) {
				return
			}
			if sim.RandomFloat("Sword Specialization") < procChance {
				icd.Use(sim)
				swordSpecSpell.Cast(sim, result.Target)
			}
		},
	})
}
