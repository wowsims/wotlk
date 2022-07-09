package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerAspectOfTheDragonhawkSpell() {
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

	actionID := core.ActionID{SpellID: 61847}
	hunter.AspectOfTheDragonhawkAura = hunter.NewTemporaryStatsAuraWrapped(
		"Aspect of the Dragonhawk",
		actionID,
		stats.Stats{
			stats.RangedAttackPower: core.TernaryFloat64(hunter.Talents.AspectMastery, 390, 300),
			stats.Dodge:             (18 + 2*float64(hunter.Talents.ImprovedAspectOfTheMonkey)) * core.DodgeRatingPerDodgeChance,
		},
		core.NeverExpires,
		func(aura *core.Aura) {
			hunter.applySharedAspectConfig(true, aura)

			if hunter.Talents.AspectMastery {
				oldOnGain := aura.OnGain
				aura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
					oldOnGain(aura, sim)
					aura.Unit.PseudoStats.DamageTakenMultiplier *= 0.95
				}
				oldOnExpire := aura.OnExpire
				aura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
					oldOnExpire(aura, sim)
					aura.Unit.PseudoStats.DamageTakenMultiplier /= 0.95
				}
			}

			aura.OnSpellHitDealt = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.ProcMask.Matches(core.ProcMaskRangedAuto) {
					return
				}

				if impHawkAura != nil && sim.RandomFloat("Imp Aspect of the Hawk") < improvedHawkProcChance {
					impHawkAura.Activate(sim)
				}
			}
		})

	hunter.AspectOfTheDragonhawk = hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.AspectOfTheDragonhawkAura.Activate(sim)
		},
	})
}

func (hunter *Hunter) registerAspectOfTheViperSpell() {
	actionID := core.ActionID{SpellID: 34074}

	damagePenalty := core.TernaryFloat64(hunter.Talents.AspectMastery, 0.6, 0.5)

	auraConfig := core.Aura{
		Label:    "Aspect of the Viper",
		ActionID: actionID,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= damagePenalty
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= damagePenalty
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
				return
			}

			// TODO: Mana gain
		},
	}
	hunter.applySharedAspectConfig(false, &auraConfig)
	hunter.AspectOfTheViperAura = hunter.RegisterAura(auraConfig)

	hunter.AspectOfTheViper = hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.AspectOfTheViperAura.Activate(sim)
		},
	})
	hunter.AspectOfTheViper.ResourceMetrics = hunter.NewManaMetrics(hunter.AspectOfTheViper.ActionID)
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
