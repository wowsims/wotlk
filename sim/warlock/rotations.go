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

	// High priority spells (does not include fillers)
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
	for i, spell := range spellBook {
		warlock.SpellsRotation[i].Spell = spell
	}

	// Calculation of spell readyness in time unit (use same order as in spellBook)
	// The associated spell will not be cast before time is 0
	// 0 : Cast ready ; core.NeverExpires : never cast
	warlock.SpellsRotation[0].CastIn = func(sim *core.Simulation) time.Duration {
		// Checking if it's the manual rotation
		if !warlock.Rotation.Corruption {
			return core.NeverExpires
		}
		// This part tracks all the damage multiplier that roll over with corruption
		// Shadow damage multipler (looking for DE)
		CurrentShadowMult := warlock.PseudoStats.ShadowDamageDealtMultiplier
		// Damage multipler (looking for TotT)
		CurrentDmgMult := warlock.PseudoStats.DamageDealtMultiplier
		// Crit rating multipler (looking for Shadow Mastery (ISB Talent) and Potion of Wild Magic)
		CurrentCritBonus := warlock.GetStat(stats.SpellCrit) + warlock.PseudoStats.BonusSpellCritRating + warlock.PseudoStats.BonusShadowCritRating +
			warlock.CurrentTarget.PseudoStats.BonusSpellCritRatingTaken
		CurrentCritMult := 1 + CurrentCritBonus/core.CritRatingPerCritChance/100*core.TernaryFloat64(warlock.Talents.Pandemic, 1, 0)
		// Combination of all multipliers
		CurrentCorruptionRolloverMult := CurrentDmgMult * CurrentShadowMult * CurrentCritMult
		// Affliction spec check
		if warlock.Talents.EverlastingAffliction > 0 {
			if (!warlock.CorruptionDot.IsActive() && (core.ShadowMasteryAura(warlock.CurrentTarget).IsActive() || warlock.Talents.ImprovedShadowBolt == 0)) ||
				// Wait for SM to be applied to cast first Corruption
				warlock.CorruptionDot.IsActive() && (CurrentCorruptionRolloverMult > warlock.CorruptionRolloverMult) {
				// If the original corruption multipliers are lower than this current time, then reapply corruption
				return 0
			} else {
				return core.NeverExpires
			}
		} else {
			return core.MaxDuration(0, warlock.CorruptionDot.RemainingDuration(sim))
		}
	}
	warlock.SpellsRotation[1].CastIn = func(sim *core.Simulation) time.Duration {
		if !(secondaryDot == proto.Warlock_Rotation_Immolate) || sim.GetRemainingDuration() < warlock.ImmolateDot.Duration/2. {
			return core.NeverExpires
		}
		return core.MaxDuration(0, warlock.ImmolateDot.RemainingDuration(sim)-warlock.ApplyCastSpeed(warlock.Immolate.DefaultCast.CastTime))
	}
	warlock.SpellsRotation[2].CastIn = func(sim *core.Simulation) time.Duration {
		if !warlock.Talents.UnstableAffliction || !(secondaryDot == proto.Warlock_Rotation_UnstableAffliction) {
			return core.NeverExpires
		}
		return core.MaxDuration(0, warlock.UnstableAfflictionDot.RemainingDuration(sim)-warlock.ApplyCastSpeed(warlock.UnstableAffliction.DefaultCast.CastTime))
	}
	warlock.SpellsRotation[3].CastIn = func(sim *core.Simulation) time.Duration {
		if !warlock.Talents.Haunt || !(specSpell == proto.Warlock_Rotation_Haunt) {
			return core.NeverExpires
		}
		hauntSBTravelTime := time.Duration(warlock.DistanceFromTarget/20) * time.Second
		hauntCastTime := warlock.ApplyCastSpeed(warlock.Haunt.DefaultCast.CastTime)
		spellCastTime := warlock.ApplyCastSpeed(core.GCDDefault)
		if sim.IsExecutePhase25() {
			spellCastTime = warlock.ApplyCastSpeed(warlock.DrainSoulDot.TickLength)
		}
		// If SE remaining duration is less than a Haunt cast time + travel time
		// (+ whichever current filler cast time so that we don't start a cast that would drop haunt)
		// and the previous cast was not Haunt or SB then cast shadow bolt so SE stacks are not lost
		KeepUpSEStacks := (warlock.PrevCastSECheck != warlock.Haunt && warlock.PrevCastSECheck != warlock.ShadowBolt &&
			warlock.ShadowEmbraceDebuffAura(warlock.CurrentTarget).RemainingDuration(sim) < hauntCastTime+hauntSBTravelTime)
		if KeepUpSEStacks && sim.GetRemainingDuration() > time.Second*10 && warlock.Haunt.IsReady(sim) {
			return 0
		} else {
			return core.MaxDuration(0, warlock.HauntDebuffAura(warlock.CurrentTarget).RemainingDuration(sim)-hauntCastTime-hauntSBTravelTime-spellCastTime)
		}
	}
	warlock.SpellsRotation[4].CastIn = func(sim *core.Simulation) time.Duration {
		if !(curse == proto.Warlock_Rotation_Doom || curse == proto.Warlock_Rotation_Agony) || warlock.CurseOfDoomDot.IsActive() || sim.GetRemainingDuration() < warlock.CurseOfAgonyDot.Duration/2 {
			return core.NeverExpires
		}
		return core.MaxDuration(0, warlock.CurseOfAgonyDot.RemainingDuration(sim))
	}
	warlock.SpellsRotation[5].CastIn = func(sim *core.Simulation) time.Duration {
		if curse != proto.Warlock_Rotation_Doom || !warlock.CurseOfDoom.IsReady(sim) || sim.GetRemainingDuration() < time.Minute {
			return core.NeverExpires
		}
		return core.MaxDuration(0, warlock.CurseOfDoomDot.RemainingDuration(sim))
	}
	warlock.SpellsRotation[6].CastIn = func(sim *core.Simulation) time.Duration {
		if !warlock.Talents.Conflagrate {
			return core.NeverExpires
		}
		if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfConflagrate) {
			return core.MaxDuration(0, warlock.Conflagrate.TimeToReady(sim))
		} else {
			return core.MaxDuration(0, warlock.ImmolateDot.RemainingDuration(sim)-warlock.ImmolateDot.TickLength)
		}
	}
	warlock.SpellsRotation[7].CastIn = func(sim *core.Simulation) time.Duration {
		if !warlock.Talents.ChaosBolt || !(specSpell == proto.Warlock_Rotation_ChaosBolt) {
			return core.NeverExpires
		}
		return core.MaxDuration(0, warlock.ChaosBolt.TimeToReady(sim))
	}

	// Priority based rotations (0 or absent means not in rotation, 1 is max)
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
		warlock.SpellsRotation[1].Priority = 3
		warlock.SpellsRotation[4].Priority = 5
		warlock.SpellsRotation[5].Priority = 2
		warlock.SpellsRotation[6].Priority = 1
		warlock.SpellsRotation[7].Priority = 4
	}

	// For Manual rotation, give spells lowest prio if user wants to experiment
	if warlock.Rotation.Corruption && warlock.SpellsRotation[0].Priority == 0 {
		warlock.SpellsRotation[0].Priority = 10
	}
	if secondaryDot == proto.Warlock_Rotation_Immolate && warlock.SpellsRotation[1].Priority == 0 {
		warlock.SpellsRotation[1].Priority = 10
	} else if secondaryDot == proto.Warlock_Rotation_UnstableAffliction && warlock.SpellsRotation[2].Priority == 0 {
		warlock.SpellsRotation[2].Priority = 10
	}
	if specSpell == proto.Warlock_Rotation_Haunt && warlock.SpellsRotation[3].Priority == 0 {
		warlock.SpellsRotation[3].Priority = 10
	} else if specSpell == proto.Warlock_Rotation_ChaosBolt && warlock.SpellsRotation[7].Priority == 0 {
		warlock.SpellsRotation[7].Priority = 10
	}
	if warlock.Talents.Conflagrate && warlock.SpellsRotation[6].Priority == 0 {
		warlock.SpellsRotation[6].Priority = 1
	}
	if curse == proto.Warlock_Rotation_Doom && warlock.SpellsRotation[5].Priority == 0 {
		warlock.SpellsRotation[5].Priority = 1
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
	curse := warlock.Rotation.Curse
	dotLag := time.Duration(10 * time.Millisecond)

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
	nextBigCD := core.NeverExpires
	for _, cd := range bigCDs {
		if cd == nil {
			continue // not on cooldown right now.
		}
		cdReadyAt := cd.Spell.ReadyAt()
		if cd.Type.Matches(core.CooldownTypeDPS) && cdReadyAt < nextBigCD {
			nextBigCD = cdReadyAt
		}
	}

	if nextBigCD-sim.CurrentTime <= 0 {
		// stop regen, start blasting
		warlock.DoingRegen = false
	}


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
	}

	if spell != nil {
		if curse != proto.Warlock_Rotation_Doom && curse != proto.Warlock_Rotation_Agony {
			if success := spell.Cast(sim, target); success {
				warlock.PrevCastSECheck = spell
				return
			}
		}
	}

	// ------------------------------------------
	// Main spells
	// ------------------------------------------

	// We're kind of trying to fit all different spec rotations in one big priority based rotation in order to let people experiment
	if filler == nil {
		switch mainSpell {
		case proto.Warlock_Rotation_ShadowBolt:
			filler = warlock.ShadowBolt
		case proto.Warlock_Rotation_Incinerate:
			filler = warlock.Incinerate
		}
	}

	// If SE remaining duration is less than a shadow bolt cast time + travel time (with a 3 second buffer to include 1 drain soul tick) and the previous cast was not haunt or SB then cast shadow bolt so SE stacks are not lost
	KeepUpSEStacksExecute := (warlock.PrevCastSECheck != warlock.Haunt && warlock.PrevCastSECheck != warlock.ShadowBolt &&
		warlock.ShadowEmbraceDebuffAura(warlock.CurrentTarget).RemainingDuration(sim).Seconds() < warlock.ApplyCastSpeed(time.Duration(warlock.ShadowBolt.DefaultCast.CastTime)).Seconds()+warlock.DistanceFromTarget/20+3)
	
	// The default filler can change because of some execute phase or proc
	if sim.IsExecutePhase25() && warlock.Talents.SoulSiphon > 0 && !KeepUpSEStacksExecute {
		// Affliction execute phase
		filler = warlock.channelCheck(sim, warlock.DrainSoulDot, 5)
	} else if warlock.DecimationAura.IsActive() {
		// Demonology execute phase
		filler = warlock.SoulFire
	} else if warlock.MoltenCoreAura.IsActive() {
		// Molten Core talent corruption proc (Demonology)
		filler = warlock.Incinerate
	}
	nextCD := core.NeverExpires
	currentCD := core.NeverExpires
	currentSpellPrio := math.MaxInt64 // Lowest priority for a filler spell
	for _, RSI := range warlock.SpellsRotation {
		currentCD = RSI.CastIn(sim)
		if currentCD < nextCD {
			nextCD = currentCD
		}
		if currentCD == 0 && (RSI.Priority < currentSpellPrio) && RSI.Spell.IsReady(sim) && RSI.Priority != 0 {
			spell = RSI.Spell
			currentSpellPrio = RSI.Priority
		}
	}
	nextCD += sim.CurrentTime
	if sim.Log != nil {
		// warlock.Log(sim, "warlock.SpellsRotation[%d]", warlock.SpellsRotation[4].CastIn(sim).Seconds())
	}

	// ------------------------------------------
	// Filler spell && Regen check
	// ------------------------------------------
	
	var ManaSpendRate float64
	var fillerCastTime time.Duration
	if warlock.Talents.SoulSiphon > 0 {
		fillerCastTime = warlock.ApplyCastSpeed(warlock.ShadowBolt.DefaultCast.CastTime)
		ManaSpendRate = warlock.ShadowBolt.BaseCost / float64(fillerCastTime.Seconds())
	} else {
		fillerCastTime = warlock.ApplyCastSpeed(filler.DefaultCast.CastTime)
		ManaSpendRate = filler.BaseCost / float64(fillerCastTime.Seconds())
	}

	if spell == nil {
		// If a CD is really close to be up, wait for it.
		if nextBigCD-sim.CurrentTime > 0 && nextBigCD-sim.CurrentTime < fillerCastTime/10 {
			warlock.WaitUntil(sim, nextBigCD)
			return
		} else if nextCD-sim.CurrentTime > 0 && nextCD-sim.CurrentTime < fillerCastTime/10 {
			// The dot lag is currently here only for UI purposes, without which the last dot tick is shown as part of the next dot cast
			warlock.WaitUntil(sim, nextCD+dotLag)
			return
		}

		var executeDuration float64
		// Estimate for desired mana needed to do affliction execute
		var DesiredManaAtExecute float64
		if warlock.Talents.Decimation > 0 {
			// We suppose that if you would want to use Soul Fire as an execute filler if and only if you have the Decimation talent.
			executeDuration = 0.35
			DesiredManaAtExecute = 0.3*sim.Duration.Seconds()*executeDuration/60
		} else if warlock.Talents.SoulSiphon > 0 {
			// We suppose that if you would want to use Drain Soul as an execute filler if and only if you have the Soul Siphon talent.
			executeDuration = 0.25
			DesiredManaAtExecute = 0.02
		}
		TotalManaAtExecute := warlock.MaxMana() * DesiredManaAtExecute
		// TotalManaAtExecute := executeDuration*sim.Duration.Seconds()/ManaSpendRate
		timeUntilOom := time.Duration((warlock.CurrentMana()-TotalManaAtExecute)/ManaSpendRate) * time.Second
		timeUntilExecute := time.Duration((sim.GetRemainingDurationPercent() - executeDuration) * float64(sim.Duration))

		if sim.Log != nil {
			warlock.Log(sim, "DesiredManaAtExecute[%d]", DesiredManaAtExecute)
		}

		if timeUntilOom < time.Second && timeUntilExecute > time.Second && warlock.CurrentManaPercent() < 0.8 {
			// If you were gonna cast a filler but are low mana, get mana instead in order not to be OOM when an important spell is coming up.
			// warlock.CurrentManaPercent() < 0.8 is here to prevent overlifetapping early in the sim since timeUntilOom could still be
			// really low since the reference is the execute time expected mana.
			warlock.LifeTapOrDarkPact(sim)
			return
		}

		// Filler
		spell = filler
	}


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
	} else {
	// Lifetap if can't cast
		warlock.LifeTapOrDarkPact(sim)
		return
	}

}
