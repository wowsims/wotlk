package core

import (
	"time"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

// Registers all consume-related effects to the Agent.
// TODO: Classic Consumes
func applyConsumeEffects(agent Agent) {
	character := agent.GetCharacter()
	consumes := character.Consumes
	if consumes == nil {
		return
	}

	if consumes.Flask != proto.Flask_FlaskUnknown {
		switch consumes.Flask {
		case proto.Flask_FlaskOfDistilledWisdom:
			character.AddStats(stats.Stats{
				stats.Mana: 2000,
			})
		case proto.Flask_FlaskOfSupremePower:
			character.AddStats(stats.Stats{
				stats.SpellPower: 150,
			})
		case proto.Flask_FlaskOfTheTitans:
			character.AddStats(stats.Stats{
				stats.Health: 1200,
			})
		case proto.Flask_FlaskOfChromaticResistance:
			character.AddStats(stats.Stats{
				stats.ArcaneResistance: 25,
				stats.FireResistance:   25,
				stats.FrostResistance:  25,
				stats.NatureResistance: 25,
				stats.ShadowResistance: 25,
			})
		}
	}

	if character.HasMHWeapon() {
		addImbueStats(character, consumes.MainHandImbue)
	}
	if character.HasOHWeapon() {
		addImbueStats(character, consumes.OffHandImbue)
	}

	if consumes.Food != proto.Food_FoodUnknown {
		switch consumes.Food {
		case proto.Food_FoodGrilledSquid:
			character.AddStats(stats.Stats{
				stats.Agility: 10,
			})
		case proto.Food_FoodSmokedDesertDumpling:
			character.AddStats(stats.Stats{
				stats.Strength: 20,
			})
		case proto.Food_FoodNightfinSoup:
			character.AddStats(stats.Stats{
				stats.MP5: 8,
			})
		case proto.Food_FoodRunnTumTuberSurprise:
			character.AddStats(stats.Stats{
				stats.Intellect: 10,
			})
		case proto.Food_FoodDirgesKickChimaerokChops:
			character.AddStats(stats.Stats{
				stats.Stamina: 25,
			})
		case proto.Food_FoodBlessedSunfruitJuice:
			character.AddStats(stats.Stats{
				stats.Spirit: 10,
			})
		case proto.Food_FoodBlessSunfruit:
			character.AddStats(stats.Stats{
				stats.Strength: 10,
			})
		}
	}

	if consumes.AgilityElixir != proto.AgilityElixir_AgilityElixirUnknown {
		switch consumes.AgilityElixir {
		case proto.AgilityElixir_ElixirOfTheMongoose:
			character.AddStats(stats.Stats{
				stats.Agility:   25,
				stats.MeleeCrit: 2 * CritRatingPerCritChance,
			})
		case proto.AgilityElixir_ElixirOfGreaterAgility:
			character.AddStats(stats.Stats{
				stats.Agility: 25,
			})
		}
	}

	if consumes.StrengthBuff != proto.StrengthBuff_StrengthBuffUnknown {
		switch consumes.StrengthBuff {
		case proto.StrengthBuff_JujuPower:
			character.AddStats(stats.Stats{
				stats.Strength: 30,
			})
		case proto.StrengthBuff_ElixirOfGiants:
			character.AddStats(stats.Stats{
				stats.Strength: 25,
			})
		}
	}

	if consumes.SpellPowerBuff != proto.SpellPowerBuff_SpellPowerBuffUnknown {
		switch consumes.SpellPowerBuff {
		case proto.SpellPowerBuff_ArcaneElixir:
			character.AddStats(stats.Stats{
				stats.SpellPower: 20,
			})
		case proto.SpellPowerBuff_GreaterArcaneElixir:
			character.AddStats(stats.Stats{
				stats.SpellPower: 35,
			})
		}
	}

	if consumes.FirePowerBuff != proto.FirePowerBuff_FirePowerBuffUnknown {
		switch consumes.FirePowerBuff {
		case proto.FirePowerBuff_ElixirOfFirepower:
			character.AddStats(stats.Stats{
				stats.FirePower: 10,
			})
		case proto.FirePowerBuff_ElixirOfGreaterFirepower:
			character.AddStats(stats.Stats{
				stats.FirePower: 40,
			})
		}
	}

	if consumes.ShadowPowerBuff {
		character.AddStats(stats.Stats{
			stats.ShadowPower: 40,
		})
	}

	if consumes.FrostPowerBuff {
		character.AddStats(stats.Stats{
			stats.FrostPower: 15,
		})
	}

	// registerPotionCD(agent, consumes)
	// registerConjuredCD(agent, consumes)
	// registerExplosivesCD(agent, consumes)
}
func addImbueStats(character *Character, imbue proto.WeaponImbue) {
	if imbue != proto.WeaponImbue_WeaponImbueUnknown {
		switch imbue {
		case proto.WeaponImbue_BrillianWizardOil:
			character.AddStats(stats.Stats{
				stats.SpellPower: 36,
				stats.SpellCrit:  1 * SpellCritRatingPerCritChance,
			})
		case proto.WeaponImbue_BrilliantManaOil:
			character.AddStats(stats.Stats{
				stats.MP5:     5,
				stats.Healing: 25,
			})
		// TODO: Classic
		// case proto.WeaponImbue_DenseSharpeningStone:
		// 	character.AddStats(stats.Stats{
		// 		stats.WeaponDamage??: 5,
		// 	})
		case proto.WeaponImbue_ElementalSharpeningStone:
			character.AddStats(stats.Stats{
				stats.MeleeCrit: 2 * CritRatingPerCritChance,
			})
		case proto.WeaponImbue_BlackfathomManaOil:
			character.AddStats(stats.Stats{
				stats.MP5:      12,
				stats.SpellHit: 2 * SpellHitRatingPerHitChance,
			})
		case proto.WeaponImbue_DenseSharpeningStone:
			character.AddStats(stats.Stats{
				stats.MeleeHit: 2 * MeleeHitRatingPerHitChance,
			})
		case proto.WeaponImbue_WildStrikes:
			buffActionID := ActionID{SpellID: 407975}
			statDep := character.NewDynamicMultiplyStat(stats.AttackPower, 1.2)

			wsBuffAura := character.GetOrRegisterAura(Aura{
				Label:     "Wild Strikes Buff",
				ActionID:  buffActionID,
				Duration:  time.Millisecond * 1500,
				MaxStacks: 2,
				OnGain: func(aura *Aura, sim *Simulation) {
					aura.Unit.EnableDynamicStatDep(sim, statDep)
				},
				OnExpire: func(aura *Aura, sim *Simulation) {
					aura.Unit.DisableDynamicStatDep(sim, statDep)
				},
			})

			var wsSpell *Spell
			icd := Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Millisecond * 1500,
			}

			MakePermanent(character.GetOrRegisterAura(Aura{
				Label: "Wild Strikes",
				OnInit: func(aura *Aura, sim *Simulation) {
					mhConfig := aura.Unit.AutoAttacks.MHConfig()
					wsSpell = character.GetOrRegisterSpell(SpellConfig{
						ActionID:         buffActionID, // temporary buff ("Windfury Attack") spell id
						SpellSchool:      mhConfig.SpellSchool,
						ProcMask:         mhConfig.ProcMask,
						Flags:            mhConfig.Flags,
						DamageMultiplier: mhConfig.DamageMultiplier,
						CritMultiplier:   mhConfig.CritMultiplier,
						ThreatMultiplier: mhConfig.ThreatMultiplier,
						ApplyEffects:     mhConfig.ApplyEffects,
					})
				},
				OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
					if spell.ProcMask.Matches(ProcMaskSuppressedExtraAttackAura) {
						return
					}
					if !result.Landed() || !spell.ProcMask.Matches(ProcMaskMelee) {
						return
					}

					if wsBuffAura.IsActive() && spell.ProcMask.Matches(ProcMaskMeleeWhiteHit) {
						wsBuffAura.RemoveStack(sim)
					}

					if !icd.IsReady(sim) {
						return
					}

					if sim.RandomFloat("Wild Strikes") > 0.2 {
						return
					}

					wsBuffAura.Activate(sim)
					wsBuffAura.SetStacks(sim, 2)
					icd.Use(sim)

					StartDelayedAction(sim, DelayedActionOptions{
						DoAt:     sim.CurrentTime + time.Millisecond*10,
						Priority: ActionPriorityAuto,
						OnAction: func(sim *Simulation) {
							wsSpell.Cast(sim, result.Target)
						},
					})
				},
			}))
		}
	}
}
