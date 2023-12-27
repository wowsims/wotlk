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

func (cat *FeralDruid) NewAPLAction(rot *core.APLRotation, config *proto.APLAction) core.APLActionImpl {
	switch config.Action.(type) {
	case *proto.APLAction_CatOptimalRotationAction:
		return cat.newActionCatOptimalRotationAction(rot, config.GetCatOptimalRotationAction())
	default:
		return nil
	}
}

type APLActionCatOptimalRotationAction struct {
	cat        *FeralDruid
	lastAction time.Duration
}

func (impl *APLActionCatOptimalRotationAction) GetInnerActions() []*core.APLAction { return nil }
func (impl *APLActionCatOptimalRotationAction) GetAPLValues() []core.APLValue      { return nil }
func (impl *APLActionCatOptimalRotationAction) Finalize(*core.APLRotation)         {}
func (impl *APLActionCatOptimalRotationAction) GetNextAction(*core.Simulation) *core.APLAction {
	return nil
}

func (cat *FeralDruid) newActionCatOptimalRotationAction(_ *core.APLRotation, config *proto.APLActionCatOptimalRotationAction) core.APLActionImpl {
	rotationOptions := &proto.FeralDruid_Rotation{
		RotationType:       config.RotationType,
		MaintainFaerieFire: true,
		UseRake:            config.UseRake,
		UseBite:            config.UseBite,
		BiteTime:           config.BiteTime,
		MangleSpam:         false,
		MaxFfDelay:         config.MaxFfDelay,
		Powerbear:          false,
		MinRoarOffset:      config.MinRoarOffset,
		RipLeeway:          config.RipLeeway,
		HotUptime:          0.0,
		FlowerWeave:        config.FlowerWeave,
		ManualParams:       config.ManualParams,
	}

	cat.setupRotation(rotationOptions)

	return &APLActionCatOptimalRotationAction{
		cat: cat,
	}
}

func (action *APLActionCatOptimalRotationAction) IsReady(sim *core.Simulation) bool {
	return sim.CurrentTime > action.lastAction
}

func (action *APLActionCatOptimalRotationAction) Execute(sim *core.Simulation) {
	cat := action.cat

	// If a melee swing resulted in an Omen proc, then schedule the
	// next player decision based on latency.
	if cat.Talents.OmenOfClarity && cat.ClearcastingAura.RemainingDuration(sim) == cat.ClearcastingAura.Duration {
		// Kick gcd loop, also need to account for any gcd 'left'
		// otherwise it breaks gcd logic
		kickTime := max(cat.NextGCDAt(), sim.CurrentTime+cat.latency)
		cat.NextRotationAction(sim, kickTime)
	}

	if cat.GCD.IsReady(sim) && (cat.rotationAction == nil || sim.CurrentTime >= cat.rotationAction.NextActionAt) {
		cat.OnGCDReady(sim)
	}

	cat.doTigersFury(sim)
	action.lastAction = sim.CurrentTime
}

func (action *APLActionCatOptimalRotationAction) Reset(*core.Simulation) {
	action.cat.usingHardcodedAPL = true
	action.lastAction = core.DurationFromSeconds(-100)
}

func (action *APLActionCatOptimalRotationAction) String() string {
	return "Execute Optimal Cat Action()"
}
