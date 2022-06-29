package warrior

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warrior *Warrior) registerHeroicStrikeSpell() {
	cost := 15.0 - float64(warrior.Talents.ImprovedHeroicStrike) - float64(warrior.Talents.FocusedRage)
	refundAmount := cost * 0.8

	warrior.HeroicStrikeOrCleave = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 29707},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeMHSpecial,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			FlatThreatBonus:  194,

			BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 176, 1, true),
			OutcomeApplier: warrior.OutcomeFuncMeleeWeaponSpecialHitAndCrit(warrior.critMultiplier(true)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
				}
			},
		}),
	})
}

func (warrior *Warrior) registerCleaveSpell() {
	cost := 20.0 - float64(warrior.Talents.FocusedRage)

	flatDamageBonus := 70 * (1 + 0.4*float64(warrior.Talents.ImprovedCleave))
	baseEffect := core.SpellEffect{
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeMHSpecial,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		FlatThreatBonus:  125,

		BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, false, flatDamageBonus, 1, true),
		OutcomeApplier: warrior.OutcomeFuncMeleeWeaponSpecialHitAndCrit(warrior.critMultiplier(true)),
	}

	numHits := core.MinInt32(2, warrior.Env.GetNumTargets())
	effects := make([]core.SpellEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		effects = append(effects, baseEffect)
		effects[i].Target = warrior.Env.GetTargetUnit(i)
	}

	warrior.HeroicStrikeOrCleave = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 25231},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})
}

func (warrior *Warrior) QueueHSOrCleave(sim *core.Simulation) {
	if warrior.CurrentRage() < warrior.HeroicStrikeOrCleave.DefaultCast.Cost {
		panic("Not enough rage for HS")
	}
	if warrior.HSOrCleaveQueueAura.IsActive() {
		return
	}
	warrior.HSOrCleaveQueueAura.Activate(sim)
	warrior.PseudoStats.DisableDWMissPenalty = true
}

func (warrior *Warrior) DequeueHSOrCleave(sim *core.Simulation) {
	warrior.HSOrCleaveQueueAura.Deactivate(sim)
	warrior.PseudoStats.DisableDWMissPenalty = false
}

// Returns true if the regular melee swing should be used, false otherwise.
func (warrior *Warrior) TryHSOrCleave(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
	if !warrior.HSOrCleaveQueueAura.IsActive() {
		return nil
	}

	if warrior.CurrentRage() < warrior.HeroicStrikeOrCleave.DefaultCast.Cost {
		warrior.DequeueHSOrCleave(sim)
		return nil
	} else if warrior.CurrentRage() < warrior.HSRageThreshold {
		if mhSwingSpell == warrior.AutoAttacks.MHAuto {
			warrior.DequeueHSOrCleave(sim)
			return nil
		}
	}

	warrior.DequeueHSOrCleave(sim)
	return warrior.HeroicStrikeOrCleave
}

func (warrior *Warrior) ShouldQueueHSOrCleave(sim *core.Simulation) bool {
	return warrior.CurrentRage() >= warrior.HSRageThreshold
}

func (warrior *Warrior) RegisterHSOrCleave(useCleave bool, rageThreshold float64) {
	if useCleave {
		warrior.registerCleaveSpell()
	} else {
		warrior.registerHeroicStrikeSpell()
	}

	warrior.HSOrCleaveQueueAura = warrior.RegisterAura(core.Aura{
		Label:    "HS Queue Aura",
		ActionID: warrior.HeroicStrikeOrCleave.ActionID,
		Duration: core.NeverExpires,
	})

	warrior.HSRageThreshold = core.MaxFloat(warrior.HeroicStrikeOrCleave.DefaultCast.Cost, rageThreshold)
}
