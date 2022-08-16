package warlock

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) LifeTapOrDarkPact(sim *core.Simulation) {
	if warlock.CurrentManaPercent() == 1 {
		panic("Life Tap or Dark Pact while full mana")
	}
	if warlock.Talents.DarkPact && warlock.Pet.CurrentMana() > warlock.GetStat(stats.SpellPower)+1200+131 {
		warlock.DarkPact.Cast(sim, warlock.CurrentTarget)
	} else {
		warlock.LifeTap.Cast(sim, warlock.CurrentTarget)
	}
}

func (warlock *Warlock) OnGCDReady(sim *core.Simulation) {
	warlock.tryUseGCD(sim)
}

func (warlock *Warlock) tryUseGCD(sim *core.Simulation) {
	var spell *core.Spell
	var filler *core.Spell
	var target = warlock.CurrentTarget
	mainSpell := warlock.Rotation.PrimarySpell
	secondaryDot := warlock.Rotation.SecondaryDot
	specSpell := warlock.Rotation.SpecSpell
	preset := warlock.Rotation.Preset
	rotationType := warlock.Rotation.Type
	curse := warlock.Rotation.Curse

	// ------------------------------------------
	// AoE (Seed)
	// ------------------------------------------
	if mainSpell == proto.Warlock_Rotation_Seed {
		if warlock.Rotation.DetonateSeed {
			if success := warlock.Seeds[0].Cast(sim, target); !success {
				warlock.LifeTapOrDarkPact(sim)
			}
			return
		}

		// If we aren't "auto popping" just put seed on and shadowbolt it.
		if !warlock.SeedDots[0].IsActive() {
			if success := warlock.Seeds[0].Cast(sim, target); success {
				return
			} else {
				warlock.LifeTapOrDarkPact(sim)
				return
			}
		}

		// If target has seed, fire a shadowbolt at main target so we start some explosions
		mainSpell = proto.Warlock_Rotation_ShadowBolt
	}

	// ------------------------------------------
	// Big CDs
	// ------------------------------------------

	bigCDs := warlock.GetMajorCooldowns()
	nextBigCD := time.Duration(math.MaxInt64)
	for _, cd := range bigCDs {
		if cd == nil {
			continue // not on cooldown right now.
		}
		cdReadyAt := cd.Spell.CD.ReadyAt()
		if cd.Type.Matches(core.CooldownTypeDPS) && cdReadyAt < nextBigCD {
			nextBigCD = cdReadyAt
		}
	}
	allCDs := []time.Duration{
		0,
		0,
		0,
	}

	if rotationType == proto.Warlock_Rotation_Affliction {
		hauntcasttime := warlock.ApplyCastSpeed(time.Millisecond * 1500)
		allCDs = []time.Duration{
			core.MaxDuration(0, time.Duration(float64(warlock.HauntDebuffAura(warlock.CurrentTarget).RemainingDuration(sim)-hauntcasttime)-float64(warlock.DistanceFromTarget)/20*1000)),
			core.MaxDuration(0, warlock.UnstableAffDot.RemainingDuration(sim)-hauntcasttime),
			core.MaxDuration(0, warlock.CurseOfAgonyDot.RemainingDuration(sim)),
		}
		if sim.Log != nil {
			//warlock.Log(sim, "Haunt[%d]", allCDs[0].Seconds())
			//warlock.Log(sim, "UA[%d]", allCDs[1].Seconds())
			//warlock.Log(sim, "SE stacks[%d]", warlock.ShadowEmbraceDebuffAura(warlock.CurrentTarget).GetStacks())
			//warlock.Log(sim, "SE time[%d]", warlock.ShadowEmbraceDebuffAura(warlock.CurrentTarget).RemainingDuration(sim).Seconds())
			//warlock.Log(sim, "travel time[%d]", float64(warlock.DistanceFromTarget)/20)
			//warlock.Log(sim, "filler time[%d]", (warlock.ApplyCastSpeed(time.Duration(warlock.ShadowBolt.DefaultCast.CastTime)).Seconds() + warlock.DistanceFromTarget/20))
		}
		nextCD := core.NeverExpires
		for _, v := range allCDs {
			if v < nextCD {
				nextCD = v
			}
		}
		nextBigCD = nextCD
	}

	// ------------------------------------------
	// Small CDs
	// ------------------------------------------
	if warlock.Talents.DemonicEmpowerment && warlock.DemonicEmpowerment.CD.IsReady(sim) && warlock.Options.Summon != proto.Warlock_Options_NoSummon {
		warlock.DemonicEmpowerment.Cast(sim, target)
	}
	if warlock.Talents.Metamorphosis && warlock.MetamorphosisAura.IsActive() &&
		warlock.ImmolationAura.CD.IsReady(sim) {
		warlock.ImmolationAura.Cast(sim, target)
	}

	// ------------------------------------------
	// Keep Glyph of Life Tap buff up
	// ------------------------------------------
	if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfLifeTap) &&
		(!warlock.GlyphOfLifeTapAura.IsActive() || warlock.GlyphOfLifeTapAura.RemainingDuration(sim) < time.Second) {
		if sim.CurrentTime < time.Second {

			// Pre-Pull Cast Shadow Bolt
			warlock.SpendMana(sim, warlock.ShadowBolt.DefaultCast.Cost, warlock.ShadowBolt.ResourceMetrics)
			warlock.ShadowBolt.SkipCastAndApplyEffects(sim, warlock.CurrentTarget)

			// Pre-pull Life Tap
			warlock.GlyphOfLifeTapAura.Activate(sim)

		} else {
			if sim.GetRemainingDuration() > time.Second*30 && rotationType == proto.Warlock_Rotation_Affliction { // more dps to not waste gcd on life tap for buff during execute unless execute is > 30 seconds
				warlock.LifeTapOrDarkPact(sim)
				return
			} else if rotationType == proto.Warlock_Rotation_Demonology || rotationType == proto.Warlock_Rotation_Destruction { // Other specs may want to keep up LT buff
				warlock.LifeTapOrDarkPact(sim)
				return
			}
		}
	}

	// ------------------------------------------
	// Curses
	// ------------------------------------------

	castCurse := func(spellToCast *core.Spell, aura *core.Aura) bool {
		if !aura.IsActive() {
			spell = spellToCast
			return true
		}
		return false
	}

	switch curse {
	case proto.Warlock_Rotation_Elements:
		castCurse(warlock.CurseOfElements, warlock.CurseOfElementsAura)
	case proto.Warlock_Rotation_Weakness:
		castCurse(warlock.CurseOfWeakness, warlock.CurseOfWeaknessAura)
	case proto.Warlock_Rotation_Tongues:
		castCurse(warlock.CurseOfTongues, warlock.CurseOfTonguesAura)
	case proto.Warlock_Rotation_Doom:
		if warlock.CurseOfDoom.CD.IsReady(sim) && sim.GetRemainingDuration() > time.Minute {
			spell = warlock.CurseOfDoom
		} else if sim.GetRemainingDuration() > time.Second*24 && !warlock.CurseOfAgonyDot.IsActive() && !warlock.CurseOfDoomDot.IsActive() {
			spell = warlock.CurseOfAgony
		}
	case proto.Warlock_Rotation_Agony:
		if rotationType == proto.Warlock_Rotation_Affliction {
			if sim.GetRemainingDuration() > time.Second*24 && allCDs[2] == 0 && allCDs[0] > 0 && allCDs[1] > 0 && warlock.CorruptionDot.IsActive() {
				spell = warlock.CurseOfAgony
			}
		} else {
			if sim.GetRemainingDuration() > time.Second*24 && !warlock.CurseOfAgonyDot.IsActive() {
				spell = warlock.CurseOfAgony
			}
		}
	}

	if spell != nil {
		if !spell.Cast(sim, target) {
			warlock.LifeTapOrDarkPact(sim)
		}
		return
	}

	// ------------------------------------------
	// Preset Rotations
	// ------------------------------------------

	// ------------------------------------------
	// Foreplay with filler
	// ------------------------------------------

	switch mainSpell {
	case proto.Warlock_Rotation_ShadowBolt:
		filler = warlock.ShadowBolt
	case proto.Warlock_Rotation_Incinerate:
		filler = warlock.Incinerate
	default:
		panic("No primary spell set")
	}

	fillerCastTime := warlock.ApplyCastSpeed(filler.DefaultCast.CastTime)
	ManaSpendRate := warlock.ShadowBolt.BaseCost / float64(fillerCastTime.Seconds()) //this is just an estimated mana spent per second
	DesiredManaAtExecute := 0.02                                                     //estimate for desired mana needed to do affliction execute
	TotalManaAtExecute := warlock.MaxMana() * DesiredManaAtExecute
	timeUntilOom := time.Duration((warlock.CurrentMana()-TotalManaAtExecute) / ManaSpendRate) * time.Second
	timeUntilExecute25 := time.Duration((sim.GetRemainingDurationPercent()- 0.25) * float64(sim.Duration))

	// If SE remaining duration is less than a shadow bolt cast time + travel time (with a 1 second buffer) and the previous cast was not haunt or SB then cast shadow bolt so SE stacks are not lost
	KeepUpSEStacks := (warlock.PrevCastSECheck != warlock.Haunt && warlock.PrevCastSECheck != warlock.ShadowBolt && warlock.ShadowEmbraceDebuffAura(warlock.CurrentTarget).RemainingDuration(sim).Seconds() < warlock.ApplyCastSpeed(time.Duration(warlock.ShadowBolt.DefaultCast.CastTime)).Seconds()+warlock.DistanceFromTarget/20+1)
	// If SE remaining duration is less than a shadow bolt cast time + travel time (with a 3 second buffer to include 1 drain soul tick) and the previous cast was not haunt or SB then cast shadow bolt so SE stacks are not lost
	KeepUpSEStacksExecute := (warlock.PrevCastSECheck != warlock.Haunt && warlock.PrevCastSECheck != warlock.ShadowBolt && warlock.ShadowEmbraceDebuffAura(warlock.CurrentTarget).RemainingDuration(sim).Seconds() < warlock.ApplyCastSpeed(time.Duration(warlock.ShadowBolt.DefaultCast.CastTime)).Seconds()+warlock.DistanceFromTarget/20+3)
	
	// This part trcks all the damage multiplier that roll over with corruption
	CurrentShadowMult := warlock.PseudoStats.ShadowDamageDealtMultiplier // Tracks the current shadow damage multipler (essentially looking for DE)
	CurrentDmgMult := warlock.PseudoStats.DamageDealtMultiplier          // Tracks the current damage multipler (essentially looking for TotT)
	CurrentCritBonus := warlock.GetStat(stats.SpellCrit) + warlock.PseudoStats.BonusSpellCritRating + warlock.PseudoStats.BonusShadowCritRating + 
		warlock.CurrentTarget.PseudoStats.BonusSpellCritRatingTaken   // Tracks the current crit rating multipler (essentially looking for Shadow Mastery (ISB))
	CurrentCritMult := 1 + CurrentCritBonus / core.CritRatingPerCritChance / 100 * core.TernaryFloat64(warlock.Talents.Pandemic, 1, 0)
	CurrentCorruptionRolloverMult := CurrentDmgMult * CurrentShadowMult * CurrentCritMult

	if sim.Log != nil {
		warlock.Log(sim, "Current Corruption Rollover Multiplier [%d]", CurrentCorruptionRolloverMult)
	}

	if preset == proto.Warlock_Rotation_Automatic {
		// ------------------------------------------
		// Affliction Rotation
		// ------------------------------------------
		if rotationType == proto.Warlock_Rotation_Affliction {
			if (CurrentCorruptionRolloverMult > warlock.CorruptionRolloverMult) || // If the original corruption multipliers are lower than this current time, then reapply corruption (also need to make sure this is some % into the fight)
				(!warlock.CorruptionDot.IsActive() && (core.ShadowMasteryAura(warlock.CurrentTarget).IsActive() || warlock.Talents.ImprovedShadowBolt == 0) && (!warlock.Haunt.CD.IsReady(sim) || allCDs[0] > 0) && allCDs[1] > 0) {
				// Cast Corruption as soon as the 5% crit debuff is up
				// Cast Corruption again when you get the execute buff (Death's Embrace)
				spell = warlock.Corruption
			} else if warlock.CorruptionDot.IsActive() && warlock.CorruptionDot.RemainingDuration(sim) < core.GCDDefault {
				// Emergency Corruption refresh just in case
				spell = warlock.DrainSoul
			} else if warlock.Talents.Haunt && warlock.Haunt.CD.IsReady(sim) && allCDs[0] == 0 && sim.GetRemainingDuration() > warlock.HauntDebuffAura(warlock.CurrentTarget).Duration/2. {
				// Keep Haunt up
				spell = warlock.Haunt
			} else if warlock.Talents.UnstableAffliction && (!warlock.Haunt.CD.IsReady(sim) || allCDs[0] > 0) && allCDs[1] == 0 && sim.GetRemainingDuration() > warlock.UnstableAffDot.Duration/2. {
				// Keep UA up
				spell = warlock.UnstableAff
			} else if sim.GetRemainingDuration() > time.Second*24 && allCDs[2] == 0 && (!warlock.Haunt.CD.IsReady(sim) || allCDs[0] > 0) && allCDs[1] > 0 && warlock.CorruptionDot.IsActive() {
				// Keep Agony up
				spell = warlock.CurseOfAgony
			} else if KeepUpSEStacks && sim.GetRemainingDuration() > time.Second*10 ||
				(core.ShadowMasteryAura(warlock.CurrentTarget).RemainingDuration(sim) < warlock.ShadowBolt.CurCast.CastTime && sim.GetRemainingDuration() > core.ShadowMasteryAura(warlock.CurrentTarget).Duration/2.) {
				// Shadow Embrace & Shadow Mastery refresh
				spell = warlock.ShadowBolt
			} else if sim.IsExecutePhase25() && !KeepUpSEStacksExecute {
				// Drain Soul execute phase
				spell = warlock.channelCheck(sim, warlock.DrainSoulDot, 5)
			}

		} else if rotationType == proto.Warlock_Rotation_Demonology {

			// ------------------------------------------
			// Demonology Rotation
			// ------------------------------------------
			if !warlock.CorruptionDot.IsActive() && core.ShadowMasteryAura(warlock.CurrentTarget).IsActive() &&
				sim.GetRemainingDuration() > warlock.CorruptionDot.Duration {
				spell = warlock.Corruption
			} else if (!warlock.ImmolateDot.IsActive() || warlock.ImmolateDot.RemainingDuration(sim) < warlock.Immolate.CurCast.CastTime) &&
				sim.GetRemainingDuration() > warlock.ImmolateDot.Duration/2. {
				spell = warlock.Immolate
			} else if core.ShadowMasteryAura(warlock.CurrentTarget).RemainingDuration(sim) < warlock.ShadowBolt.CurCast.CastTime && sim.GetRemainingDuration() > core.ShadowMasteryAura(warlock.CurrentTarget).Duration/2. {
				// Shadow Mastery refresh
				spell = warlock.ShadowBolt
			} else if warlock.DecimationAura.IsActive() {
				// Demonology execute phase
				spell = warlock.SoulFire
			} else if warlock.MoltenCoreAura.IsActive() {
				// Corruption proc
				spell = warlock.Incinerate
			}
		} else if rotationType == proto.Warlock_Rotation_Destruction {

			// ------------------------------------------
			// Destruction Rotation
			// ------------------------------------------
			if warlock.Talents.Shadowburn && sim.GetRemainingDuration() < 2*time.Second && warlock.Shadowburn.CD.IsReady(sim) {
				// TODO: ^ maybe use a better heuristic then a static 2s for using our finishers
				spell = warlock.Shadowburn
			} else if warlock.CanConflagrate(sim) && (warlock.ImmolateDot.TickCount > warlock.ImmolateDot.NumberOfTicks-2 || warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfConflagrate)) {
				spell = warlock.Conflagrate
			} else if (!warlock.ImmolateDot.IsActive() || warlock.ImmolateDot.RemainingDuration(sim) < warlock.Immolate.CurCast.CastTime) &&
				sim.GetRemainingDuration() > warlock.ImmolateDot.Duration/2. {
				spell = warlock.Immolate
			} else if warlock.Talents.ChaosBolt && warlock.ChaosBolt.CD.IsReady(sim) {
				spell = warlock.ChaosBolt
			}
		}
	}

	// ------------------------------------------
	// Manual Rotation
	// ------------------------------------------

	// ------------------------------------------
	// Main spells
	// ------------------------------------------

	// We're kind of trying to fit all different spec rotations in one big priority based rotation in order to let people experiment

	if preset == proto.Warlock_Rotation_Manual {
		if warlock.Rotation.Corruption &&
			(!warlock.CorruptionDot.IsActive() && (core.ShadowMasteryAura(warlock.CurrentTarget).IsActive() || warlock.Talents.ImprovedShadowBolt == 0) ||
				sim.IsExecutePhase35() && time.Duration(warlock.CorruptionDot.TickCount)*warlock.CorruptionDot.TickLength > sim.CurrentTime) {
			// Cast Corruption as soon as the 5% crit debuff is up if you have the talent
			// Cast Corruption again when you get the execute buff (Death's Embrace)
			spell = warlock.Corruption
		} else if warlock.CanConflagrate(sim) && (warlock.ImmolateDot.TickCount > warlock.ImmolateDot.NumberOfTicks-2 || warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfConflagrate)) {
			// Cast Conflagrate at last Immo tick or on CD if you have the glyph
			spell = warlock.Conflagrate
		} else if secondaryDot == proto.Warlock_Rotation_Immolate && (!warlock.ImmolateDot.IsActive() || warlock.ImmolateDot.RemainingDuration(sim) < warlock.Immolate.CurCast.CastTime) &&
			sim.GetRemainingDuration() > warlock.ImmolateDot.Duration/2. {
			// Refresh Immolate when it is gonna fade but not if the fight is ending
			spell = warlock.Immolate
		} else if warlock.Talents.UnstableAffliction && secondaryDot == proto.Warlock_Rotation_UnstableAffliction &&
			(!warlock.UnstableAffDot.IsActive() || warlock.UnstableAffDot.RemainingDuration(sim) < warlock.UnstableAff.CurCast.CastTime) &&
			sim.GetRemainingDuration() > warlock.UnstableAffDot.Duration {
			// Refresh Unstable when it is gonna fade but not if the fight is ending
			spell = warlock.UnstableAff
		} else if warlock.Talents.Haunt && specSpell == proto.Warlock_Rotation_Haunt && warlock.Haunt.CD.IsReady(sim) && !warlock.HauntDebuffAura(warlock.CurrentTarget).IsActive() {
			// Refresh Haunt Debuff
			spell = warlock.Haunt
		} else if warlock.Talents.ShadowEmbrace > 0 && warlock.ShadowEmbraceDebuffAura(warlock.CurrentTarget).RemainingDuration(sim) < warlock.ShadowBolt.CurCast.CastTime+core.GCDDefault ||
			warlock.Talents.ImprovedShadowBolt > 0 && core.ShadowMasteryAura(warlock.CurrentTarget).RemainingDuration(sim) < warlock.ShadowBolt.CurCast.CastTime {
			// Shadow Embrace & Shadow Mastery refresh
			spell = warlock.ShadowBolt
		} else if warlock.DecimationAura.IsActive() {
			// Spam Soulfire if you have the Decimation buff (Demonology execute phase)
			spell = warlock.SoulFire
		} else if warlock.MoltenCoreAura.IsActive() {
			// Spam Incinerate if you have the Molten Core buff (procs off Corruption ticks)
			spell = warlock.Incinerate
		} else if warlock.Talents.ChaosBolt && specSpell == proto.Warlock_Rotation_ChaosBolt && warlock.ChaosBolt.CD.IsReady(sim) {
			spell = warlock.ChaosBolt
		} else if sim.IsExecutePhase25() && warlock.Talents.SoulSiphon > 0 {
			// Drain Soul execute phase for Affliction
			spell = warlock.channelCheck(sim, warlock.DrainSoulDot, 5)
		}
	}
	// ------------------------------------------
	// Regen check
	// ------------------------------------------
	// If big CD coming up and we don't have enough mana for it, lifetap
	// Also, never do a big regen in the last few seconds of the fight.
	if !warlock.DoingRegen && nextBigCD-sim.CurrentTime < time.Second*6 && sim.GetRemainingDuration() > time.Second*30 {
		if warlock.CurrentManaPercent() < 0.6 {
			warlock.DoingRegen = true
		}
	}

	if warlock.DoingRegen {
		if nextBigCD-sim.CurrentTime < time.Second*2 {
			// stop regen, start blasting
			warlock.DoingRegen = false
		} else {
			warlock.LifeTapOrDarkPact(sim)
			if warlock.CurrentManaPercent() > 0.6 {
				warlock.DoingRegen = false
			}
			return
		}
	}

	// ------------------------------------------
	// Filler spell
	// ------------------------------------------
	if spell == nil {
		if timeUntilOom < 5*time.Second && timeUntilExecute25 > time.Second {
			// If you were gonna cast a filler but are low mana, get mana instead in order not to be OOM when an important spell is coming up
			warlock.LifeTapOrDarkPact(sim)
			return
		} else {
			// Filler
			if nextBigCD-sim.CurrentTime > 0 && (nextBigCD-sim.CurrentTime)*15 < fillerCastTime {
				warlock.WaitUntil(sim, sim.CurrentTime+nextBigCD)
				return
			} else {
				spell = filler
			}
		}
	}

	// ------------------------------------------
	// Spell casting
	// ------------------------------------------

	if success := spell.Cast(sim, target); success {
		warlock.PrevCastSECheck = spell
		if spell == warlock.Corruption {
		// We are recording the current rollover power of corruption
			warlock.CorruptionRolloverMult = CurrentCorruptionRolloverMult
		}
		return
	}

	// Lifetap if nothing else
	if warlock.CurrentManaPercent() < 0.8 {
		warlock.LifeTapOrDarkPact(sim)
		return
	}

	// If we get here, something's wrong
}
