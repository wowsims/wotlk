package sim

import (
	_ "github.com/wowsims/tbc/sim/common"
	"github.com/wowsims/tbc/sim/druid/balance"
	"github.com/wowsims/tbc/sim/druid/feral"
	feralTank "github.com/wowsims/tbc/sim/druid/tank"
	_ "github.com/wowsims/tbc/sim/encounters"
	"github.com/wowsims/tbc/sim/hunter"
	"github.com/wowsims/tbc/sim/mage"
	protectionPaladin "github.com/wowsims/tbc/sim/paladin/protection"
	"github.com/wowsims/tbc/sim/paladin/retribution"
	"github.com/wowsims/tbc/sim/priest/shadow"
	"github.com/wowsims/tbc/sim/priest/smite"
	"github.com/wowsims/tbc/sim/rogue"
	"github.com/wowsims/tbc/sim/shaman/elemental"
	"github.com/wowsims/tbc/sim/shaman/enhancement"
	"github.com/wowsims/tbc/sim/warlock"
	dpsWarrior "github.com/wowsims/tbc/sim/warrior/dps"
	protectionWarrior "github.com/wowsims/tbc/sim/warrior/protection"
)

var registered = false

func RegisterAll() {
	if registered {
		return
	}
	registered = true

	balance.RegisterBalanceDruid()
	feral.RegisterFeralDruid()
	feralTank.RegisterFeralTankDruid()
	elemental.RegisterElementalShaman()
	enhancement.RegisterEnhancementShaman()
	hunter.RegisterHunter()
	mage.RegisterMage()
	shadow.RegisterShadowPriest()
	rogue.RegisterRogue()
	dpsWarrior.RegisterDpsWarrior()
	protectionWarrior.RegisterProtectionWarrior()
	retribution.RegisterRetributionPaladin()
	protectionPaladin.RegisterProtectionPaladin()
	smite.RegisterSmitePriest()
	warlock.RegisterWarlock()
}
