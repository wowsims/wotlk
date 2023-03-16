package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (rogue *Rogue) registerShivSpell() {
	baseCost := 20.0
	if ohWeapon := rogue.GetOHWeapon(); ohWeapon != nil {
		baseCost = rogue.costModifier(20 + 10*ohWeapon.SwingSpeed)
	}

	rogue.Shiv = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 5938},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeOHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagBuilder,

		EnergyCost: core.EnergyCostOptions{
			Cost: baseCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: (1 +
			0.02*float64(rogue.Talents.FindWeakness) +
			core.TernaryFloat64(rogue.Talents.SurpriseAttacks, 0.1, 0)) * rogue.dwsMultiplier(),
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())

				switch rogue.Options.OhImbue {
				case proto.Rogue_Options_DeadlyPoison:
					rogue.DeadlyPoison.Cast(sim, target)
				case proto.Rogue_Options_InstantPoison:
					rogue.InstantPoison[ShivProc].Cast(sim, target)
				case proto.Rogue_Options_WoundPoison:
					rogue.WoundPoison[ShivProc].Cast(sim, target)
				}
			}
		},
	})
}
