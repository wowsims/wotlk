package hunter

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (hunter *Hunter) registerAspectOfTheHawkSpell() {
	var impHawkAura *core.Aura
	const improvedHawkProcChance = 0.1
	if hunter.Talents.ImprovedAspectOfTheHawk > 0 {
		improvedHawkBonus := 1 + 0.03*float64(hunter.Talents.ImprovedAspectOfTheHawk)
		impHawkAura = hunter.GetOrRegisterAura(core.Aura{
			Label:    "Improved Aspect of the Hawk",
			ActionID: core.ActionID{SpellID: 19556},
			Duration: time.Second * 12,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.RangedSpeedMultiplier *= improvedHawkBonus
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.RangedSpeedMultiplier /= improvedHawkBonus
			},
		})
	}

	actionID := core.ActionID{SpellID: 27044}
	hunter.AspectOfTheHawkAura = hunter.NewTemporaryStatsAuraWrapped("Aspect of the Hawk", actionID, stats.Stats{stats.RangedAttackPower: 155}, core.NeverExpires, func(aura *core.Aura) {
		hunter.applySharedAspectConfig(true, aura)
		aura.OnSpellHitDealt = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.ProcMask.Matches(core.ProcMaskRangedAuto) {
				return
			}

			if impHawkAura != nil && sim.RandomFloat("Imp Aspect of the Hawk") < improvedHawkProcChance {
				impHawkAura.Activate(sim)
			}
		}
	})

	baseCost := 140.0
	hunter.AspectOfTheHawk = hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.AspectOfTheHawkAura.Activate(sim)
		},
	})
}

func (hunter *Hunter) registerAspectOfTheViperSpell() {
	actionID := core.ActionID{SpellID: 34074}
	auraConfig := core.Aura{
		Label:    "Aspect of the Viper",
		ActionID: actionID,
	}
	hunter.applySharedAspectConfig(false, &auraConfig)
	hunter.AspectOfTheViperAura = hunter.RegisterAura(auraConfig)

	baseCost := 40.0
	hunter.AspectOfTheViper = hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.AspectOfTheViperAura.Activate(sim)
		},
	})
}

func (hunter *Hunter) applySharedAspectConfig(isHawk bool, aura *core.Aura) {
	if isHawk != (hunter.Rotation.ViperStartManaPercent >= 1) {
		aura.OnReset = func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		}
	}

	aura.Tag = "Aspect"
	aura.Priority = 1
	aura.Duration = core.NeverExpires

	oldOnGain := aura.OnGain
	if oldOnGain == nil {
		aura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
			hunter.currentAspect = aura
		}
	} else {
		aura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
			oldOnGain(aura, sim)
			hunter.currentAspect = aura
		}
	}
}
