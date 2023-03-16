package rogue

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (rogue *Rogue) OnEnergyGain(sim *core.Simulation) {
	if rogue.Talents.Mutilate && sim.GetNumTargets() <= 3 {
		rogue.OnCanAct(sim)
		return
	}
	if rogue.KillingSpreeAura.IsActive() {
		rogue.DoNothing()
		return
	}
	rogue.TryUseCooldowns(sim)
	if rogue.GCD.IsReady(sim) {
		if rogue.Talents.HonorAmongThieves > 0 {
			rogue.doSubtletyRotation(sim)
		} else {
			rogue.rotation(sim)
		}
	}
}

func (rogue *Rogue) OnGCDReady(sim *core.Simulation) {
	if rogue.Talents.Mutilate && sim.GetNumTargets() <= 3 {
		rogue.OnCanAct(sim)
		return
	}
	if rogue.Talents.HonorAmongThieves > 0 && sim.GetNumTargets() <= 3 {
		rogue.OnCanActSubtlety(sim)
		return
	}
	if rogue.KillingSpreeAura.IsActive() {
		rogue.DoNothing()
		return
	}
	rogue.TryUseCooldowns(sim)
	if rogue.IsWaitingForEnergy() {
		rogue.DoNothing()
		return
	}
	if rogue.GCD.IsReady(sim) {
		rogue.rotation(sim)
	}
}

type rogueRotationItem struct {
	ExpiresAt            time.Duration
	MinimumBuildDuration time.Duration
	MaximumBuildDuration time.Duration
	PrioIndex            int
}

type roguePriorityItem struct {
	Aura               *core.Aura
	CastCount          int32
	EnergyCost         float64
	PoolAmount         float64
	GetDuration        func(*Rogue, int32) time.Duration
	GetSpell           func(*Rogue, int32) *core.Spell
	IsFiller           bool
	MaximumComboPoints int32
	MaxCasts           int32
	MinimumComboPoints int32
}

type shouldCastRotationItemResult int32

const (
	ShouldNotCast shouldCastRotationItemResult = iota
	ShouldBuild
	ShouldCast
	ShouldWait
)

func (rogue *Rogue) energyToBuild(points int32) float64 {
	costPerBuilder := rogue.Builder.DefaultCast.Cost

	buildersNeeded := math.Ceil(float64(points) / float64(rogue.BuilderPoints))
	return buildersNeeded * costPerBuilder
}

func (rogue *Rogue) timeToBuild(_ *core.Simulation, points int32, builderPoints int32, eps float64, finisherCost float64) time.Duration {
	energyNeeded := rogue.energyToBuild(points) + finisherCost
	secondsNeeded := energyNeeded / eps
	globalsNeeded := math.Ceil(float64(points)/float64(builderPoints)) + 1
	// Return greater of the time it takes to use the globals and the time it takes to build the energy
	return core.MaxDuration(time.Second*time.Duration(secondsNeeded), time.Second*time.Duration(globalsNeeded))
}

func (rogue *Rogue) shouldCastNextRotationItem(sim *core.Simulation, eps float64) shouldCastRotationItemResult {
	if len(rogue.rotationItems) == 0 {
		panic("Empty rotation")
	}
	currentEnergy := rogue.CurrentEnergy()
	comboPoints := rogue.ComboPoints()
	currentTime := sim.CurrentTime
	item := rogue.rotationItems[0]
	prio := rogue.priorityItems[item.PrioIndex]
	tte := item.ExpiresAt - currentTime
	clippingThreshold := time.Second * 2
	timeUntilNextGCD := rogue.GCD.TimeToReady(sim)

	// Check to see if a higher prio item will expire
	if len(rogue.rotationItems) >= 2 {
		timeElapsed := time.Second * 1
		for _, nextItem := range rogue.rotationItems[1:] {
			if nextItem.ExpiresAt <= currentTime+timeElapsed {
				return ShouldNotCast
			}
			timeElapsed += nextItem.MinimumBuildDuration
		}
	}

	// Expires before next GCD
	if tte <= timeUntilNextGCD {
		if comboPoints >= prio.MinimumComboPoints && currentEnergy >= (prio.EnergyCost+prio.PoolAmount) {
			return ShouldCast
		} else if comboPoints < prio.MinimumComboPoints && currentEnergy >= rogue.Builder.DefaultCast.Cost {
			return ShouldBuild
		} else {
			return ShouldWait
		}
	}
	if comboPoints >= prio.MaximumComboPoints { // Don't need CP
		// Cast
		if tte <= clippingThreshold && currentEnergy >= (prio.EnergyCost+prio.PoolAmount) {
			return ShouldCast
		}
		// Pool energy
		if tte <= clippingThreshold && currentEnergy < (prio.EnergyCost+prio.PoolAmount) {
			return ShouldWait
		}
		// We have time to squeeze in another spell
		if tte > item.MinimumBuildDuration {
			// Find the first lower prio item that can be cast and use it
			for lpi, lowerPrio := range rogue.priorityItems[item.PrioIndex+1:] {
				if comboPoints > lowerPrio.MinimumComboPoints && currentEnergy > lowerPrio.EnergyCost && lowerPrio.MaxCasts == 0 {
					rogue.rotationItems = append([]rogueRotationItem{
						{ExpiresAt: currentTime, PrioIndex: lpi + item.PrioIndex + 1},
					}, rogue.rotationItems...)
					return ShouldCast
				}
			}
		}
		// Overcap CP with builder
		if rogue.timeToBuild(sim, 1, rogue.BuilderPoints, eps, prio.EnergyCost+prio.PoolAmount) <= tte && currentEnergy >= rogue.Builder.DefaultCast.Cost {
			return ShouldBuild
		}
	} else if comboPoints < prio.MinimumComboPoints { // Need CP
		if currentEnergy >= rogue.Builder.DefaultCast.Cost {
			return ShouldBuild
		} else {
			return ShouldWait
		}
	} else { // Between MinimumComboPoints and MaximumComboPoints
		if currentEnergy >= prio.EnergyCost+prio.PoolAmount && tte <= timeUntilNextGCD {
			return ShouldCast
		}
		ttb := rogue.timeToBuild(sim, 1, 2, eps, prio.EnergyCost+prio.PoolAmount-currentEnergy)
		if currentEnergy >= rogue.Builder.DefaultCast.Cost && tte > ttb {
			return ShouldBuild
		}
	}
	return ShouldWait
}

func (rogue *Rogue) planRotation(sim *core.Simulation) []rogueRotationItem {
	var rotationItems []rogueRotationItem
	eps := rogue.getExpectedEnergyPerSecond()
	for pi, prio := range rogue.priorityItems {
		if prio.MaxCasts > 0 && prio.CastCount >= prio.MaxCasts {
			continue
		}
		maxCP := prio.MaximumComboPoints
		for maxCP > 0 && prio.GetDuration(rogue, maxCP)+sim.CurrentTime > sim.Duration {
			maxCP--
		}
		var expiresAt time.Duration
		if prio.Aura != nil {
			expiresAt = prio.Aura.ExpiresAt()
		} else if prio.MaxCasts == 1 {
			expiresAt = sim.CurrentTime // TODO looks fishy, repeated expiresAt = sim.CurrentTime
		} else {
			expiresAt = sim.CurrentTime
		}
		minimumBuildDuration := rogue.timeToBuild(sim, prio.MinimumComboPoints, rogue.BuilderPoints, eps, prio.EnergyCost)
		maximumBuildDuration := rogue.timeToBuild(sim, maxCP, rogue.BuilderPoints, eps, prio.EnergyCost)
		rotationItems = append(rotationItems, rogueRotationItem{
			ExpiresAt:            expiresAt,
			MaximumBuildDuration: maximumBuildDuration,
			MinimumBuildDuration: minimumBuildDuration,
			PrioIndex:            pi,
		})
	}

	currentTime := sim.CurrentTime
	comboPoints := rogue.ComboPoints()
	currentEnergy := rogue.CurrentEnergy()

	var prioStack []rogueRotationItem
	for _, item := range rotationItems {
		if item.ExpiresAt >= sim.Duration {
			continue
		}
		prio := rogue.priorityItems[item.PrioIndex]
		maxBuildAt := item.ExpiresAt - item.MaximumBuildDuration
		if prio.Aura == nil {
			timeValueOfResources := time.Duration((float64(comboPoints)*rogue.Builder.DefaultCast.Cost/float64(rogue.BuilderPoints) + currentEnergy) / eps)
			maxBuildAt = currentTime - item.MaximumBuildDuration - timeValueOfResources
		}
		if currentTime < maxBuildAt {
			// Put it on the to cast stack
			prioStack = append(prioStack, item)
			if prio.MinimumComboPoints > 0 {
				comboPoints = 0
			}
			currentTime += item.MaximumBuildDuration
		} else {
			cpUsed := core.MaxInt32(0, prio.MinimumComboPoints-comboPoints)
			energyUsed := core.MaxFloat(0, prio.EnergyCost-currentEnergy)
			minBuildTime := rogue.timeToBuild(sim, cpUsed, rogue.BuilderPoints, eps, energyUsed)
			if currentTime+minBuildTime <= item.ExpiresAt || !prio.IsFiller {
				prioStack = append(prioStack, item)
				currentTime = core.MaxDuration(item.ExpiresAt, currentTime+minBuildTime)
				currentEnergy = 0
				if prio.MinimumComboPoints > 0 {
					comboPoints = 0
				}
			} else if len(prioStack) < 1 || (prio.Aura != nil && !prio.Aura.IsActive() && !prio.IsFiller) || prio.MaxCasts == 1 {
				// Plan to cast it as soon as possible
				prioStack = append(prioStack, item)
				currentTime += item.MinimumBuildDuration
				currentEnergy = 0
				if prio.MinimumComboPoints > 0 {
					comboPoints = 0
				}
			}
		}
	}

	// Reverse
	for i, j := 0, len(prioStack)-1; i < j; i, j = i+1, j-1 {
		prioStack[i], prioStack[j] = prioStack[j], prioStack[i]
	}

	return prioStack
}

func (rogue *Rogue) setPriorityItems(sim *core.Simulation) {
	rogue.Builder = rogue.SinisterStrike
	rogue.BuilderPoints = 1
	if rogue.PrimaryTalentTree == AssassinTree {
		if rogue.CanMutilate() {
			rogue.Builder = rogue.Mutilate
			rogue.BuilderPoints = 2
		}
		rogue.setupAssassinationRotation(sim)
	}
	if rogue.PrimaryTalentTree == SubtletyTree {
		rogue.setSubtletyBuilder(sim)
		rogue.setupSubtletyRotation(sim)
	}
	isMultiTarget := sim.GetNumTargets() >= 3
	// Slice and Dice
	rogue.priorityItems = rogue.priorityItems[:0]

	sliceAndDice := roguePriorityItem{
		MinimumComboPoints: 1,
		MaximumComboPoints: 5,
		Aura:               rogue.SliceAndDiceAura,
		EnergyCost:         rogue.SliceAndDice.DefaultCast.Cost,
		GetDuration: func(rogue *Rogue, cp int32) time.Duration {
			return rogue.sliceAndDiceDurations[cp]
		},
		GetSpell: func(rogue *Rogue, cp int32) *core.Spell {
			return rogue.SliceAndDice
		},
	}
	if isMultiTarget {
		if rogue.Rotation.MultiTargetSliceFrequency != proto.Rogue_Rotation_Never {
			sliceAndDice.MinimumComboPoints = core.MaxInt32(1, rogue.Rotation.MinimumComboPointsMultiTargetSlice)
			if rogue.Rotation.MultiTargetSliceFrequency == proto.Rogue_Rotation_Once {
				sliceAndDice.MaxCasts = 1
			}
			rogue.priorityItems = append(rogue.priorityItems, sliceAndDice)
		}
	} else {
		rogue.priorityItems = append(rogue.priorityItems, sliceAndDice)
	}

	// Expose Armor
	if rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Maintain ||
		rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Once {
		minPoints := int32(1)
		maxCasts := int32(0)
		if rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Once {
			minPoints = rogue.Rotation.MinimumComboPointsExposeArmor
			maxCasts = 1
		}
		rogue.priorityItems = append(rogue.priorityItems, roguePriorityItem{
			MaxCasts:           maxCasts,
			MaximumComboPoints: 5,
			MinimumComboPoints: minPoints,
			Aura:               rogue.ExposeArmorAuras.Get(rogue.CurrentTarget),
			EnergyCost:         rogue.ExposeArmor.DefaultCast.Cost,
			GetDuration: func(rogue *Rogue, cp int32) time.Duration {
				return rogue.exposeArmorDurations[cp]
			},
			GetSpell: func(rogue *Rogue, cp int32) *core.Spell {
				return rogue.ExposeArmor
			},
		})
	}

	// Hunger for Blood
	if rogue.Talents.HungerForBlood {
		rogue.priorityItems = append(rogue.priorityItems, roguePriorityItem{
			MaximumComboPoints: 0,
			Aura:               rogue.HungerForBloodAura,
			EnergyCost:         rogue.HungerForBlood.DefaultCast.Cost,
			GetDuration: func(rogue *Rogue, cp int32) time.Duration {
				return rogue.HungerForBloodAura.Duration
			},
			GetSpell: func(rogue *Rogue, cp int32) *core.Spell {
				return rogue.HungerForBlood
			},
		})
	}

	// Dummy priority to enable CDs
	rogue.priorityItems = append(rogue.priorityItems, roguePriorityItem{
		MaxCasts:           1,
		MaximumComboPoints: 0,
		GetDuration: func(rogue *Rogue, cp int32) time.Duration {
			return 0
		},
		GetSpell: func(rogue *Rogue, cp int32) *core.Spell {
			if rogue.allMCDsDisabled {
				for _, mcd := range rogue.GetMajorCooldowns() {
					mcd.Enable()
				}
				rogue.allMCDsDisabled = false
			}
			return nil
		},
	})

	// Rupture
	rupture := roguePriorityItem{
		MinimumComboPoints: 3,
		MaximumComboPoints: 5,
		Aura:               rogue.Rupture.CurDot().Aura,
		EnergyCost:         rogue.Rupture.DefaultCast.Cost,
		GetDuration: func(rogue *Rogue, cp int32) time.Duration {
			return rogue.RuptureDuration(cp)
		},
		GetSpell: func(rogue *Rogue, cp int32) *core.Spell {
			return rogue.Rupture
		},
	}

	// Eviscerate
	eviscerate := roguePriorityItem{
		MinimumComboPoints: 1,
		MaximumComboPoints: 5,
		EnergyCost:         rogue.Eviscerate.DefaultCast.Cost,
		GetDuration: func(rogue *Rogue, cp int32) time.Duration {
			return 0
		},
		GetSpell: func(rogue *Rogue, cp int32) *core.Spell {
			return rogue.Eviscerate
		},
	}

	if isMultiTarget {
		rogue.priorityItems = append(rogue.priorityItems, roguePriorityItem{
			MaximumComboPoints: 0,
			EnergyCost:         rogue.FanOfKnives.DefaultCast.Cost,
			GetSpell: func(rogue *Rogue, i int32) *core.Spell {
				return rogue.FanOfKnives
			},
		})

	} else if rogue.Talents.MasterPoisoner > 0 || rogue.Talents.CutToTheChase > 0 {
		// Envenom
		envenom := roguePriorityItem{
			MinimumComboPoints: 1,
			MaximumComboPoints: 5,
			Aura:               rogue.EnvenomAura,
			EnergyCost:         rogue.Envenom.DefaultCast.Cost,
			GetDuration: func(rogue *Rogue, cp int32) time.Duration {
				return rogue.EnvenomAura.Duration
			},
			GetSpell: func(rogue *Rogue, cp int32) *core.Spell {
				return rogue.Envenom
			},
		}
		switch rogue.Rotation.AssassinationFinisherPriority {
		case proto.Rogue_Rotation_EnvenomRupture:
			envenom.MinimumComboPoints = core.MaxInt32(1, rogue.Rotation.MinimumComboPointsPrimaryFinisher)
			rogue.priorityItems = append(rogue.priorityItems, envenom)
			rupture.MinimumComboPoints = rogue.Rotation.MinimumComboPointsSecondaryFinisher
			rupture.IsFiller = true
			if rupture.MinimumComboPoints > 0 && rupture.MinimumComboPoints <= 5 {
				rogue.priorityItems = append(rogue.priorityItems, rupture)
			}
		case proto.Rogue_Rotation_RuptureEnvenom:
			rupture.MinimumComboPoints = core.MaxInt32(1, rogue.Rotation.MinimumComboPointsPrimaryFinisher)
			rogue.priorityItems = append(rogue.priorityItems, rupture)
			envenom.MinimumComboPoints = rogue.Rotation.MinimumComboPointsSecondaryFinisher
			envenom.IsFiller = true
			if envenom.MinimumComboPoints > 0 && envenom.MinimumComboPoints <= 5 {
				rogue.priorityItems = append(rogue.priorityItems, envenom)
			}
		}
		eviscerate.IsFiller = true
		eviscerate.MinimumComboPoints = 1
		rogue.priorityItems = append(rogue.priorityItems, eviscerate)
	} else {
		if rogue.PrimaryTalentTree == CombatTree {
			switch rogue.Rotation.CombatFinisherPriority {
			case proto.Rogue_Rotation_RuptureEviscerate:
				rupture.MinimumComboPoints = core.MaxInt32(1, rogue.Rotation.MinimumComboPointsPrimaryFinisher)
				rogue.priorityItems = append(rogue.priorityItems, rupture)
				eviscerate.MinimumComboPoints = core.MaxInt32(1, rogue.Rotation.MinimumComboPointsSecondaryFinisher)
				eviscerate.IsFiller = true
				rogue.priorityItems = append(rogue.priorityItems, eviscerate)
			case proto.Rogue_Rotation_EviscerateRupture:
				eviscerate.MinimumComboPoints = core.MaxInt32(1, rogue.Rotation.MinimumComboPointsPrimaryFinisher)
				rogue.priorityItems = append(rogue.priorityItems, eviscerate)
				rupture.MinimumComboPoints = rogue.Rotation.MinimumComboPointsSecondaryFinisher
				rupture.IsFiller = true
				if rupture.MinimumComboPoints > 0 && rupture.MinimumComboPoints <= 5 {
					rogue.priorityItems = append(rogue.priorityItems, rupture)
				}
			}
		} else {
			switch rogue.Rotation.SubtletyFinisherPriority {
			case proto.Rogue_Rotation_SubtletyEviscerate:
				rupture.MinimumComboPoints = core.MaxInt32(1, rogue.Rotation.MinimumComboPointsPrimaryFinisher)
				rogue.priorityItems = append(rogue.priorityItems, rupture)
				eviscerate.MinimumComboPoints = core.MaxInt32(1, rogue.Rotation.MinimumComboPointsSecondaryFinisher)
				eviscerate.IsFiller = true
				rogue.priorityItems = append(rogue.priorityItems, eviscerate)
			case proto.Rogue_Rotation_SubtletyEnvenom:
				eviscerate.MinimumComboPoints = core.MaxInt32(1, rogue.Rotation.MinimumComboPointsPrimaryFinisher)
				rogue.priorityItems = append(rogue.priorityItems, eviscerate)
				rupture.MinimumComboPoints = rogue.Rotation.MinimumComboPointsSecondaryFinisher
				rupture.IsFiller = true
				if rupture.MinimumComboPoints > 0 && rupture.MinimumComboPoints <= 5 {
					rogue.priorityItems = append(rogue.priorityItems, rupture)
				}
			}
		}
	}
	rogue.rotationItems = rogue.planRotation(sim)
}

func (rogue *Rogue) rotation(sim *core.Simulation) {
	if len(rogue.rotationItems) < 1 {
		panic("Rotation is empty")
	}
	eps := rogue.getExpectedEnergyPerSecond()
	shouldCast := rogue.shouldCastNextRotationItem(sim, eps)
	item := rogue.rotationItems[0]
	prio := rogue.priorityItems[item.PrioIndex]

	switch shouldCast {
	case ShouldNotCast:
		rogue.rotationItems = rogue.rotationItems[1:]
		rogue.rotation(sim)
	case ShouldBuild:
		spell := rogue.Builder
		if spell == nil || spell.Cast(sim, rogue.CurrentTarget) {
			if rogue.GCD.IsReady(sim) {
				rogue.rotation(sim)
			}
		} else {
			panic("Unexpected builder cast failure")
		}
	case ShouldCast:
		spell := prio.GetSpell(rogue, rogue.ComboPoints())
		if spell == nil || spell.Cast(sim, rogue.CurrentTarget) {
			rogue.priorityItems[item.PrioIndex].CastCount += 1
			rogue.rotationItems = rogue.planRotation(sim)
			if rogue.GCD.IsReady(sim) {
				rogue.rotation(sim)
			}
		} else {
			panic("Unexpected cast failure")
		}
	case ShouldWait:
		desiredEnergy := 100.0
		if rogue.ComboPoints() == 5 {
			desiredEnergy = prio.EnergyCost
		} else {
			if rogue.CurrentEnergy() < prio.EnergyCost && rogue.ComboPoints() >= prio.MinimumComboPoints {
				desiredEnergy = prio.EnergyCost
			} else if rogue.ComboPoints() < 5 {
				desiredEnergy = rogue.Builder.DefaultCast.Cost
			}
		}
		cdAvailableTime := time.Second * 10
		if sim.CurrentTime > cdAvailableTime {
			cdAvailableTime = core.NeverExpires
		}
		nextExpiration := cdAvailableTime
		for _, otherItem := range rogue.rotationItems {
			if otherItem.ExpiresAt < nextExpiration {
				nextExpiration = otherItem.ExpiresAt
			}
		}
		neededEnergy := desiredEnergy - rogue.CurrentEnergy()
		energyAvailableTime := time.Second*time.Duration(neededEnergy/eps) + 1*time.Second
		energyAt := sim.CurrentTime + energyAvailableTime
		if energyAt < nextExpiration {
			rogue.WaitForEnergy(sim, desiredEnergy)
		} else if nextExpiration > sim.CurrentTime {
			rogue.WaitUntil(sim, nextExpiration)
		} else {
			rogue.DoNothing()
		}
	}
}
