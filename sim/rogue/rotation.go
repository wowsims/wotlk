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

func (rogue *Rogue) SetPriorityList() {
	rogue.PriorityList = make([]RoguePriority, 0)
	rogue.PriorityList = append(rogue.PriorityList, RoguePriority{
		Aura:                 rogue.SliceAndDiceAura,
		CastCount:            0,
		Frequency:            proto.Rogue_Rotation_Maintain,
		GeneratedComboPoints: 0,
		MinComboPoints:       1,
		Spell: func(rogue *Rogue) *core.Spell {
			return rogue.SliceAndDice[rogue.ComboPoints()]
		},
	})
	if rogue.Rotation.ExposeArmorFrequency != proto.Rogue_Rotation_Never {
		exposeCP := int32(1)
		if rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Once {
			exposeCP = rogue.Rotation.MinimumComboPointsExposeArmor
		}
		rogue.PriorityList = append(rogue.PriorityList, RoguePriority{
			Aura:                 rogue.ExposeArmorAura,
			CastCount:            0,
			Frequency:            rogue.Rotation.ExposeArmorFrequency,
			GeneratedComboPoints: 0,
			MinComboPoints:       exposeCP,
			Spell: func(rogue *Rogue) *core.Spell {
				return rogue.ExposeArmor[rogue.ComboPoints()]
			},
		})
	}
	mutilatePrio := RoguePriority{
		Aura:                 nil,
		CastCount:            0,
		Frequency:            proto.Rogue_Rotation_Build,
		GeneratedComboPoints: 2,
		MinComboPoints:       0,
		Spell: func(rogue *Rogue) *core.Spell {
			return rogue.Mutilate
		},
	}
	sinisterStrikePrio := RoguePriority{
		Aura:                 nil,
		CastCount:            0,
		Frequency:            proto.Rogue_Rotation_Build,
		GeneratedComboPoints: 1,
		MinComboPoints:       0,
		Spell: func(rogue *Rogue) *core.Spell {
			return rogue.SinisterStrike
		},
	}
	envenomPrio := RoguePriority{
		Aura:                 rogue.EnvenomAura,
		CastCount:            0,
		Frequency:            proto.Rogue_Rotation_Maintain,
		GeneratedComboPoints: 0,
		MinComboPoints:       3,
		Spell: func(rogue *Rogue) *core.Spell {
			return rogue.Envenom[rogue.ComboPoints()]
		},
	}
	rupturePrio := RoguePriority{
		Aura:                 rogue.RuptureDot.Aura,
		CastCount:            0,
		Frequency:            proto.Rogue_Rotation_Maintain,
		GeneratedComboPoints: 0,
		MinComboPoints:       3,
		Spell: func(rogue *Rogue) *core.Spell {
			return rogue.Rupture[rogue.ComboPoints()]
		},
	}
	evisceratePrio := RoguePriority{
		Aura:                 nil,
		CastCount:            0,
		Frequency:            proto.Rogue_Rotation_Build,
		GeneratedComboPoints: 0,
		MinComboPoints:       3,
		Spell: func(rogue *Rogue) *core.Spell {
			return rogue.Eviscerate[rogue.ComboPoints()]
		},
	}
	if rogue.Talents.Mutilate {
		rogue.PriorityList = append(rogue.PriorityList, mutilatePrio)
		switch rogue.Rotation.AssassinationFinisherPriority {
		case proto.Rogue_Rotation_EnvenomRupture:
			envenomPrio.MinComboPoints = rogue.Rotation.MinimumComboPointsPrimaryFinisher
			rupturePrio.MinComboPoints = rogue.Rotation.MinimumComboPointsSecondaryFinisher
			rupturePrio.Frequency = proto.Rogue_Rotation_Fill
			rogue.PriorityList = append(rogue.PriorityList, envenomPrio)
			rogue.PriorityList = append(rogue.PriorityList, rupturePrio)
		case proto.Rogue_Rotation_RuptureEnvenom:
			rupturePrio.MinComboPoints = rogue.Rotation.MinimumComboPointsPrimaryFinisher
			envenomPrio.MinComboPoints = rogue.Rotation.MinimumComboPointsSecondaryFinisher
			envenomPrio.Frequency = proto.Rogue_Rotation_Fill
			rogue.PriorityList = append(rogue.PriorityList, rupturePrio)
			rogue.PriorityList = append(rogue.PriorityList, envenomPrio)
		}
	} else {
		rogue.PriorityList = append(rogue.PriorityList, sinisterStrikePrio)
		switch rogue.Rotation.CombatFinisherPriority {
		case proto.Rogue_Rotation_RuptureEviscerate:
			rupturePrio.MinComboPoints = rogue.Rotation.MinimumComboPointsPrimaryFinisher
			evisceratePrio.MinComboPoints = rogue.Rotation.MinimumComboPointsSecondaryFinisher
			evisceratePrio.Frequency = proto.Rogue_Rotation_Fill
			rogue.PriorityList = append(rogue.PriorityList, rupturePrio)
			rogue.PriorityList = append(rogue.PriorityList, evisceratePrio)
		case proto.Rogue_Rotation_EviscerateRupture:
			evisceratePrio.MinComboPoints = rogue.Rotation.MinimumComboPointsPrimaryFinisher
			rupturePrio.MinComboPoints = rogue.Rotation.MinimumComboPointsSecondaryFinisher
			rupturePrio.Frequency = proto.Rogue_Rotation_Fill
			rogue.PriorityList = append(rogue.PriorityList, evisceratePrio)
			rogue.PriorityList = append(rogue.PriorityList, rupturePrio)
		}
	}
	rogue.PriorityList = append(rogue.PriorityList, RoguePriority{
		Aura:                 nil,
		CastCount:            0,
		Frequency:            proto.Rogue_Rotation_Build,
		GeneratedComboPoints: 0,
		MinComboPoints:       0,
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
	Aura                 *core.Aura
	CastCount            int32
	Frequency            proto.Rogue_Rotation_Frequency
	GeneratedComboPoints int32
	MinComboPoints       int32
	Spell                func(rogue *Rogue) *core.Spell
}

func (prio *RoguePriority) ShouldCast(rogue *Rogue, sim *core.Simulation) bool {
	if prio.Frequency == proto.Rogue_Rotation_Never {
		return false
	}
	if prio.Frequency == proto.Rogue_Rotation_Once && prio.CastCount >= 1 {
		return false
	}
	if prio.Frequency == proto.Rogue_Rotation_Build && prio.GeneratedComboPoints+rogue.ComboPoints() > 5 {
		return false
	}
	if prio.Frequency == proto.Rogue_Rotation_Maintain && RemainingAuraDuration(sim, prio.Aura) >= 0 {
		return false
	}
	if prio.Frequency == proto.Rogue_Rotation_Fill && (prio.Aura != nil && RemainingAuraDuration(sim, prio.Aura) >= 0) {
		return false
	}
	return true
}

func (prio *RoguePriority) NeedsCombatPoints(rogue *Rogue) bool {
	return rogue.ComboPoints() < prio.MinComboPoints
}

func (prio *RoguePriority) NeedsEnergy(rogue *Rogue) bool {
	spell := prio.Spell(rogue)
	if spell == nil {
		return false
	}
	return rogue.CurrentEnergy() < spell.DefaultCast.Cost
}

func (rogue *Rogue) getPriority(sim *core.Simulation) *RoguePriority {
	var nextPrio *RoguePriority
	for idx, prio := range rogue.PriorityList {
		if prio.ShouldCast(rogue, sim) {
			needsCP := prio.NeedsCombatPoints(rogue)
			needsEnergy := prio.NeedsEnergy(rogue)
			if !needsCP && !needsEnergy && prio.Frequency != proto.Rogue_Rotation_Fill {
				return &rogue.PriorityList[idx]
			}
			if nextPrio == nil && !needsCP {
				nextPrio = &rogue.PriorityList[idx]
			}
		}
	}
	return nextPrio
}
