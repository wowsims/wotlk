package balance

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (moonkin *BalanceDruid) OnGCDReady(sim *core.Simulation) {
	moonkin.tryUseGCD(sim)
}

func (moonkin *BalanceDruid) tryUseGCD(sim *core.Simulation) {
	if moonkin.useSurplusRotation {
		moonkin.manaTracker.Update(sim, moonkin.GetCharacter())

		// If we have enough mana to burn, use the surplus rotation.
		if moonkin.manaTracker.ProjectedManaSurplus(sim, moonkin.GetCharacter()) {
			moonkin.actRotation(sim, moonkin.surplusRotation)
		} else {
			moonkin.actRotation(sim, moonkin.primaryRotation)
		}
	} else {
		moonkin.actRotation(sim, moonkin.primaryRotation)
	}
}

func (moonkin *BalanceDruid) actRotation(sim *core.Simulation, rotation proto.BalanceDruid_Rotation) {
	// Activate shared druid behaviors
	// Use Rebirth at the beginning of the fight if flagged in rotation settings
	// Potentially allow options for "Time of cast" in future or default cast like 1 min into fight
	// Currently just casts at the beginning of encounter (with all CDs popped)
	if moonkin.useBattleRes && moonkin.TryRebirth(sim) {
		return
	}

	target := moonkin.CurrentTarget

	var spell *core.Spell
	// TODO: add starfall always

	if moonkin.ShouldFaerieFire(sim) {
		spell = moonkin.FaerieFire
	} else if moonkin.ShouldCastHurricane(sim, rotation) {
		spell = moonkin.Hurricane
	} else if moonkin.ShouldCastInsectSwarm(sim, target, rotation) {
		spell = moonkin.InsectSwarm
	} else if moonkin.ShouldCastMoonfire(sim, target, rotation) {
		spell = moonkin.Moonfire
	} else {
		spell = moonkin.Starfire
		// TODO: Check for eclipse to decide if starfire or wrath
	}

	if success := spell.Cast(sim, target); !success {
		moonkin.WaitForMana(sim, spell.CurCast.Cost)
	}
}
