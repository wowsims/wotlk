package rogue

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (rogue *Rogue) makeKillingSpreedWeaponSwingEffect(isMh bool) core.SpellEffect {
	var procMask core.ProcMask
	var baseMultiplier float64
	var hand core.Hand
	if isMh {
		procMask = core.ProcMaskMeleeMHSpecial
		baseMultiplier = 1
		hand = core.MainHand
	} else {
		procMask = core.ProcMaskMeleeOHSpecial
		baseMultiplier = 1 + 0.1*float64(rogue.Talents.DualWieldSpecialization)
		hand = core.OffHand
	}
	return core.SpellEffect{
		ProcMask: procMask,
		DamageMultiplier: (1 +
			0.02*float64(rogue.Talents.FindWeakness)) * baseMultiplier,
		ThreatMultiplier: 1,
		BaseDamage:       core.BaseDamageConfigMeleeWeapon(hand, true, 0, true),
		OutcomeApplier:   rogue.OutcomeFuncMeleeWeaponSpecialHitAndCrit(rogue.MeleeCritMultiplier(isMh, false)),
	}
}
func (rogue *Rogue) registerKillingSpreeSpell() {
	mhWeaponSwing := rogue.GetOrRegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 51690, Tag: 1},
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(rogue.makeKillingSpreedWeaponSwingEffect(true)),
	})
	ohWeaponSwing := rogue.GetOrRegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 51690, Tag: 2},
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(rogue.makeKillingSpreedWeaponSwingEffect(false)),
	})
	rogue.KillingSpreeAura = rogue.RegisterAura(core.Aura{
		Label:    "Killing Spree",
		ActionID: core.ActionID{SpellID: 51690},
		Duration: time.Second * 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.DamageDealtMultiplier *= 1.2
			// This first attack could/should? be implemented as an immediate tick
			// but there is currently an issue causing the periodic action
			// to only fire once when this flag is set
			mhWeaponSwing.Cast(sim, rogue.CurrentTarget)
			ohWeaponSwing.Cast(sim, rogue.CurrentTarget)
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:          time.Millisecond * 500,
				NumTicks:        4,
				TickImmediately: false,
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
			bladeFlurry := rogue.GetMajorCooldown(BladeFlurryActionID)
			if bladeFlurry != nil && bladeFlurry.IsReady(sim) {
				return false
			}
			if rogue.CurrentEnergy() > 60 || (rogue.CurrentEnergy() > 30 && rogue.AdrenalineRushAura.IsActive()) {
				return false
			}
			return true
		},
	})
}
