package mage

import (
	"github.com/wowsims/tbc/sim/core"
)

// Winters Chill has a separate hit check from frostbolt, so it needs its own spell.
func (mage *Mage) registerWintersChillSpell() {
	effect := core.SpellEffect{
		ProcMask:            core.ProcMaskEmpty,
		BonusSpellHitRating: float64(mage.Talents.ElementalPrecision) * 1 * core.SpellHitRatingPerHitChance,
		ThreatMultiplier:    1,
		OutcomeApplier:      mage.OutcomeFuncMagicHit(),
	}

	if mage.Talents.WintersChill > 0 {
		wcAura := mage.CurrentTarget.GetAura(core.WintersChillAuraLabel)
		if wcAura == nil {
			wcAura = core.WintersChillAura(mage.CurrentTarget, 0)
		}

		effect.OnSpellHitDealt = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				wcAura.Activate(sim)
				wcAura.AddStack(sim)
			}
		}
	}

	mage.WintersChill = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 28595},
		SpellSchool: core.SpellSchoolFrost,
		Flags:       SpellFlagMage,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (mage *Mage) applyWintersChill() {
	if mage.Talents.WintersChill == 0 {
		return
	}

	procChance := float64(mage.Talents.WintersChill) / 5.0

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
				if procChance != 1.0 && sim.RandomFloat("Winters Chill") > procChance {
					return
				}

				mage.WintersChill.Cast(sim, spellEffect.Target)
			}
		},
	})
}
