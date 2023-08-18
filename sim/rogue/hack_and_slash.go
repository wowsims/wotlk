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
			hackAndSlashSpell = rogue.GetOrRegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 13964},
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskMeleeMHAuto,
				Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

				DamageMultiplier: rogue.AutoAttacks.MHConfig.DamageMultiplier,
				CritMultiplier:   rogue.MeleeCritMultiplier(false),
				ThreatMultiplier: rogue.AutoAttacks.MHConfig.ThreatMultiplier,

				ApplyEffects: rogue.AutoAttacks.MHConfig.ApplyEffects,
			})
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
