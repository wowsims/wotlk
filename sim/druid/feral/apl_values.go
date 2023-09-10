package feral

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (cat *FeralDruid) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_CatExcessEnergy:
		return cat.newValueCatExcessEnergy(rot, config.GetCatExcessEnergy())
	default:
		return nil
	}
}

type APLValueCatExcessEnergy struct {
	core.DefaultAPLValueImpl
	cat *FeralDruid
}

func (cat *FeralDruid) newValueCatExcessEnergy(rot *core.APLRotation, config *proto.APLValueCatExcessEnergy) core.APLValue {
	return &APLValueCatExcessEnergy{
		cat: cat,
	}
}
func (value *APLValueCatExcessEnergy) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCatExcessEnergy) GetFloat(sim *core.Simulation) float64 {
	cat := value.cat
	pendingPool := PoolingActions{}
	pendingPool.create(4)

	curCp := cat.ComboPoints()
	simTimeRemain := sim.GetRemainingDuration()
	rakeDot := cat.Rake.CurDot()
	ripDot := cat.Rip.CurDot()
	mangleRefreshPending := cat.bleedAura.IsActive() && cat.bleedAura.RemainingDuration(sim) < (simTimeRemain-time.Second)
	endThresh := time.Second * 10

	if ripDot.IsActive() && (ripDot.RemainingDuration(sim) < simTimeRemain-endThresh) && curCp == 5 {
		ripCost := core.Ternary(cat.berserkExpectedAt(sim, ripDot.ExpiresAt()), cat.Rip.DefaultCast.Cost*0.5, cat.Rip.DefaultCast.Cost)
		pendingPool.addAction(ripDot.ExpiresAt(), ripCost)
		cat.ripRefreshPending = true
	}
	if rakeDot.IsActive() && (rakeDot.RemainingDuration(sim) < simTimeRemain-rakeDot.Duration) {
		rakeCost := core.Ternary(cat.berserkExpectedAt(sim, rakeDot.ExpiresAt()), cat.Rake.DefaultCast.Cost*0.5, cat.Rake.DefaultCast.Cost)
		pendingPool.addAction(rakeDot.ExpiresAt(), rakeCost)
	}
	if mangleRefreshPending {
		mangleCost := core.Ternary(cat.berserkExpectedAt(sim, cat.bleedAura.ExpiresAt()), cat.MangleCat.DefaultCast.Cost*0.5, cat.MangleCat.DefaultCast.Cost)
		pendingPool.addAction(cat.bleedAura.ExpiresAt(), mangleCost)
	}
	if cat.SavageRoarAura.IsActive() {
		roarCost := core.Ternary(cat.berserkExpectedAt(sim, cat.SavageRoarAura.ExpiresAt()), cat.SavageRoar.DefaultCast.Cost*0.5, cat.SavageRoar.DefaultCast.Cost)
		pendingPool.addAction(cat.SavageRoarAura.ExpiresAt(), roarCost)
	}

	pendingPool.sort()

	floatingEnergy := pendingPool.calcFloatingEnergy(cat, sim)
	return cat.CurrentEnergy() - floatingEnergy
}
func (value *APLValueCatExcessEnergy) String() string {
	return "Cat Excess Energy()"
}
