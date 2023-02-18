package wotlk

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type StatCDFactory func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration)

func init() {
	// Keep these separated by stat, ordered by item ID within each group.

	// Wraps factory functions so that only the first item is included in tests.
	testFirstOnly := func(factory StatCDFactory) StatCDFactory {
		first := true
		return func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
			if first {
				first = false
				factory(itemID, bonus, duration, cooldown)
			} else {
				core.AddEffectsToTest = false
				factory(itemID, bonus, duration, cooldown)
				core.AddEffectsToTest = true
			}
		}
	}

	newHasteActive := testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
		core.NewSimpleStatOffensiveTrinketEffect(itemID, stats.Stats{stats.MeleeHaste: bonus, stats.SpellHaste: bonus}, duration, cooldown)
	})
	newHasteActive(36972, 256, time.Second*20, time.Minute*2) // Tome of Arcane Phenomena
	//newHasteActive(37558, 122, time.Second*20, time.Minute*2) // Tidal Boon
	//newHasteActive(37560, 124, time.Second*20, time.Minute*2) // Vial of Renewal
	//newHasteActive(37562, 140, time.Second*20, time.Minute*2) // Fury of the Crimson Drake
	//newHasteActive(38070, 148, time.Second*20, time.Minute*2) // Foresight's Anticipation
	//newHasteActive(38258, 140, time.Second*20, time.Minute*2) // Sailor's Knotted Charm
	//newHasteActive(38259, 140, time.Second*20, time.Minute*2) // First Mate's Pocketwatch
	newHasteActive(38764, 208, time.Second*20, time.Minute*2) // Rune of Finite Variation
	newHasteActive(40531, 491, time.Second*20, time.Minute*2) // Mark of Norgannon
	newHasteActive(43836, 212, time.Second*20, time.Minute*2) // Thorny Rose Brooch
	newHasteActive(45466, 457, time.Second*20, time.Minute*2) // Scale of Fates
	newHasteActive(46088, 375, time.Second*20, time.Minute*2) // Platinum Disks of Swiftness
	newHasteActive(48722, 512, time.Second*20, time.Minute*2) // Shard of the Crystal Heart

	newAttackPowerActive := testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
		core.NewSimpleStatOffensiveTrinketEffect(itemID, stats.Stats{stats.AttackPower: bonus, stats.RangedAttackPower: bonus}, duration, cooldown)
	})
	//newAttackPowerActive(35937, 328, time.Second*20, time.Minute*2)  // Braxley's Backyard Moonshine
	//newAttackPowerActive(36871, 280, time.Second*20, time.Minute*2)  // Fury of the Encroaching Storm
	newAttackPowerActive(37166, 670, time.Second*20, time.Minute*2) // Sphere of Red Dragon's Blood
	//newAttackPowerActive(37556, 248, time.Second*20, time.Minute*2)  // Talisman of the Tundra
	//newAttackPowerActive(37557, 304, time.Second*20, time.Minute*2)  // Warsong's Fervor
	//newAttackPowerActive(38080, 264, time.Second*20, time.Minute*2)  // Automated Weapon Coater
	//newAttackPowerActive(38081, 280, time.Second*20, time.Minute*2)  // Scarab of Isanoth
	newAttackPowerActive(38761, 248, time.Second*20, time.Minute*2)  // Talon of Hatred
	newAttackPowerActive(39257, 670, time.Second*20, time.Minute*2)  // Loatheb's Shadow
	newAttackPowerActive(44014, 432, time.Second*15, time.Minute*2)  // Fezzik's Pocketwatch
	newAttackPowerActive(45263, 905, time.Second*20, time.Minute*2)  // Wrathstone
	newAttackPowerActive(46086, 752, time.Second*20, time.Minute*2)  // Platinum Disks of Battle
	newAttackPowerActive(47734, 1024, time.Second*20, time.Minute*2) // Mark of Supremacy

	newSpellPowerActive := testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
		core.NewSimpleStatOffensiveTrinketEffect(itemID, stats.Stats{stats.SpellPower: bonus}, duration, cooldown)
	})
	//newSpellPowerActive(35935, 178, time.Second*20, time.Minute*2) // Infused Coldstone Rune
	//newSpellPowerActive(36872, 173, time.Second*20, time.Minute*2) // Mender of the Oncoming Dawn
	//newSpellPowerActive(36874, 183, time.Second*20, time.Minute*2) // Horn of the Herald
	//newSpellPowerActive(37555, 149, time.Second*20, time.Minute*2) // Warsong's Wrath
	newSpellPowerActive(37844, 346, time.Second*20, time.Minute*2) // Winged Talisman
	newSpellPowerActive(37873, 346, time.Second*20, time.Minute*2) // Mark of the War Prisoner
	//newSpellPowerActive(38073, 120, time.Second*15, time.Minute*2) // Will of the Red Dragonflight
	//newSpellPowerActive(38213, 149, time.Second*20, time.Minute*2) // Harbringer's Wrath
	//newSpellPowerActive(38527, 183, time.Second*20, time.Minute*2) // Strike of the Seas
	newSpellPowerActive(38760, 145, time.Second*20, time.Minute*2) // Mendicant's Charm
	newSpellPowerActive(38762, 145, time.Second*20, time.Minute*2) // Insignia of Bloody Fire
	newSpellPowerActive(38765, 202, time.Second*20, time.Minute*2) // Rune of Infinite Power
	newSpellPowerActive(39811, 183, time.Second*20, time.Minute*2) // Badge of the Infiltrator
	newSpellPowerActive(39819, 145, time.Second*20, time.Minute*2) // Bloodbinder's Runestone
	newSpellPowerActive(39821, 145, time.Second*20, time.Minute*2) // Spiritist's Focus
	newSpellPowerActive(42395, 292, time.Second*20, time.Minute*2) // Figurine - Twilight Serpent
	newSpellPowerActive(43837, 281, time.Second*20, time.Minute*2) // Soflty Glowing Orb
	newSpellPowerActive(44013, 281, time.Second*20, time.Minute*2) // Cannoneer's Fuselighter
	newSpellPowerActive(44015, 281, time.Second*20, time.Minute*2) // Cannoneer's Morale
	newSpellPowerActive(45148, 534, time.Second*20, time.Minute*2) // Living Flame
	newSpellPowerActive(45292, 431, time.Second*20, time.Minute*2) // Energy Siphon
	newSpellPowerActive(46087, 440, time.Second*20, time.Minute*2) // Platinum Disks of Sorcery
	newSpellPowerActive(48724, 599, time.Second*20, time.Minute*2) // Talisman of Resurgence
	newSpellPowerActive(50357, 716, time.Second*20, time.Minute*2) // Maghia's Misguided Quill

	newArmorPenActive := testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
		core.NewSimpleStatOffensiveTrinketEffect(itemID, stats.Stats{stats.ArmorPenetration: bonus}, duration, cooldown)
	})
	newArmorPenActive(37723, 291, time.Second*20, time.Minute*2) // Incisor Fragment

	newHealthActive := testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
		core.NewSimpleStatDefensiveTrinketEffect(itemID, stats.Stats{stats.Health: bonus}, duration, cooldown)
	})
	newHealthActive(37638, 3025, time.Second*15, time.Minute*3) // Offering of Sacrifice
	newHealthActive(39292, 3025, time.Second*15, time.Minute*3) // Repelling Charge
	newHealthActive(42128, 3385, time.Second*15, time.Minute*3) // Battlemaster's Hostility
	newHealthActive(42129, 3385, time.Second*15, time.Minute*3) // Battlemaster's Accuracy
	newHealthActive(42130, 3385, time.Second*15, time.Minute*3) // Battlemaster's Avidity
	newHealthActive(42131, 3385, time.Second*15, time.Minute*3) // Battlemaster's Conviction
	newHealthActive(42132, 3385, time.Second*15, time.Minute*3) // Battlemaster's Bravery
	newHealthActive(42133, 4608, time.Second*15, time.Minute*3) // Battlemaster's Fury
	newHealthActive(42134, 4608, time.Second*15, time.Minute*3) // Battlemaster's Precision
	newHealthActive(42135, 4608, time.Second*15, time.Minute*3) // Battlemaster's Vivacity
	newHealthActive(42136, 4608, time.Second*15, time.Minute*3) // Battlemaster's Rage
	newHealthActive(42137, 4608, time.Second*15, time.Minute*3) // Battlemaster's Ruination
	newHealthActive(47080, 4610, time.Second*15, time.Minute*3) // Satrina's Impeding Scarab
	newHealthActive(47088, 5186, time.Second*15, time.Minute*3) // Satrina's Impeding Scarab H
	newHealthActive(47290, 4610, time.Second*15, time.Minute*3) // Juggernaut's Vitality
	newHealthActive(47451, 5186, time.Second*15, time.Minute*3) // Juggernaut's Vitality H
	newHealthActive(50235, 4104, time.Second*15, time.Minute*3) // Ick's Rotting Thumb

	newArmorActive := testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
		core.NewSimpleStatDefensiveTrinketEffect(itemID, stats.Stats{stats.Armor: bonus}, duration, cooldown)
	})
	newArmorActive(36993, 3570, time.Second*20, time.Minute*2) // Seal of the Pantheon
	newArmorActive(45313, 5448, time.Second*20, time.Minute*2) // Furnace Stone

	newBlockValueActive := testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
		//core.NewSimpleStatDefensiveTrinketEffect(itemID, stats.Stats{stats.BlockValue: bonus}, duration, cooldown)
		// Hack for Lavanthor's Talisman Shared CD being shorter than its effect
		core.NewSimpleStatItemActiveEffect(itemID, stats.Stats{stats.BlockValue: bonus}, duration, cooldown, func(character *core.Character) core.Cooldown {
			return core.Cooldown{
				Timer:    character.GetDefensiveTrinketCD(),
				Duration: time.Second * 20,
			}
		}, nil)
	})
	newBlockValueActive(37872, 440, time.Second*40, time.Minute*2) // Lavanthor's Talisman

	newDodgeActive := testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
		core.NewSimpleStatDefensiveTrinketEffect(itemID, stats.Stats{stats.Dodge: bonus}, duration, cooldown)
	})
	newDodgeActive(40257, 455, time.Second*20, time.Minute*2) // Defender's Code
	newDodgeActive(40683, 335, time.Second*20, time.Minute*2) // Valor Medal of the First War
	newDodgeActive(44063, 300, time.Second*10, time.Minute*1) // Figurine - Monarch Crab
	newDodgeActive(45158, 457, time.Second*20, time.Minute*2) // Heart of Iron
	newDodgeActive(49080, 335, time.Second*20, time.Minute*2) // Brawler's Souvenir
	newDodgeActive(47735, 512, time.Second*20, time.Minute*2) // Glyph of Indomitability

	newParryActive := testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
		core.NewSimpleStatDefensiveTrinketEffect(itemID, stats.Stats{stats.Parry: bonus}, duration, cooldown)
	})
	newParryActive(40372, 375, time.Second*20, time.Minute*2) // Rune of Repulsion
	newParryActive(46021, 402, time.Second*20, time.Minute*2) // Royal Seal of King Llane

	newSpiritActive := testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
		core.NewSimpleStatDefensiveTrinketEffect(itemID, stats.Stats{stats.Spirit: bonus}, duration, cooldown)
	})
	newSpiritActive(38763, 184, time.Second*20, time.Minute*2) // Futuresight Rune
	newSpiritActive(39388, 336, time.Second*20, time.Minute*2) // Spirit-World Glass

	newResistsActive := testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
		core.NewSimpleStatDefensiveTrinketEffect(itemID, stats.Stats{stats.ArcaneResistance: bonus, stats.FireResistance: bonus, stats.FrostResistance: bonus, stats.NatureResistance: bonus, stats.ShadowResistance: bonus}, duration, cooldown)
	})
	newResistsActive(50361, 239, time.Second*10, time.Minute*1) // Sindragona's Flawless Fang
	newResistsActive(50364, 268, time.Second*10, time.Minute*1) // Sindragona's Flawless Fang H
}
