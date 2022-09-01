package mage

import (
	"github.com/wowsims/wotlk/sim/core"
)

// Winters Chill has a separate hit check from frostbolt, so it needs its own spell.
func (mage *Mage) registerWintersChillSpell() {
	if mage.Talents.WintersChill == 0 {
		return
	}

	wcAura := mage.CurrentTarget.GetAura(core.WintersChillAuraLabel)
	if wcAura == nil {
		wcAura = core.WintersChillAura(mage.CurrentTarget, 0)
	}

	mage.WintersChill = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 28595},
		SpellSchool: core.SpellSchoolFrost,
		Flags:       SpellFlagMage,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			BonusHitRating:   0,
			ThreatMultiplier: 1,
			OutcomeApplier:   mage.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					wcAura.Activate(sim)
					if wcAura.IsActive() {
						wcAura.AddStack(sim)
					}
				}
			},
		}),
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
