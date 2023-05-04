package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

const FanOfKnivesSpellID int32 = 51723

func (rogue *Rogue) makeFanOfKnivesWeaponHitSpell(isMH bool) *core.Spell {
	var procMask core.ProcMask
	var weaponMultiplier float64
	if isMH {
		weaponMultiplier = core.TernaryFloat64(rogue.HasDagger(core.MainHand), 1.05, 0.7)
		procMask = core.ProcMaskMeleeMHSpecial
	} else {
		weaponMultiplier = core.TernaryFloat64(rogue.HasDagger(core.OffHand), 1.05, 0.7)
		weaponMultiplier *= rogue.dwsMultiplier()
		procMask = core.ProcMaskMeleeOHSpecial
	}

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: FanOfKnivesSpellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagColdBlooded,

		DamageMultiplier: weaponMultiplier * (1 +
			0.02*float64(rogue.Talents.FindWeakness) +
			core.TernaryFloat64(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfFanOfKnives), 0.2, 0.0)),
		CritMultiplier:   rogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,
	})
}

func (rogue *Rogue) registerFanOfKnives() {
	mhSpell := rogue.makeFanOfKnivesWeaponHitSpell(true)
	ohSpell := rogue.makeFanOfKnivesWeaponHitSpell(false)
	results := make([]*core.SpellResult, len(rogue.Env.Encounter.TargetUnits))

	rogue.FanOfKnives = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: FanOfKnivesSpellID},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		EnergyCost: core.EnergyCostOptions{
			Cost: 50,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			// Calc and apply all OH hits first, because MH hits can benefit from an OH felstriker proc.
			for i, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := ohSpell.Unit.OHWeaponDamage(sim, ohSpell.MeleeAttackPower())
				baseDamage *= sim.Encounter.AOECapMultiplier()
				results[i] = ohSpell.CalcDamage(sim, aoeTarget, baseDamage, ohSpell.OutcomeMeleeWeaponSpecialHitAndCrit)
			}
			for i := range sim.Encounter.TargetUnits {
				ohSpell.DealDamage(sim, results[i])
			}

			for i, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := mhSpell.Unit.MHWeaponDamage(sim, mhSpell.MeleeAttackPower())
				baseDamage *= sim.Encounter.AOECapMultiplier()
				results[i] = mhSpell.CalcDamage(sim, aoeTarget, baseDamage, mhSpell.OutcomeMeleeWeaponSpecialHitAndCrit)
			}
			for i := range sim.Encounter.TargetUnits {
				mhSpell.DealDamage(sim, results[i])
			}
		},
	})
}
