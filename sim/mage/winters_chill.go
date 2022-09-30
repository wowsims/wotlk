package mage

import (
	"github.com/wowsims/wotlk/sim/core"
)

// Winters Chill has a separate hit check from frostbolt, so it needs its own spell.
func (mage *Mage) registerWintersChillSpell() {
	if mage.Talents.WintersChill == 0 {
		return
	}

	wcAuras := make([]*core.Aura, mage.Env.GetNumTargets())
	for _, target := range mage.Env.Encounter.Targets {
		wcAuras[target.Index] = core.WintersChillAura(&target.Unit, 0)
	}

	mage.WintersChill = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 28595},
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       SpellFlagMage,

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, 0, spell.OutcomeMagicHit)
			spell.DealDamage(sim, &result)

			if result.Landed() {
				aura := wcAuras[target.Index]
				aura.Activate(sim)
				if aura.IsActive() {
					aura.AddStack(sim)
				}
			}
		},
	})
}

func (mage *Mage) applyWintersChill() {
	if mage.Talents.WintersChill == 0 {
		return
	}

	procChance := float64(mage.Talents.WintersChill) / 3

	mage.RegisterAura(core.Aura{
		Label:    "Winters Chill Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}

			if spell.SpellSchool == core.SpellSchoolFrost && spell != mage.WintersChill {
				if procChance == 1.0 || sim.RandomFloat("Winters Chill") < procChance {
					mage.WintersChill.Cast(sim, spellEffect.Target)
				}
			}
		},
	})
}
