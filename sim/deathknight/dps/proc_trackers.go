package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

type ProcTracker struct {
	id          int32
	aura        *core.Aura
	didActivate bool
	expiresAt   time.Duration
}

type UnholyRotation struct {
	dk            *DpsDeathknight
	lastCastSpell *core.Spell

	ffFirst bool
	syncFF  bool

	syncTimeFF time.Duration

	recastedFF bool
	recastedBP bool

	procTrackers []*ProcTracker
}

func (ur *UnholyRotation) addProc(id int32, label string) bool {
	if !ur.dk.HasAura(label) {
		return false
	}
	ur.procTrackers = append(ur.procTrackers, &ProcTracker{
		id:          id,
		didActivate: false,
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

func (ur *UnholyRotation) Reset(sim *core.Simulation) {
	ur.syncFF = false

	ur.syncTimeFF = 0

	ur.recastedFF = false
	ur.recastedBP = false

	ur.resetProcTrackers()
}

func (dk *DpsDeathknight) initProcTrackers() {
	dk.ur.procTrackers = make([]*ProcTracker, 0)

	dk.ur.addProc(54999, "Hyperspeed Acceleration")
	dk.ur.addProc(26297, "Berserking (Troll)")
	dk.ur.addProc(59620, "Berserking MH Proc")
	dk.ur.addProc(59620, "Berserking OH Proc")

	dk.ur.addProc(59626, "Black Magic Proc")
	dk.ur.addProc(53344, "Rune Of The Fallen Crusader Proc")
	dk.ur.addProc(55379, "Thundering Skyflare Diamond Proc")

	dk.ur.addProc(42987, "DMC Greatness Strength Proc")
	dk.ur.addProc(47115, "Deaths Verdict Strength Proc")
	dk.ur.addProc(47131, "Deaths Verdict H Strength Proc")
	dk.ur.addProc(47303, "Deaths Choice Strength Proc")
	dk.ur.addProc(47464, "Deaths Choice H Strength Proc")

	dk.ur.addProc(71484, "Deathbringer's Will Strength Proc")
	dk.ur.addProc(71486, "Deathbringer's Will AP Proc")
	dk.ur.addProc(71492, "Deathbringer's Will Haste Proc")

	dk.ur.addProc(71561, "Deathbringer's Will H Strength Proc")
	dk.ur.addProc(71558, "Deathbringer's Will H AP Proc")
	dk.ur.addProc(71560, "Deathbringer's Will H Haste Proc")

	dk.ur.addProc(37390, "Meteorite Whetstone Proc")
	dk.ur.addProc(39229, "Embrace of the Spider Proc")
	dk.ur.addProc(40684, "Mirror of Truth Proc")
	dk.ur.addProc(40767, "Sonic Booster Proc")
	dk.ur.addProc(43573, "Tears of Bitter Anguish Proc")
	dk.ur.addProc(44308, "Signet of Edward the Odd Proc")
	dk.ur.addProc(44914, "Anvil of Titans Proc")
	dk.ur.addProc(45286, "Pyrite Infuser Proc")
	dk.ur.addProc(45522, "Blood of the Old God Proc")
	dk.ur.addProc(45609, "Comet's Trail Proc")
	dk.ur.addProc(45866, "Elemental Focus Stone Proc")
	dk.ur.addProc(47214, "Banner of Victory Proc")
	dk.ur.addProc(49074, "Coren's Chromium Coaster Proc")
	dk.ur.addProc(50342, "Whispering Fanged Skull Proc")
	dk.ur.addProc(50343, "Whispering Fanged Skull H Proc")
	dk.ur.addProc(50401, "Ashen Band of Unmatched Vengeance Proc")
	dk.ur.addProc(50402, "Ashen Band of Endless Vengeance Proc")
	dk.ur.addProc(52571, "Ashen Band of Unmatched Might Proc")
	dk.ur.addProc(52572, "Ashen Band of Endless Might Proc")
	dk.ur.addProc(54569, "Sharpened Twilight Scale Proc")
	dk.ur.addProc(54590, "Sharpened Twilight Scale H Proc")
}
