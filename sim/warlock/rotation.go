package warlock

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const humanReactionTime = 150 * time.Millisecond

func (warlock *Warlock) setupCooldowns(sim *core.Simulation) {
	// TODO: also need to consider shared CDs that cause delays, like double on-use trinkets

	// check if waiting `waitDuration` amount of time results in fewer cooldown usages
	retainUses := func(timeLeft time.Duration, spellCD time.Duration, waitDuration time.Duration) bool {
		return math.Floor(float64(timeLeft)/float64(spellCD)+1) ==
			math.Floor(float64(timeLeft-waitDuration)/float64(spellCD)+1)
	}

	// TODO: find a way of getting the duration directly from the spell instead
	durMap := make(map[core.ActionID]time.Duration)
	if warlock.MetamorphosisAura != nil {
		durMap[core.ActionID{SpellID: 47241}] = warlock.MetamorphosisAura.Duration
	}
	durMap[core.ActionID{SpellID: 33697}] = 15 * time.Second
	durMap[core.ActionID{SpellID: 54758}] = 12 * time.Second
	durMap[core.ActionID{SpellID: 10060}] = 15 * time.Second
	durMap[core.ActionID{ItemID: 40211}] = 15 * time.Second
	durMap[core.ActionID{ItemID: 40212}] = 15 * time.Second
	durMap[core.ActionID{ItemID: 45466}] = 20 * time.Second
	durMap[core.ActionID{ItemID: 45148}] = 20 * time.Second
	durMap[core.ActionID{ItemID: 37873}] = 20 * time.Second

	ignoredCDs := make(map[core.ActionID]struct{})
	ignoredCDs[core.ActionID{ItemID: 42641}] = struct{}{}       // sapper
	ignoredCDs[core.ActionID{ItemID: 41119}] = struct{}{}       // saronite bomb
	ignoredCDs[core.ActionID{ItemID: 40536}] = struct{}{}       // explosive decoy
	ignoredCDs[core.BloodlustActionID.WithTag(-1)] = struct{}{} // don't mess with BL
	if warlock.Inferno != nil {
		ignoredCDs[warlock.Inferno.ActionID] = struct{}{}
	}

	var executeActive func() bool
	var executePhase time.Duration
	if warlock.Talents.Decimation > 0 {
		// approximation since we don't know exactly when decimation will be active
		executePhase = time.Duration(float64(sim.Duration)*(1.0-sim.Encounter.ExecuteProportion_35)) + 3*time.Second
		executeActive = func() bool { return warlock.DecimationAura.IsActive() }
	} else if warlock.Talents.Haunt {
		executePhase = time.Duration(float64(sim.Duration) * (1.0 - sim.Encounter.ExecuteProportion_25))
		executeActive = func() bool { return sim.IsExecutePhase25() }
	} else {
		executePhase = time.Duration(0)
		executeActive = func() bool { return true }
	}

	lustCD := warlock.GetMajorCooldownIgnoreTag(core.BloodlustActionID)
	for _, cd := range warlock.GetMajorCooldowns() {
		if _, ignored := ignoredCDs[cd.Spell.ActionID]; ignored {
			continue
		}

		spellCD := core.MaxDuration(cd.Spell.CD.Duration, cd.Spell.SharedCD.Duration)
		runTime := time.Duration(float64(durMap[cd.Spell.ActionID]) * 0.75)
		spell := cd.Spell

		cd.ShouldActivate = func(sim *core.Simulation, character *core.Character) bool {
			timeLeft := sim.GetRemainingDuration() - runTime
			timeUntilExecute := core.MaxDuration(0, executePhase-sim.CurrentTime)

			// if time until execute is less than the CD AND remaining time minus time till execute gives
			// the same amount of uses as remaining time alone then delay
			if !executeActive() && timeUntilExecute < spellCD+runTime &&
				retainUses(timeLeft, spellCD, timeUntilExecute) {
				return false
			}

			if warlock.Talents.Metamorphosis && spell.ActionID != warlock.Metamorphosis.ActionID {
				metaCD := warlock.GetMajorCooldown(warlock.Metamorphosis.ActionID)
				if !warlock.MetamorphosisAura.IsActive() && metaCD.TimeToNextCast(sim) < spellCD+runTime &&
					retainUses(timeLeft, spellCD, metaCD.TimeToNextCast(sim)) {
					return false
				}
			}

			if lustCD != nil && !character.HasActiveAuraWithTag(core.BloodlustAuraTag) &&
				lustCD.TimeToNextCast(sim) < spellCD+runTime && retainUses(timeLeft, spellCD,
				lustCD.TimeToNextCast(sim)) {
				return false
			}

			if spell.ActionID.SameActionIgnoreTag(core.PowerInfusionActionID) &&
				(character.HasActiveAuraWithTag(core.BloodlustAuraTag) || lustCD.TimeToNextCast(sim) < runTime) {
				return false // don't use PI while lust is active or it would overlap
			}

			return true
		}
	}
}

func (warlock *Warlock) calcRelativeCorruptionInc(target *core.Unit) float64 {
	dot := warlock.Corruption.Dot(target)
	snapshotCrit := dot.SnapshotCritChance
	snapshotDmg := dot.SnapshotAttackerMultiplier * (snapshotCrit*(warlock.Corruption.CritMultiplier-1) + 1)

	attackTable := warlock.AttackTables[target.UnitIndex]
	curCrit := warlock.Corruption.SpellCritChance(target)
	curDmg := dot.Spell.AttackerDamageMultiplier(attackTable) * (curCrit*(warlock.Corruption.CritMultiplier-1) + 1)

	return curDmg / snapshotDmg
}

func aclAppendSimple(acl []ActionCondition, spell *core.Spell, cond func(sim *core.Simulation) (
	bool, *core.Unit, string)) []ActionCondition {
	return append(acl, ActionCondition{
		Spell: spell,
		Condition: func(sim *core.Simulation) (ACLaction, *core.Unit, string) {
			if cond, target, reason := cond(sim); cond {
				return ACLCast, target, reason
			} else {
				return ACLNext, nil, reason
			}
		},
	})
}

func (warlock *Warlock) defineRotation() {
	acl := warlock.acl
	mainTarget := warlock.CurrentTarget // assumed to be the first element in the target list
	hauntTravel := time.Duration(float64(time.Second) * warlock.DistanceFromTarget / warlock.Haunt.MissileSpeed)
	critDebuffCat := warlock.GetEnemyExclusiveCategories(core.SpellCritEffectCategory).Get(mainTarget)

	allUnits := warlock.Env.Encounter.TargetUnits
	if mainTarget != allUnits[0] {
		panic("CurrentTarget assumption violated")
	}

	var multidotTargets, uaDotTargets []*core.Unit
	multidotCount := core.MinInt(len(allUnits), 3)
	if warlock.Rotation.Type == proto.Warlock_Rotation_Affliction {
		// up to 3 targets: multidot, no seed
		// 4 targets: corruption+UA 3x, seed on 4th; possibly only 1x UA since it's close in value
		// 5 targets: corruption x3, UA 1x, seed
		// 6 targets: corruption x2, UA 1x, seed; only 1x corruption + UA is close in value
		// 7-9 targets: corruption x1, no UA, seed
		// 10+ targets: no corruption anymore probably
		uaCount := core.MinInt(len(allUnits), 3)

		if len(allUnits) > 4 {
			uaCount = 1
		}
		if len(allUnits) == 6 {
			multidotCount = 2
		} else if len(allUnits) > 6 {
			uaCount = 0
			multidotCount = core.TernaryInt(len(allUnits) > 9, 0, 1)
		}

		uaDotTargets = allUnits[:uaCount]
	} else if warlock.Rotation.Type == proto.Warlock_Rotation_Destruction {
		multidotCount = core.MinInt(len(allUnits), 4)
	}
	multidotTargets = allUnits[:multidotCount]

	if warlock.Talents.DemonicEmpowerment && warlock.Options.Summon != proto.Warlock_Options_NoSummon {
		acl = aclAppendSimple(acl, warlock.DemonicEmpowerment, func(sim *core.Simulation) (bool, *core.Unit, string) {
			return !warlock.Rotation.UseInfernal || warlock.Inferno.IsReady(sim), mainTarget, ""
		})
	}

	if warlock.Talents.Metamorphosis {
		acl = aclAppendSimple(acl, warlock.ImmolationAura, func(sim *core.Simulation) (bool, *core.Unit, string) {
			// TODO: potentially wait for procs
			return true, nil, ""
		})
	}

	// only handles deliberate overrides of the primary spell
	if warlock.Rotation.PrimarySpell == proto.Warlock_Rotation_Seed {
		acl = aclAppendSimple(acl, warlock.Seed, func(sim *core.Simulation) (bool, *core.Unit, string) {
			return warlock.Rotation.DetonateSeed || !warlock.Seed.Dot(mainTarget).IsActive(), mainTarget, ""
		})
	}

	if warlock.Talents.Conflagrate {
		acl = aclAppendSimple(acl, warlock.Conflagrate, func(sim *core.Simulation) (bool, *core.Unit, string) {
			return warlock.Immolate.Dot(mainTarget).IsActive(), mainTarget, ""
		})
	}

	if warlock.Talents.Haunt && warlock.Rotation.SpecSpell == proto.Warlock_Rotation_Haunt {
		curIndex := len(acl)

		acl = aclAppendSimple(acl, warlock.Haunt, func(sim *core.Simulation) (bool, *core.Unit, string) {
			// no need for haunt until dots are up, mostly relevant in the opener
			if !warlock.Corruption.Dot(mainTarget).IsActive() && !warlock.UnstableAffliction.Dot(mainTarget).IsActive() {
				return false, nil, ""
			}

			if !warlock.Haunt.CD.IsReady(sim) {
				return false, nil, ""
			}

			if sim.GetRemainingDuration() < 5*time.Second {
				return false, nil, ""
			}

			castTime := warlock.Haunt.CastTime()
			_, nextActionTime := warlock.getAlternativeAction(sim, curIndex)
			hauntRem := warlock.HauntDebuffAuras.Get(mainTarget).RemainingDuration(sim)

			// 250ms of leeway in case haste buffs run out
			return hauntRem-castTime-hauntTravel < nextActionTime+250*time.Millisecond, mainTarget, ""
		})

		acl = aclAppendSimple(acl, warlock.LifeTap, func(sim *core.Simulation) (bool, *core.Unit, string) {
			val := warlock.ShadowBolt.DefaultCast.Cost

			if sim.IsExecutePhase25() {
				dsDot := warlock.DrainSoul.CurDot()
				if dsDot.IsActive() && dsDot.NumTicksRemaining(sim) >= 1 {
					return false, nil, "" // continuing to channel drain soul doesn't cost us any mana
				}

				val = warlock.UnstableAffliction.DefaultCast.Cost // highest mana cost spell outside SB
			}
			val += warlock.Haunt.DefaultCast.Cost

			if warlock.CurrentMana() > val || sim.GetRemainingDuration() > 5*time.Second {
				return false, nil, ""
			}

			return true, nil, "Casting life tap to not drop haunt"
		})
	}

	// refresh corruption with shadow bolt if it's running out
	if warlock.Talents.EverlastingAffliction == 5 && len(allUnits) > 1 {
		travel := time.Duration(float64(time.Second) * warlock.DistanceFromTarget / warlock.ShadowBolt.MissileSpeed)
		curIndex := len(acl)

		acl = aclAppendSimple(acl, warlock.ShadowBolt, func(sim *core.Simulation) (bool, *core.Unit, string) {
			type targetRem struct {
				target *core.Unit
				rem    time.Duration
			}
			targets := make([]targetRem, 0, len(sim.Encounter.TargetUnits))
			for _, target := range sim.Encounter.TargetUnits {
				// if there's already an shadowbolt on the way then skip
				if warlock.corrRefreshList[target.UnitIndex] >= sim.CurrentTime-travel {
					continue
				}

				// same when we can't refresh in time
				if warlock.Corruption.Dot(target).RemainingDuration(sim) < travel+warlock.ShadowBolt.CastTime() {
					continue
				}

				// assuming haunt doesn't drop, which it shouldn't, corruption will already be refreshed
				if target == mainTarget && warlock.HauntDebuffAuras.Get(target).RemainingDuration(sim) <
					warlock.Corruption.Dot(target).RemainingDuration(sim) {
					continue
				}

				rem := core.MinDuration(warlock.Corruption.Dot(target).RemainingDuration(sim),
					warlock.ShadowEmbraceAuras.Get(target).RemainingDuration(sim))
				targets = append(targets, targetRem{rem: rem, target: target})
			}
			sort.Slice(targets, func(i, j int) bool { return targets[i].rem < targets[j].rem })

			// we know that the only higher priority action is haunt, thus the only 2 things we need to
			// consider outside of shadow bolts is haunt and mana
			nextSpell, timeAdvance := warlock.getAlternativeAction(sim, curIndex)
			sbCastTime := warlock.ShadowBolt.EffectiveCastTime()
			timeAdvance += sbCastTime
			recast := false
			// shadow trance proc will only speed up one cast
			if warlock.ShadowBolt.CastTimeMultiplier == 0 {
				// somewhat hacky, breaks if CastTimeMultiplier is ever changed by anything else
				sbCastTime = time.Duration(float64(warlock.ShadowBolt.DefaultCast.CastTime) * warlock.CastSpeed)
				sbCastTime = core.MaxDuration(sbCastTime, warlock.SpellGCD())

				if nextSpell == warlock.ShadowBolt {
					timeAdvance += sbCastTime - warlock.ShadowBolt.EffectiveCastTime()
				}
			}
			mana := warlock.CurrentMana() - warlock.ShadowBolt.DefaultCast.Cost
			consideredHaunt := false
			for _, ele := range targets {
				if mana < warlock.ShadowBolt.DefaultCast.Cost {
					timeAdvance += warlock.LifeTap.EffectiveCastTime()
					mana += 10000.0 // we only need 1 life tap, so the exact value doesn't matter
				}
				mana -= warlock.ShadowBolt.DefaultCast.Cost

				if !consideredHaunt && timeAdvance+warlock.Haunt.CastTime()+hauntTravel >=
					warlock.HauntDebuffAuras.Get(mainTarget).RemainingDuration(sim) {
					timeAdvance += warlock.Haunt.EffectiveCastTime()
					mana -= warlock.Haunt.DefaultCast.Cost
					consideredHaunt = true
				}

				// some extra time to accommodate haste buffs running out
				if timeAdvance+travel+250*time.Millisecond >= ele.rem {
					recast = true
					break
				}

				timeAdvance += sbCastTime + 50*time.Millisecond
			}

			if recast {
				return true, targets[0].target, ""
			} else {
				return false, nil, ""
			}
		})
	}

	if warlock.Rotation.Corruption && warlock.Talents.EverlastingAffliction > 0 {
		acl = aclAppendSimple(acl, warlock.Corruption, func(sim *core.Simulation) (bool, *core.Unit, string) {
			// TODO: wait for all targets SB debuff?
			if !critDebuffCat.AnyActive() &&
				warlock.Talents.ImprovedShadowBolt > 0 && sim.CurrentTime < 25 {
				return false, nil, ""
			}

			reason := ""
			for _, target := range multidotTargets {
				if !warlock.Corruption.Dot(target).IsActive() {
					return true, target, ""
				}

				// check if reapplying corruption is worthwhile
				relDmgInc := warlock.calcRelativeCorruptionInc(target)
				snapshotDmg := warlock.Corruption.ExpectedDamageFromCurrentSnapshot(sim, target)
				snapshotDmg *= float64(sim.GetRemainingDuration()) / float64(warlock.Corruption.Dot(target).TickPeriod())
				snapshotDmg *= (relDmgInc - 1)
				snapshotDmg -= warlock.Corruption.ExpectedDamageFromCurrentSnapshot(sim, target)

				reason = fmt.Sprintf("Relative Corruption Inc: [%.2f], expected dmg gain: [%.2f]",
					relDmgInc, snapshotDmg)

				if relDmgInc > 1.15 || snapshotDmg > 10000 {
					return true, target, reason
				}
			}

			return false, nil, reason
		})
	}

	prefCurse := warlock.CurseOfAgony.CurDot().Aura
	switch warlock.Rotation.Curse {
	case proto.Warlock_Rotation_Elements:
		prefCurse = warlock.CurseOfElementsAuras.Get(mainTarget)
		acl = aclAppendSimple(acl, warlock.CurseOfElements, func(sim *core.Simulation) (bool, *core.Unit, string) {
			return warlock.CurseOfElementsAuras.Get(mainTarget).RemainingDuration(sim) < 3*time.Second, mainTarget, ""
		})
	case proto.Warlock_Rotation_Weakness:
		prefCurse = warlock.CurseOfWeaknessAuras.Get(mainTarget)
		acl = aclAppendSimple(acl, warlock.CurseOfWeakness, func(sim *core.Simulation) (bool, *core.Unit, string) {
			return warlock.CurseOfWeaknessAuras.Get(mainTarget).RemainingDuration(sim) < 3*time.Second, mainTarget, ""
		})
	case proto.Warlock_Rotation_Tongues:
		prefCurse = warlock.CurseOfTonguesAuras.Get(mainTarget)
		acl = aclAppendSimple(acl, warlock.CurseOfTongues, func(sim *core.Simulation) (bool, *core.Unit, string) {
			return warlock.CurseOfTonguesAuras.Get(mainTarget).RemainingDuration(sim) < 3*time.Second, mainTarget, ""
		})
	}

	if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfLifeTap) {
		acl = aclAppendSimple(acl, warlock.LifeTap, func(sim *core.Simulation) (bool, *core.Unit, string) {
			// try to keep up the buff for the entire execute phase if possible
			expiresAt := core.MaxDuration(0, warlock.GlyphOfLifeTapAura.RemainingDuration(sim))
			if sim.GetRemainingDuration() <= 40*time.Second &&
				expiresAt+10*time.Second < sim.GetRemainingDuration() &&
				warlock.CurrentManaPercent() < 0.35 {
				return true, nil, "Casting life tap to keep up GoLT (40s till EOF)"
			} else if sim.GetRemainingDuration() <= 55*time.Second {
				return false, nil, ""
			}

			if warlock.GlyphOfLifeTapAura.RemainingDuration(sim) > 1*time.Second ||
				sim.GetRemainingDuration() <= 10*time.Second {
				return false, nil, ""
			}

			return true, nil, "Casting life tap to keep up GoLT"
		})
	}

	if warlock.Talents.UnstableAffliction && warlock.Rotation.SecondaryDot == proto.Warlock_Rotation_UnstableAffliction {
		acl = aclAppendSimple(acl, warlock.UnstableAffliction, func(sim *core.Simulation) (bool, *core.Unit, string) {
			castTime := warlock.UnstableAffliction.CastTime()
			for _, target := range uaDotTargets {
				if warlock.UnstableAffliction.Dot(target).RemainingDuration(sim)-castTime <= 0 &&
					sim.GetRemainingDuration() >= 9*time.Second+castTime {
					return true, target, ""
				}
			}

			return false, nil, ""
		})
	}

	if len(allUnits) > len(multidotTargets) {
		acl = aclAppendSimple(acl, warlock.Seed, func(sim *core.Simulation) (bool, *core.Unit, string) {
			for _, target := range sim.Encounter.TargetUnits {
				// avoid mainTarget as we may want to corruption that later
				if !warlock.Corruption.Dot(target).IsActive() && target != mainTarget {
					return true, target, ""
				}
			}
			panic("No viable seed target found")
		})
	}

	// TODO: automatically determine based on haunt/SE?
	if warlock.Rotation.Curse == proto.Warlock_Rotation_Doom {
		acl = aclAppendSimple(acl, warlock.CurseOfDoom, func(sim *core.Simulation) (bool, *core.Unit, string) {
			return warlock.CurseOfDoom.Dot(mainTarget).RemainingDuration(sim) <= 0 &&
				sim.GetRemainingDuration() >= 60*time.Second, mainTarget, ""
		})
	}

	if warlock.Rotation.Corruption && warlock.Talents.EverlastingAffliction <= 0 {
		acl = aclAppendSimple(acl, warlock.Corruption, func(sim *core.Simulation) (bool, *core.Unit, string) {
			for _, target := range multidotTargets {
				dot := warlock.Corruption.Dot(target)
				if dot.IsActive() {
					continue
				}

				tickLen := dot.TickLength
				if dot.AffectedByCastSpeed {
					tickLen = warlock.ApplyCastSpeed(tickLen)
				}

				if sim.GetRemainingDuration() >= 4*tickLen {
					return true, target, ""
				}
			}
			return false, nil, ""
		})
	}

	if warlock.Rotation.Curse == proto.Warlock_Rotation_Agony || warlock.Rotation.Curse == proto.Warlock_Rotation_Doom {
		tickHeuristic := core.TernaryDuration(warlock.Talents.Haunt, 16*time.Second, 22*time.Second)

		acl = aclAppendSimple(acl, warlock.CurseOfAgony, func(sim *core.Simulation) (bool, *core.Unit, string) {
			for _, target := range multidotTargets {
				if !warlock.CurseOfDoom.Dot(target).IsActive() && !warlock.CurseOfAgony.
					Dot(target).IsActive() && sim.GetRemainingDuration() >= tickHeuristic {
					return true, target, ""
				}
			}

			return false, nil, ""
		})
	}

	if !warlock.Talents.UnstableAffliction && warlock.Rotation.SecondaryDot == proto.Warlock_Rotation_Immolate {
		tickHeuristic := core.TernaryDuration(warlock.Talents.Conflagrate, 6*time.Second, 12*time.Second)

		acl = aclAppendSimple(acl, warlock.Immolate, func(sim *core.Simulation) (bool, *core.Unit, string) {
			castTime := warlock.Immolate.CastTime()
			for _, target := range multidotTargets {
				if warlock.Immolate.Dot(target).RemainingDuration(sim)-castTime <= 0 &&
					sim.GetRemainingDuration() >= tickHeuristic+castTime {
					return true, target, ""
				}
			}
			return false, nil, ""
		})
	}

	if warlock.Talents.ChaosBolt {
		acl = aclAppendSimple(acl, warlock.ChaosBolt, func(sim *core.Simulation) (bool, *core.Unit, string) {
			return true, mainTarget, ""
		})
	}

	if warlock.Talents.Haunt {
		function := func(sim *core.Simulation) (ACLaction, *core.Unit, string) {
			dsDot := warlock.DrainSoul.CurDot()
			if !sim.IsExecutePhase25() {
				return ACLNext, nil, ""
			}

			if !dsDot.IsActive() || dsDot.TimeUntilNextTick(sim) < dsDot.TickPeriod()-humanReactionTime {
				return ACLCast, mainTarget, ""
			}

			if warlock.Corruption.CurDot().RemainingDuration(sim) < dsDot.TickPeriod() {
				return ACLRecast, mainTarget, "Recasting drain soul to not let corruption drop"
			}

			// check if recasting drain soul is worthwhile

			// check when UA, CoA and haunt have to be refreshed, respectively
			uaRefresh := warlock.UnstableAffliction.Dot(mainTarget).RemainingDuration(sim) -
				warlock.UnstableAffliction.CastTime()

			curseRefresh := core.MaxDuration(prefCurse.RemainingDuration(sim),
				warlock.CurseOfDoom.CurDot().RemainingDuration(sim)) - warlock.CurseOfAgony.CastTime()

			hauntRefresh := warlock.HauntDebuffAuras.Get(mainTarget).RemainingDuration(sim) -
				warlock.Haunt.CastTime() - hauntTravel

			timeUntilRefresh := core.MinDuration(uaRefresh, curseRefresh)

			// the amount of ticks we have left, assuming we continue channeling
			ticksLeft := int(timeUntilRefresh/dsDot.TickPeriod()) + 1
			ticksLeft = core.MinInt(ticksLeft, int(hauntRefresh/dsDot.TickPeriod()))
			ticksLeft = core.MinInt(ticksLeft, dsDot.NumTicksRemaining(sim))

			// amount of ticks we'd get assuming we recast drain soul
			recastTicks := int(timeUntilRefresh/warlock.ApplyCastSpeed(dsDot.TickLength)) + 1
			recastTicks = core.MinInt(recastTicks, int(hauntRefresh/warlock.ApplyCastSpeed(dsDot.TickLength)))
			recastTicks = core.MinInt(recastTicks, int(dsDot.NumberOfTicks))

			if ticksLeft <= 0 || recastTicks <= 0 {
				return ACLCast, mainTarget, ""
			}

			snapshotDmg := warlock.DrainSoul.ExpectedDamageFromCurrentSnapshot(sim, mainTarget) * float64(ticksLeft)
			recastDmg := warlock.DrainSoul.ExpectedDamage(sim, mainTarget) * float64(recastTicks)
			snapshotDPS := snapshotDmg / (float64(ticksLeft) * dsDot.TickPeriod().Seconds())
			recastDps := recastDmg / (float64(recastTicks)*warlock.ApplyCastSpeed(dsDot.TickLength).Seconds() +
				humanReactionTime.Seconds())

			if recastDps > snapshotDPS {
				return ACLRecast, mainTarget, fmt.Sprintf("Recasting drain soul, %.2f (%d) > %.2f (%d)",
					recastDps, recastTicks, snapshotDPS, ticksLeft)
			}

			// TODO: if number of ticks left < number of ticks until we need to recast dots/haunt
			// and some proc effect falls off before the next tick, check if recasting is a DPS gain

			return ACLCast, mainTarget, ""
		}

		acl = append(acl, ActionCondition{
			Spell:     warlock.DrainSoul,
			Condition: function,
		})
	}

	if warlock.Talents.Decimation > 0 {
		acl = aclAppendSimple(acl, warlock.SoulFire, func(sim *core.Simulation) (bool, *core.Unit, string) {
			return warlock.DecimationAura.IsActive(), mainTarget, ""
		})
	}

	if warlock.Talents.MoltenCore > 0 {
		acl = aclAppendSimple(acl, warlock.Incinerate, func(sim *core.Simulation) (bool, *core.Unit, string) {
			return warlock.MoltenCoreAura.IsActive(), mainTarget, ""
		})
	}

	if warlock.Rotation.PrimarySpell == proto.Warlock_Rotation_Incinerate {
		acl = aclAppendSimple(acl, warlock.Incinerate, func(sim *core.Simulation) (bool, *core.Unit, string) {
			return true, mainTarget, ""
		})
	}

	acl = aclAppendSimple(acl, warlock.ShadowBolt, func(sim *core.Simulation) (bool, *core.Unit, string) {
		return true, mainTarget, ""
	})

	if warlock.Talents.DarkPact {
		acl = aclAppendSimple(acl, warlock.DarkPact, func(sim *core.Simulation) (bool, *core.Unit, string) {
			// if pet has enough mana, prefer dark pact over life tap
			return warlock.Pet.CurrentMana() > warlock.GetStat(stats.SpellPower)+1200+131, nil, ""
		})
	}

	acl = aclAppendSimple(acl, warlock.LifeTap, func(sim *core.Simulation) (bool, *core.Unit, string) {
		return true, nil, ""
	})

	warlock.acl = acl
}

func aclNextAction(sim *core.Simulation, acl []ActionCondition, skipIndex int) (*core.Spell, bool) {
	for _, ac := range acl[skipIndex+1:] {
		if action, _, _ := ac.Condition(sim); action != ACLNext && ac.Spell.IsReady(sim) {
			return ac.Spell, action == ACLRecast
		}
	}

	panic("ACL list exhausted but no match found")
}

// Returns the spell and casttime of the alternative action we'd take, if we skip skipIndex
func (warlock *Warlock) getAlternativeAction(sim *core.Simulation, skipIndex int) (*core.Spell, time.Duration) {
	var nextSpellTime time.Duration
	nextSpell, recast := aclNextAction(sim, warlock.acl, skipIndex)

	if nextSpell == warlock.DrainSoul {
		if recast || !nextSpell.CurDot().IsActive() {
			nextSpellTime = warlock.ApplyCastSpeed(nextSpell.CurDot().TickLength) + humanReactionTime
		} else {
			nextSpellTime = nextSpell.CurDot().TickPeriod() + humanReactionTime
		}
	} else {
		nextSpellTime = nextSpell.EffectiveCastTime()
	}

	return nextSpell, core.MaxDuration(core.GCDMin, nextSpellTime)
}

func (warlock *Warlock) OnGCDReady(sim *core.Simulation) {
	if warlock.Options.Summon != proto.Warlock_Options_NoSummon && warlock.Talents.DemonicKnowledge > 0 {
		// TODO: investigate a better way of handling this like a "reverse inheritance" for pets.
		bonus := (warlock.Pet.GetStat(stats.Stamina) + warlock.Pet.GetStat(stats.Intellect)) * (0.04 * float64(warlock.Talents.DemonicKnowledge))
		if bonus != warlock.petStmBonusSP {
			warlock.AddStatDynamic(sim, stats.SpellPower, bonus-warlock.petStmBonusSP)
			warlock.petStmBonusSP = bonus
		}
	}

	if warlock.Talents.DemonicPact > 0 && sim.CurrentTime != 0 && warlock.Pet != nil {
		dpspCurrent := warlock.DemonicPactAura.ExclusiveEffects[0].Priority
		currentTimeJump := sim.CurrentTime.Seconds() - warlock.PreviousTime.Seconds()

		if currentTimeJump > 0 {
			warlock.DPSPAggregate += dpspCurrent * currentTimeJump
			warlock.Metrics.UpdateDpasp(dpspCurrent * currentTimeJump)
		}

		if sim.Log != nil {
			warlock.Log(sim, "[Info] Demonic Pact spell power bonus average [%.0f]", warlock.DPSPAggregate/sim.CurrentTime.Seconds())
		}

		warlock.PreviousTime = sim.CurrentTime
	}

	for _, ac := range warlock.acl {
		action, target, reason := ac.Condition(sim)
		if reason != "" && sim.Log != nil {
			warlock.Log(sim, "[Info] %s\n", reason)
		}
		if action == ACLNext || !ac.Spell.IsReady(sim) {
			continue
		}

		// TODO: find a more general way of dealing with channeling spells, but for now this is fine since drain
		// soul is the only one being used anyway
		if action == ACLRecast {
			if ac.Spell != warlock.DrainSoul {
				panic("Trying to recast unknown spell")
			}
			warlock.DrainSoul.Dot(target).Cancel(sim)
		}

		if warlock.DrainSoul.CurDot().IsActive() {
			if ac.Spell != warlock.DrainSoul && warlock.DrainSoul.CurDot().TickCount != 0 {
				warlock.DrainSoul.CurDot().Cancel(sim)
			} else {
				warlock.WaitUntil(sim, sim.CurrentTime+warlock.DrainSoul.CurDot().TimeUntilNextTick(sim)+humanReactionTime)
				return
			}
		}

		castTime := ac.Spell.CastTime()
		if success := ac.Spell.Cast(sim, target); success {
			// track shadowbolts "in the air" that haven't refreshed corruption yet
			if ac.Spell == warlock.ShadowBolt || ac.Spell == warlock.Haunt {
				warlock.corrRefreshList[target.UnitIndex] = sim.CurrentTime + castTime
			}

			if !warlock.GCD.IsReady(sim) {
				// after-GCD actions
				if ac.Spell == warlock.Corruption && warlock.ItemSwap.IsEnabled() && warlock.ItemSwap.IsSwapped() {
					warlock.ItemSwap.SwapItems(sim, []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand,
						proto.ItemSlot_ItemSlotOffHand, proto.ItemSlot_ItemSlotRanged}, true)
				}

				return
			}
		} else if warlock.CurrentMana() < ac.Spell.DefaultCast.Cost {
			// TODO: this will only cast life tap right now
			if success := warlock.acl[len(warlock.acl)-1].Spell.Cast(sim, nil); !success {
				panic("Failed to cast life tap / dark pact")
			}
			return
		}
	}

	panic("ACL list exhausted but no match found")
}
