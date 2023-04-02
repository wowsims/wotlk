package warlock

import (
	"math"
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

func aclAppendSimple(acl []ActionCondition, spell *core.Spell, cond func(sim *core.Simulation) (bool, *core.Unit)) []ActionCondition {
	return append(acl, ActionCondition{
		Spell: spell,
		Condition: func(sim *core.Simulation) (ACLaction, *core.Unit) {
			if cond, target := cond(sim); cond {
				return ACLCast, target
			} else {
				return ACLNext, target
			}
		},
	})
}

func (warlock *Warlock) defineRotation() {
	acl := warlock.acl
	warlock.skipList = make(map[int]struct{})
	mainTarget := warlock.CurrentTarget
	hauntTravel := time.Duration(float64(time.Second) * warlock.DistanceFromTarget / warlock.Haunt.MissileSpeed)

	if warlock.Talents.DemonicEmpowerment && warlock.Options.Summon != proto.Warlock_Options_NoSummon {
		acl = aclAppendSimple(acl, warlock.DemonicEmpowerment, func(sim *core.Simulation) (bool, *core.Unit) {
			return !warlock.Rotation.UseInfernal || warlock.Inferno.IsReady(sim), mainTarget
		})
	}

	if warlock.Talents.Metamorphosis {
		acl = aclAppendSimple(acl, warlock.ImmolationAura, func(sim *core.Simulation) (bool, *core.Unit) {
			// TODO: potentially wait for procs
			return true, nil
		})
	}

	// TODO: the real AoE rotation is way more complicated than this, this is really just a stub
	if warlock.Rotation.PrimarySpell == proto.Warlock_Rotation_Seed {
		acl = aclAppendSimple(acl, warlock.Seed, func(sim *core.Simulation) (bool, *core.Unit) {
			return warlock.Rotation.DetonateSeed || !warlock.Seed.Dot(mainTarget).IsActive(), mainTarget
		})
	}

	if warlock.Talents.Conflagrate {
		acl = aclAppendSimple(acl, warlock.Conflagrate, func(sim *core.Simulation) (bool, *core.Unit) {
			return warlock.Immolate.Dot(mainTarget).IsActive(), mainTarget
		})
	}

	if warlock.Talents.Haunt && warlock.Rotation.SpecSpell == proto.Warlock_Rotation_Haunt {
		curIndex := len(acl)

		acl = aclAppendSimple(acl, warlock.Haunt, func(sim *core.Simulation) (bool, *core.Unit) {
			// no need for haunt until dots are up, mostly relevant in the opener
			if !warlock.Corruption.Dot(mainTarget).IsActive() && !warlock.UnstableAffliction.Dot(mainTarget).IsActive() {
				return false, nil
			}

			if !warlock.Haunt.CD.IsReady(sim) {
				return false, nil
			}

			if sim.GetRemainingDuration() < 5*time.Second {
				return false, nil
			}

			castTime := warlock.Haunt.CastTime()
			nextActionTime := warlock.getNextActionTime(sim, curIndex)
			hauntRem := warlock.HauntDebuffAuras.Get(mainTarget).RemainingDuration(sim)

			// 250ms of leeway in case haste buffs run out
			return hauntRem-castTime-hauntTravel < nextActionTime+250*time.Millisecond, mainTarget
		})

		acl = aclAppendSimple(acl, warlock.LifeTap, func(sim *core.Simulation) (bool, *core.Unit) {
			val := warlock.ShadowBolt.DefaultCast.Cost

			if sim.IsExecutePhase25() {
				dsDot := warlock.DrainSoul.CurDot()
				if dsDot.IsActive() && dsDot.NumTicksRemaining(sim) >= 1 {
					return false, nil // continuing to channel drain soul doesn't cost us any mana
				}

				val = warlock.UnstableAffliction.DefaultCast.Cost // highest mana cost spell outside SB
			}
			val += warlock.Haunt.DefaultCast.Cost

			if warlock.CurrentMana() > val || sim.GetRemainingDuration() > 5*time.Second {
				return false, nil
			}

			if sim.Log != nil && len(warlock.skipList) == 0 {
				warlock.Log(sim, "[Info] Casting life tap to not drop haunt")
			}

			return true, nil
		})
	}

	if warlock.Rotation.Corruption && warlock.Talents.EverlastingAffliction > 0 {
		acl = aclAppendSimple(acl, warlock.Corruption, func(sim *core.Simulation) (bool, *core.Unit) {
			if !warlock.CritDebuffCategory.AnyActive() &&
				warlock.Talents.ImprovedShadowBolt > 0 && sim.CurrentTime < 25 {
				return false, nil
			}

			if !warlock.Corruption.Dot(mainTarget).IsActive() {
				return true, mainTarget
			}

			// check if reapplying corruption is a worthwhile
			relDmgInc := warlock.calcRelativeCorruptionInc(mainTarget)
			snapshotDmg := warlock.Corruption.ExpectedDamageFromCurrentSnapshot(sim, mainTarget)
			snapshotDmg *= float64(sim.GetRemainingDuration()) / float64(warlock.Corruption.Dot(mainTarget).TickPeriod())
			snapshotDmg *= (relDmgInc - 1)

			if sim.Log != nil && len(warlock.skipList) == 0 {
				warlock.Log(sim, "[Info] Relative Corruption Inc: [%.2f], expected dmg gain: [%.2f]",
					relDmgInc, snapshotDmg)
			}

			if relDmgInc > 1.15 || snapshotDmg > 6000 {
				return true, mainTarget
			}

			return false, nil
		})
	}

	prefCurse := warlock.CurseOfAgony
	switch warlock.Rotation.Curse {
	case proto.Warlock_Rotation_Elements:
		prefCurse = warlock.CurseOfElements
		acl = aclAppendSimple(acl, warlock.CurseOfElements, func(sim *core.Simulation) (bool, *core.Unit) {
			return warlock.CurseOfElementsAuras.Get(mainTarget).RemainingDuration(sim) < 3*time.Second, mainTarget
		})
	case proto.Warlock_Rotation_Weakness:
		prefCurse = warlock.CurseOfWeakness
		acl = aclAppendSimple(acl, warlock.CurseOfWeakness, func(sim *core.Simulation) (bool, *core.Unit) {
			return warlock.CurseOfWeaknessAuras.Get(mainTarget).RemainingDuration(sim) < 3*time.Second, mainTarget
		})
	case proto.Warlock_Rotation_Tongues:
		prefCurse = warlock.CurseOfTongues
		acl = aclAppendSimple(acl, warlock.CurseOfTongues, func(sim *core.Simulation) (bool, *core.Unit) {
			return warlock.CurseOfTonguesAuras.Get(mainTarget).RemainingDuration(sim) < 3*time.Second, mainTarget
		})
	}

	if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfLifeTap) {
		acl = aclAppendSimple(acl, warlock.LifeTap, func(sim *core.Simulation) (bool, *core.Unit) {
			// try to keep up the buff for the entire execute phase if possible
			expiresAt := core.MaxDuration(0, warlock.GlyphOfLifeTapAura.RemainingDuration(sim))
			if sim.GetRemainingDuration() <= 40*time.Second &&
				expiresAt+10*time.Second < sim.GetRemainingDuration() &&
				warlock.CurrentManaPercent() < 0.35 {
				if sim.Log != nil && len(warlock.skipList) == 0 {
					warlock.Log(sim, "[Info] Casting life tap to keep up GoLT (40s till EOF)")
				}

				return true, nil
			} else if sim.GetRemainingDuration() <= 55*time.Second {
				return false, nil
			}

			if warlock.GlyphOfLifeTapAura.RemainingDuration(sim) > 1*time.Second ||
				sim.GetRemainingDuration() <= 10*time.Second {
				return false, nil
			}

			if sim.Log != nil && len(warlock.skipList) == 0 {
				warlock.Log(sim, "[Info] Casting life tap to keep up GoLT")
			}

			return true, nil
		})
	}

	if warlock.Talents.UnstableAffliction && warlock.Rotation.SecondaryDot == proto.Warlock_Rotation_UnstableAffliction {
		acl = aclAppendSimple(acl, warlock.UnstableAffliction, func(sim *core.Simulation) (bool, *core.Unit) {
			castTime := warlock.UnstableAffliction.CastTime()
			if warlock.UnstableAffliction.Dot(mainTarget).RemainingDuration(sim)-castTime <= 0 &&
				sim.GetRemainingDuration() >= 9*time.Second+castTime {
				return true, mainTarget
			}

			return false, nil
		})
	}

	// TODO: automatically determine based on haunt/SE?
	if warlock.Rotation.Curse == proto.Warlock_Rotation_Doom {
		acl = aclAppendSimple(acl, warlock.CurseOfDoom, func(sim *core.Simulation) (bool, *core.Unit) {
			return warlock.CurseOfDoom.Dot(mainTarget).RemainingDuration(sim) <= 0 &&
				sim.GetRemainingDuration() >= 60*time.Second, mainTarget
		})
	}

	if warlock.Rotation.Curse == proto.Warlock_Rotation_Agony || warlock.Rotation.Curse == proto.Warlock_Rotation_Doom {
		tickHeuristic := core.TernaryDuration(warlock.Talents.Haunt, 16*time.Second, 22*time.Second)

		acl = aclAppendSimple(acl, warlock.CurseOfAgony, func(sim *core.Simulation) (bool, *core.Unit) {
			if !warlock.CurseOfDoom.Dot(mainTarget).IsActive() && !warlock.CurseOfAgony.
				Dot(mainTarget).IsActive() && sim.GetRemainingDuration() >= tickHeuristic {
				return true, mainTarget
			}

			return false, nil
		})
	}

	if warlock.Rotation.Corruption && warlock.Talents.EverlastingAffliction <= 0 {
		acl = aclAppendSimple(acl, warlock.Corruption, func(sim *core.Simulation) (bool, *core.Unit) {
			dot := warlock.Corruption.Dot(mainTarget)
			if dot.IsActive() {
				return false, nil
			}

			tickLen := dot.TickLength
			if dot.AffectedByCastSpeed {
				tickLen = warlock.ApplyCastSpeed(tickLen)
			}

			return sim.GetRemainingDuration() >= 4*tickLen, mainTarget
		})
	}

	if !warlock.Talents.UnstableAffliction && warlock.Rotation.SecondaryDot == proto.Warlock_Rotation_Immolate {
		tickHeuristic := core.TernaryDuration(warlock.Talents.Conflagrate, 6*time.Second, 12*time.Second)

		acl = aclAppendSimple(acl, warlock.Immolate, func(sim *core.Simulation) (bool, *core.Unit) {
			castTime := warlock.Immolate.CastTime()
			return warlock.Immolate.Dot(mainTarget).RemainingDuration(sim)-castTime <= 0 &&
				sim.GetRemainingDuration() >= tickHeuristic+castTime, mainTarget
		})
	}

	if warlock.Talents.ChaosBolt {
		acl = aclAppendSimple(acl, warlock.ChaosBolt, func(sim *core.Simulation) (bool, *core.Unit) {
			return true, mainTarget
		})
	}

	if warlock.Talents.Haunt {
		function := func(sim *core.Simulation) (ACLaction, *core.Unit) {
			dsDot := warlock.DrainSoul.CurDot()
			if !sim.IsExecutePhase25() {
				return ACLNext, nil
			}

			if !dsDot.IsActive() || dsDot.TimeUntilNextTick(sim) < dsDot.TickPeriod()-humanReactionTime {
				return ACLCast, mainTarget
			}

			if warlock.Corruption.CurDot().RemainingDuration(sim) < dsDot.TickPeriod() {
				if sim.Log != nil && len(warlock.skipList) == 0 {
					warlock.Log(sim, "[Info] Recasting drain soul to not let corruption drop")
				}
				return ACLRecast, mainTarget // recast to not let corruption drop
			}

			// check if recasting drain soul is worthwhile

			// check when UA, CoA and haunt have to be refreshed, respectively
			uaRefresh := warlock.UnstableAffliction.Dot(mainTarget).RemainingDuration(sim) -
				warlock.UnstableAffliction.CastTime()

			curseRefresh := core.MaxDuration(prefCurse.CurDot().RemainingDuration(sim),
				warlock.CurseOfDoom.CurDot().RemainingDuration(sim)) - prefCurse.CastTime()

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
				return ACLCast, mainTarget
			}

			snapshotDmg := warlock.DrainSoul.ExpectedDamageFromCurrentSnapshot(sim, mainTarget) * float64(ticksLeft)
			recastDmg := warlock.DrainSoul.ExpectedDamage(sim, mainTarget) * float64(recastTicks)
			snapshotDPS := snapshotDmg / (float64(ticksLeft) * dsDot.TickPeriod().Seconds())
			recastDps := recastDmg / (float64(recastTicks)*warlock.ApplyCastSpeed(dsDot.TickLength).Seconds() +
				humanReactionTime.Seconds())

			if recastDps > snapshotDPS {
				if sim.Log != nil && len(warlock.skipList) == 0 {
					warlock.Log(sim, "[Info] Recasting drain soul, %.2f (%d) > %.2f (%d)\n",
						recastDps, recastTicks, snapshotDPS, ticksLeft)
				}

				return ACLRecast, mainTarget
			}

			// TODO: if number of ticks left < number of ticks until we need to recast dots/haunt
			// and some proc effect falls off before the next tick, check if recasting is a DPS gain

			return ACLCast, mainTarget
		}

		acl = append(acl, ActionCondition{
			Spell:     warlock.DrainSoul,
			Condition: function,
		})
	}

	if warlock.Talents.Decimation > 0 {
		acl = aclAppendSimple(acl, warlock.SoulFire, func(sim *core.Simulation) (bool, *core.Unit) {
			return warlock.DecimationAura.IsActive(), mainTarget
		})
	}

	if warlock.Talents.MoltenCore > 0 {
		acl = aclAppendSimple(acl, warlock.Incinerate, func(sim *core.Simulation) (bool, *core.Unit) {
			return warlock.MoltenCoreAura.IsActive(), mainTarget
		})
	}

	if warlock.Talents.Emberstorm > 0 {
		acl = aclAppendSimple(acl, warlock.Incinerate, func(sim *core.Simulation) (bool, *core.Unit) {
			return true, mainTarget
		})
	}

	acl = aclAppendSimple(acl, warlock.ShadowBolt, func(sim *core.Simulation) (bool, *core.Unit) {
		return true, mainTarget
	})

	if warlock.Talents.DarkPact {
		acl = aclAppendSimple(acl, warlock.DarkPact, func(sim *core.Simulation) (bool, *core.Unit) {
			// if pet has enough mana, prefer dark pact over life tap
			return warlock.Pet.CurrentMana() > warlock.GetStat(stats.SpellPower)+1200+131, nil
		})
	}

	acl = aclAppendSimple(acl, warlock.LifeTap, func(sim *core.Simulation) (bool, *core.Unit) { return true, nil })

	warlock.acl = acl
}

func aclNextAction(sim *core.Simulation, acl []ActionCondition, skipList map[int]struct{}, skipIndex int) (*core.Spell, bool) {
	skipList[skipIndex] = struct{}{}
	for i, ac := range acl {
		if _, contains := skipList[i]; contains {
			continue
		}

		if action, _ := ac.Condition(sim); action != ACLNext && ac.Spell.IsReady(sim) {
			delete(skipList, skipIndex)
			return ac.Spell, action == ACLRecast
		}
	}

	panic("ACL list exhausted but no match found")
}

// time the next action will take, until we are ready to cast something else again
func (warlock *Warlock) getNextActionTime(sim *core.Simulation, skipIndex int) time.Duration {
	var nextSpellTime time.Duration
	nextSpell, recast := aclNextAction(sim, warlock.acl, warlock.skipList, skipIndex)

	if nextSpell == warlock.DrainSoul {
		if recast || !nextSpell.CurDot().IsActive() {
			nextSpellTime = warlock.ApplyCastSpeed(nextSpell.CurDot().TickLength) + humanReactionTime
		} else {
			nextSpellTime = nextSpell.CurDot().TickPeriod() + humanReactionTime
		}
	} else {
		nextSpellTime = nextSpell.EffectiveCastTime()
	}

	return core.MaxDuration(core.GCDMin, nextSpellTime)
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
		action, target := ac.Condition(sim)
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

		if success := ac.Spell.Cast(sim, target); success {
			if !warlock.GCD.IsReady(sim) {
				return
			}

			// TODO: if the reason we failed to cast something is that we have not enough mana, we may want
			// to just tap. On the other hand maybe falling through and casting the next best thing
			// sometimes has value?
		}
	}

	panic("ACL list exhausted but no match found")
}
