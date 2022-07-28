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
	rotationType := warlock.Rotation.Type
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
	case proto.Warlock_Rotation_Doom:
		if warlock.CurseOfDoom.CD.IsReady(sim) && sim.GetRemainingDuration() > time.Minute {
			spell = warlock.CurseOfDoom
		} else if sim.GetRemainingDuration() > time.Second*24 && !warlock.CurseOfAgonyDot.IsActive() && !warlock.CurseOfDoomDot.IsActive() {
			spell = warlock.CurseOfAgony
		}
	case proto.Warlock_Rotation_Agony:
		if sim.GetRemainingDuration() > time.Second*24 && !warlock.CurseOfAgonyDot.IsActive() {
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
	// Preset Rotations
	// ------------------------------------------
	if preset == proto.Warlock_Rotation_Automatic {
		// ------------------------------------------
		// Affliction Rotation
		// ------------------------------------------
		if rotationType == proto.Warlock_Rotation_Affliction {
			if !warlock.CorruptionDot.IsActive() && (core.ShadowMasteryAura(warlock.CurrentTarget).IsActive() || warlock.Talents.ImprovedShadowBolt == 0) ||
				sim.IsExecutePhase35() && time.Duration(warlock.CorruptionDot.TickCount)*warlock.CorruptionDot.TickLength > sim.CurrentTime {
				// Cast Corruption as soon as the 5% crit debuff is up
				// Cast Corruption again when you get the execute buff (Death's Embrace)
				spell = warlock.Corruption
			} else if warlock.CorruptionDot.IsActive() && warlock.CorruptionDot.RemainingDuration(sim) < core.GCDDefault {
				// Emergency Corruption refresh just in case
				spell = warlock.DrainSoul
			} else if warlock.Talents.UnstableAffliction && (!warlock.UnstableAffDot.IsActive() || warlock.UnstableAffDot.RemainingDuration(sim) < warlock.UnstableAff.CurCast.CastTime) &&
				sim.GetRemainingDuration() > warlock.UnstableAffDot.Duration {
				// Keep UA up
				spell = warlock.UnstableAff
			} else if warlock.Talents.Haunt && warlock.Haunt.CD.IsReady(sim) && 2*sim.GetRemainingDuration() > warlock.HauntAura.Duration && warlock.CorruptionDot.IsActive() {
				// Keep Haunt up
				spell = warlock.Haunt
			} else if sim.IsExecutePhase25() {
				// Drain Soul execute phase
				spell = warlock.channelCheck(sim, warlock.DrainSoulDot, 5)
			} else {
				// Filler
				spell = warlock.ShadowBolt
			}
		} else if rotationType == proto.Warlock_Rotation_Demonology {

			// ------------------------------------------
			// Demonology Rotation
			// ------------------------------------------
			if !warlock.CorruptionDot.IsActive() {
				spell = warlock.Corruption
			} else if !warlock.ImmolateDot.IsActive() || warlock.ImmolateDot.RemainingDuration(sim) < warlock.Immolate.CurCast.CastTime {
				spell = warlock.Immolate
			} else if warlock.DecimationAura.IsActive() {
				spell = warlock.SoulFire
			} else if warlock.MoltenCoreAura.IsActive() {
				spell = warlock.Incinerate
			} else {
				spell = warlock.ShadowBolt
			}
		} else if rotationType == proto.Warlock_Rotation_Destruction {

			// ------------------------------------------
			// Destruction Rotation
			// ------------------------------------------
			if warlock.CanConflagrate(sim) && (warlock.ImmolateDot.TickCount > warlock.ImmolateDot.NumberOfTicks-2 || warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfConflagrate)) {
				spell = warlock.Conflagrate
			} else if !warlock.CorruptionDot.IsActive() {
				spell = warlock.Corruption
			} else if !warlock.ImmolateDot.IsActive() || warlock.ImmolateDot.RemainingDuration(sim) < warlock.Immolate.CurCast.CastTime {
				spell = warlock.Immolate
			} else if warlock.Talents.ChaosBolt && warlock.ChaosBolt.CD.IsReady(sim) {
				spell = warlock.ChaosBolt
			} else {
				spell = warlock.Incinerate
			}
		}
	} else {
		preset = proto.Warlock_Rotation_Manual
		warlock.Rotation.Preset = proto.Warlock_Rotation_Manual
	}

	// ------------------------------------------
	// Manual Rotation
	// ------------------------------------------

	// ------------------------------------------
	// Main spells
	// ------------------------------------------
	if preset == proto.Warlock_Rotation_Manual {
		if warlock.Rotation.Corruption && (!warlock.CorruptionDot.IsActive() && core.ShadowMasteryAura(warlock.CurrentTarget).IsActive() ||
			sim.IsExecutePhase35() && time.Duration(warlock.CorruptionDot.TickCount)*warlock.CorruptionDot.TickLength > sim.CurrentTime) {
			// Cast Corruption as soon as the 5% crit debuff is up
			// Cast Corruption again when you get the execute buff (Death's Embrace)
			spell = warlock.Corruption
		} else if warlock.CanConflagrate(sim) && (warlock.ImmolateDot.TickCount > warlock.ImmolateDot.NumberOfTicks-2 || warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfConflagrate)) {
			spell = warlock.Conflagrate
		} else if secondaryDot == proto.Warlock_Rotation_Immolate && (!warlock.ImmolateDot.IsActive() || warlock.ImmolateDot.RemainingDuration(sim) < warlock.Immolate.CurCast.CastTime) {
			spell = warlock.Immolate
		} else if warlock.Talents.UnstableAffliction && secondaryDot == proto.Warlock_Rotation_UnstableAffliction &&
			(!warlock.UnstableAffDot.IsActive() || warlock.UnstableAffDot.RemainingDuration(sim) < warlock.UnstableAff.CurCast.CastTime) &&
			sim.GetRemainingDuration() > warlock.UnstableAffDot.Duration {
			spell = warlock.UnstableAff
		} else if warlock.Talents.Decimation > 0 && warlock.DecimationAura.IsActive() {
			spell = warlock.SoulFire
		} else if warlock.Talents.MoltenCore > 0 && warlock.MoltenCoreAura.IsActive() {
			spell = warlock.Incinerate
		} else if warlock.Talents.ChaosBolt && specSpell == proto.Warlock_Rotation_ChaosBolt && warlock.ChaosBolt.CD.IsReady(sim) {
			spell = warlock.ChaosBolt
		} else if warlock.Talents.Haunt && specSpell == proto.Warlock_Rotation_Haunt && warlock.Haunt.CD.IsReady(sim) && !warlock.HauntAura.IsActive() {
			spell = warlock.Haunt
		} else if sim.IsExecutePhase20() {
			// Drain Soul execute phase
			spell = warlock.channelCheck(sim, warlock.DrainSoulDot, 5)
		} else {
			switch mainSpell {
			case proto.Warlock_Rotation_ShadowBolt:
				spell = warlock.ShadowBolt
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
