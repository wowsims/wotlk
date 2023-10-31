package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warrior *Warrior) registerWhirlwindSpell() {
	actionID := core.ActionID{SpellID: 1680}
	numHits := min(4, warrior.Env.GetNumTargets())
	results := make([]*core.SpellResult, numHits)

	if warrior.AutoAttacks.IsDualWielding && warrior.GetOHWeapon().WeaponType != proto.WeaponType_WeaponTypeStaff &&
		warrior.GetOHWeapon().WeaponType != proto.WeaponType_WeaponTypePolearm {
		warrior.WhirlwindOH = warrior.RegisterSpell(core.SpellConfig{
			ActionID:    actionID.WithTag(1),
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty, // whirlwind offhand hits usually don't proc auras
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete | SpellFlagWhirlwindOH,

			DamageMultiplier: 1 *
				(1 + 0.02*float64(warrior.Talents.UnendingFury) + 0.1*float64(warrior.Talents.ImprovedWhirlwind)) *
				(1 + 0.05*float64(warrior.Talents.DualWieldSpecialization)),
			CritMultiplier:   warrior.critMultiplier(oh),
			ThreatMultiplier: 1.25,
		})
	}

	warrior.Whirlwind = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBloodsurge | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost: 25 - float64(warrior.Talents.FocusedRage),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: core.TernaryDuration(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfWhirlwind), time.Second*8, time.Second*10),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BerserkerStance)
		},

		DamageMultiplier: 1 *
			(1 + 0.02*float64(warrior.Talents.UnendingFury) + 0.1*float64(warrior.Talents.ImprovedWhirlwind)),
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1.25,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := 0 +
					spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
				results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			curTarget = target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				spell.DealDamage(sim, results[hitIndex])
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			if warrior.WhirlwindOH != nil {
				curTarget = target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					baseDamage := 0 +
						spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
						spell.BonusWeaponDamage()
					results[hitIndex] = warrior.WhirlwindOH.CalcDamage(sim, curTarget, baseDamage, warrior.WhirlwindOH.OutcomeMeleeWeaponSpecialHitAndCrit)

					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}

				curTarget = target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					warrior.WhirlwindOH.DealDamage(sim, results[hitIndex])
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
			}
		},
	})
}

func (warrior *Warrior) CanWhirlwind(sim *core.Simulation) bool {
	return warrior.StanceMatches(BerserkerStance) && warrior.CurrentRage() >= warrior.Whirlwind.DefaultCast.Cost && warrior.Whirlwind.IsReady(sim)
}
