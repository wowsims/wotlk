package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerAspectOfTheDragonhawkSpell() {
	var impHawkAura *core.Aura
	const improvedHawkProcChance = 0.1
	if hunter.Talents.ImprovedAspectOfTheHawk > 0 {
		improvedHawkBonus := 1 +
			0.03*float64(hunter.Talents.ImprovedAspectOfTheHawk) +
			core.TernaryFloat64(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfTheHawk), 0.06, 0)

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
				if spell != hunter.AutoAttacks.RangedAuto {
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

	gronnstalkerMultiplier := core.TernaryFloat64(hunter.HasSetBonus(ItemSetGronnstalker, 2), 1.25, 1)
	baseManaRegen := 0.02 * hunter.BaseMana *
		core.TernaryFloat64(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfAspectOfTheViper), 1.1, 1) *
		gronnstalkerMultiplier
	manaPerRangedHit := baseManaRegen * hunter.AutoAttacks.Ranged.SwingSpeed * gronnstalkerMultiplier
	manaPerMHHit := baseManaRegen * hunter.AutoAttacks.MH.SwingSpeed * gronnstalkerMultiplier
	manaPerOHHit := baseManaRegen * hunter.AutoAttacks.OH.SwingSpeed * gronnstalkerMultiplier
	var tickPA *core.PendingAction

	hasCryptstalker4pc := hunter.HasSetBonus(ItemSetCryptstalkerBattlegear, 4)

	auraConfig := core.Aura{
		Label:    "Aspect of the Viper",
		ActionID: actionID,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= damagePenalty
			if hasCryptstalker4pc {
				aura.Unit.PseudoStats.RangedSpeedMultiplier *= 1.2
			}

			tickPA = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second * 3,
				OnAction: func(sim *core.Simulation) {
					hunter.AddMana(sim, 0.04*hunter.MaxMana(), hunter.AspectOfTheViper.ResourceMetrics, false)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= damagePenalty
			if hasCryptstalker4pc {
				aura.Unit.PseudoStats.RangedSpeedMultiplier /= 1.2
			}
			tickPA.Cancel(sim)
			tickPA = nil
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.ProcMask.Matches(core.ProcMaskRanged) {
				hunter.AddMana(sim, manaPerRangedHit, hunter.AspectOfTheViper.ResourceMetrics, false)
			} else if spellEffect.ProcMask.Matches(core.ProcMaskMeleeMH) {
				hunter.AddMana(sim, manaPerMHHit, hunter.AspectOfTheViper.ResourceMetrics, false)
			} else if spellEffect.ProcMask.Matches(core.ProcMaskMeleeOH) {
				hunter.AddMana(sim, manaPerOHHit, hunter.AspectOfTheViper.ResourceMetrics, false)
			}
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
