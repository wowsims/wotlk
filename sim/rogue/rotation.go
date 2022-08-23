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

func RemainingAuraDuration(sim *core.Simulation, aura *core.Aura) time.Duration {
	if aura == nil {
		return 0
	}
	return core.MaxDuration(0, aura.RemainingDuration(sim))
}

func (rogue *Rogue) ExpectedEnergyGain(sim *core.Simulation, duration time.Duration) float64 {
	adrenalineRushDuration := time.Second * 0
	if rogue.Talents.AdrenalineRush {
		adrenalineRushDuration = core.MinDuration(core.MaxDuration(rogue.AdrenalineRushAura.RemainingDuration(sim), 0), duration)
	}
	offHandWeaponSpeed := rogue.AutoAttacks.OH.SwingSpeed
	baseEnergyDuration := duration - adrenalineRushDuration
	adrenalineRushEnergyPerSecond := 12.5
	baseEnergyPerSecond := 10.0
	guaranteedEnergy := (adrenalineRushDuration.Seconds() * adrenalineRushEnergyPerSecond) + (baseEnergyDuration.Seconds() * baseEnergyPerSecond)
	combatPotencyEnergyPerSwing := float64(rogue.Talents.CombatPotency) * 3 * 0.2
	expectedCombatPotencyEnergy := (duration.Seconds() / offHandWeaponSpeed) * combatPotencyEnergyPerSwing
	return guaranteedEnergy + expectedCombatPotencyEnergy
}

func (rogue *Rogue) ExpectedComboPoints(sim *core.Simulation, duration time.Duration, builderEnergyCost float64, finisherEnergyCost float64) float64 {
	energyGained := rogue.ExpectedEnergyGain(sim, duration)
	availableEnergy := energyGained + rogue.CurrentEnergy() - finisherEnergyCost
	currentComboPoints := float64(rogue.ComboPoints())
	expectedComboPoints := currentComboPoints
	availableDuration := duration
	for availableDuration > 0 && ((currentComboPoints >= 4.5 && availableEnergy >= finisherEnergyCost) ||
		(currentComboPoints < 4.5 && availableEnergy >= builderEnergyCost)) {
		if currentComboPoints >= 4.5 {
			currentComboPoints -= 4.4
			availableEnergy -= finisherEnergyCost
		} else {
			currentComboPoints += 1
			availableEnergy -= builderEnergyCost
			expectedComboPoints += 1
		}
		availableDuration -= 1.0
	}
	return expectedComboPoints
}

func (rogue *Rogue) rotation(sim *core.Simulation) {
	rogue.updatePlan(sim)
	spell := rogue.CurrentPriority.GetSpell(sim, rogue)
	if spell == nil && (rogue.CurrentPriority.IsBuilder || rogue.CurrentPriority.IsFinisher) {
		panic("Builders and Finishers should always have a spell")
	} else if spell == nil {
		if rogue.CurrentPriority.IncrementIfNil {
			rogue.CurrentPriority.CastCount += 1
		}
		rogue.CurrentPriority = rogue.CurrentPriority.Prev
	} else {
		castSucceeded := false
		if rogue.CurrentPriority.IsFinisher {
			// Pool energy?
			if rogue.CurrentEnergy() < 100 && RemainingAuraDuration(sim, rogue.CurrentPriority.Aura) > 0 {
				castSucceeded = false
			} else {
				castSucceeded = spell.Cast(sim, rogue.CurrentTarget)
			}
		} else {
			castSucceeded = spell.Cast(sim, rogue.CurrentTarget)
		}
		if castSucceeded {
			rogue.CurrentPriority.CastCount += 1
			rogue.CurrentPriority = rogue.CurrentPriority.Prev
		}
	}
	if rogue.GCD.IsReady(sim) {
		rogue.DoNothing()
	}
}

type RoguePriority struct {
	Prev                 *RoguePriority
	CastCount            int32
	MaxCasts             int32
	IncrementIfNil       bool
	ComboPointsConsumed  float64
	Index                int
	MinimumComboPoints   float64
	IsBuilder            bool
	IsFinisher           bool
	Label                string
	Aura                 *core.Aura
	EnergyCost           float64
	ComboPointsGenerated float64
	GetSpell             func(*core.Simulation, *Rogue) *core.Spell
}

func (rogue *Rogue) SetMultiTargetPriorityList() {
	index := 0
	rogue.PriorityList = make([]RoguePriority, 0)

	if rogue.Rotation.MultiTargetSliceFrequency != proto.Rogue_Rotation_Never {

		// Slice And Dice
		sliceAndDice := RoguePriority{
			Index:              index,
			MinimumComboPoints: 1,
			IsBuilder:          false,
			IsFinisher:         true,
			Label:              "Slice and Dice",
			Aura:               rogue.SliceAndDiceAura,
			EnergyCost:         rogue.SliceAndDice[1].DefaultCast.Cost,
			GetSpell: func(s *core.Simulation, r *Rogue) *core.Spell {
				comboPoints := rogue.ComboPoints()
				if comboPoints < 1 {
					return nil
				}
				return rogue.SliceAndDice[comboPoints]
			},
		}
		if rogue.Rotation.MultiTargetSliceFrequency == proto.Rogue_Rotation_Once {
			sliceAndDice.MaxCasts = 1
			sliceAndDice.MinimumComboPoints = core.MaxFloat(1.0, float64(rogue.Rotation.MinimumComboPointsMultiTargetSlice))
		}
		rogue.PriorityList = append(rogue.PriorityList, sliceAndDice)
		index += 1

		if rogue.CanMutilate() {
			mutilate := RoguePriority{
				MinimumComboPoints:   0,
				IsBuilder:            true,
				IsFinisher:           false,
				Label:                "Mutilate",
				ComboPointsGenerated: 2,
				EnergyCost:           rogue.Mutilate.DefaultCast.Cost,
				GetSpell: func(s *core.Simulation, r *Rogue) *core.Spell {
					return rogue.Mutilate
				},
			}
			rogue.PriorityList = append(rogue.PriorityList, mutilate)
			index += 1
		} else {
			sinisterStrike := RoguePriority{
				MinimumComboPoints:   0,
				IsBuilder:            true,
				IsFinisher:           false,
				Label:                "Sinister Strike",
				ComboPointsGenerated: 1,
				EnergyCost:           rogue.SinisterStrike.DefaultCast.Cost,
				GetSpell: func(s *core.Simulation, r *Rogue) *core.Spell {
					return rogue.SinisterStrike
				},
			}
			rogue.PriorityList = append(rogue.PriorityList, sinisterStrike)
			index += 1
		}
	}
	if rogue.Talents.HungerForBlood {
		hungerForBlood := RoguePriority{
			IsBuilder:  false,
			IsFinisher: false,
			Label:      "Hunger for Blood",
			EnergyCost: rogue.HungerForBlood.DefaultCast.Cost,
			Aura:       rogue.HungerForBloodAura,
			GetSpell: func(s *core.Simulation, r *Rogue) *core.Spell {
				return rogue.HungerForBlood
			},
		}
		rogue.PriorityList = append(rogue.PriorityList, hungerForBlood)
		index += 1
	}

	// Dummy thing to enable cooldowns
	cooldownEnabler := RoguePriority{
		Index:          index,
		Label:          "Enable Cooldowns",
		MaxCasts:       1,
		IncrementIfNil: true,
		GetSpell: func(s *core.Simulation, r *Rogue) *core.Spell {
			if rogue.disabledMCDs != nil {
				rogue.EnableAllCooldowns(rogue.disabledMCDs)
				rogue.disabledMCDs = nil
			}
			return nil
		},
	}
	rogue.PriorityList = append(rogue.PriorityList, cooldownEnabler)
	index += 1

	fanOfKnives := RoguePriority{
		Label:      "Fan of Knives",
		IsBuilder:  false,
		IsFinisher: false,
		EnergyCost: rogue.FanOfKnives.DefaultCast.Cost,
		GetSpell: func(s *core.Simulation, r *Rogue) *core.Spell {
			return rogue.FanOfKnives
		},
	}
	rogue.PriorityList = append(rogue.PriorityList, fanOfKnives)
	index += 1
}

func (rogue *Rogue) SetPriorityList(sim *core.Simulation) {
	if rogue.CanMutilate() {
		rogue.Builder = rogue.Mutilate
	} else {
		rogue.Builder = rogue.SinisterStrike
	}
	if sim.GetNumTargets() > 3 {
		rogue.SetMultiTargetPriorityList()
	} else {
		rogue.SetStandardPriorityList()
	}
}

func (rogue *Rogue) SetStandardPriorityList() {
	index := 0
	rogue.PriorityList = make([]RoguePriority, 0)

	// Slice And Dice
	sliceAndDice := RoguePriority{
		Index:              index,
		MinimumComboPoints: 1,
		IsBuilder:          false,
		IsFinisher:         true,
		Label:              "Slice and Dice",
		Aura:               rogue.SliceAndDiceAura,
		EnergyCost:         rogue.SliceAndDice[1].DefaultCast.Cost,
		GetSpell: func(s *core.Simulation, r *Rogue) *core.Spell {
			comboPoints := rogue.ComboPoints()
			if comboPoints < 1 {
				return nil
			}
			return rogue.SliceAndDice[comboPoints]
		},
	}
	rogue.PriorityList = append(rogue.PriorityList, sliceAndDice)
	index += 1

	// Expose Armor
	exposeArmor := RoguePriority{
		MinimumComboPoints: 1,
		IsBuilder:          false,
		IsFinisher:         true,
		Label:              "Expose Armor",
		Aura:               rogue.ExposeArmorAura,
		EnergyCost:         rogue.ExposeArmor[1].DefaultCast.Cost,
		GetSpell: func(s *core.Simulation, r *Rogue) *core.Spell {
			comboPoints := rogue.ComboPoints()
			if comboPoints < 1 {
				return nil
			}
			return rogue.ExposeArmor[comboPoints]
		},
	}
	if rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Maintain {
		exposeArmor.Index = index
		rogue.PriorityList = append(rogue.PriorityList, exposeArmor)
		index += 1
	} else if rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Once {
		exposeArmor.Index = index
		exposeArmor.MaxCasts = 1
		exposeArmor.MinimumComboPoints = float64(rogue.Rotation.MinimumComboPointsExposeArmor)
		rogue.PriorityList = append(rogue.PriorityList, exposeArmor)
		index += 1
	}

	if rogue.Talents.HungerForBlood {
		hungerForBlood := RoguePriority{
			IsBuilder:  false,
			IsFinisher: false,
			Label:      "Hunger for Blood",
			EnergyCost: rogue.HungerForBlood.DefaultCast.Cost,
			Aura:       rogue.HungerForBloodAura,
			GetSpell: func(s *core.Simulation, r *Rogue) *core.Spell {
				return rogue.HungerForBlood
			},
		}
		rogue.PriorityList = append(rogue.PriorityList, hungerForBlood)
		index += 1
	}

	// Dummy thing to enable cooldowns
	cooldownEnabler := RoguePriority{
		Index:          index,
		Label:          "Enable Cooldowns",
		MaxCasts:       1,
		IncrementIfNil: true,
		GetSpell: func(s *core.Simulation, r *Rogue) *core.Spell {
			if rogue.disabledMCDs != nil {
				rogue.EnableAllCooldowns(rogue.disabledMCDs)
				rogue.disabledMCDs = nil
			}
			return nil
		},
	}
	rogue.PriorityList = append(rogue.PriorityList, cooldownEnabler)
	index += 1

	rupture := RoguePriority{
		MinimumComboPoints: 3,
		IsBuilder:          false,
		IsFinisher:         true,
		Label:              "Rupture",
		Aura:               rogue.RuptureDot.Aura,
		EnergyCost:         rogue.Rupture[1].DefaultCast.Cost,
		GetSpell: func(s *core.Simulation, r *Rogue) *core.Spell {
			comboPoints := rogue.ComboPoints()
			if comboPoints < 1 {
				return nil
			}
			return rogue.Rupture[comboPoints]
		},
	}

	eviscerate := RoguePriority{
		MinimumComboPoints: 3,
		IsBuilder:          false,
		IsFinisher:         true,
		Label:              "Eviscerate",
		EnergyCost:         rogue.Eviscerate[1].DefaultCast.Cost,
		GetSpell: func(s *core.Simulation, r *Rogue) *core.Spell {
			comboPoints := rogue.ComboPoints()
			if comboPoints < 1 {
				return nil
			}
			return rogue.Eviscerate[comboPoints]
		},
	}

	if rogue.Talents.MasterPoisoner > 0 || rogue.Talents.CutToTheChase > 0 {
		envenom := RoguePriority{
			MinimumComboPoints: 3,
			IsBuilder:          false,
			IsFinisher:         true,
			Label:              "Envenom",
			Aura:               rogue.EnvenomAura,
			EnergyCost:         rogue.Envenom[1].DefaultCast.Cost,
			GetSpell: func(s *core.Simulation, r *Rogue) *core.Spell {
				comboPoints := rogue.ComboPoints()
				if comboPoints < 1 {
					return nil
				}
				return rogue.Envenom[comboPoints]
			},
		}
		switch rogue.Rotation.AssassinationFinisherPriority {
		case proto.Rogue_Rotation_EnvenomRupture:
			envenom.Index = index
			envenom.MinimumComboPoints = float64(rogue.Rotation.MinimumComboPointsPrimaryFinisher)
			rogue.PriorityList = append(rogue.PriorityList, envenom)
			index += 1
			rupture.Index = index
			rupture.MinimumComboPoints = float64(rogue.Rotation.MinimumComboPointsSecondaryFinisher)
			rogue.PriorityList = append(rogue.PriorityList, rupture)
			index += 1
		case proto.Rogue_Rotation_RuptureEnvenom:
			rupture.Index = index
			rupture.MinimumComboPoints = float64(rogue.Rotation.MinimumComboPointsPrimaryFinisher)
			rogue.PriorityList = append(rogue.PriorityList, rupture)
			index += 1
			envenom.Index = index
			envenom.MinimumComboPoints = float64(rogue.Rotation.MinimumComboPointsSecondaryFinisher)
			rogue.PriorityList = append(rogue.PriorityList, envenom)
			index += 1
		}
	} else {
		switch rogue.Rotation.CombatFinisherPriority {
		case proto.Rogue_Rotation_RuptureEviscerate:
			rupture.Index = index
			rupture.MinimumComboPoints = float64(rogue.Rotation.MinimumComboPointsPrimaryFinisher)
			rogue.PriorityList = append(rogue.PriorityList, rupture)
			index += 1
			eviscerate.Index = index
			eviscerate.MinimumComboPoints = float64(rogue.Rotation.MinimumComboPointsSecondaryFinisher)
			rogue.PriorityList = append(rogue.PriorityList, eviscerate)
			index += 1
		case proto.Rogue_Rotation_EviscerateRupture:
			eviscerate.Index = index
			eviscerate.MinimumComboPoints = float64(rogue.Rotation.MinimumComboPointsPrimaryFinisher)
			rogue.PriorityList = append(rogue.PriorityList, eviscerate)
			index += 1
			rupture.Index = index
			rupture.MinimumComboPoints = float64(rogue.Rotation.MinimumComboPointsSecondaryFinisher)
			rogue.PriorityList = append(rogue.PriorityList, rupture)
			index += 1
		}
	}

	if rogue.CanMutilate() {
		mutilate := RoguePriority{
			MinimumComboPoints:   0,
			IsBuilder:            true,
			IsFinisher:           false,
			Label:                "Mutilate",
			ComboPointsGenerated: 2,
			EnergyCost:           rogue.Mutilate.DefaultCast.Cost,
			GetSpell: func(s *core.Simulation, r *Rogue) *core.Spell {
				return rogue.Mutilate
			},
		}
		rogue.PriorityList = append(rogue.PriorityList, mutilate)
		index += 1
	} else {
		sinisterStrike := RoguePriority{
			MinimumComboPoints:   0,
			IsBuilder:            true,
			IsFinisher:           false,
			Label:                "Sinister Strike",
			ComboPointsGenerated: 1,
			EnergyCost:           rogue.SinisterStrike.DefaultCast.Cost,
			GetSpell: func(s *core.Simulation, r *Rogue) *core.Spell {
				return rogue.SinisterStrike
			},
		}
		rogue.PriorityList = append(rogue.PriorityList, sinisterStrike)
		index += 1
	}
}

func (rogue *Rogue) MinimalComboPointsNeededForPlan() float64 {
	prio := rogue.CurrentPriority
	if prio == nil {
		return 0
	}
	cpNeeded := prio.MinimumComboPoints
	for prio.Prev != nil {
		prio = prio.Prev
		cpNeeded += prio.MinimumComboPoints
	}
	return cpNeeded
}

func (rogue *Rogue) IdealComboPointsNeededForPlan() float64 {
	prio := rogue.CurrentPriority
	if prio == nil {
		return 0
	}
	cpNeeded := prio.ComboPointsConsumed
	for prio.Prev != nil {
		prio = prio.Prev
		cpNeeded += prio.ComboPointsConsumed
	}
	return cpNeeded
}

func (rogue *Rogue) DurationOfCurrentPlan(sim *core.Simulation) time.Duration {
	prio := rogue.CurrentPriority
	if prio == nil {
		return 0
	}
	for prio.Prev != nil {
		prio = prio.Prev
	}
	return RemainingAuraDuration(sim, prio.Aura)
}

func (rogue *Rogue) DurationUntilNextExpiration(sim *core.Simulation) (bool, time.Duration) {
	prio := rogue.CurrentPriority
	hasExpiration := false
	duration := time.Second * math.MaxInt32
	for prio != nil {
		if prio.Aura != nil {
			hasExpiration = true
			duration = core.MinDuration(duration, RemainingAuraDuration(sim, prio.Aura))
		}
		prio = prio.Prev
	}
	return hasExpiration, duration
}

func (rogue *Rogue) verifyPlan(sim *core.Simulation) (bool, float64) {
	prio := rogue.CurrentPriority
	pointsSpent := 0.0
	error := 0.0
	for prio != nil {
		if prio.Aura != nil {
			pointsNeeded := prio.ComboPointsConsumed + pointsSpent
			neededBy := RemainingAuraDuration(sim, prio.Aura)
			expectedPoints := rogue.ExpectedComboPoints(sim, neededBy, rogue.Builder.DefaultCast.Cost, prio.EnergyCost)
			if expectedPoints < pointsNeeded {
				delta := pointsNeeded - expectedPoints
				if delta > prio.ComboPointsConsumed-prio.MinimumComboPoints {
					return false, error
				}
				// We can reduce consumed points to make up the difference
				pointsNeeded -= delta
				error += delta
			}
			pointsSpent += pointsNeeded
		}
		prio = prio.Prev
	}
	return true, error
}

func (rogue *Rogue) HasPriorityWithInactiveAura(sim *core.Simulation) bool {
	prio := rogue.CurrentPriority
	for prio != nil {
		if prio.Aura != nil && !prio.Aura.IsActive() {
			return true
		}
		prio = prio.Prev
	}
	return false
}

func (rogue *Rogue) updatePlan(sim *core.Simulation) {
	currentIndex := 0
	builderEnergyCost := rogue.Builder.DefaultCast.Cost

	// Verify current plan
	if rogue.CurrentPriority != nil {
		validPlan, pointDelta := rogue.verifyPlan(sim)
		if !validPlan {
			rogue.CurrentPriority = nil
		} else if pointDelta > 0 {
			if rogue.CurrentPriority.ComboPointsConsumed-pointDelta >= rogue.CurrentPriority.MinimumComboPoints {
				rogue.CurrentPriority.ComboPointsConsumed -= pointDelta
			} else {
				rogue.CurrentPriority = nil
			}
		}
		if rogue.CurrentPriority != nil && rogue.CurrentPriority.IsBuilder {
			return
		}
	}
	if rogue.CurrentPriority != nil {
		currentIndex = rogue.CurrentPriority.Index + 1
	}

	// Add new priorities
	freezeFinishers := false
	for nextIndex, prio := range rogue.PriorityList[currentIndex:] {
		prioIndex := currentIndex + nextIndex

		// Filter out cast-limited priorities
		if prio.MaxCasts > 0 && prio.CastCount >= prio.MaxCasts {
			continue
		}
		if prio.IsFinisher && !freezeFinishers {
			if rogue.HasPriorityWithInactiveAura(sim) {
				freezeFinishers = true
				continue
			}
			expectedPoints := 0.0
			idealConsumedPoints := 0.0
			minimalConsumedPoints := 0.0
			energyCost := prio.EnergyCost
			if rogue.CurrentPriority != nil {
				minimalConsumedPoints = rogue.MinimalComboPointsNeededForPlan()
				idealConsumedPoints = rogue.IdealComboPointsNeededForPlan()
				expectedPoints = rogue.ExpectedComboPoints(sim, rogue.DurationOfCurrentPlan(sim), builderEnergyCost, energyCost)
			} else {
				expectedPoints = rogue.ExpectedComboPoints(sim, RemainingAuraDuration(sim, prio.Aura), builderEnergyCost, energyCost)
			}
			hasInactiveAura := prio.Aura != nil && !prio.Aura.IsActive()
			canPrioritizeIdeal := expectedPoints >= (idealConsumedPoints+prio.MinimumComboPoints) && (prio.Aura == nil || prio.Aura.IsActive())
			shouldPrioritize := expectedPoints >= (minimalConsumedPoints+prio.MinimumComboPoints) && hasInactiveAura
			if rogue.CurrentPriority == nil || canPrioritizeIdeal || shouldPrioritize {
				nextPrio := &rogue.PriorityList[prioIndex]
				nextPrio.Prev = rogue.CurrentPriority
				pointsToUse := 0.0
				if canPrioritizeIdeal {
					pointsToUse = core.MinFloat(expectedPoints-idealConsumedPoints, 5)
				} else {
					pointsToUse = core.MaxFloat(prio.MinimumComboPoints, float64(rogue.ComboPoints()))
				}
				nextPrio.ComboPointsConsumed = pointsToUse
				if nextPrio.ComboPointsConsumed <= 0 {
					panic("Consumed points of finishers should always be greater than 0")
				}
				rogue.CurrentPriority = nextPrio
			} else {
				freezeFinishers = true
			}
		}
		if !prio.IsBuilder && !prio.IsFinisher && (rogue.CurrentPriority == nil || rogue.CurrentPriority.Aura.IsActive()) {
			// Hunger for Blood
			if prio.Aura != nil && !prio.Aura.IsActive() {
				nextPrio := &rogue.PriorityList[prioIndex]
				nextPrio.Prev = rogue.CurrentPriority
				rogue.CurrentPriority = nextPrio
			} else if prio.Aura == nil {
				nextPrio := &rogue.PriorityList[prioIndex]
				nextPrio.Prev = rogue.CurrentPriority
				rogue.CurrentPriority = nextPrio
			}
		}
		if prio.IsBuilder {
			if rogue.CurrentPriority == nil || !rogue.CurrentPriority.IsFinisher {
				continue
			}
			if float64(rogue.ComboPoints()) < core.MinFloat(5, rogue.CurrentPriority.ComboPointsConsumed) {
				nextPrio := &rogue.PriorityList[prioIndex]
				nextPrio.Prev = rogue.CurrentPriority
				rogue.CurrentPriority = nextPrio
			}
		}
	}
}
