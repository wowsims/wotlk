package rogue

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (rogue *Rogue) registerKillingSpreeSpell() {
	mhWeaponSwing := rogue.GetOrRegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 51690, Tag: 1}, // actual spellID is 57841
		SpellSchool:      core.SpellSchoolPhysical,
		ProcMask:         core.ProcMaskMeleeMHSpecial,
		Flags:            core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		DamageMultiplier: 1 + 0.02*float64(rogue.Talents.FindWeakness),
		CritMultiplier:   rogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})
	ohWeaponSwing := rogue.GetOrRegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 51690, Tag: 2}, // actual spellID is 57842
		SpellSchool:      core.SpellSchoolPhysical,
		ProcMask:         core.ProcMaskMeleeOHSpecial,
		Flags:            core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		DamageMultiplier: (1 + 0.02*float64(rogue.Talents.FindWeakness)) * rogue.dwsMultiplier(),
		CritMultiplier:   rogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0 +
				spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})
	rogue.KillingSpreeAura = rogue.RegisterAura(core.Aura{
		Label:    "Killing Spree",
		ActionID: core.ActionID{SpellID: 51690},
		Duration: time.Second*2 + 1,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.DamageDealtMultiplier *= 1.2
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:          time.Millisecond * 500,
				NumTicks:        5,
				TickImmediately: true,
				OnAction: func(s *core.Simulation) {
					targetCount := sim.GetNumTargets()
					target := rogue.CurrentTarget
					if targetCount > 1 {
						newUnitIndex := int32(math.Ceil(float64(targetCount)*sim.RandomFloat("Killing Spree"))) - 1
						target = sim.GetTargetUnit(newUnitIndex)
					}
					mhWeaponSwing.Cast(sim, target)
					ohWeaponSwing.Cast(sim, target)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.DamageDealtMultiplier /= 1.2
		},
	})
	killingSpreeSpell := rogue.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 51690},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Minute*2 - core.TernaryDuration(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfKillingSpree), time.Second*45, 0),
			},
		},

		ApplyEffects: func(sim *core.Simulation, u *core.Unit, s2 *core.Spell) {
			rogue.KillingSpreeAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    killingSpreeSpell,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityLow,
		ShouldActivate: func(sim *core.Simulation, c *core.Character) bool {
			if bf := rogue.GetMajorCooldown(BladeFlurryActionID); bf != nil && bf.IsReady(sim) {
				return false
			}
			if rogue.CurrentEnergy() > 60 || (rogue.CurrentEnergy() > 30 && rogue.AdrenalineRushAura.IsActive()) {
				return false
			}
			return true
		},
	})
}
