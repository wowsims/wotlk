package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warlock *Warlock) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_WarlockShouldRecastDrainSoul:
		return warlock.newValueWarlockShouldRecastDrainSoul(rot, config.GetWarlockShouldRecastDrainSoul())
	case *proto.APLValue_WarlockShouldRefreshCorruption:
		return warlock.newValueWarlockShouldRefreshCorruption(rot, config.GetWarlockShouldRefreshCorruption())
	default:
		return nil
	}
}

type APLValueWarlockShouldRecastDrainSoul struct {
	core.DefaultAPLValueImpl
	warlock *Warlock
}

func (warlock *Warlock) newValueWarlockShouldRecastDrainSoul(rot *core.APLRotation, config *proto.APLValueWarlockShouldRecastDrainSoul) core.APLValue {
	return &APLValueWarlockShouldRecastDrainSoul{
		warlock: warlock,
	}
}
func (value *APLValueWarlockShouldRecastDrainSoul) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueWarlockShouldRecastDrainSoul) GetBool(sim *core.Simulation) bool {
	warlock := value.warlock

	// Assert that we're currently channeling Drain Soul.
	if warlock.ChanneledDot == nil || warlock.ChanneledDot.Spell != warlock.DrainSoul {
		return false
	}

	uaRefresh := warlock.UnstableAffliction.CurDot().RemainingDuration(sim) -
		warlock.UnstableAffliction.CastTime()

	curseRefresh := max(
		warlock.CurseOfAgony.CurDot().RemainingDuration(sim),
		warlock.CurseOfDoom.CurDot().RemainingDuration(sim),
		warlock.CurseOfElementsAuras.Get(warlock.CurrentTarget).RemainingDuration(sim),
		warlock.CurseOfTonguesAuras.Get(warlock.CurrentTarget).RemainingDuration(sim),
		warlock.CurseOfWeaknessAuras.Get(warlock.CurrentTarget).RemainingDuration(sim),
	) - warlock.CurseOfAgony.CastTime()

	hauntRefresh := 1000 * time.Second
	if warlock.HauntDebuffAuras != nil {
		hauntRefresh = warlock.HauntDebuffAuras.Get(warlock.CurrentTarget).RemainingDuration(sim) -
			warlock.Haunt.CastTime() -
			warlock.Haunt.TravelTime()
	}

	timeUntilRefresh := min(uaRefresh, curseRefresh)

	// the amount of ticks we have left, assuming we continue channeling
	dsDot := warlock.ChanneledDot
	ticksLeft := int(timeUntilRefresh/dsDot.TickPeriod()) + 1
	ticksLeft = min(ticksLeft, int(hauntRefresh/dsDot.TickPeriod()))
	ticksLeft = min(ticksLeft, dsDot.NumTicksRemaining(sim))

	// amount of ticks we'd get assuming we recast drain soul
	recastTicks := int(timeUntilRefresh/warlock.ApplyCastSpeed(dsDot.TickLength)) + 1
	recastTicks = min(recastTicks, int(hauntRefresh/warlock.ApplyCastSpeed(dsDot.TickLength)))
	recastTicks = min(recastTicks, int(dsDot.NumberOfTicks))

	if ticksLeft <= 0 || recastTicks <= 0 {
		return false
	}

	snapshotDmg := warlock.DrainSoul.ExpectedTickDamageFromCurrentSnapshot(sim, warlock.CurrentTarget) * float64(ticksLeft)
	recastDmg := warlock.DrainSoul.ExpectedTickDamage(sim, warlock.CurrentTarget) * float64(recastTicks)
	snapshotDPS := snapshotDmg / (time.Duration(ticksLeft) * dsDot.TickPeriod()).Seconds()
	recastDps := recastDmg / (time.Duration(recastTicks)*warlock.ApplyCastSpeed(dsDot.TickLength) + warlock.ChannelClipDelay).Seconds()

	//if sim.Log != nil {
	//	warlock.Log(sim, "Should Recast Drain Soul Calc: %.2f (%d) > %.2f (%d)", recastDps, recastTicks, snapshotDPS, ticksLeft)
	//}
	return recastDps > snapshotDPS
}
func (value *APLValueWarlockShouldRecastDrainSoul) String() string {
	return "Warlock Should Recast Drain Soul()"
}

type APLValueWarlockShouldRefreshCorruption struct {
	core.DefaultAPLValueImpl
	warlock *Warlock
	target  core.UnitReference
}

func (warlock *Warlock) newValueWarlockShouldRefreshCorruption(rot *core.APLRotation, config *proto.APLValueWarlockShouldRefreshCorruption) core.APLValue {
	target := rot.GetTargetUnit(config.TargetUnit)
	if target.Get() == nil {
		return nil
	}

	return &APLValueWarlockShouldRefreshCorruption{
		warlock: warlock,
		target:  target,
	}
}
func (value *APLValueWarlockShouldRefreshCorruption) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueWarlockShouldRefreshCorruption) GetBool(sim *core.Simulation) bool {
	warlock := value.warlock
	target := value.target.Get()

	dot := warlock.Corruption.Dot(target)
	if !dot.IsActive() {
		return true
	}

	// check if reapplying corruption is worthwhile
	snapshotCrit := dot.SnapshotCritChance
	snapshotMult := dot.SnapshotAttackerMultiplier * (snapshotCrit*(warlock.Corruption.CritMultiplier-1) + 1)

	attackTable := warlock.AttackTables[target.UnitIndex]
	curCrit := warlock.Corruption.SpellCritChance(target)
	curDmg := dot.Spell.AttackerDamageMultiplier(attackTable) * (curCrit*(warlock.Corruption.CritMultiplier-1) + 1)

	relDmgInc := curDmg / snapshotMult

	snapshotDmg := warlock.Corruption.ExpectedTickDamageFromCurrentSnapshot(sim, target)
	snapshotDmg *= float64(sim.GetRemainingDuration()) / float64(dot.TickPeriod())
	snapshotDmg *= (relDmgInc - 1)
	snapshotDmg -= warlock.Corruption.ExpectedTickDamageFromCurrentSnapshot(sim, target)

	//if sim.Log != nil {
	//	warlock.Log(sim, "Relative Corruption Inc: [%.2f], expected dmg gain: [%.2f]", relDmgInc, snapshotDmg)
	//}

	return relDmgInc > 1.15 || snapshotDmg > 10000
}
func (value *APLValueWarlockShouldRefreshCorruption) String() string {
	return "Warlock Should Refresh Corruption()"
}
