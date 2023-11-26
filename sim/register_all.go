package sim

import (
	_ "github.com/wowsims/classic/sim/common"
	// "github.com/wowsims/classic/sim/druid/balance"
	// "github.com/wowsims/classic/sim/druid/feral"
	// restoDruid "github.com/wowsims/classic/sim/druid/restoration"
	// feralTank "github.com/wowsims/classic/sim/druid/tank"
	// _ "github.com/wowsims/classic/sim/encounters"
	// "github.com/wowsims/classic/sim/hunter"
	"github.com/wowsims/classic/sim/mage"
	// holyPaladin "github.com/wowsims/classic/sim/paladin/holy"
	// protectionPaladin "github.com/wowsims/classic/sim/paladin/protection"
	// "github.com/wowsims/classic/sim/paladin/retribution"
	// healingPriest "github.com/wowsims/classic/sim/priest/healing"
	"github.com/wowsims/classic/sim/priest/shadow"
	// "github.com/wowsims/classic/sim/rogue"
	// "github.com/wowsims/classic/sim/shaman/elemental"
	// "github.com/wowsims/classic/sim/shaman/enhancement"
	// restoShaman "github.com/wowsims/classic/sim/shaman/restoration"
	// "github.com/wowsims/classic/sim/warlock"
	dpsWarrior "github.com/wowsims/classic/sim/warrior/dps"
	// protectionWarrior "github.com/wowsims/classic/sim/warrior/protection"
)

var registered = false

func RegisterAll() {
	if registered {
		return
	}
	registered = true

	// balance.RegisterBalanceDruid()
	// feral.RegisterFeralDruid()
	// feralTank.RegisterFeralTankDruid()
	// restoDruid.RegisterRestorationDruid()
	// elemental.RegisterElementalShaman()
	// enhancement.RegisterEnhancementShaman()
	// restoShaman.RegisterRestorationShaman()
	// hunter.RegisterHunter()
	mage.RegisterMage()
	// healingPriest.RegisterHealingPriest()
	shadow.RegisterShadowPriest()
	// rogue.RegisterRogue()
	dpsWarrior.RegisterDpsWarrior()
	// protectionWarrior.RegisterProtectionWarrior()
	// holyPaladin.RegisterHolyPaladin()
	// protectionPaladin.RegisterProtectionPaladin()
	// retribution.RegisterRetributionPaladin()
	// warlock.RegisterWarlock()
}
