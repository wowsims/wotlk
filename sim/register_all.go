package sim

import (
	_ "github.com/wowsims/wotlk/sim/common"
	dpsDeathKnight "github.com/wowsims/wotlk/sim/deathknight/dps"
	tankDeathKnight "github.com/wowsims/wotlk/sim/deathknight/tank"
	"github.com/wowsims/wotlk/sim/druid/balance"
	"github.com/wowsims/wotlk/sim/druid/feral"
	feralTank "github.com/wowsims/wotlk/sim/druid/tank"
	_ "github.com/wowsims/wotlk/sim/encounters"
	"github.com/wowsims/wotlk/sim/hunter"
	"github.com/wowsims/wotlk/sim/mage"
	protectionPaladin "github.com/wowsims/wotlk/sim/paladin/protection"
	"github.com/wowsims/wotlk/sim/paladin/retribution"
	healingPriest "github.com/wowsims/wotlk/sim/priest/healing"
	"github.com/wowsims/wotlk/sim/priest/shadow"
	"github.com/wowsims/wotlk/sim/priest/smite"
	"github.com/wowsims/wotlk/sim/rogue"
	"github.com/wowsims/wotlk/sim/shaman/elemental"
	"github.com/wowsims/wotlk/sim/shaman/enhancement"
	"github.com/wowsims/wotlk/sim/warlock"
	dpsWarrior "github.com/wowsims/wotlk/sim/warrior/dps"
	protectionWarrior "github.com/wowsims/wotlk/sim/warrior/protection"
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
	healingPriest.RegisterHealingPriest()
	shadow.RegisterShadowPriest()
	smite.RegisterSmitePriest()
	rogue.RegisterRogue()
	dpsWarrior.RegisterDpsWarrior()
	protectionWarrior.RegisterProtectionWarrior()
	retribution.RegisterRetributionPaladin()
	protectionPaladin.RegisterProtectionPaladin()
	warlock.RegisterWarlock()
	dpsDeathKnight.RegisterDpsDeathknight()
	tankDeathKnight.RegisterTankDeathknight()
}
