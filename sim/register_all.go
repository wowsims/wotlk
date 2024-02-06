package sim

import (
	_ "github.com/wowsims/sod/sim/common"
	"github.com/wowsims/sod/sim/druid/balance"
	dpsRogue "github.com/wowsims/sod/sim/rogue"

	"github.com/wowsims/sod/sim/druid/feral"
	// restoDruid "github.com/wowsims/sod/sim/druid/restoration"
	// feralTank "github.com/wowsims/sod/sim/druid/tank"
	// _ "github.com/wowsims/sod/sim/encounters"
	"github.com/wowsims/sod/sim/hunter"
	"github.com/wowsims/sod/sim/mage"

	// holyPaladin "github.com/wowsims/sod/sim/paladin/holy"
	// protectionPaladin "github.com/wowsims/sod/sim/paladin/protection"
	// "github.com/wowsims/sod/sim/paladin/retribution"
	// healingPriest "github.com/wowsims/sod/sim/priest/healing"
	"github.com/wowsims/sod/sim/priest/shadow"
	// "github.com/wowsims/sod/sim/shaman/elemental"
	// "github.com/wowsims/sod/sim/shaman/enhancement"
	// restoShaman "github.com/wowsims/sod/sim/shaman/restoration"
	dpsWarlock "github.com/wowsims/sod/sim/warlock/dps"
	tankWarlock "github.com/wowsims/sod/sim/warlock/tank"
	dpsWarrior "github.com/wowsims/sod/sim/warrior/dps"
	// protectionWarrior "github.com/wowsims/sod/sim/warrior/protection"
)

var registered = false

func RegisterAll() {
	if registered {
		return
	}
	registered = true

	balance.RegisterBalanceDruid()
	feral.RegisterFeralDruid()
	// feralTank.RegisterFeralTankDruid()
	// restoDruid.RegisterRestorationDruid()
	// elemental.RegisterElementalShaman()
	// enhancement.RegisterEnhancementShaman()
	// restoShaman.RegisterRestorationShaman()
	hunter.RegisterHunter()
	mage.RegisterMage()
	// healingPriest.RegisterHealingPriest()
	shadow.RegisterShadowPriest()
	dpsRogue.RegisterRogue()
	dpsWarrior.RegisterDpsWarrior()
	// protectionWarrior.RegisterProtectionWarrior()
	// holyPaladin.RegisterHolyPaladin()
	// protectionPaladin.RegisterProtectionPaladin()
	// retribution.RegisterRetributionPaladin()
	dpsWarlock.RegisterDpsWarlock()
	tankWarlock.RegisterTankWarlock()
}
