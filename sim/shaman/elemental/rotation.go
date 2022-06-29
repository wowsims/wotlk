package elemental

import (
	"time"

	"github.com/wowsims/tbc/sim/common"
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (eleShaman *ElementalShaman) GetPresimOptions(_ proto.Player) *core.PresimOptions {
	return eleShaman.rotation.GetPresimOptions()
}

func (eleShaman *ElementalShaman) OnGCDReady(sim *core.Simulation) {
	eleShaman.tryUseGCD(sim)
}

func (eleShaman *ElementalShaman) OnManaTick(sim *core.Simulation) {
	if eleShaman.FinishedWaitingForManaAndGCDReady(sim) {
		eleShaman.tryUseGCD(sim)
	}
}

func (eleShaman *ElementalShaman) tryUseGCD(sim *core.Simulation) {
	if eleShaman.TryDropTotems(sim) {
		return
	}

	eleShaman.rotation.DoAction(eleShaman, sim)
	//actionSuccessful := newAction.Cast(sim)
	//if actionSuccessful {
	//	eleShaman.rotation.OnActionAccepted(eleShaman, sim, newAction)
	//} else {
	//	// Only way for a shaman spell to fail is due to mana cost.
	//	// Wait until we have enough mana to cast.
	//	eleShaman.WaitForMana(sim, newAction.GetManaCost())
	//}
}

// Picks which attacks / abilities the Shaman does.
type Rotation interface {
	GetPresimOptions() *core.PresimOptions

	// Returns the action this rotation would like to take next.
	DoAction(*ElementalShaman, *core.Simulation)

	// Returns this rotation to its initial state. Called before each Sim iteration.
	Reset(*ElementalShaman, *core.Simulation)
}

// ################################################################
//                              LB ONLY
// ################################################################
type LBOnlyRotation struct {
}

func (rotation *LBOnlyRotation) DoAction(eleShaman *ElementalShaman, sim *core.Simulation) {
	if !eleShaman.LightningBolt.Cast(sim, eleShaman.CurrentTarget) {
		eleShaman.WaitForMana(sim, eleShaman.LightningBolt.CurCast.Cost)
	}
}

func (rotation *LBOnlyRotation) Reset(eleShaman *ElementalShaman, sim *core.Simulation) {}
func (rotation *LBOnlyRotation) GetPresimOptions() *core.PresimOptions                  { return nil }

func NewLBOnlyRotation() *LBOnlyRotation {
	return &LBOnlyRotation{}
}

// ################################################################
//                             CL ON CD
// ################################################################
type CLOnCDRotation struct {
}

func (rotation *CLOnCDRotation) DoAction(eleShaman *ElementalShaman, sim *core.Simulation) {
	var spell *core.Spell
	if eleShaman.ChainLightning.IsReady(sim) {
		spell = eleShaman.ChainLightning
	} else {
		spell = eleShaman.LightningBolt
	}

	if !spell.Cast(sim, eleShaman.CurrentTarget) {
		eleShaman.WaitForMana(sim, spell.CurCast.Cost)
	}
}

func (rotation *CLOnCDRotation) Reset(eleShaman *ElementalShaman, sim *core.Simulation) {}
func (rotation *CLOnCDRotation) GetPresimOptions() *core.PresimOptions                  { return nil }

func NewCLOnCDRotation() *CLOnCDRotation {
	return &CLOnCDRotation{}
}

// ################################################################
//                          FIXED ROTATION
// ################################################################
type FixedRotation struct {
	numLBsPerCL       int32
	numLBsSinceLastCL int32
}

func (rotation *FixedRotation) DoAction(eleShaman *ElementalShaman, sim *core.Simulation) {
	var spell *core.Spell
	if rotation.numLBsSinceLastCL < rotation.numLBsPerCL {
		spell = eleShaman.LightningBolt
		rotation.numLBsSinceLastCL++
	} else if eleShaman.ChainLightning.IsReady(sim) {
		spell = eleShaman.ChainLightning
		rotation.numLBsSinceLastCL = 0
	} else if eleShaman.HasTemporarySpellCastSpeedIncrease() {
		// If we have a temporary haste effect (like bloodlust or quags eye) then
		// we should add LB casts instead of waiting
		spell = eleShaman.LightningBolt
		rotation.numLBsSinceLastCL++
	}

	if spell == nil {
		common.NewWaitAction(sim, &eleShaman.Unit, eleShaman.ChainLightning.TimeToReady(sim), common.WaitReasonRotation).Cast(sim)
	} else {
		if !spell.Cast(sim, eleShaman.CurrentTarget) {
			eleShaman.WaitForMana(sim, spell.CurCast.Cost)
		}
	}
}

func (rotation *FixedRotation) Reset(eleShaman *ElementalShaman, sim *core.Simulation) {
	rotation.numLBsSinceLastCL = rotation.numLBsPerCL // This lets us cast CL first
}

func (rotation *FixedRotation) GetPresimOptions() *core.PresimOptions { return nil }

func NewFixedRotation(numLBsPerCL int32) *FixedRotation {
	return &FixedRotation{
		numLBsPerCL: numLBsPerCL,
	}
}

// ################################################################
//                          CL ON CLEARCAST
// ################################################################
type CLOnClearcastRotation struct {
	// Whether the second-to-last spell procced clearcasting
	prevPrevCastProccedCC bool
}

func (rotation *CLOnClearcastRotation) DoAction(eleShaman *ElementalShaman, sim *core.Simulation) {
	var spell *core.Spell
	if !eleShaman.ChainLightning.IsReady(sim) || !rotation.prevPrevCastProccedCC {
		spell = eleShaman.LightningBolt
	} else {
		spell = eleShaman.ChainLightning
	}

	if !spell.Cast(sim, eleShaman.CurrentTarget) {
		eleShaman.WaitForMana(sim, spell.CurCast.Cost)
	} else {
		rotation.prevPrevCastProccedCC = eleShaman.ClearcastingAura.GetStacks() == 2
	}
}

func (rotation *CLOnClearcastRotation) Reset(eleShaman *ElementalShaman, sim *core.Simulation) {
	rotation.prevPrevCastProccedCC = true // Lets us cast CL first
}

func (rotation *CLOnClearcastRotation) GetPresimOptions() *core.PresimOptions { return nil }

func NewCLOnClearcastRotation() *CLOnClearcastRotation {
	return &CLOnClearcastRotation{}
}

// ################################################################
//                             ADAPTIVE
// ################################################################
type AdaptiveRotation struct {
	manaTracker common.ManaSpendingRateTracker

	hasClearcasting bool
	baseRotation    Rotation // The rotation used most of the time
	surplusRotation Rotation // The rotation used when we have extra mana
}

func (rotation *AdaptiveRotation) DoAction(eleShaman *ElementalShaman, sim *core.Simulation) {
	didLB := false
	if sim.GetNumTargets() == 1 {
		sp := eleShaman.GetStat(stats.NatureSpellPower) + eleShaman.GetStat(stats.SpellPower)
		lb := ((612 + (sp * 0.794)) * 1.2) / (2 * eleShaman.CastSpeed)
		cl := ((786 + (sp * 0.651)) * 1.0666) / core.MaxFloat((1.5*eleShaman.CastSpeed), 1)
		if eleShaman.has4pT6 {
			lb *= 1.05
		}
		if lb+10 >= cl {
			eleShaman.LightningBolt.Cast(sim, eleShaman.CurrentTarget)
			didLB = true
		}
	}

	if !didLB {
		// If we have enough mana to burn, use the surplus rotation.
		if rotation.manaTracker.ProjectedManaSurplus(sim, eleShaman.GetCharacter()) {
			rotation.surplusRotation.DoAction(eleShaman, sim)
		} else {
			rotation.baseRotation.DoAction(eleShaman, sim)
		}
	}

	rotation.manaTracker.Update(sim, eleShaman.GetCharacter())
}

func (rotation *AdaptiveRotation) Reset(eleShaman *ElementalShaman, sim *core.Simulation) {
	rotation.manaTracker.Reset()
	rotation.baseRotation.Reset(eleShaman, sim)
	rotation.surplusRotation.Reset(eleShaman, sim)
}

func (rotation *AdaptiveRotation) GetPresimOptions() *core.PresimOptions {
	return &core.PresimOptions{
		SetPresimPlayerOptions: func(player *proto.Player) {
			player.Spec.(*proto.Player_ElementalShaman).ElementalShaman.Rotation.Type = proto.ElementalShaman_Rotation_CLOnClearcast
		},

		OnPresimResult: func(presimResult proto.UnitMetrics, iterations int32, duration time.Duration) bool {
			if !rotation.hasClearcasting {
				rotation.baseRotation = NewLBOnlyRotation()
				rotation.surplusRotation = NewCLOnCDRotation()
			} else {
				if float64(presimResult.SecondsOomAvg) >= 0.03*duration.Seconds() {
					rotation.baseRotation = NewLBOnlyRotation()
					rotation.surplusRotation = NewCLOnClearcastRotation()
				} else {
					rotation.baseRotation = NewCLOnClearcastRotation()
					rotation.surplusRotation = NewCLOnCDRotation()
				}
			}
			return true
		},
	}
}

func NewAdaptiveRotation(talents *proto.ShamanTalents) *AdaptiveRotation {
	return &AdaptiveRotation{
		hasClearcasting: talents.ElementalFocus,
		manaTracker:     common.NewManaSpendingRateTracker(),
	}
}

// A single action that an Agent can take.
type AgentAction interface {
	GetActionID() core.ActionID

	// TODO: Maybe change this to 'ResourceCost'
	// Amount of mana required to perform the action.
	GetManaCost() float64

	// Do the action. Returns whether the action was successful. An unsuccessful
	// action indicates that the prerequisites, like resource cost, were not met.
	Cast(sim *core.Simulation) bool
}
