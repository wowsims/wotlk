package warlock

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) OnGCDReady(sim *core.Simulation) {
	warlock.tryUseGCD(sim)
}

func (warlock *Warlock) tryUseGCD(sim *core.Simulation) {
	var spell *core.Spell
	var target = warlock.CurrentTarget
	mainSpell := warlock.Rotation.PrimarySpell
	secondaryDot := warlock.Rotation.SecondaryDot
	specSpell := warlock.Rotation.SpecSpell
	preset := warlock.Rotation.Preset
	curse := warlock.Rotation.Curse

	// ------------------------------------------
	// AoE (Seed)
	// ------------------------------------------
	if mainSpell == proto.Warlock_Rotation_Seed {
		if warlock.Rotation.DetonateSeed {
			if success := warlock.Seeds[0].Cast(sim, target); !success {
				warlock.LifeTap.Cast(sim, target)
			}
			return
		}

		// If we aren't "auto popping" just put seed on and shadowbolt it.
		if !warlock.SeedDots[0].IsActive() {
			if success := warlock.Seeds[0].Cast(sim, target); success {
				return
			} else {
				warlock.LifeTap.Cast(sim, target)
				return
			}
		}

		// If target has seed, fire a shadowbolt at main target so we start some explosions
		mainSpell = proto.Warlock_Rotation_Shadowbolt
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

	// ------------------------------------------
	// Regen check
	// ------------------------------------------
	// If big CD coming up and we don't have enough mana for it, lifetap
	// Also, never do a big regen in the last few seconds of the fight.
	if !warlock.DoingRegen && nextBigCD-sim.CurrentTime < time.Second*5 && sim.GetRemainingDuration() > time.Second*30 {
		if warlock.GetStat(stats.SpellPower) > warlock.GetInitialStat(stats.SpellPower) || warlock.HasTemporarySpellCastSpeedIncrease() {
			// never start regen if you have boosted sp or boosted cast speed
		} else if warlock.CurrentManaPercent() < 0.2 {
			warlock.DoingRegen = true
		}
	}

	if warlock.DoingRegen {
		if nextBigCD-sim.CurrentTime < time.Second*2 {
			// stop regen, start blasting
			warlock.DoingRegen = false
		} else {
			warlock.LifeTap.Cast(sim, target)
			if warlock.CurrentManaPercent() > 0.6 {
				warlock.DoingRegen = false
			}
			return
		}
	}


	// ------------------------------------------
	// Small CDs
	// ------------------------------------------
	if warlock.Talents.DemonicEmpowerment && warlock.DemonicEmpowerment.CD.IsReady(sim) {
		warlock.DemonicEmpowerment.Cast(sim, target)
	}

	// ------------------------------------------
	// Keep Glyph of Life Tap buff up
	// ------------------------------------------
	if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfLifeTap) && !warlock.GlyphOfLifeTapAura.IsActive() {
		warlock.LifeTap.Cast(sim, target)
		return
	}

	// ------------------------------------------
	// Preset Rotations
	// ------------------------------------------
	if preset == proto.Warlock_Rotation_Automatic {
		// ------------------------------------------
		// Affliction Rotation
		// ------------------------------------------
		if warlock.Talents.Haunt {
			if !warlock.CurseOfAgonyDot.IsActive() {
				spell = warlock.CurseOfAgony
			} else if !warlock.CorruptionDot.IsActive() && sim.GetRemainingDuration() > time.Second*24 {
				spell = warlock.Corruption
			} else if warlock.CorruptionDot.IsActive() && warlock.CorruptionDot.TickCount > warlock.CorruptionDot.NumberOfTicks-2 {
				spell = warlock.Shadowbolt
			} else if !warlock.UnstableAffDot.IsActive() {
				spell = warlock.UnstableAff
			} else if !warlock.HauntAura.IsActive() && warlock.Haunt.CD.IsReady(sim) {
				spell = warlock.Haunt
			} else {
				spell = warlock.Shadowbolt
			}
		} else if warlock.Talents.Metamorphosis {

		// ------------------------------------------
		// Demonology Rotation
		// ------------------------------------------
			if warlock.CurseOfDoom.CD.IsReady(sim) && sim.GetRemainingDuration() > time.Minute {
				spell = warlock.CurseOfDoom
			} else if sim.GetRemainingDuration() > time.Second*24 && !warlock.CurseOfAgonyDot.IsActive() && !warlock.CurseOfDoomDot.IsActive() {
				// Can't cast agony until we are at end and both agony and doom are not ticking.
				spell = warlock.CurseOfAgony
			} else if !warlock.CorruptionDot.IsActive() {
				spell = warlock.Corruption
			} else if !warlock.ImmolateDot.IsActive() {
				spell = warlock.Immolate
			} else if warlock.DecimationAura.IsActive() {
				spell = warlock.SoulFire
			} else if warlock.MoltenCoreAura.IsActive() {
				spell = warlock.Incinerate
			} else {
				spell = warlock.Shadowbolt
			}
		} else if warlock.Talents.ChaosBolt {

		// ------------------------------------------
		// Destruction Rotation
		// ------------------------------------------
			if warlock.CurseOfDoom.CD.IsReady(sim) && sim.GetRemainingDuration() > time.Minute {
				spell = warlock.CurseOfDoom
			} else if sim.GetRemainingDuration() > time.Second*24 && !warlock.CurseOfAgonyDot.IsActive() && !warlock.CurseOfDoomDot.IsActive() {
				// Can't cast agony until we are at end and both agony and doom are not ticking.
				spell = warlock.CurseOfAgony
			} else if warlock.CanConflagrate(sim) && (warlock.ImmolateDot.TickCount > warlock.ImmolateDot.NumberOfTicks-2 || warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfConflagrate)) {
				spell = warlock.Conflagrate
			} else if !warlock.CorruptionDot.IsActive() {
				spell = warlock.Corruption
			} else if !warlock.ImmolateDot.IsActive() {
				spell = warlock.Immolate
			} else if warlock.ChaosBolt.CD.IsReady(sim) {
				spell = warlock.ChaosBolt
			} else {
				spell = warlock.Incinerate
			}
		} else { preset = proto.Warlock_Rotation_Manual }
	}

	// ------------------------------------------
	// Manual Rotation
	// ------------------------------------------
	if preset == proto.Warlock_Rotation_Manual {

		// ------------------------------------------
		// Curses (priority)
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
		default:
			fallthrough
		case proto.Warlock_Rotation_Doom:
			if sim.GetRemainingDuration() < time.Minute {
				// Can't cast agony until we are at end and both agony and doom are not ticking.
				if sim.GetRemainingDuration() > time.Second*24 && !warlock.CurseOfAgonyDot.IsActive() && !warlock.CurseOfDoomDot.IsActive() {
					spell = warlock.CurseOfAgony
				}
			} else if warlock.CurseOfDoom.CD.IsReady(sim) && !warlock.CurseOfDoomDot.IsActive() {
				spell = warlock.CurseOfDoom
			}
		case proto.Warlock_Rotation_Agony:
			if !warlock.CurseOfAgonyDot.IsActive() {
				spell = warlock.CurseOfAgony
			}
		}
		if spell != nil {
			if !spell.Cast(sim, target) {
				warlock.LifeTap.Cast(sim, target)
			}
			return
		}

		// ------------------------------------------
		// Main spells
		// ------------------------------------------
		if warlock.Talents.ChaosBolt && specSpell == proto.Warlock_Rotation_ChaosBolt && warlock.ChaosBolt.CD.IsReady(sim) {
			spell = warlock.ChaosBolt
		} else if warlock.Talents.Haunt && specSpell == proto.Warlock_Rotation_Haunt && warlock.Haunt.CD.IsReady(sim) && !warlock.HauntAura.IsActive() {
			spell = warlock.Haunt
		} else if warlock.Rotation.Corruption && !warlock.CorruptionDot.IsActive() {
			spell = warlock.Corruption
		} else if warlock.CanConflagrate(sim) && (warlock.ImmolateDot.TickCount > warlock.ImmolateDot.NumberOfTicks-2 || warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfConflagrate)) {
			spell = warlock.Conflagrate
		} else if warlock.Talents.UnstableAffliction && secondaryDot == proto.Warlock_Rotation_UnstableAffliction && !warlock.UnstableAffDot.IsActive() {
			spell = warlock.UnstableAff
		} else if secondaryDot == proto.Warlock_Rotation_Immolate && !warlock.ImmolateDot.IsActive() {
			spell = warlock.Immolate
		} else {
			switch mainSpell {
			case proto.Warlock_Rotation_Shadowbolt:
				spell = warlock.Shadowbolt
			case proto.Warlock_Rotation_Incinerate:
				spell = warlock.Incinerate
			default:
				panic("No primary spell set")
			}
		}

	} 


	// ------------------------------------------
	// Spell casting
	// ------------------------------------------

	if success := spell.Cast(sim, target); success {
		return
	}

	// Lifetap if nothing else
	if warlock.CurrentManaPercent() < 0.8 {
		warlock.LifeTap.Cast(sim, target)
		return
	}

	// If we get here, something's wrong
}
