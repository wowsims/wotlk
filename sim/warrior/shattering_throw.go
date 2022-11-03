package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) RegisterShatteringThrowCD() {
	cost := 25 - float64(warrior.Talents.FocusedRage)

	ShatteringThrowSpell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 64382},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagChanneled | core.SpellFlagApplyArmorReduction,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     cost,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Minute * 5,
			},
			IgnoreHaste: true,
		},
		DamageMultiplier: 1,
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.5 * spell.MeleeAttackPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
			core.ShatteringThrowAura(target).Activate(sim)
			warrior.AutoAttacks.DelayMainhandMeleeUntil(sim, warrior.AutoAttacks.MainhandSwingAt+warrior.AutoAttacks.MainhandSwingSpeed())

			// To desync same speed weapon
			if warrior.AutoAttacks.MainhandSwingSpeed() == warrior.AutoAttacks.OffhandSwingSpeed() {
				warrior.AutoAttacks.DelayOffhandMeleeUntil(sim, warrior.AutoAttacks.OffhandSwingAt+warrior.AutoAttacks.OffhandSwingSpeed()+warrior.AutoAttacks.OffhandSwingSpeed()/2)
			} else {
				warrior.AutoAttacks.DelayOffhandMeleeUntil(sim, warrior.AutoAttacks.OffhandSwingAt+warrior.AutoAttacks.OffhandSwingSpeed())
			}
			warrior.disableHsCleaveUntil = sim.CurrentTime + spell.DefaultCast.CastTime + time.Millisecond*10
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: ShatteringThrowSpell,
		Type:  core.CooldownTypeDPS,
		ActivationFactory: func(sim *core.Simulation) core.CooldownActivation {
			return func(sim *core.Simulation, character *core.Character) {
				if !warrior.StanceMatches(BattleStance) {
					if !warrior.BattleStance.IsReady(sim) {
						return
					}
					warrior.BattleStance.Cast(sim, nil)
				}
				if warrior.CurrentRage() < cost {
					return
				}
				ShatteringThrowSpell.Cast(sim, character.CurrentTarget)
			}
		},
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return true
		},
	})
}
