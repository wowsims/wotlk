package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) registerShivSpell() {
	if !rogue.HasRune(proto.RogueRune_RuneShiv) {
		return
	}

	baseCost := 20.0
	if ohWeapon := rogue.GetOHWeapon(); ohWeapon != nil {
		baseCost = rogue.costModifier(baseCost + 10*ohWeapon.SwingSpeed)
	}

	rogue.Shiv = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 424799},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeOHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagBuilder | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   baseCost,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: (1 * rogue.dwsMultiplier()),
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			// TODO: Cannot Miss
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

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
