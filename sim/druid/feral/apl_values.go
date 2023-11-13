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
	case *proto.APLValue_CatNewSavageRoarDuration:
		return cat.newValueCatNewSavageRoarDuration(rot, config.GetCatNewSavageRoarDuration())
	default:
		return nil
	}
}

type APLValueCatExcessEnergy struct {
	core.DefaultAPLValueImpl
	cat *FeralDruid
}

func (cat *FeralDruid) newValueCatExcessEnergy(_ *core.APLRotation, _ *proto.APLValueCatExcessEnergy) core.APLValue {
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

	simTimeRemain := sim.GetRemainingDuration()
	if ripDot := cat.Rip.CurDot(); ripDot.IsActive() && ripDot.RemainingDuration(sim) < simTimeRemain-time.Second*10 && cat.ComboPoints() == 5 {
		ripCost := core.Ternary(cat.berserkExpectedAt(sim, ripDot.ExpiresAt()), cat.Rip.DefaultCast.Cost*0.5, cat.Rip.DefaultCast.Cost)
		pendingPool.addAction(ripDot.ExpiresAt(), ripCost)
		cat.ripRefreshPending = true
	}
	if rakeDot := cat.Rake.CurDot(); rakeDot.IsActive() && rakeDot.RemainingDuration(sim) < simTimeRemain-rakeDot.Duration {
		rakeCost := core.Ternary(cat.berserkExpectedAt(sim, rakeDot.ExpiresAt()), cat.Rake.DefaultCast.Cost*0.5, cat.Rake.DefaultCast.Cost)
		pendingPool.addAction(rakeDot.ExpiresAt(), rakeCost)
	}
	if cat.bleedAura.IsActive() && cat.bleedAura.RemainingDuration(sim) < simTimeRemain-time.Second {
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

type APLValueCatNewSavageRoarDuration struct {
	core.DefaultAPLValueImpl
	cat *FeralDruid
}

func (cat *FeralDruid) newValueCatNewSavageRoarDuration(_ *core.APLRotation, _ *proto.APLValueCatNewSavageRoarDuration) core.APLValue {
	return &APLValueCatNewSavageRoarDuration{
		cat: cat,
	}
}
func (value *APLValueCatNewSavageRoarDuration) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueCatNewSavageRoarDuration) GetDuration(_ *core.Simulation) time.Duration {
	cat := value.cat
	return cat.SavageRoarDurationTable[cat.ComboPoints()]
}
func (value *APLValueCatNewSavageRoarDuration) String() string {
	return "New Savage Roar Duration()"
}
