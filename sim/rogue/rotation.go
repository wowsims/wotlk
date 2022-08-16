package rogue

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (rogue *Rogue) OnEnergyGain(sim *core.Simulation) {
	if rogue.KillingSpreeAura.IsActive() {
		rogue.DoNothing()
		return
	}
	rogue.TryUseCooldowns(sim)
	if rogue.GCD.IsReady(sim) {
		rogue.rotation(sim)
	}
}

func (rogue *Rogue) OnGCDReady(sim *core.Simulation) {
	if rogue.KillingSpreeAura.IsActive() {
		rogue.DoNothing()
		return
	}
	rogue.rotation(sim)
}

func (rogue *Rogue) rotation(sim *core.Simulation) {
	if rogue.CurrentPriority == nil {
		rogue.CurrentPriority = rogue.getPriority(sim)
	}
	if rogue.CurrentPriority == nil {
		sim.Log("Problem")
	}
	spell := rogue.CurrentPriority.Spell(rogue)
	if spell == nil {
		rogue.CurrentPriority = nil
	} else if spell.Cast(sim, rogue.CurrentTarget) {
		rogue.CurrentPriority.CastCount += 1
		rogue.CurrentPriority = nil
	}
	if rogue.GCD.IsReady(sim) {
		rogue.DoNothing()
	}
}

func MakeConditionForFrequency(frequency proto.Rogue_Rotation_Frequency, aura *core.Aura) func(*core.Simulation, *Rogue, *RoguePriority) bool {
	var shouldCast func(*core.Simulation, *Rogue, *RoguePriority) bool
	switch frequency {
	case proto.Rogue_Rotation_Never:
		shouldCast = func(s *core.Simulation, r *Rogue, p *RoguePriority) bool { return false }
	case proto.Rogue_Rotation_Once:
		shouldCast = func(s *core.Simulation, r *Rogue, p *RoguePriority) bool { return p.CastCount < 1 }
	case proto.Rogue_Rotation_Build:
		shouldCast = func(s *core.Simulation, r *Rogue, p *RoguePriority) bool {
			return p.GeneratedComboPoints+r.ComboPoints() <= 5
		}
	case proto.Rogue_Rotation_Maintain:
		shouldCast = func(s *core.Simulation, r *Rogue, p *RoguePriority) bool {
			if aura == nil {
				return true
			}
			return RemainingAuraDuration(s, aura) <= 0
		}
	case proto.Rogue_Rotation_Fill:
		shouldCast = func(s *core.Simulation, r *Rogue, p *RoguePriority) bool {
			return true
		}
	}
	return shouldCast
}

func (rogue *Rogue) SetMultiTargetPriorityList() {
	rogue.PriorityList = make([]RoguePriority, 0)
	sliceCP := int32(1)
	if rogue.Rotation.MultiTargetSliceFrequency == proto.Rogue_Rotation_Once {
		sliceCP = rogue.Rotation.MinimumComboPointsMultiTargetSlice
	}
	sliceAndDice := RoguePriority{
		MinComboPoints: sliceCP,
		ShouldCast:     MakeConditionForFrequency(rogue.Rotation.MultiTargetSliceFrequency, rogue.SliceAndDiceAura),
		Spell:          func(r *Rogue) *core.Spell { return r.SliceAndDice[r.ComboPoints()] },
	}
	rogue.PriorityList = append(rogue.PriorityList, sliceAndDice)
	sliceAndDiceIndex := 0

	builderCondition := MakeConditionForFrequency(proto.Rogue_Rotation_Build, nil)
	if rogue.Rotation.MultiTargetSliceFrequency == proto.Rogue_Rotation_Once {
		oldCondition := builderCondition
		builderCondition = func(s *core.Simulation, r *Rogue, rp *RoguePriority) bool {
			return oldCondition(s, r, rp) && r.PriorityList[sliceAndDiceIndex].CastCount < 1
		}
	}
	// Mutilate
	rogue.PriorityList = append(rogue.PriorityList, RoguePriority{
		GeneratedComboPoints: 2,
		ShouldCast: func(s *core.Simulation, r *Rogue, rp *RoguePriority) bool {
			return r.Talents.Mutilate && builderCondition(s, r, rp)
		},
		Spell: func(rogue *Rogue) *core.Spell {
			return rogue.Mutilate
		},
	})
	// Sinister Strike
	rogue.PriorityList = append(rogue.PriorityList, RoguePriority{
		GeneratedComboPoints: 1,
		ShouldCast: func(s *core.Simulation, r *Rogue, rp *RoguePriority) bool {
			return !r.Talents.Mutilate && builderCondition(s, r, rp)
		},
		Spell: func(rogue *Rogue) *core.Spell {
			return rogue.SinisterStrike
		},
	})
	// Hunger for Blood
	rogue.PriorityList = append(rogue.PriorityList, RoguePriority{
		ShouldCast: func(s *core.Simulation, r *Rogue, rp *RoguePriority) bool {
			return r.Talents.HungerForBlood && MakeConditionForFrequency(proto.Rogue_Rotation_Maintain, rogue.HungerForBloodAura)(s, r, rp)
		},
		Spell: func(rogue *Rogue) *core.Spell {
			return rogue.HungerForBlood
		},
	})
	// Enable CDs
	rogue.PriorityList = append(rogue.PriorityList, RoguePriority{
		ShouldCast: func(sim *core.Simulation, rogue *Rogue, prio *RoguePriority) bool {
			return rogue.disabledMCDs != nil
		},
		Spell: func(rogue *Rogue) *core.Spell {
			if rogue.disabledMCDs != nil {
				rogue.EnableAllCooldowns(rogue.disabledMCDs)
				rogue.disabledMCDs = nil
			}
			return nil
		},
	})
	// Fan of Knives
	rogue.PriorityList = append(rogue.PriorityList, RoguePriority{
		Spell: func(rogue *Rogue) *core.Spell {
			return rogue.FanOfKnives
		},
	})
}

func (rogue *Rogue) SetPriorityList() {
	rogue.PriorityList = make([]RoguePriority, 0)
	// Slice and Dice
	rogue.PriorityList = append(rogue.PriorityList, RoguePriority{
		MinComboPoints: 1,
		ShouldCast:     MakeConditionForFrequency(proto.Rogue_Rotation_Maintain, rogue.SliceAndDiceAura),
		Spell: func(rogue *Rogue) *core.Spell {
			return rogue.SliceAndDice[rogue.ComboPoints()]
		},
	})
	// Expose armor
	exposeCP := int32(1)
	if rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Once {
		exposeCP = rogue.Rotation.MinimumComboPointsExposeArmor
	}
	rogue.PriorityList = append(rogue.PriorityList, RoguePriority{
		MinComboPoints: exposeCP,
		ShouldCast:     MakeConditionForFrequency(rogue.Rotation.ExposeArmorFrequency, rogue.ExposeArmorAura),
		Spell: func(rogue *Rogue) *core.Spell {
			return rogue.ExposeArmor[rogue.ComboPoints()]
		},
	})
	builderCondition := MakeConditionForFrequency(proto.Rogue_Rotation_Build, nil)
	// Mutilate
	rogue.PriorityList = append(rogue.PriorityList, RoguePriority{
		GeneratedComboPoints: 2,
		ShouldCast: func(s *core.Simulation, r *Rogue, rp *RoguePriority) bool {
			return r.Talents.Mutilate && builderCondition(s, r, rp)
		},
		Spell: func(rogue *Rogue) *core.Spell {
			return rogue.Mutilate
		},
	})
	// Sinister Strike
	rogue.PriorityList = append(rogue.PriorityList, RoguePriority{
		GeneratedComboPoints: 1,
		ShouldCast: func(s *core.Simulation, r *Rogue, rp *RoguePriority) bool {
			return !r.Talents.Mutilate && builderCondition(s, r, rp)
		},
		Spell: func(rogue *Rogue) *core.Spell {
			return rogue.SinisterStrike
		},
	})

	// Envenom
	envenomPrio := RoguePriority{
		MinComboPoints: 3,
		ShouldCast:     MakeConditionForFrequency(proto.Rogue_Rotation_Maintain, rogue.EnvenomAura),
		Spell: func(rogue *Rogue) *core.Spell {
			return rogue.Envenom[rogue.ComboPoints()]
		},
	}
	// Rupture
	rupturePrio := RoguePriority{
		MinComboPoints: 3,
		ShouldCast:     MakeConditionForFrequency(proto.Rogue_Rotation_Maintain, rogue.RuptureDot.Aura),
		Spell: func(rogue *Rogue) *core.Spell {
			return rogue.Rupture[rogue.ComboPoints()]
		},
	}
	// Evis
	evisceratePrio := RoguePriority{
		MinComboPoints: 3,
		ShouldCast:     MakeConditionForFrequency(proto.Rogue_Rotation_Fill, nil),
		Spell: func(rogue *Rogue) *core.Spell {
			return rogue.Eviscerate[rogue.ComboPoints()]
		},
	}
	if rogue.Talents.Mutilate {
		switch rogue.Rotation.AssassinationFinisherPriority {
		case proto.Rogue_Rotation_EnvenomRupture:
			envenomPrio.MinComboPoints = rogue.Rotation.MinimumComboPointsPrimaryFinisher
			rupturePrio.MinComboPoints = rogue.Rotation.MinimumComboPointsSecondaryFinisher
			rogue.PriorityList = append(rogue.PriorityList, envenomPrio)
			dependencyIdx := len(rogue.PriorityList) - 1
			oldCondition := rupturePrio.ShouldCast
			rupturePrio.ShouldCast = func(sim *core.Simulation, rogue *Rogue, prio *RoguePriority) bool {
				return !rogue.PriorityList[dependencyIdx].ShouldCast(sim, rogue, prio) && oldCondition(sim, rogue, prio)
			}
			rupturePrio.IsFiller = true
			rogue.PriorityList = append(rogue.PriorityList, rupturePrio)
		case proto.Rogue_Rotation_RuptureEnvenom:
			rupturePrio.MinComboPoints = rogue.Rotation.MinimumComboPointsPrimaryFinisher
			envenomPrio.MinComboPoints = rogue.Rotation.MinimumComboPointsSecondaryFinisher
			rogue.PriorityList = append(rogue.PriorityList, rupturePrio)
			dependencyIdx := len(rogue.PriorityList) - 1
			oldCondition := envenomPrio.ShouldCast
			envenomPrio.ShouldCast = func(sim *core.Simulation, rogue *Rogue, prio *RoguePriority) bool {
				return !rogue.PriorityList[dependencyIdx].ShouldCast(sim, rogue, prio) && oldCondition(sim, rogue, prio)
			}
			envenomPrio.IsFiller = true
			rogue.PriorityList = append(rogue.PriorityList, envenomPrio)
		}
	} else {
		switch rogue.Rotation.CombatFinisherPriority {
		case proto.Rogue_Rotation_RuptureEviscerate:
			rupturePrio.MinComboPoints = rogue.Rotation.MinimumComboPointsPrimaryFinisher
			evisceratePrio.MinComboPoints = rogue.Rotation.MinimumComboPointsSecondaryFinisher
			rogue.PriorityList = append(rogue.PriorityList, rupturePrio)
			dependencyIdx := len(rogue.PriorityList) - 1
			oldCondition := evisceratePrio.ShouldCast
			evisceratePrio.ShouldCast = func(sim *core.Simulation, rogue *Rogue, prio *RoguePriority) bool {
				return !rogue.PriorityList[dependencyIdx].ShouldCast(sim, rogue, prio) && oldCondition(sim, rogue, prio)
			}
			evisceratePrio.IsFiller = true
			rogue.PriorityList = append(rogue.PriorityList, evisceratePrio)
		case proto.Rogue_Rotation_EviscerateRupture:
			evisceratePrio.MinComboPoints = rogue.Rotation.MinimumComboPointsPrimaryFinisher
			rupturePrio.MinComboPoints = rogue.Rotation.MinimumComboPointsSecondaryFinisher
			rupturePrio.IsFiller = true
			rogue.PriorityList = append(rogue.PriorityList, rupturePrio)
			rogue.PriorityList = append(rogue.PriorityList, evisceratePrio)
		}
	}
	// Hunger for Blood
	rogue.PriorityList = append(rogue.PriorityList, RoguePriority{
		ShouldCast: func(s *core.Simulation, r *Rogue, rp *RoguePriority) bool {
			return r.Talents.HungerForBlood && MakeConditionForFrequency(proto.Rogue_Rotation_Maintain, rogue.HungerForBloodAura)(s, r, rp)
		},
		Spell: func(rogue *Rogue) *core.Spell {
			return rogue.HungerForBlood
		},
	})
	// Enable CDs
	rogue.PriorityList = append(rogue.PriorityList, RoguePriority{
		Spell: func(rogue *Rogue) *core.Spell {
			if rogue.disabledMCDs != nil {
				rogue.EnableAllCooldowns(rogue.disabledMCDs)
				rogue.disabledMCDs = nil
			}
			return nil
		},
	})
}

func RemainingAuraDuration(sim *core.Simulation, aura *core.Aura) time.Duration {
	if aura == nil {
		return time.Second * math.MaxInt32
	}
	remainingDuration := aura.RemainingDuration(sim)
	if remainingDuration < sim.GetRemainingDuration() {
		return remainingDuration
	} else {
		return time.Second * math.MaxInt32
	}
}

type RoguePriority struct {
	CastCount            int32
	GeneratedComboPoints int32
	IsFiller             bool
	MinComboPoints       int32
	ShouldCast           func(sim *core.Simulation, rogue *Rogue, prio *RoguePriority) bool
	Spell                func(rogue *Rogue) *core.Spell
}

func (prio *RoguePriority) NeedsComboPoints(rogue *Rogue) int32 {
	return prio.MinComboPoints - rogue.ComboPoints()
}

func (prio *RoguePriority) NeedsEnergy(rogue *Rogue) float64 {
	spell := prio.Spell(rogue)
	if spell == nil {
		return 0
	}
	return spell.DefaultCast.Cost - rogue.CurrentEnergy()
}

func (rogue *Rogue) getPriority(sim *core.Simulation) *RoguePriority {
	var nextPrio *RoguePriority
	for idx, prio := range rogue.PriorityList {
		if prio.ShouldCast == nil || prio.ShouldCast(sim, rogue, &prio) {
			needsCP := prio.NeedsComboPoints(rogue)
			needsEnergy := prio.NeedsEnergy(rogue)
			if nextPrio == nil || !prio.IsFiller {
				if needsCP <= 0 && needsEnergy <= 0 {
					return &rogue.PriorityList[idx]
				}
				if needsCP <= 0 {
					nextPrio = &rogue.PriorityList[idx]
				}
			} else {
				if needsCP <= 0 && needsEnergy <= 0 {
					if nextPrio.NeedsComboPoints(rogue) < 1 && prio.MinComboPoints < 1 {
						return &rogue.PriorityList[idx]
					}
				}
			}
		}
	}
	return nextPrio
}
