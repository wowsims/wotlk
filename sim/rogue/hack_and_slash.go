package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (rogue *Rogue) registerHackAndSlash(mask core.ProcMask) {
	if rogue.Talents.HackAndSlash < 1 || mask == core.ProcMaskUnknown {
		return
	}
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
				Flags:       core.SpellFlagMeleeMetrics,

				DamageMultiplier: rogue.AutoAttacks.MHConfig.DamageMultiplier,
				ThreatMultiplier: rogue.AutoAttacks.MHConfig.ThreatMultiplier,

				ApplyEffects: core.ApplyEffectFuncDirectDamage(rogue.AutoAttacks.MHEffect),
			})
		},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}
			if !spellEffect.ProcMask.Matches(mask) {
				return
			}
			if !icd.IsReady(sim) {
				return
			}
			if sim.RandomFloat("Sword Specialization") > procChance {
				return
			}
			icd.Use(sim)
			hackAndSlashSpell.Cast(sim, spellEffect.Target)
		},
	})
}
