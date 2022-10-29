package warlock

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warlock *Warlock) staticAdditiveDamageMultiplier(actionID core.ActionID, spellSchool core.SpellSchool, isPeriodic bool) float64 {
	// actionID spellbook
	actionID_ShadowBolt := core.ActionID{SpellID: 47809}
	actionID_Corruption := core.ActionID{SpellID: 47813}
	actionID_Seed := core.ActionID{SpellID: 47836}
	actionID_UnstableAffliction := core.ActionID{SpellID: 47843}
	actionID_Incinerate := core.ActionID{SpellID: 47838}
	actionID_Immolate := core.ActionID{SpellID: 47811}
	actionID_Conflagrate := core.ActionID{SpellID: 17962}
	actionID_CurseOfAgony := core.ActionID{SpellID: 47864}
	actionID_CurseOfDoom := core.ActionID{SpellID: 47867}
	// actionID_SoulFire := core.ActionID{SpellID: 47825}
	// actionID_ChaosBolt := core.ActionID{SpellID: 59172}
	// actionID_Haunt := core.ActionID{SpellID: 59164}
	// actionID_DrainSoul := core.ActionID{SpellID: 47855}

	// Aura bonuses are treated separately as they function like normal multipliers
	additiveDamageMultiplier := 1.

	// Additive Multipliers
	// Weapon Imbues
	if (isPeriodic && warlock.Options.WeaponImbue == proto.Warlock_Options_GrandSpellstone && !(actionID == actionID_CurseOfAgony) && !(actionID == actionID_CurseOfDoom)) ||
		(!isPeriodic && warlock.Options.WeaponImbue == proto.Warlock_Options_GrandFirestone) {
		additiveDamageMultiplier += 0.01
	}

	// Talent & Glyphs Bonuses
	if spellSchool == core.SpellSchoolShadow {
		additiveDamageMultiplier += 0.03 * float64(warlock.Talents.ShadowMastery)
	} else if spellSchool == core.SpellSchoolFire {
		additiveDamageMultiplier += 0.03 * float64(warlock.Talents.Emberstorm)
	}

	if actionID == actionID_CurseOfAgony || actionID == actionID_Corruption || actionID == actionID_Seed {
		additiveDamageMultiplier += 0.01 * float64(warlock.Talents.Contagion)
	}

	if warlock.Talents.SiphonLife && (actionID == actionID_UnstableAffliction || actionID == actionID_Corruption || (actionID == actionID_Seed && isPeriodic)) {
		additiveDamageMultiplier += 0.05
	}

	if actionID == actionID_ShadowBolt {
		additiveDamageMultiplier += 0.02 * float64(warlock.Talents.ImprovedShadowBolt)
	}

	if actionID == actionID_Incinerate && warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfIncinerate) {
		additiveDamageMultiplier += 0.05
	}

	if actionID == actionID_CurseOfAgony {
		additiveDamageMultiplier += 0.05 * float64(warlock.Talents.ImprovedCurseOfAgony)
	}

	if actionID == actionID_Corruption {
		additiveDamageMultiplier += 0.02 * float64(warlock.Talents.ImprovedCorruption)
	}

	if actionID.SameActionIgnoreTag(actionID_Immolate) || actionID.SameActionIgnoreTag(actionID_Conflagrate) {
		additiveDamageMultiplier += 0.1 * float64(warlock.Talents.ImprovedImmolate)
	}

	if (actionID.SameActionIgnoreTag(actionID_Immolate) && isPeriodic) || actionID.SameActionIgnoreTag(actionID_Conflagrate) {
		additiveDamageMultiplier += 0.03 * float64(warlock.Talents.Aftermath)
		if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfImmolate) {
			additiveDamageMultiplier += 0.1
		}
	}

	//  Tier Sets Bonuses (additive)
	if warlock.HasSetBonus(ItemSetMaleficRaiment, 4) {
		if actionID == actionID_ShadowBolt || actionID == actionID_Incinerate {
			additiveDamageMultiplier += 0.06
		}
	}

	if warlock.HasSetBonus(ItemSetDeathbringerGarb, 2) {
		if actionID.SameActionIgnoreTag(actionID_Immolate) {
			additiveDamageMultiplier += 0.1
		}
		if actionID == actionID_UnstableAffliction {
			additiveDamageMultiplier += 0.2
		}
	}

	if warlock.HasSetBonus(ItemSetGuldansRegalia, 4) {
		if actionID.SameActionIgnoreTag(actionID_Immolate) || actionID == actionID_Corruption || actionID == actionID_UnstableAffliction {
			additiveDamageMultiplier += 0.1
		}
	}

	return additiveDamageMultiplier
}
