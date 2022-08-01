package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (mage *Mage) registerManaGemsCD() {
	if mage.Consumes.DefaultConjured != proto.Conjured_ConjuredMageManaSapphire {
		return
	}

	actionID := core.ActionID{ItemID: 33312}
	manaMetrics := mage.NewManaMetrics(actionID)

	minManaEmeraldGain := 2340.0
	maxManaEmeraldGain := 2460.0
	minManaSapphireGain := 3330.0
	maxManaSapphireGain := 3500.0
	manaEmeraldGainRange := maxManaEmeraldGain - minManaEmeraldGain
	manaSapphireGainRange := maxManaSapphireGain - minManaSapphireGain

	var remainingManaGems int
	mage.RegisterResetEffect(func(sim *core.Simulation) {
		remainingManaGems = 6
	})

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if remainingManaGems <= 3 {
				// Mana Emerald: Restores 2340 to 2460 mana. (2 Min Cooldown)
				manaGain := minManaEmeraldGain + (sim.RandomFloat("Mana Gem") * manaEmeraldGainRange)
				mage.AddMana(sim, manaGain, manaMetrics, true)
			} else {
				// Mana Sapphire: Restores 3330 to 3500 mana. (2 Min Cooldown)
				manaGain := minManaSapphireGain + (sim.RandomFloat("Mana Gem") * manaSapphireGainRange)
				mage.AddMana(sim, manaGain, manaMetrics, true)
			}

			remainingManaGems--
			if remainingManaGems == 0 {
				// Disable this cooldown since we're out of emeralds.
				mage.DisableMajorCooldown(actionID)
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell:    spell,
		Priority: core.CooldownPriorityDefault,
		Type:     core.CooldownTypeMana,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return remainingManaGems != 0
		},
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			// Only pop if we have less than the max mana provided by the gem minus 1mp5 tick.
			totalRegen := character.ManaRegenPerSecondWhileCasting() * 5
			maxManaGain := maxManaEmeraldGain
			if remainingManaGems <= 3 {
				maxManaGain = maxManaEmeraldGain
			}
			if character.MaxMana()-(character.CurrentMana()+totalRegen) < maxManaGain {
				return false
			}

			return true
		},
	})
}
