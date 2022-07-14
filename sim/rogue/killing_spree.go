package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (rogue *Rogue) makeKillingSpreeAttackSpell() *core.Spell {
	baseEffectMH := core.SpellEffect{
		ProcMask:         core.ProcMaskMeleeMHSpecial,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BaseDamage:       core.BaseDamageConfigMeleeWeapon(core.MainHand, true, 0, 1, true),
		OutcomeApplier:   rogue.OutcomeFuncMeleeWeaponSpecialHitAndCrit(rogue.MeleeCritMultiplier(true, true)),
	}
	baseEffectMH.Target = rogue.CurrentTarget
	baseEffectOH := core.SpellEffect{
		ProcMask:         core.ProcMaskMeleeOHSpecial,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BaseDamage:       core.BaseDamageConfigMeleeWeapon(core.OffHand, true, 0, 1+0.05*float64(rogue.Talents.DualWieldSpecialization), true),
		OutcomeApplier:   rogue.OutcomeFuncMeleeWeaponSpecialHitAndCrit(rogue.MeleeCritMultiplier(true, true)),
	}
	baseEffectOH.Target = rogue.CurrentTarget
	killingSpreeAttackSpell := rogue.GetOrRegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 51690, Tag: 1},
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDamageMultiple([]core.SpellEffect{baseEffectMH, baseEffectOH}),
	})
	return killingSpreeAttackSpell
}
func (rogue *Rogue) registerKillingSpreeSpell() {
	attackSpell := rogue.makeKillingSpreeAttackSpell()
	rogue.KillingSpreeAura = rogue.RegisterAura(core.Aura{
		Label:    "Killing Spree",
		ActionID: core.ActionID{SpellID: 51690},
		Duration: time.Second * 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.DamageDealtMultiplier *= 1.2
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:          time.Millisecond * 500,
				NumTicks:        5,
				TickImmediately: true,
				OnAction: func(s *core.Simulation) {
					attackSpell.Cast(sim, rogue.CurrentTarget)
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
				Duration: time.Minute*2 - core.TernaryDuration(rogue.HasGlyph(int32(proto.RogueMajorGlyph_GlyphOfKillingSpree)), time.Second*45, 0),
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
	})
}
