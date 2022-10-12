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
		ActionID:     core.ActionID{SpellID: 5938},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeOHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | SpellFlagBuilder,
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

		DamageMultiplier: (1 +
			0.02*float64(rogue.Talents.FindWeakness) +
			core.TernaryFloat64(rogue.Talents.SurpriseAttacks, 0.1, 0)) * (1 + 0.1*float64(rogue.Talents.DualWieldSpecialization)),
		CritMultiplier:   rogue.MeleeCritMultiplier(false, true),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.OffHand, true, 0, false),
			OutcomeApplier: rogue.OutcomeFuncMeleeSpecialNoBlockDodgeParryNoCrit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())

					switch rogue.Options.OhImbue {
					case proto.Rogue_Options_DeadlyPoison:
						rogue.DeadlyPoison.Cast(sim, spellEffect.Target)
					case proto.Rogue_Options_InstantPoison:
						rogue.InstantPoison[ShivProc].Cast(sim, spellEffect.Target)
					case proto.Rogue_Options_WoundPoison:
						rogue.WoundPoison[ShivProc].Cast(sim, spellEffect.Target)
					}
				}
			},
		}),
	})
}
