package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (rogue *Rogue) registerShivSpell() {
	cost := 20.0
	if rogue.GetOHWeapon() != nil {
		cost = 20 + 10*rogue.GetOHWeapon().SwingSpeed
	}

	rogue.Shiv = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 5938},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagBuilder | core.SpellFlagCannotBeDodged,

		ResourceType: stats.Energy,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
			ModifyCast:  rogue.CastModifier,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeOHSpecial,
			DamageMultiplier: 1 +
				0.02*float64(rogue.Talents.FindWeakness) +
				core.TernaryFloat64(rogue.Talents.SurpriseAttacks, 0.1, 0),
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMeleeWeapon(core.OffHand, true, 0, 1+0.1*float64(rogue.Talents.DualWieldSpecialization), false),
			OutcomeApplier:   rogue.OutcomeFuncMeleeSpecialHitAndCrit(rogue.MeleeCritMultiplier(false, true)),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())

					switch rogue.Options.OhImbue {
					case proto.Rogue_Options_DeadlyPoison:
						rogue.DeadlyPoison.Cast(sim, spellEffect.Target)
					case proto.Rogue_Options_InstantPoison:
						rogue.InstantPoison[2].Cast(sim, spellEffect.Target)
					}
				}
			},
		}),
	})
}
