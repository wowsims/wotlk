package warlock

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) defineRotation() {
	rotationType := warlock.Rotation.Type
	curse := warlock.Rotation.Curse
	secondaryDot := warlock.Rotation.SecondaryDot
	specSpell := warlock.Rotation.SpecSpell

	// Spells (does not include fillers)
 	spellBook := [...]*core.Spell{
 		warlock.Corruption,
 		warlock.Immolate,
 		warlock.UnstableAffliction,
 		warlock.Haunt,
 		warlock.CurseOfAgony,
 		warlock.CurseOfDoom,
 		warlock.Conflagrate,
 		warlock.ChaosBolt,
 	}

 	warlock.SpellsRotation = make([]SpellRotation, len(spellBook))
	for i, spell := range spellBook{
		warlock.SpellsRotation[i].Spell = spell
	}

	// Next Cast readyness in time unit (use same order as in spellBook)
	// 0 : Cast ready ; core.NeverExpires : never cast
	warlock.SpellsRotation[0].CastIn = func(sim *core.Simulation) time.Duration {
		if !warlock.Rotation.Corruption {
			return core.NeverExpires
		}
		// This part tracks all the damage multiplier that roll over with corruption
		CurrentShadowMult := warlock.PseudoStats.ShadowDamageDealtMultiplier // Tracks the current shadow damage multipler (essentially looking for DE)
		CurrentDmgMult := warlock.PseudoStats.DamageDealtMultiplier          // Tracks the current damage multipler (essentially looking for TotT)
		CurrentCritBonus := warlock.GetStat(stats.SpellCrit) + warlock.PseudoStats.BonusSpellCritRating + warlock.PseudoStats.BonusShadowCritRating +
			warlock.CurrentTarget.PseudoStats.BonusSpellCritRatingTaken 	 // Tracks the current crit rating multipler (essentially looking for Shadow Mastery (ISB))
		CurrentCritMult := 1 + CurrentCritBonus/core.CritRatingPerCritChance/100*core.TernaryFloat64(warlock.Talents.Pandemic, 1, 0)
		CurrentCorruptionRolloverMult := CurrentDmgMult * CurrentShadowMult * CurrentCritMult
		if warlock.Talents.EverlastingAffliction > 0 && ((CurrentCorruptionRolloverMult > warlock.CorruptionRolloverMult) ||
			// If the original corruption multipliers are lower than this current time, then reapply corruption (also need to make sure this is some % into the fight)
			(!warlock.CorruptionDot.IsActive() && (core.ShadowMasteryAura(warlock.CurrentTarget).IsActive() || warlock.Talents.ImprovedShadowBolt == 0))) {
			// Wait for SM to be applied to cast first Corruption
			return 0
		} else {
			return core.MaxDuration(0,warlock.CorruptionDot.RemainingDuration(sim))
		}
	}
	warlock.SpellsRotation[1].CastIn = func(sim *core.Simulation) time.Duration {
		if !(secondaryDot == proto.Warlock_Rotation_Immolate) || sim.GetRemainingDuration() < warlock.ImmolateDot.Duration/2. {
			return core.NeverExpires
		}
		return core.MaxDuration(0,warlock.ImmolateDot.RemainingDuration(sim)-warlock.ApplyCastSpeed(warlock.Immolate.DefaultCast.CastTime))
	}
	warlock.SpellsRotation[2].CastIn = func(sim *core.Simulation) time.Duration {
		if !warlock.Talents.UnstableAffliction || !(secondaryDot == proto.Warlock_Rotation_UnstableAffliction) {
			return core.NeverExpires
		}
		return core.MaxDuration(0,warlock.UnstableAfflictionDot.RemainingDuration(sim)-warlock.ApplyCastSpeed(warlock.UnstableAffliction.DefaultCast.CastTime))
	}
	warlock.SpellsRotation[3].CastIn = func(sim *core.Simulation) time.Duration {
		if !warlock.Talents.Haunt || !(specSpell == proto.Warlock_Rotation_Haunt) {
			return core.NeverExpires
		}
		hauntTravelTime := time.Duration(warlock.DistanceFromTarget/20) * time.Second
		hauntCastTime := warlock.ApplyCastSpeed(warlock.Haunt.DefaultCast.CastTime)
		spellCastTime := warlock.ApplyCastSpeed(core.GCDDefault)
		if sim.IsExecutePhase25() {
			spellCastTime = warlock.ApplyCastSpeed(warlock.DrainSoulDot.TickLength)
		}
		// If SE remaining duration is less than a shadow bolt cast time + travel time (with a 1 second buffer) and the previous cast was not haunt or SB then cast shadow bolt so SE stacks are not lost
		KeepUpSEStacks := (warlock.PrevCastSECheck != warlock.Haunt && warlock.PrevCastSECheck != warlock.ShadowBolt && warlock.ShadowEmbraceDebuffAura(warlock.CurrentTarget).RemainingDuration(sim) < hauntCastTime+hauntTravelTime+spellCastTime)
		if KeepUpSEStacks && sim.GetRemainingDuration() > time.Second*10 && warlock.Haunt.IsReady(sim) {
			return 0
		} else {
			return core.MaxDuration(0,warlock.HauntDebuffAura(warlock.CurrentTarget).RemainingDuration(sim)-hauntCastTime-hauntTravelTime)
		}
	}
	warlock.SpellsRotation[4].CastIn = func(sim *core.Simulation) time.Duration {
		if !(curse == proto.Warlock_Rotation_Doom || curse == proto.Warlock_Rotation_Agony) || warlock.CurseOfDoomDot.IsActive() || sim.GetRemainingDuration() < warlock.CurseOfAgonyDot.Duration/2 {
			return core.NeverExpires
		}
		return core.MaxDuration(0,warlock.CurseOfAgonyDot.RemainingDuration(sim))
	}
	warlock.SpellsRotation[5].CastIn = func(sim *core.Simulation) time.Duration {
		if curse != proto.Warlock_Rotation_Doom || !warlock.CurseOfDoom.IsReady(sim) || sim.GetRemainingDuration() < time.Minute {
			return core.NeverExpires
		}
		return core.MaxDuration(0,warlock.CurseOfDoomDot.RemainingDuration(sim))
	}
	warlock.SpellsRotation[6].CastIn = func(sim *core.Simulation) time.Duration {
		if !warlock.Talents.Conflagrate {
			return core.NeverExpires
		}
		if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfConflagrate) {
			return core.MaxDuration(0,warlock.SpellsRotation[6].Spell.TimeToReady(sim))
		} else {
			return core.MaxDuration(0,warlock.ImmolateDot.RemainingDuration(sim)-warlock.ImmolateDot.TickLength)
		}
	}
	warlock.SpellsRotation[7].CastIn = func(sim *core.Simulation) time.Duration {
		if !warlock.Talents.ChaosBolt || !(specSpell == proto.Warlock_Rotation_ChaosBolt)  {
			return core.NeverExpires
		}
		return core.MaxDuration(0,warlock.ChaosBolt.TimeToReady(sim))
	}

	// Priority based rotations (0 means not in rotation, 1 is max
	if rotationType == proto.Warlock_Rotation_Affliction {
		warlock.SpellsRotation[0].Priority = 1
		warlock.SpellsRotation[2].Priority = 2
		warlock.SpellsRotation[3].Priority = 3
		warlock.SpellsRotation[4].Priority = 4
	} else if rotationType == proto.Warlock_Rotation_Demonology {
		warlock.SpellsRotation[0].Priority = 2
		warlock.SpellsRotation[1].Priority = 3
		warlock.SpellsRotation[4].Priority = 4
		warlock.SpellsRotation[5].Priority = 1
	} else if rotationType == proto.Warlock_Rotation_Destruction {
		warlock.SpellsRotation[1].Priority = 2
		warlock.SpellsRotation[4].Priority = 5
		warlock.SpellsRotation[5].Priority = 1
		warlock.SpellsRotation[6].Priority = 3
		warlock.SpellsRotation[7].Priority = 4
	}
}

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
	preset := warlock.Rotation.Preset
	rotationType := warlock.Rotation.Type
	curse := warlock.Rotation.Curse
	dotLag := time.Duration(10*time.Millisecond)

	// ------------------------------------------
	// Data
	// ------------------------------------------
	if warlock.Talents.DemonicPact > 0 && sim.CurrentTime != 0 {
		// We are integrating the Demonic Pact SP bonus over the course of the simulation to get the average
		warlock.DPSPAverage *= float64(warlock.PreviousTime)
		warlock.DPSPAverage += core.DemonicPactAura(warlock.GetCharacter(), 0).Priority * float64(sim.CurrentTime-warlock.PreviousTime)
		warlock.DPSPAverage /= float64(sim.CurrentTime)
		warlock.PreviousTime = sim.CurrentTime
	}

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

	if nextBigCD - sim.CurrentTime <= 0 {
		// stop regen, start blasting
		warlock.DoingRegen = false
	}

	allCDs := []time.Duration{
		0,
		0,
		0,
	}

	nextCD := core.NeverExpires
	immolateCastTime := warlock.ApplyCastSpeed(warlock.Immolate.DefaultCast.CastTime)
	hauntTravelTime := time.Duration(float64(warlock.DistanceFromTarget)/20) * time.Second
	if rotationType == proto.Warlock_Rotation_Affliction {
		hauntCastTime := warlock.ApplyCastSpeed(warlock.Haunt.DefaultCast.CastTime)
		UACastTime := warlock.ApplyCastSpeed(warlock.UnstableAffliction.DefaultCast.CastTime)
		allCDs = []time.Duration{
			core.MaxDuration(0, warlock.HauntDebuffAura(warlock.CurrentTarget).RemainingDuration(sim)-hauntCastTime-hauntTravelTime),
			core.MaxDuration(0, warlock.UnstableAfflictionDot.RemainingDuration(sim)-UACastTime),
			core.MaxDuration(0, warlock.CurseOfAgonyDot.RemainingDuration(sim)),
		}
	} else if rotationType == proto.Warlock_Rotation_Demonology {
		allCDs = []time.Duration{
			core.MaxDuration(0, warlock.CorruptionDot.RemainingDuration(sim)),
			core.MaxDuration(0, warlock.ImmolateDot.RemainingDuration(sim)-immolateCastTime),
			core.MaxDuration(0, warlock.CurseOfDoomDot.RemainingDuration(sim)),
		}
	} else if rotationType == proto.Warlock_Rotation_Destruction {
		allCDs = []time.Duration{
			core.MaxDuration(0, warlock.ImmolateDot.RemainingDuration(sim)-immolateCastTime),
			core.MaxDuration(0, warlock.CurseOfDoomDot.RemainingDuration(sim)),
			core.NeverExpires,
		}
	}
	for _, v := range allCDs {
		if v < nextCD {
			nextCD = v
		}
	}
	nextCD += sim.CurrentTime

	// ------------------------------------------
	// Small CDs
	// ------------------------------------------
	if warlock.Talents.DemonicEmpowerment && warlock.DemonicEmpowerment.IsReady(sim) && warlock.Options.Summon != proto.Warlock_Options_NoSummon {
		warlock.DemonicEmpowerment.Cast(sim, target)
	}
	if warlock.Talents.Metamorphosis && warlock.MetamorphosisAura.IsActive() &&
		warlock.ImmolationAura.IsReady(sim) {
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
			if sim.GetRemainingDuration() > time.Second*30 {
 				// More dps to not waste gcd on life tap for buff during execute unless execute is > 30 seconds
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
		if preset == proto.Warlock_Rotation_Automatic {
			if warlock.CurseOfDoom.CD.IsReady(sim) && sim.GetRemainingDuration() > time.Minute {
				spell = warlock.CurseOfDoom
			} else if sim.GetRemainingDuration() > warlock.CurseOfAgonyDot.Duration/2 && !warlock.CurseOfAgonyDot.IsActive() && !warlock.CurseOfDoomDot.IsActive() {
				spell = warlock.CurseOfAgony
			}
		}
	case proto.Warlock_Rotation_Agony:
		if preset == proto.Warlock_Rotation_Automatic {
			if rotationType == proto.Warlock_Rotation_Affliction {
				if sim.GetRemainingDuration() > warlock.CurseOfAgonyDot.Duration/2 && allCDs[2] == 0 && (!warlock.Haunt.CD.IsReady(sim) || allCDs[0] > 0) && allCDs[1] > 0 && warlock.CorruptionDot.IsActive() {
					spell = warlock.CurseOfAgony
				}
			} else {
				if sim.GetRemainingDuration() > warlock.CurseOfAgonyDot.Duration/2 && !warlock.CurseOfAgonyDot.IsActive() {
					spell = warlock.CurseOfAgony
				}
			}
		}
	}

	if spell != nil {
		if preset == proto.Warlock_Rotation_Automatic || (curse != proto.Warlock_Rotation_Doom && curse != proto.Warlock_Rotation_Agony) {
			if !spell.Cast(sim, target) {
				warlock.LifeTapOrDarkPact(sim)
			}
			return
		}
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
	timeUntilOom := time.Duration((warlock.CurrentMana()-TotalManaAtExecute)/ManaSpendRate) * time.Second
	timeUntilExecute25 := time.Duration((sim.GetRemainingDurationPercent() - 0.25) * float64(sim.Duration))

	// If SE remaining duration is less than a shadow bolt cast time + travel time (with a 1 second buffer) and the previous cast was not haunt or SB then cast shadow bolt so SE stacks are not lost
	KeepUpSEStacks := (warlock.PrevCastSECheck != warlock.Haunt && warlock.PrevCastSECheck != warlock.ShadowBolt && warlock.ShadowEmbraceDebuffAura(warlock.CurrentTarget).RemainingDuration(sim).Seconds() < warlock.ApplyCastSpeed(time.Duration(warlock.ShadowBolt.DefaultCast.CastTime)).Seconds()+warlock.DistanceFromTarget/20+1)
	// If SE remaining duration is less than a shadow bolt cast time + travel time (with a 3 second buffer to include 1 drain soul tick) and the previous cast was not haunt or SB then cast shadow bolt so SE stacks are not lost
	KeepUpSEStacksExecute := (warlock.PrevCastSECheck != warlock.Haunt && warlock.PrevCastSECheck != warlock.ShadowBolt && warlock.ShadowEmbraceDebuffAura(warlock.CurrentTarget).RemainingDuration(sim).Seconds() < warlock.ApplyCastSpeed(time.Duration(warlock.ShadowBolt.DefaultCast.CastTime)).Seconds()+warlock.DistanceFromTarget/20+3)

	// This part tracks all the damage multiplier that roll over with corruption
	CurrentShadowMult := warlock.PseudoStats.ShadowDamageDealtMultiplier // Tracks the current shadow damage multipler (essentially looking for DE)
	CurrentDmgMult := warlock.PseudoStats.DamageDealtMultiplier          // Tracks the current damage multipler (essentially looking for TotT)
	CurrentCritBonus := warlock.GetStat(stats.SpellCrit) + warlock.PseudoStats.BonusSpellCritRating + warlock.PseudoStats.BonusShadowCritRating +
		warlock.CurrentTarget.PseudoStats.BonusSpellCritRatingTaken // Tracks the current crit rating multipler (essentially looking for Shadow Mastery (ISB))
	CurrentCritMult := 1 + CurrentCritBonus/core.CritRatingPerCritChance/100*core.TernaryFloat64(warlock.Talents.Pandemic, 1, 0)
	CurrentCorruptionRolloverMult := CurrentDmgMult * CurrentShadowMult * CurrentCritMult

	if sim.Log != nil {
		if warlock.Talents.EverlastingAffliction > 0 {
			warlock.Log(sim, "[Info] Initial Corruption Rollover Multiplier [%.2f]", warlock.CorruptionRolloverMult)
			warlock.Log(sim, "[Info] Current Corruption Rollover Multiplier [%.2f]", CurrentCorruptionRolloverMult)
		}
		if warlock.Talents.DemonicPact > 0 {
			warlock.Log(sim, "[Info] Demonic Pact Spell Power Average [%.0f]", warlock.DPSPAverage)
		}
	}

	SBTravelTime := time.Duration(float64(warlock.DistanceFromTarget)/20) * time.Second

	if preset == proto.Warlock_Rotation_Automatic {
		// ------------------------------------------
		// Affliction Rotation
		// ------------------------------------------
		if rotationType == proto.Warlock_Rotation_Affliction {
			if (CurrentCorruptionRolloverMult > warlock.CorruptionRolloverMult) && warlock.Talents.EverlastingAffliction > 0 ||
				// If the original corruption multipliers are lower than this current time, then reapply corruption (also need to make sure this is some % into the fight)
				(!warlock.CorruptionDot.IsActive() && (core.ShadowMasteryAura(warlock.CurrentTarget).IsActive() || warlock.Talents.ImprovedShadowBolt == 0)) {
				// Cast Corruption as soon as the 5% crit debuff is up
				// Cast Corruption again when you get the execute buff (Death's Embrace)
				spell = warlock.Corruption
			} else if warlock.CorruptionDot.IsActive() && warlock.CorruptionDot.RemainingDuration(sim) < warlock.ApplyCastSpeed(core.GCDDefault) {
				// Emergency Corruption refresh just in case
				spell = warlock.DrainSoul
			} else if warlock.Talents.Haunt && warlock.Haunt.CD.IsReady(sim) && allCDs[0] == 0 && sim.GetRemainingDuration() > warlock.HauntDebuffAura(warlock.CurrentTarget).Duration/2. {
				// Keep Haunt up
				spell = warlock.Haunt
			} else if warlock.Talents.UnstableAffliction && allCDs[1] == 0 && sim.GetRemainingDuration() > warlock.UnstableAfflictionDot.Duration/2. {
				// Keep UA up
				spell = warlock.UnstableAffliction
			} else if sim.GetRemainingDuration() > warlock.CurseOfAgonyDot.Duration/2. &&
			allCDs[2] == 0 && (!warlock.Haunt.CD.IsReady(sim) || allCDs[0] > 0) && allCDs[1] > 0 && warlock.CorruptionDot.IsActive() {
				// Keep Agony up
				spell = warlock.CurseOfAgony
			} else if KeepUpSEStacks && sim.GetRemainingDuration() > time.Second*10 ||
				(core.ShadowMasteryAura(warlock.CurrentTarget).RemainingDuration(sim) < warlock.ShadowBolt.CurCast.CastTime + SBTravelTime && sim.GetRemainingDuration() > core.ShadowMasteryAura(warlock.CurrentTarget).Duration/2.) {
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
				sim.GetRemainingDuration() > warlock.CorruptionDot.Duration/2. {
				spell = warlock.Corruption
			} else if (!warlock.ImmolateDot.IsActive() || warlock.ImmolateDot.RemainingDuration(sim) < warlock.Immolate.CurCast.CastTime) &&
				sim.GetRemainingDuration() > warlock.ImmolateDot.Duration/2. {
				spell = warlock.Immolate
			} else if core.ShadowMasteryAura(warlock.CurrentTarget).RemainingDuration(sim) < warlock.ShadowBolt.CurCast.CastTime + SBTravelTime &&
				sim.GetRemainingDuration() > core.ShadowMasteryAura(warlock.CurrentTarget).Duration/2. {
				// Shadow Mastery refresh
				spell = warlock.ShadowBolt
			} else if warlock.DecimationAura.IsActive() {
				// Demonology execute phase
				filler = warlock.SoulFire
			} else if warlock.MoltenCoreAura.IsActive() {
				// Corruption proc
				filler = warlock.Incinerate
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
		if sim.IsExecutePhase25() && warlock.Talents.SoulSiphon > 0 {
			// Affliction execute phase
			filler = warlock.channelCheck(sim, warlock.DrainSoulDot, 5)
		} else if warlock.DecimationAura.IsActive() {
			// Demonology execute phase
			filler = warlock.SoulFire
		} else if warlock.MoltenCoreAura.IsActive() {
			// Molten Core talent corruption proc (Demonology)
			filler = warlock.Incinerate
		}
		currentSpellPrio := math.MaxInt64 // Lowest priority for a filler spell
		for _, RSI := range warlock.SpellsRotation {
			if RSI.CastIn(sim) == 0 && (RSI.Priority < currentSpellPrio) && RSI.Spell.IsReady(sim) && RSI.Priority != 0 {
				spell = RSI.Spell
				currentSpellPrio = RSI.Priority
			}
		}
		if sim.Log != nil {
			// warlock.Log(sim, "warlock.SpellsRotation[%d]", warlock.SpellsRotation[4].CastIn(sim).Seconds())
		}
	}

	// ------------------------------------------
	// Filler spell && Regen check
	// ------------------------------------------

	if spell == nil {
		// If a CD is really close to be up, wait for it.
		if nextBigCD - sim.CurrentTime > 0 && nextBigCD - sim.CurrentTime < fillerCastTime/10 {
			warlock.WaitUntil(sim, nextBigCD)
			return
		} else if nextCD - sim.CurrentTime > 0 && nextCD - sim.CurrentTime < fillerCastTime/10 {
			warlock.WaitUntil(sim, nextCD+dotLag)
			return
		} else if timeUntilOom < 5*time.Second && timeUntilExecute25 > time.Second {
			// If you were gonna cast a filler but are low mana, get mana instead in order not to be OOM when an important spell is coming up.
			warlock.LifeTapOrDarkPact(sim)
			return
		} else if !warlock.DoingRegen && nextBigCD - sim.CurrentTime < time.Second*6 && sim.GetRemainingDuration() > time.Second*30 {
			// If big CD coming up and we don't have enough mana for it, lifetap
			// Also, never do a big regen in the last few seconds of the fight.
			// TODO: Specify regen goals depending on CD
			if warlock.CurrentManaPercent() < 0.2 {
				warlock.DoingRegen = true
			}
		}

		if warlock.DoingRegen {
			warlock.LifeTapOrDarkPact(sim)
			if warlock.CurrentManaPercent() > 0.2 {
				warlock.DoingRegen = false
			}
			return
		}

		// Filler
		spell = filler
	}


	// ------------------------------------------
	// Spell casting
	// ------------------------------------------

	if success := spell.Cast(sim, target); success {
		warlock.PrevCastSECheck = spell
		if spell == warlock.Corruption && warlock.Talents.EverlastingAffliction > 0 {
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
