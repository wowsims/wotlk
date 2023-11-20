package core

import (
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
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

	if consumes.WeaponBuff != proto.WeaponBuff_WeaponBuffUnknown {
		switch consumes.WeaponBuff {
		case proto.WeaponBuff_BrillianWizardOil:
			character.AddStats(stats.Stats{
				stats.SpellPower: 36,
				stats.SpellCrit:  1 * SpellCritRatingPerCritChance,
			})
		case proto.WeaponBuff_BrilliantManaOil:
			character.AddStats(stats.Stats{
				stats.MP5:     5,
				stats.Healing: 25,
			})
		// TODO: Classic
		// case proto.WeaponBuff_DenseSharpeningStone:
		// 	character.AddStats(stats.Stats{
		// 		stats.WeaponDamage??: 5,
		// 	})
		case proto.WeaponBuff_ElementalSharpeningStone:
			character.AddStats(stats.Stats{
				stats.MeleeCrit: 2 * CritRatingPerCritChance,
			})
		}
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
				stats.Agility: 25,
			})
		}
	}

	if consumes.SpellPowerBuff {
		character.AddStats(stats.Stats{
			stats.SpellPower: 35,
		})
	}

	if consumes.ShadowPowerBuff {
		character.AddStats(stats.Stats{
			stats.ShadowPower: 40,
		})
	}

	if consumes.FirePowerBuff {
		character.AddStats(stats.Stats{
			stats.FirePower: 40,
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
