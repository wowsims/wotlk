/* README(2022-08):
Main Author: Glandalf (Discord : Glandalf#0679)
Co-Author's: Ketesh (Ketesh#8103),
			Linelo aka "The good bad guy" (Discord: Linelo#3958),
			Pötiküvi(Discord: Pötiküvi#7506)

This file (rotations.go) contains the logic behind how the sim chooses a spell at a given time.
There are two rotation types, Manual & Automatic.
	Automatic predetermines the spell priorities according to tested and theorycrafted information.
	Manual lets users decide their spell priorities and which spells they cast, allowing for further experimentation.


*/

// importing dependencies
package warlock

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

/*
In this section of the code, we will be predefinining some intermediary functions.
These act as a setup for the tryUseGCD() function.
*/
/*	Roll Multiplier & Evaluation
 This part tracks all the damage multiplier that roll over with corruption.
 Everlasting Affliction talent allows you to "Roll" snapshot values for DoT's, carrying their benefits beyond their buff time on you.
  Ex: If you have a 6 seconds Tricks on you with %10 damage increase, you can have your corruption "roll" with that buff indefinately.

These variables are used to estimate how good the roll will be, and determine if refreshing the corruption again will be a DPS increase.
*/

func (warlock *Warlock) corruptionTracker() float64 {
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
	return CurrentDmgMult * CurrentShadowMult * CurrentCritMult
}

func (warlock *Warlock) defineRotation() {
	rotationType := warlock.Rotation.Type
	curse := warlock.Rotation.Curse
	secondaryDot := warlock.Rotation.SecondaryDot
	specSpell := warlock.Rotation.SpecSpell

	/* The warlock rotaitonal spells come mostly in three shapes
	Filler Spells: are the bottom of the spell hierarchy, you cast this whenever you don't have any other "Priorities"
		Ex: Shadowbolt as Affliction or Demo, Incinerate as Destro are your regular fillers.
		Ex: Soul Siphon as Affliction, Soulfire as Demo/Hybrid are your Execute fillers. (More onto this later)
	Priority Spells: Are things you want to aggressively cast whenever the situation calls for it.
		Ex: When Conflag comes off CD for Destro, When UA expires on the target as Affliction, Incinerate when you have MC procs as Demo.
	Regen Spells: Are cast when the best thing to do is regen, either for when you want to prepare a big burst and top your mana reserves, or there is nothing better to cast.
		Ex: Lifetap / DarkPact
	Although the relationship between these are not as simple as I've put out to be, the code below is all about determining which spell is best to cast in a given moment in time,
	*/

	spellBook := [...]*core.Spell{ //These are all of your possible "Priority Spells " being stored in an array.
		warlock.Corruption,
		warlock.Immolate,
		warlock.UnstableAffliction,
		warlock.Haunt,
		warlock.CurseOfAgony,
		warlock.CurseOfDoom,
		warlock.Conflagrate,
		warlock.ChaosBolt,
	}

	warlock.SpellsRotation = make([]SpellRotation, len(spellBook)) // an array containing Spell Rotation structs
	/*
		type SpellRotation struct {
		Spell    *core.Spell //The spell in question
		CastIn   CastReadyness // If the time to cast the spell is right. This is measured in time units. If this is 0, the spell is ready to go!
		Priority int // Priority of the spell. A metric used to compare... you guessed right, the priority of the spell. Higher is cast first.
		}
	*/

	for i, spell := range spellBook { //Setup the spells, match them with their spellBook entries
		warlock.SpellsRotation[i].Spell = spell
	}

	/*
		Now, each spell in the spellbook above, will be evaluated one by one, and will be given a CastIn value, that determines when each spell will be ready to cast.
		return 0: means cast it now, it's ready!
		return core.NeverExpires : This value is used to completely disable a spell, like when the talent for it is missing.
		If a spell is never ready to cast, it won't ever get casted.
	*/

	//Starting off: Corruption
	warlock.SpellsRotation[0].CastIn = func(sim *core.Simulation) time.Duration {

		if !warlock.Rotation.Corruption { //If corruption is not added to the rotation, it will never be cast.
			return core.NeverExpires
		}

		// Affliction spec check
		if warlock.Talents.EverlastingAffliction > 0 {
			if (!warlock.CorruptionDot.IsActive() && (core.ShadowMasteryAura(warlock.CurrentTarget).IsActive() || warlock.Talents.ImprovedShadowBolt == 0)) ||
				// Wait for SM to be applied to cast first Corruption
				warlock.CorruptionDot.IsActive() && (warlock.corruptionTracker() > warlock.CorruptionRolloverPower) {
				// If the active corruption multipliers are lower than the ones for a potential new corruption, then reapply corruption
				return 0
			} else {
				return core.NeverExpires //Never will be cast
			}
		} else {
			return core.MaxDuration(0, warlock.CorruptionDot.RemainingDuration(sim)) // Will be "ready to cast in this many seconds"
		} //This is due to not having EA on, you will have to manually reapply Corr.
	}
	//Immolate
	warlock.SpellsRotation[1].CastIn = func(sim *core.Simulation) time.Duration {
		if !(secondaryDot == proto.Warlock_Rotation_Immolate) || sim.GetRemainingDuration() < warlock.ImmolateDot.Duration/2. {
			return core.NeverExpires
		}
		return core.MaxDuration(0, warlock.ImmolateDot.RemainingDuration(sim)-warlock.ApplyCastSpeed(warlock.Immolate.DefaultCast.CastTime))
		//This return is used as "the time left to refresh the spell"
		//It's remaining duration - time it will take to reapply the spell, when it is 0, you optimally benefit from everytick, while restoring the debuff the milisecond it falls off.
	}
	//UA
	warlock.SpellsRotation[2].CastIn = func(sim *core.Simulation) time.Duration {
		if !warlock.Talents.UnstableAffliction || !(secondaryDot == proto.Warlock_Rotation_UnstableAffliction) {
			return core.NeverExpires
		}
		return core.MaxDuration(0, warlock.UnstableAfflictionDot.RemainingDuration(sim)-warlock.ApplyCastSpeed(warlock.UnstableAffliction.DefaultCast.CastTime))
	}
	/*Haunt
	Haunt is different than all your other DoT's, reason being, it dynamicly amplifies other DoT's for it's duration.
	Meaning, all your other dots,you let them tick to their last second and elapse, and then reapply as soon as possible.
	In Haunt, you want to maximize uptime and that it never falls off.
	It also shares debuff duration with Shadow Embrace, which stacks off to 3, and is a huge dps loss when dropped. */
	warlock.SpellsRotation[3].CastIn = func(sim *core.Simulation) time.Duration {
		if !warlock.Talents.Haunt || !(specSpell == proto.Warlock_Rotation_Haunt) {
			return core.NeverExpires
		}
		hauntSBTravelTime := time.Duration(warlock.DistanceFromTarget/20) * time.Second //Haunt is a projectile, so if you are at range, you have to account distance from target
		hauntCastTime := warlock.ApplyCastSpeed(warlock.Haunt.DefaultCast.CastTime)
		spellCastTime := warlock.ApplyCastSpeed(core.GCDDefault)
		if sim.IsExecutePhase25() {
			spellCastTime = warlock.ApplyCastSpeed(warlock.DrainSoulDot.TickLength)
		}
		return core.MaxDuration(0, warlock.HauntDebuffAura(warlock.CurrentTarget).RemainingDuration(sim)-hauntCastTime-hauntSBTravelTime-spellCastTime)
		//Since Haunt's unique behavior, this return is the "Leeway" you have for the spell. Meaning, if this hits below 0, you are too late and haunt dropped off.
		//On the other hand, reapplying this when not 0, but say 0.5 or 1, is not a tick loss as it is for other dots.
	}
	//Curse of Agony
	warlock.SpellsRotation[4].CastIn = func(sim *core.Simulation) time.Duration {
		if !(curse == proto.Warlock_Rotation_Doom || curse == proto.Warlock_Rotation_Agony) || warlock.CurseOfDoomDot.IsActive() || sim.GetRemainingDuration() < warlock.CurseOfAgonyDot.Duration/2 {
			return core.NeverExpires
		}
		return core.MaxDuration(0, warlock.CurseOfAgonyDot.RemainingDuration(sim))
	}
	//Curse of Doom
	warlock.SpellsRotation[5].CastIn = func(sim *core.Simulation) time.Duration {
		if curse != proto.Warlock_Rotation_Doom || !warlock.CurseOfDoom.IsReady(sim) || sim.GetRemainingDuration() < time.Minute {
			return core.NeverExpires
		}
		return core.MaxDuration(0, warlock.CurseOfDoomDot.RemainingDuration(sim))
	}
	//Conflagrate
	warlock.SpellsRotation[6].CastIn = func(sim *core.Simulation) time.Duration {
		if !warlock.Talents.Conflagrate || !warlock.ImmolateDot.IsActive() {
			return core.NeverExpires
		}

		if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfConflagrate) { //This glyph makes Conflag not consume the Immolate
			// Cast on CD
			return core.MaxDuration(0, warlock.Conflagrate.TimeToReady(sim))
		} else {
			// Cast at the end of an Immolate
			return core.MaxDuration(core.MaxDuration(0, warlock.ImmolateDot.RemainingDuration(sim)-warlock.ImmolateDot.TickLength), warlock.Conflagrate.TimeToReady(sim))
		}
	}
	//Chaos Bolt
	warlock.SpellsRotation[7].CastIn = func(sim *core.Simulation) time.Duration {
		if !warlock.Talents.ChaosBolt || !(specSpell == proto.Warlock_Rotation_ChaosBolt) {
			return core.NeverExpires
		}
		return core.MaxDuration(0, warlock.ChaosBolt.TimeToReady(sim))
		//for spells that have a set CD, and not a DoT, this return simply becomes the current CD on your spell.
	}

	// Rotation Presets: These are the presets representing theorycrafted rotation priorities.
	// We use this variable to distinguish between two spells that are ready at the same time, highest prio is cast first.
	// Value Legend: 0 / Absent = Not cast, 1 is highest prio, 10 is lowest prio.
	if rotationType == proto.Warlock_Rotation_Affliction {
		warlock.SpellsRotation[0].Priority = 1 //Corruption
		warlock.SpellsRotation[2].Priority = 2 //UA
		warlock.SpellsRotation[3].Priority = 3 //Haunt
		warlock.SpellsRotation[4].Priority = 4 //Curse of Agony
	} else if rotationType == proto.Warlock_Rotation_Demonology {
		warlock.SpellsRotation[5].Priority = 1 // Curse of Doom
		warlock.SpellsRotation[0].Priority = 2 // Corruption
		warlock.SpellsRotation[1].Priority = 3 // Immolate
		warlock.SpellsRotation[4].Priority = 4 // Curse of Agony
	} else if rotationType == proto.Warlock_Rotation_Destruction {
		warlock.SpellsRotation[6].Priority = 1 // Conflagrate
		warlock.SpellsRotation[5].Priority = 2 // Curse of Doom
		warlock.SpellsRotation[1].Priority = 3 // Immolate
		warlock.SpellsRotation[7].Priority = 4 // Chaos Bolt
		warlock.SpellsRotation[4].Priority = 5 // Curse of Agony
	}

	//Manual Rotation Feature:
	//This part sets every castable spells prio to the lowest value of 10, to later let the user reorder them.
	//CAUTION:This section is not yet implemented in the UI and is WIP.
	if warlock.Rotation.Corruption && warlock.SpellsRotation[0].Priority == 0 {
		warlock.SpellsRotation[0].Priority = 10
	} else if !warlock.Rotation.Corruption && warlock.SpellsRotation[0].Priority != 0 {
		warlock.SpellsRotation[0].Priority = 0
	}
	if secondaryDot == proto.Warlock_Rotation_Immolate && warlock.SpellsRotation[1].Priority == 0 {
		warlock.SpellsRotation[1].Priority = 10
		warlock.SpellsRotation[2].Priority = 0
	} else if secondaryDot == proto.Warlock_Rotation_UnstableAffliction && warlock.SpellsRotation[2].Priority == 0 {
		warlock.SpellsRotation[1].Priority = 0
		warlock.SpellsRotation[2].Priority = 10
	} else if secondaryDot == proto.Warlock_Rotation_NoSecondaryDot {
		warlock.SpellsRotation[1].Priority = 0
		warlock.SpellsRotation[2].Priority = 0
	}

	if specSpell == proto.Warlock_Rotation_Haunt && warlock.SpellsRotation[3].Priority == 0 {
		warlock.SpellsRotation[3].Priority = 10
		warlock.SpellsRotation[7].Priority = 0
	} else if specSpell == proto.Warlock_Rotation_ChaosBolt && warlock.SpellsRotation[7].Priority == 0 {
		warlock.SpellsRotation[3].Priority = 0
		warlock.SpellsRotation[7].Priority = 10
	} else if specSpell == proto.Warlock_Rotation_NoSpecSpell {
		warlock.SpellsRotation[3].Priority = 0
		warlock.SpellsRotation[7].Priority = 0
	}
	if warlock.Talents.Conflagrate && warlock.SpellsRotation[6].Priority == 0 {
		warlock.SpellsRotation[6].Priority = 1
	}
	if curse == proto.Warlock_Rotation_Doom && warlock.SpellsRotation[5].Priority == 0 {
		warlock.SpellsRotation[5].Priority = 1
	} else if curse != proto.Warlock_Rotation_Doom && warlock.SpellsRotation[5].Priority != 0 {
		warlock.SpellsRotation[5].Priority = 0
	}
}

/*
	At the end of this function, every spell in our arsenal is given;
	 A castIn() value, that determines when the spell will be needed to/will be ready to recast again.
	 A Priority value, that lets us order which one of the "Ready Spells(aka. castIn = 0)" we will be using.
*/

// Regen Spells: Casts the regen spell that will give you the most mana, includes a error whenever we cast pact on full mana.
func (warlock *Warlock) LifeTapOrDarkPact(sim *core.Simulation) {
	if warlock.CurrentManaPercent() == 1 {
		panic("Life Tap or Dark Pact while full mana")
	}
	if warlock.Talents.DarkPact && warlock.Pet.CurrentMana() > warlock.GetStat(stats.SpellPower)+1200+131 { //Evaluates based on your SP, if DP or LT will give you the highest mana.
		warlock.DarkPact.Cast(sim, warlock.CurrentTarget)
	} else {
		warlock.LifeTap.Cast(sim, warlock.CurrentTarget)
	}
}

// This function is an intermediary, it is used when sim has a GCD ready, not much to see here.
func (warlock *Warlock) OnGCDReady(sim *core.Simulation) {
	warlock.tryUseGCD(sim)
}

//preparation function definitions ends::

/*
This function is the way we execute the main functionality of this entire script.
All of the previously implemented functions come together in this function.
Function takes the Warlock character that's used to model the client behavior, and returns the "modified" simulation state.
Might sound complicated, worry not, things will get better.
*/
func (warlock *Warlock) tryUseGCD(sim *core.Simulation) {
	var spell *core.Spell                          //the variable we'll be returning to the sim as our final decision
	var filler *core.Spell                         //the filler spell we'll store, we will cast this whenever we have all our priorities in check
	var target = warlock.CurrentTarget             //our current target
	mainSpell := warlock.Rotation.PrimarySpell     // our primary spell
	curse := warlock.Rotation.Curse                // our curse of choice
	dotLag := time.Duration(10 * time.Millisecond) // the lag time for dots, a small value that allows us to gap two dots properly

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
	//For aoe situations, sets your main spell as Seed.
	//This is currently a WIP.
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

	bigCDs := warlock.GetMajorCooldowns() // all of our major CD's, things like pots, racials, metamorphing power rangers, you name it.
	nextBigCD := core.NeverExpires        // just setting the highest possible value for convenience in declaration
	for _, cd := range bigCDs {           //a loop that iterates over all possible CD's, and orders them based on their time to get ready.
		if cd == nil {
			continue // not on cooldown right now.
		}
		cdReadyAt := cd.Spell.ReadyAt()                                     //Cooldown will be ready in cdReadyAt.
		if cd.Type.Matches(core.CooldownTypeDPS) && cdReadyAt < nextBigCD { //If the cooldown is a DPS cooldown, nextBigCD will be this cd
			nextBigCD = cdReadyAt
		}
	}

	// ------------------------------------------
	// Small CDs (Cast on CD)
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
	if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfLifeTap) && // This glyph gives you a buff to SP when you cast Life Tap, and we want this on at all times.
		(!warlock.GlyphOfLifeTapAura.IsActive() || warlock.GlyphOfLifeTapAura.RemainingDuration(sim) < time.Second) {
		if sim.CurrentTime < time.Second {

			// Pre-Pull Cast Shadow Bolt
			warlock.SpendMana(sim, warlock.ShadowBolt.DefaultCast.Cost, warlock.ShadowBolt.ResourceMetrics)
			warlock.ShadowBolt.SkipCastAndApplyEffects(sim, warlock.CurrentTarget)

			// Pre-pull Life Tap
			warlock.GlyphOfLifeTapAura.Activate(sim)

			//These lines emulate you pre-casting a shadowbolt and having Life Tap on
			//TODO: Illustration of Dragon Soul stacking to 10 with Life Funnel.
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

	castDebuffCurse := func(spellToCast *core.Spell, aura *core.Aura) bool {
		if !aura.IsActive() {
			spell = spellToCast
			return true
		}
		return false
	} // a simple function for our debuff curses, we want to apply them immediately as they drop, disregarding any priority as responsible debuff slaves.

	switch curse {
	case proto.Warlock_Rotation_Elements:
		castDebuffCurse(warlock.CurseOfElements, warlock.CurseOfElementsAura)
	case proto.Warlock_Rotation_Weakness:
		castDebuffCurse(warlock.CurseOfWeakness, warlock.CurseOfWeaknessAura)
	case proto.Warlock_Rotation_Tongues:
		castDebuffCurse(warlock.CurseOfTongues, warlock.CurseOfTonguesAura)
	}

	if spell != nil {
		if curse != proto.Warlock_Rotation_Doom && curse != proto.Warlock_Rotation_Agony {
			if success := spell.Cast(sim, target); success {
				return
			}
		}
	}

	// ------------------------------------------
	// Main spells
	// ------------------------------------------

	// We're kind of trying to fit all different spec rotations in one big priority based rotation in order to let people experiment
	if filler == nil { // There are two main fillers in warlocks arsenal, SB for Affliction & Demo, Incinerate for Destro
		switch mainSpell {
		case proto.Warlock_Rotation_ShadowBolt:
			filler = warlock.ShadowBolt
		case proto.Warlock_Rotation_Incinerate:
			filler = warlock.Incinerate
		}
	}

	/* Execute Phase: Warlock, with it's Demo/Affliction/Hybrid specs, rely on execute mechanics;
	Affliction:
		Death's Embrace : When the target is at or below %35 hp, all shadow damage done is %12 increased flat.
		Drain Soul : Spell has a built in mechanic, that does 4 Times it's normal damage, if the target is at or below %25 HP.
	Demo & Hybrid:
		Decimation: Hitting a target at or below %35 hp, will make your Soul Fire spell %40 faster, and cost no Soul Shards
	*/
	if sim.IsExecutePhase25() && warlock.Talents.SoulSiphon > 0 {
		// Drain Soul phase, Soul Siphon is an affliction talent, so if below %25, and you have siphon talented, sim assumes you are Affliction.
		filler = warlock.channelCheck(sim, warlock.DrainSoulDot, 5) //This function checks if you are channeling DS or not. Returns continuing channeling or casting DS actions accordingly.
	} else if warlock.DecimationAura.IsActive() { //Molten Core, buffs Incinerate and Soul Fire, however, since below %35 you will spam SF, you don't need to check Molten Core.
		// Demo & Hybrid execute phase
		filler = warlock.SoulFire
	} else if warlock.MoltenCoreAura.IsActive() { //Now you have to alternate to Incinerate from Shadow Bolts, so you will check for MC's
		// Molten Core talent corruption proc (Demonology)
		filler = warlock.Incinerate
	}
	//The loop that filters the currently ready Priority Spells, and decides on which to cast based on Priority.
	nextSpell := core.NeverExpires               // declaration convention, don't worry, just breathe
	currentSpell := core.NeverExpires            // in..... and out, breath in fully....
	currentSpellPrio := math.MaxInt64            // Lowest priority for a filler spell, oh btw, remember to keep breathing
	for _, RSI := range warlock.SpellsRotation { // For all the spells off cooldown (aka. castIn = 0, check the explanations in warlock.SpellsReady if this makes no sense )
		currentSpell = RSI.CastIn(sim)
		if currentSpell < nextSpell {
			nextSpell = currentSpell
		}
		if currentSpell == 0 && (RSI.Priority < currentSpellPrio) && RSI.Spell.IsReady(sim) && RSI.Priority != 0 {
			spell = RSI.Spell
			currentSpellPrio = RSI.Priority
		} // find and cast the highest prio
	}
	nextSpell += sim.CurrentTime
	if sim.Log != nil {
		// warlock.Log(sim, "warlock.SpellsRotation[%d]", warlock.SpellsRotation[4].CastIn(sim).Seconds())
	}

	// ------------------------------------------
	// Filler spell && Regen check
	// ------------------------------------------
	//We decide if we can cast our fillers without repercussions,
	var ManaSpendRate float64
	var fillerCastTime time.Duration
	if warlock.Talents.SoulSiphon > 0 { //SoulSiphon >0 is an affliction check.
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
		} else if nextSpell-sim.CurrentTime > 0 && nextSpell-sim.CurrentTime < fillerCastTime/10 {
			// The dot lag is currently here only for UI purposes, without which the last dot tick is shown as part of the next dot cast
			warlock.WaitUntil(sim, nextSpell+dotLag)
			return
		}

		var executeDuration float64
		// Estimate for desired mana needed to do affliction execute
		var DesiredManaAtExecute float64
		if warlock.Talents.Decimation > 0 {
			// We suppose that if you would want to use Soul Fire as an execute filler if and only if you have the Decimation talent.
			executeDuration = 0.35
			DesiredManaAtExecute = 0.3 * sim.Duration.Seconds() * executeDuration / 60
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

		// After all the previous checks, if everything checks out, you are free to cast filler.
		spell = filler
	}

	// This part tracks all the damage multiplier that roll over with corruption
	PotentialCorruptionRolloverPower := warlock.corruptionTracker()
	if sim.Log != nil {
		if warlock.Talents.EverlastingAffliction > 0 {
			warlock.Log(sim, "[Info] Active Corruption rollover power [%.2f]", warlock.CorruptionRolloverPower)
			warlock.Log(sim, "[Info] Potential Corruption rollover power [%.2f]", PotentialCorruptionRolloverPower)
		}
		if warlock.Talents.DemonicPact > 0 {
			warlock.Log(sim, "[Info] Demonic Pact spell power bonus average [%.0f]", warlock.DPSPAverage)
		}
	}

	// ------------------------------------------
	// Spell casting
	// ------------------------------------------
	if success := spell.Cast(sim, target); success {
		if spell == warlock.Corruption && warlock.Talents.EverlastingAffliction > 0 {
			// We are recording the current rollover power of corruption
			warlock.CorruptionRolloverPower = PotentialCorruptionRolloverPower
		}
		return
	} else {
		// Regen Cast if you can't cast anything else.
		warlock.LifeTapOrDarkPact(sim)
		return
	}

}
