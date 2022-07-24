package warlock

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warlock *Warlock) dynamicMultiplier(sim *core.Simulation, spell *core.Spell, _ *core.SpellEffect) float64 {
	dynamicMultiplier:= 1.0

	// Execute Multipliers
	if sim.IsExecutePhase20() && spell == warlock.DrainSoul {
		dynamicMultiplier *= 4.0
	}
	if sim.IsExecutePhase35() && spell.SpellSchool == core.SpellSchoolShadow {
		dynamicMultiplier += 0.04*float64(warlock.Talents.DeathsEmbrace)
	}

	// Normal Multipliers
	if spell == warlock.DrainSoul {
		afflictionSpellNumber := core.TernaryFloat64(warlock.DrainSoulDot.IsActive(), 1, 0) + //core.TernaryFloat64(warlock.ConflagrateDot.IsActive(), 1, 0) +
			core.TernaryFloat64(warlock.CorruptionDot.IsActive(), 1, 0) + //core.TernaryFloat64(warlock.SeedDots.IsActive(), 1, 0) +
			core.TernaryFloat64(warlock.CurseOfDoomDot.IsActive(), 1, 0) + core.TernaryFloat64(warlock.CurseOfAgonyDot.IsActive(), 1, 0) +
			core.TernaryFloat64(warlock.UnstableAffDot.IsActive(), 1, 0) + core.TernaryFloat64(warlock.ImmolateDot.IsActive(), 1, 0)
		dynamicMultiplier *= 1 + 0.03*float64(warlock.Talents.SoulSiphon) * core.MinFloat(3, afflictionSpellNumber)
	}
	return dynamicMultiplier
}

func (warlock *Warlock) staticAdditiveDamageMultiplier(actionID core.ActionID, spellSchool core.SpellSchool, IsPeriodic bool) float64 {
	// actionID spellbook
	actionID_Incinerate := core.ActionID{SpellID: 47838}
	actionID_ShadowBolt := core.ActionID{SpellID: 47809}
	actionID_UnstableAffliction := core.ActionID{SpellID: 47843}
	actionID_Immolate := core.ActionID{SpellID: 47811}
	actionID_Conflagrate := core.ActionID{SpellID: 17962}
	actionID_CurseOfAgony := core.ActionID{SpellID: 47864}
	actionID_Corruption := core.ActionID{SpellID: 47813}
	actionID_Seed := core.ActionID{SpellID: 47836}
	// actionID_SoulFire := core.ActionID{SpellID: 47825}
	// actionID_ChaosBolt := core.ActionID{SpellID: 59172}
	// actionID_CurseOfDoom := core.ActionID{SpellID: 47867}
	// actionID_Haunt := core.ActionID{SpellID: 59164}
	// actionID_DrainSoul := core.ActionID{SpellID: 47855}

	// Aura bonuses are treated separately as they function like normal multipliers
	staticAdditiveDamageMultiplier:= 1.0

	// Additive Multipliers
	// Weapon Imbues
	if (IsPeriodic && warlock.Options.WeaponImbue == proto.Warlock_Options_GrandSpellstone) ||
		(!IsPeriodic && warlock.Options.WeaponImbue == proto.Warlock_Options_GrandFirestone) {
		staticAdditiveDamageMultiplier += 0.01
	}

	// Talent & Glyphs Bonuses
	if spellSchool == core.SpellSchoolShadow {
		staticAdditiveDamageMultiplier += 0.03 * float64(warlock.Talents.ShadowMastery)
	} else if spellSchool == core.SpellSchoolFire {
		staticAdditiveDamageMultiplier += 0.03*float64(warlock.Talents.Emberstorm)
	}

	if actionID == actionID_CurseOfAgony || actionID == actionID_Corruption || actionID == actionID_Seed {
		staticAdditiveDamageMultiplier += 0.01*float64(warlock.Talents.Contagion)
	}

	if warlock.Talents.SiphonLife && (actionID == actionID_UnstableAffliction || actionID == actionID_Corruption || actionID == actionID_Seed) {
		staticAdditiveDamageMultiplier += 0.05
	}

	if actionID == actionID_ShadowBolt {
		staticAdditiveDamageMultiplier += 0.05*float64(warlock.Talents.ImprovedShadowBolt)
	}

	if actionID == actionID_Incinerate && warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfIncinerate) {
		staticAdditiveDamageMultiplier += 0.05
	}

	if actionID == actionID_CurseOfAgony {
		staticAdditiveDamageMultiplier += 0.05*float64(warlock.Talents.ImprovedCurseOfAgony)
	}

	if actionID == actionID_Corruption {
		staticAdditiveDamageMultiplier += 0.02*float64(warlock.Talents.ImprovedCorruption)
	}

	if actionID == actionID_Immolate || actionID == actionID_Conflagrate {
		staticAdditiveDamageMultiplier += 0.1 * float64(warlock.Talents.ImprovedImmolate)
	}

	if (actionID == actionID_Immolate && IsPeriodic) || actionID == actionID_Conflagrate {
		staticAdditiveDamageMultiplier += 0.03*float64(warlock.Talents.Aftermath)
		if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfImmolate) {
			staticAdditiveDamageMultiplier += 0.1
		}
	}

	//  Tier Sets Bonuses (additive)
	if warlock.HasSetBonus(ItemSetMaleficRaiment, 4) {
		if actionID == actionID_ShadowBolt || actionID == actionID_Incinerate {
			staticAdditiveDamageMultiplier += 0.06
		}
	}

	if warlock.HasSetBonus(ItemSetDeathbringerGarb, 2) {
		if actionID == actionID_Immolate {
			staticAdditiveDamageMultiplier += 0.1
		}
		if actionID == actionID_UnstableAffliction {
			staticAdditiveDamageMultiplier += 0.2
		}
	}

	if warlock.HasSetBonus(ItemSetGuldansRegalia, 4) {
		if actionID == actionID_Immolate || actionID == actionID_Corruption || actionID == actionID_UnstableAffliction {
			staticAdditiveDamageMultiplier += 0.1
		}
	}

	return staticAdditiveDamageMultiplier
}

// func (warlock *Warlock) spellDamageMultiplierHelper(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) float64 {
// 	// Aura bonuses are treated separately as they function like normal multipliers
// 	additiveDamageMultiplier:= 1.0
// 	executeDamageMultiplier:= 1.0
// 	normalMultiplier:= 1.0

// 	// Additive Multipliers
// 	// Weapon Imbues
// 	if (spellEffect.IsPeriodic && warlock.Options.WeaponImbue == proto.Warlock_Options_GrandSpellstone) ||
// 		(!spellEffect.IsPeriodic && warlock.Options.WeaponImbue == proto.Warlock_Options_GrandFirestone) {
// 		additiveDamageMultiplier += 0.01
// 	}

// 	// Talent & Glyphs Bonuses
// 	if spell.SpellSchool == core.SpellSchoolShadow {
// 		additiveDamageMultiplier += 0.03 * float64(warlock.Talents.ShadowMastery)
// 	} else if spell.SpellSchool == core.SpellSchoolFire {
// 		additiveDamageMultiplier += 0.03*float64(warlock.Talents.Emberstorm)
// 	}

// 	if spell == warlock.CurseOfAgony || spell == warlock.Corruption || spell == warlock.Seeds[0] {
// 		additiveDamageMultiplier += 0.01*float64(warlock.Talents.Contagion)
// 	}

// 	if warlock.Talents.SiphonLife && (spell == warlock.UnstableAff || spell == warlock.Corruption || spell == warlock.Seeds[0]) {
// 		additiveDamageMultiplier += 0.05
// 	}

// 	if spell == warlock.ShadowBolt {
// 		additiveDamageMultiplier += 0.05*float64(warlock.Talents.ImprovedShadowBolt)
// 	}

// 	if spell == warlock.Incinerate && warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfIncinerate) {
// 		additiveDamageMultiplier += 0.05
// 	}

// 	if (spell == warlock.Incinerate || spell == warlock.ChaosBolt) && warlock.ImmolateDot.IsActive() {
// 		additiveDamageMultiplier += 0.02*float64(warlock.Talents.FireAndBrimstone)
// 	}

// 	if spell == warlock.CurseOfAgony {
// 		additiveDamageMultiplier += 0.05*float64(warlock.Talents.ImprovedCurseOfAgony)
// 	}

// 	if spell == warlock.Corruption {
// 		additiveDamageMultiplier += 0.02*float64(warlock.Talents.ImprovedCorruption)
// 	}

// 	if spell == warlock.Immolate || spell == warlock.Conflagrate {
// 		additiveDamageMultiplier += 0.1 * float64(warlock.Talents.ImprovedImmolate)
// 	}

// 	if (spell == warlock.Immolate && spellEffect.IsPeriodic) || spell == warlock.Conflagrate {
// 		additiveDamageMultiplier += 0.03*float64(warlock.Talents.Aftermath)
// 		if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfImmolate) {
// 			additiveDamageMultiplier += 0.1
// 		}
// 	}

// 	//  Tier Bonuses (additive)
// 	if warlock.HasSetBonus(ItemSetMaleficRaiment, 4) {
// 		if spell == warlock.ShadowBolt || spell == warlock.Incinerate {
// 			additiveDamageMultiplier += 0.06
// 		}
// 	}

// 	if warlock.HasSetBonus(ItemSetDeathbringerGarb, 2) {
// 		if spell == warlock.Immolate {
// 			additiveDamageMultiplier += 0.1
// 		}
// 		if spell == warlock.UnstableAff {
// 			additiveDamageMultiplier += 0.2
// 		}
// 	}

// 	if warlock.HasSetBonus(ItemSetGuldansRegalia, 4) {
// 		if spell == warlock.Immolate || spell == warlock.Corruption || spell == warlock.UnstableAff {
// 			additiveDamageMultiplier += 0.1
// 		}
// 	}

// 	// Execute Multipliers
// 	if sim.IsExecutePhase20() && spell == warlock.DrainSoul {
// 		executeDamageMultiplier *= 4.0
// 	}
// 	if sim.IsExecutePhase35() && spell.SpellSchool == core.SpellSchoolShadow {
// 		executeDamageMultiplier += 0.04*float64(warlock.Talents.DeathsEmbrace)
// 	}

// 	// Normal Multipliers
// 	if spell == warlock.DrainSoul {
// 		afflictionSpellNumber := core.TernaryFloat64(warlock.DrainSoulDot.IsActive(), 1, 0) + //core.TernaryFloat64(warlock.ConflagrateDot.IsActive(), 1, 0) +
// 			core.TernaryFloat64(warlock.CorruptionDot.IsActive(), 1, 0) + //core.TernaryFloat64(warlock.SeedDots.IsActive(), 1, 0) +
// 			core.TernaryFloat64(warlock.CurseOfDoomDot.IsActive(), 1, 0) + core.TernaryFloat64(warlock.CurseOfAgonyDot.IsActive(), 1, 0) +
// 			core.TernaryFloat64(warlock.UnstableAffDot.IsActive(), 1, 0) + core.TernaryFloat64(warlock.ImmolateDot.IsActive(), 1, 0)
// 		normalMultiplier *= 1 + 0.03*float64(warlock.Talents.SoulSiphon) * core.MinFloat(3, afflictionSpellNumber)
// 	}
	
// 	return additiveDamageMultiplier * executeDamageMultiplier * normalMultiplier
// }
