package paladin

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (paladin *Paladin) registerSealOfCommandSpellAndAura() {
	/*
	 * Seal of Command is an Spell/Aura that when active makes the paladin capable of procing
	 * 2 different SpellIDs depending on a paladin's casted spell or melee swing.
	 *
	 * SpellID 20467 (Judgement of Command):
	 *   - Procs off of any "Primary" Judgement (JoL, JoW, JoJ).
	 *   - Cannot miss or be dodged/parried.
	 *   - Deals hybrid AP/SP damage.
	 *   - Crits off of a melee modifier.
	 *
	 * SpellID 20424 (Seal of Command):
	 *   - Procs off of any melee special ability, or white hit.
	 *   - If the ability is SINGLE TARGET, it hits up to 2 extra targets.
	 *   - Deals hybrid AP/SP damage * current weapon speed.
	 *   - Crits off of a melee modifier.
	 *   - CAN MISS, BE DODGED/PARRIED/BLOCKED.
	 */

	numHits := min(3, paladin.Env.GetNumTargets()) // primary target + 2 others
	results := make([]*core.SpellResult, numHits)

	onJudgementProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20467}, // Judgement of Command
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskMeleeSpecial,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagSecondaryJudgement,

		BonusCritRating: (6 * float64(paladin.Talents.Fanaticism) * core.CritRatingPerCritChance) +
			(core.TernaryFloat64(paladin.HasSetBonus(ItemSetTuralyonsBattlegear, 4), 5, 0) * core.CritRatingPerCritChance),

		DamageMultiplier: 1 *
			(1 + paladin.getItemSetLightswornBattlegearBonus4() +
				paladin.getMajorGlyphOfJudgementBonus() + paladin.getTalentTheArtOfWarBonus()) *
			(1 + paladin.getTalentTwoHandedWeaponSpecializationBonus()),
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			mhWeaponDamage := 0 +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()
			baseDamage := 0.19*mhWeaponDamage +
				0.08*spell.MeleeAttackPower() +
				0.13*spell.SpellPower()

			// Secondary Judgements cannot miss if the Primary Judgement hit, only roll for crit.
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	onSpecialOrSwingActionID := core.ActionID{SpellID: 20424}
	onSpecialOrSwingProcCleave := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    onSpecialOrSwingActionID,
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics,

		DamageMultiplier: 1 *
			(1 + paladin.getItemSetLightswornBattlegearBonus4()) *
			(1 + paladin.getTalentTwoHandedWeaponSpecializationBonus()) *
			0.36, // Only 36% of weapon damage.
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := 0 +
					spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()

				results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			curTarget = target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				spell.DealDamage(sim, results[hitIndex])
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
		},
	})

	onSpecialOrSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    onSpecialOrSwingActionID,
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty, // unlike SoV, SoC crits don't proc Vengeance
		Flags:       core.SpellFlagMeleeMetrics,

		DamageMultiplier: 1 *
			(1 + paladin.getItemSetLightswornBattlegearBonus4()) *
			(1 + paladin.getTalentTwoHandedWeaponSpecializationBonus()) *
			0.36, // Only 36% of weapon damage.
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0 +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})

	var glyphManaMetrics *core.ResourceMetrics
	glyphManaGain := .08 * paladin.BaseMana
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfSealOfCommand) {
		glyphManaMetrics = paladin.NewManaMetrics(core.ActionID{ItemID: 41094})
	}

	// Seal of Command aura.
	auraActionID := core.ActionID{SpellID: 20375}
	paladin.SealOfCommandAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Command",
		Tag:      "Seal",
		ActionID: auraActionID,
		Duration: SealDuration,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if glyphManaMetrics != nil && spell.Flags.Matches(SpellFlagPrimaryJudgement) {
				paladin.AddMana(sim, glyphManaGain, glyphManaMetrics)
			}

			// Don't proc on misses or our own procs.
			if !result.Landed() || spell == onJudgementProc || spell.SameAction(onSpecialOrSwingActionID) {
				return
			}

			// Differ between judgements and other melee abilities.
			if spell.Flags.Matches(SpellFlagPrimaryJudgement) {
				onJudgementProc.Cast(sim, result.Target)
				if paladin.Talents.JudgementsOfTheJust > 0 {
					// Special JoJ talent behavior, procs swing seal on judgements
					// For SoC this is a cleave.
					onSpecialOrSwingProcCleave.Cast(sim, result.Target)
				}
			} else if spell.IsMelee() {
				// Temporary check to avoid AOE double procing.
				if spell.SpellID == paladin.HammerOfTheRighteous.SpellID || spell.SpellID == paladin.DivineStorm.SpellID {
					onSpecialOrSwingProc.Cast(sim, result.Target)
				} else {
					onSpecialOrSwingProcCleave.Cast(sim, result.Target)
				}
			}
		},
	})

	aura := paladin.SealOfCommandAura
	paladin.SealOfCommand = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    auraActionID, // Seal of Command self buff.
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.14,
			Multiplier: 1 - 0.02*float64(paladin.Talents.Benediction),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if paladin.CurrentSeal != nil {
				paladin.CurrentSeal.Deactivate(sim)
			}
			paladin.CurrentSeal = aura
			paladin.CurrentSeal.Activate(sim)
		},
	})
}
