package mage

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (mage *Mage) registerManaGemsCD() {
	if mage.Consumes.DefaultConjured != proto.Conjured_ConjuredMageManaEmerald {
		return
	}

	manaMetrics := mage.NewManaMetrics(core.MageManaGemMCDActionID)

	var serpentCoilAura *core.Aura
	if mage.HasTrinketEquipped(SerpentCoilBraidID) {
		serpentCoilAura = mage.NewTemporaryStatsAura("Serpent Coil Braid", core.ActionID{ItemID: SerpentCoilBraidID}, stats.Stats{stats.SpellPower: 225}, time.Second*15)
	}

	manaMultiplier := 1.0
	minManaEmeraldGain := 2340.0
	maxManaEmeraldGain := 2460.0
	minManaRubyGain := 1073.0
	maxManaRubyGain := 1127.0
	if serpentCoilAura != nil {
		manaMultiplier = 1.25
		minManaEmeraldGain *= manaMultiplier
		maxManaEmeraldGain *= manaMultiplier
		minManaRubyGain *= manaMultiplier
		maxManaRubyGain *= manaMultiplier
	}
	manaEmeraldGainRange := maxManaEmeraldGain - minManaEmeraldGain
	manaRubyGainRange := maxManaRubyGain - minManaRubyGain

	var remainingManaGems int
	mage.RegisterResetEffect(func(sim *core.Simulation) {
		remainingManaGems = 4
	})

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: core.MageManaGemMCDActionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.GetConjuredCD(),
				Duration: time.Minute * 2,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if remainingManaGems == 1 {
				// Mana Ruby: Restores 1073 to 1127 mana. (2 Min Cooldown)
				manaGain := minManaRubyGain + (sim.RandomFloat("Mana Gem") * manaRubyGainRange)
				mage.AddMana(sim, manaGain, manaMetrics, true)
			} else {
				// Mana Emerald: Restores 2340 to 2460 mana. (2 Min Cooldown)
				manaGain := minManaEmeraldGain + (sim.RandomFloat("Mana Gem") * manaEmeraldGainRange)
				mage.AddMana(sim, manaGain, manaMetrics, true)
			}

			if serpentCoilAura != nil {
				serpentCoilAura.Activate(sim)
			}

			remainingManaGems--
			if remainingManaGems == 0 {
				// Disable this cooldown since we're out of emeralds.
				mage.DisableMajorCooldown(core.MageManaGemMCDActionID)
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
			if remainingManaGems == 1 {
				maxManaGain = maxManaRubyGain
			}
			if character.MaxMana()-(character.CurrentMana()+totalRegen) < maxManaGain {
				return false
			}

			return true
		},
	})
}
