package warlock

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

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
	staticAdditiveDamageMultiplier := 1.0

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
		staticAdditiveDamageMultiplier += 0.03 * float64(warlock.Talents.Emberstorm)
	}

	if actionID == actionID_CurseOfAgony || actionID == actionID_Corruption || actionID == actionID_Seed {
		staticAdditiveDamageMultiplier += 0.01 * float64(warlock.Talents.Contagion)
	}

	if warlock.Talents.SiphonLife && (actionID == actionID_UnstableAffliction || actionID == actionID_Corruption || (actionID == actionID_Seed && IsPeriodic)) {
		staticAdditiveDamageMultiplier += 0.05
	}

	if actionID == actionID_ShadowBolt {
		staticAdditiveDamageMultiplier += 0.02 * float64(warlock.Talents.ImprovedShadowBolt)
	}

	if actionID == actionID_Incinerate && warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfIncinerate) {
		staticAdditiveDamageMultiplier += 0.05
	}

	if actionID == actionID_CurseOfAgony {
		staticAdditiveDamageMultiplier += 0.05 * float64(warlock.Talents.ImprovedCurseOfAgony)
	}

	if actionID == actionID_Corruption {
		staticAdditiveDamageMultiplier += 0.02 * float64(warlock.Talents.ImprovedCorruption)
	}

	if actionID == actionID_Immolate || actionID == actionID_Conflagrate {
		staticAdditiveDamageMultiplier += 0.1 * float64(warlock.Talents.ImprovedImmolate)
	}

	if (actionID == actionID_Immolate && IsPeriodic) || actionID == actionID_Conflagrate {
		staticAdditiveDamageMultiplier += 0.03 * float64(warlock.Talents.Aftermath)
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
