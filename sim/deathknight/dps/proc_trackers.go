package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

type ProcTracker struct {
	id          int32
	aura        *core.Aura
	didActivate bool
	isActive    bool
	expiresAt   time.Duration
}

func (ur *UnholyRotation) addProc(id int32, label string, isActive bool) bool {
	if !ur.dk.HasAura(label) {
		return false
	}
	ur.procTrackers = append(ur.procTrackers, &ProcTracker{
		id:          id,
		didActivate: false,
		isActive:    isActive,
		expiresAt:   -1,
		aura:        ur.dk.GetAura(label),
	})
	return true
}

func (ur *UnholyRotation) resetProcTrackers() {
	for _, procTracker := range ur.procTrackers {
		procTracker.didActivate = false
		procTracker.expiresAt = -1
	}
}

func (dk *DpsDeathknight) initProcTrackers() {
	dk.ur.procTrackers = make([]*ProcTracker, 0)

	dk.ur.addProc(40211, "Potion of Speed", true)
	dk.ur.addProc(54999, "Hyperspeed Acceleration", true)
	dk.ur.addProc(26297, "Berserking (Troll)", true)
	dk.ur.addProc(33697, "Blood Fury", true)

	dk.ur.addProc(53344, "Rune Of The Fallen Crusader Proc", false)
	dk.ur.addProc(55379, "Thundering Skyflare Diamond Proc", false)
	dk.ur.addProc(59620, "Berserking MH Proc", false)
	dk.ur.addProc(59620, "Berserking OH Proc", false)
	dk.ur.addProc(59626, "Black Magic Proc", false)

	dk.ur.addProc(42987, "DMC Greatness Strength Proc", false)

	dk.ur.addProc(47115, "Deaths Verdict Strength Proc", false)
	dk.ur.addProc(47131, "Deaths Verdict H Strength Proc", false)
	dk.ur.addProc(47303, "Deaths Choice Strength Proc", false)
	dk.ur.addProc(47464, "Deaths Choice H Strength Proc", false)

	dk.ur.addProc(71484, "Deathbringer's Will Strength Proc", false)
	dk.ur.addProc(71492, "Deathbringer's Will Haste Proc", false)
	dk.ur.addProc(71561, "Deathbringer's Will H Strength Proc", false)
	dk.ur.addProc(71560, "Deathbringer's Will H Haste Proc", false)

	dk.ur.addProc(37390, "Meteorite Whetstone Proc", false)
	dk.ur.addProc(39229, "Embrace of the Spider Proc", false)
	dk.ur.addProc(40684, "Mirror of Truth Proc", false)
	dk.ur.addProc(40767, "Sonic Booster Proc", false)
	dk.ur.addProc(43573, "Tears of Bitter Anguish Proc", false)
	dk.ur.addProc(44308, "Signet of Edward the Odd Proc", false)
	dk.ur.addProc(44914, "Anvil of Titans Proc", false)
	dk.ur.addProc(45286, "Pyrite Infuser Proc", false)
	dk.ur.addProc(45522, "Blood of the Old God Proc", false)
	dk.ur.addProc(45609, "Comet's Trail Proc", false)
	dk.ur.addProc(45866, "Elemental Focus Stone Proc", false)
	dk.ur.addProc(47214, "Banner of Victory Proc", false)
	dk.ur.addProc(49074, "Coren's Chromium Coaster Proc", false)
	dk.ur.addProc(50342, "Whispering Fanged Skull Proc", false)
	dk.ur.addProc(50343, "Whispering Fanged Skull H Proc", false)
	dk.ur.addProc(50401, "Ashen Band of Unmatched Vengeance Proc", false)
	dk.ur.addProc(50402, "Ashen Band of Endless Vengeance Proc", false)
	dk.ur.addProc(52571, "Ashen Band of Unmatched Might Proc", false)
	dk.ur.addProc(52572, "Ashen Band of Endless Might Proc", false)
	dk.ur.addProc(54569, "Sharpened Twilight Scale Proc", false)
	dk.ur.addProc(54590, "Sharpened Twilight Scale H Proc", false)
}

func (dk *DpsDeathknight) setupGargoyleCooldowns() {
	dk.ur.majorCds = make([]*core.MajorCooldown, 0)

	// hyperspeed accelerators
	dk.gargoyleCooldownSync(core.ActionID{SpellID: 54758}, false)

	// berserking (troll)
	dk.gargoyleCooldownSync(core.ActionID{SpellID: 26297}, false)

	// blood fury (orc)
	dk.gargoyleCooldownSync(core.ActionID{SpellID: 33697}, false)

	// potion of speed
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 40211}, true)

	// active ap trinkets
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 35937}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 36871}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 37166}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 37556}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 37557}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 38080}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 38081}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 38761}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 39257}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 45263}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 46086}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 47734}, false)

	// active haste trinkets
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 36972}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 37558}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 37560}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 37562}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 38070}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 38258}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 38259}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 38764}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 40531}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 43836}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 45466}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 46088}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 48722}, false)
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 50260}, false)
}

func (dk *DpsDeathknight) gargoyleCooldownSync(actionID core.ActionID, isPotion bool) {
	if dk.Character.HasMajorCooldown(actionID) {
		majorCd := dk.Character.GetMajorCooldown(actionID)
		dk.ur.majorCds = append(dk.ur.majorCds, majorCd)

		majorCd.ShouldActivate = func(sim *core.Simulation, character *core.Character) bool {
			return dk.ur.activatingGargoyle || (dk.SummonGargoyle.CD.TimeToReady(sim) > majorCd.Spell.CD.Duration && !isPotion) || dk.SummonGargoyle.CD.ReadyAt() > dk.Env.Encounter.Duration
		}
	}
}
