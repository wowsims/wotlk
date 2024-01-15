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
				aura.Unit.MultiplyRangedSpeed(sim, improvedHawkBonus)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.MultiplyRangedSpeed(sim, 1/improvedHawkBonus)
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
			if hunter.Talents.AspectMastery {
				aura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.DamageTakenMultiplier *= 0.95
				})
				aura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.DamageTakenMultiplier /= 0.95
				})
			}

			aura.OnSpellHitDealt = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell != hunter.AutoAttacks.RangedAuto() {
					return
				}

				if impHawkAura != nil && sim.RandomFloat("Imp Aspect of the Hawk") < improvedHawkProcChance {
					impHawkAura.Activate(sim)
				}
			}
		})
	hunter.applySharedAspectConfig(true, hunter.AspectOfTheDragonhawkAura)

	hunter.AspectOfTheDragonhawk = hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.AspectOfTheDragonhawkAura.Activate(sim)
		},
	})
}

func (hunter *Hunter) registerAspectOfTheViperSpell() {
	actionID := core.ActionID{SpellID: 34074}

	damagePenalty := core.TernaryFloat64(hunter.Talents.AspectMastery, 0.6, 0.5)

	baseManaRegenMultiplier := 0.01 *
		core.TernaryFloat64(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfAspectOfTheViper), 1.1, 1) *
		core.TernaryFloat64(hunter.HasSetBonus(ItemSetGronnstalker, 2), 1.25, 1)
	manaPerRangedHitMultiplier := baseManaRegenMultiplier * hunter.AutoAttacks.Ranged().SwingSpeed
	manaPerMHHitMultiplier := baseManaRegenMultiplier * hunter.AutoAttacks.MH().SwingSpeed
	manaPerOHHitMultiplier := baseManaRegenMultiplier * hunter.AutoAttacks.OH().SwingSpeed
	var tickPA *core.PendingAction

	hasCryptstalker4pc := hunter.HasSetBonus(ItemSetCryptstalkerBattlegear, 4)

	auraConfig := core.Aura{
		Label:    "Aspect of the Viper",
		ActionID: actionID,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= damagePenalty
			if hasCryptstalker4pc {
				aura.Unit.MultiplyRangedSpeed(sim, 1.2)
			}

			tickPA = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second * 3,
				OnAction: func(sim *core.Simulation) {
					hunter.AddMana(sim, 0.04*hunter.MaxMana(), hunter.AspectOfTheViper.ResourceMetrics)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= damagePenalty
			if hasCryptstalker4pc {
				aura.Unit.MultiplyRangedSpeed(sim, 1/1.2)
			}
			tickPA.Cancel(sim)
			tickPA = nil
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskRanged) {
				hunter.AddMana(sim, manaPerRangedHitMultiplier*hunter.MaxMana(), hunter.AspectOfTheViper.ResourceMetrics)
			} else if spell.ProcMask.Matches(core.ProcMaskMeleeMH) {
				hunter.AddMana(sim, manaPerMHHitMultiplier*hunter.MaxMana(), hunter.AspectOfTheViper.ResourceMetrics)
			} else if spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
				hunter.AddMana(sim, manaPerOHHitMultiplier*hunter.MaxMana(), hunter.AspectOfTheViper.ResourceMetrics)
			}
		},
	}
	hunter.AspectOfTheViperAura = hunter.RegisterAura(auraConfig)
	hunter.applySharedAspectConfig(false, hunter.AspectOfTheViperAura)

	hunter.AspectOfTheViper = hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.AspectOfTheViperAura.Activate(sim)
		},
	})
	hunter.AspectOfTheViper.ResourceMetrics = hunter.NewManaMetrics(hunter.AspectOfTheViper.ActionID)
}

func (hunter *Hunter) applySharedAspectConfig(isHawk bool, aura *core.Aura) {
	if isHawk {
		aura.OnReset = func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		}
	}

	aura.Duration = core.NeverExpires
	aura.NewExclusiveEffect("Aspect", true, core.ExclusiveEffect{})
}
