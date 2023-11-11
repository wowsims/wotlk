package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (rogue *Rogue) registerHackAndSlash(mask core.ProcMask) {
	// https://wotlk.wowhead.com/spell=13964/sword-specialization, proc mask = 20.
	var hackAndSlashSpell *core.Spell
	icd := core.Cooldown{
		Timer:    rogue.NewTimer(),
		Duration: time.Millisecond * 500,
	}
	procChance := 0.01 * float64(rogue.Talents.HackAndSlash)

	rogue.RegisterAura(core.Aura{
		Label:    "Hack and Slash",
		Duration: core.NeverExpires,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			config := *rogue.AutoAttacks.MHConfig()
			config.ActionID = core.ActionID{SpellID: 13964}
			hackAndSlashSpell = rogue.GetOrRegisterSpell(config)
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
			if sim.RandomFloat("Hack and Slash") < procChance {
				icd.Use(sim)
				hackAndSlashSpell.Cast(sim, result.Target)
			}
		},
	})
}
