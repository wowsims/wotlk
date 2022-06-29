package smite

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
)

// TODO: probably do something different instead of making it global?
const (
	mbidx int = iota
	swdidx
	vtidx
	swpidx
)

func (spriest *SmitePriest) OnGCDReady(sim *core.Simulation) {
	spriest.tryUseGCD(sim)
}

func (spriest *SmitePriest) OnManaTick(sim *core.Simulation) {
	if spriest.FinishedWaitingForManaAndGCDReady(sim) {
		spriest.tryUseGCD(sim)
	}
}

func (spriest *SmitePriest) tryUseGCD(sim *core.Simulation) {

	// Calculate higher SW:P uptime if using HF
	swpRemaining := spriest.ShadowWordPainDot.RemainingDuration(sim)

	// smite cast time, talent assumed
	smiteCastTime := spriest.ApplyCastSpeed(time.Millisecond * 2000)

	// holy fire cast time
	hfCastTime := spriest.ApplyCastSpeed(time.Millisecond * 3000)

	var spell *core.Spell
	// Always attempt to keep SW:P up if its down
	if !spriest.ShadowWordPainDot.IsActive() {
		spell = spriest.ShadowWordPain
		// Favor star shards for NE if off cooldown first
	} else if spriest.rotation.UseStarshards && spriest.Starshards.IsReady(sim) {
		spell = spriest.Starshards
		// Allow for undead to use devouring plague off CD
	} else if spriest.rotation.UseDevPlague && spriest.DevouringPlague.IsReady(sim) {
		spell = spriest.DevouringPlague
		// If setting enabled, throw mind blast into our rotation off CD
	} else if spriest.rotation.UseMindBlast && spriest.MindBlast.IsReady(sim) {
		spell = spriest.MindBlast
		// If setting enabled, cast Shadow Word: Death on cooldown
	} else if spriest.rotation.UseShadowWordDeath && spriest.ShadowWordDeath.IsReady(sim) {
		spell = spriest.ShadowWordDeath
		// Consider HF if SWP will fall off after 1 smite but before 2 smites from now finishes
		//	and swp falls off after hf finishes (assumption never worth clipping)
	} else if spriest.rotation.RotationType == proto.SmitePriest_Rotation_HolyFireWeave && swpRemaining > smiteCastTime && swpRemaining < hfCastTime {
		spell = spriest.HolyFire
		// Base filler spell is smite
	} else {
		spell = spriest.Smite
	}

	if success := spell.Cast(sim, spriest.CurrentTarget); !success {
		spriest.WaitForMana(sim, spell.CurCast.Cost)
	}
}
