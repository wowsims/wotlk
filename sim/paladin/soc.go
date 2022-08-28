package paladin

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
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

	baseModifiers := Multiplicative{
		Additive{paladin.getItemSetLightswornBattlegearBonus4()},
		Additive{paladin.getTalentTwoHandedWeaponSpecializationBonus()},
	}
	baseMultiplier := baseModifiers.Get()

	judgementModifiers := append(baseModifiers.Clone(),
		Additive{paladin.getMajorGlyphOfJudgementBonus(), paladin.getTalentTheArtOfWarBonus()},
	)
	judgementMultiplier := judgementModifiers.Get()

	onJudgementProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20467}, // Judgement of Command
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagSecondaryJudgement,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeOrRangedSpecial,
			DamageMultiplier: judgementMultiplier,
			ThreatMultiplier: 1,

			BonusCritRating: (6 * float64(paladin.Talents.Fanaticism) * core.CritRatingPerCritChance) +
				(core.TernaryFloat64(paladin.HasSetBonus(ItemSetTuralyonsBattlegear, 4) || paladin.HasSetBonus(ItemSetLiadrinsBattlegear, 4), 5, 0) * core.CritRatingPerCritChance),
			BaseDamage: core.WrapBaseDamageConfig(core.BaseDamageConfigMeleeWeapon(
				core.MainHand,
				false,
				0,
				1.0,
				(0.19),
				true,
			), func(oldCalculator core.BaseDamageCalculator) core.BaseDamageCalculator {
				return func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					scaling := hybridScaling{
						AP: 0.08,
						SP: 0.13,
					}
					return oldCalculator(sim, hitEffect, spell) + (scaling.AP * hitEffect.MeleeAttackPower(spell.Unit)) + (scaling.SP * hitEffect.SpellPower(spell.Unit, spell))
				}
			}),
			OutcomeApplier: paladin.OutcomeFuncMeleeSpecialCritOnly(paladin.MeleeCritMultiplier()), // Secondary Judgements cannot miss if the Primary Judgement hit, only roll for crit.
		}),
	})

	baseEffect := core.SpellEffect{
		ProcMask:         core.ProcMaskEmpty,
		DamageMultiplier: baseMultiplier,
		ThreatMultiplier: 1,
		BaseDamage: core.BaseDamageConfigMeleeWeapon(
			core.MainHand,
			false,
			0,
			1,
			(0.36),
			true,
		),
		OutcomeApplier: paladin.OutcomeFuncMeleeSpecialHitAndCrit(paladin.MeleeCritMultiplier()),
	}

	numHits := core.MinInt32(3, paladin.Env.GetNumTargets()) // primary target + 2 others
	effects := make([]core.SpellEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		mhEffect := baseEffect
		mhEffect.Target = paladin.Env.GetTargetUnit(i)
		effects = append(effects, mhEffect)
	}

	onSpecialOrSwingActionID := core.ActionID{SpellID: 20424}
	onSpecialOrSwingProcCleave := paladin.RegisterSpell(core.SpellConfig{
		ActionID:     onSpecialOrSwingActionID, // Seal of Command damage bonus for single target spells.
		SpellSchool:  core.SpellSchoolHoly,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})

	onSpecialOrSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:     onSpecialOrSwingActionID, // Seal of Command damage bonus for cleaves.
		SpellSchool:  core.SpellSchoolHoly,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(baseEffect),
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

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if glyphManaMetrics != nil && spell.Flags.Matches(SpellFlagPrimaryJudgement) {
				paladin.AddMana(sim, glyphManaGain, glyphManaMetrics, false)
			}

			// Don't proc on misses or our own procs.
			if !spellEffect.Landed() || spell.SpellID == onJudgementProc.SpellID || spell.SpellID == onSpecialOrSwingProc.SpellID {
				return
			}

			// Differ between judgements and other melee abilities.
			if spell.Flags.Matches(SpellFlagPrimaryJudgement) {
				onJudgementProc.Cast(sim, spellEffect.Target)
				if paladin.Talents.JudgementsOfTheJust > 0 {
					// Special JoJ talent behavior, procs swing seal on judgements
					// For SoC this is a cleave.
					onSpecialOrSwingProcCleave.Cast(sim, spellEffect.Target)
				}
			} else {
				if spellEffect.IsMelee() {
					// Temporary check to avoid AOE double procing.
					if spell.SpellID == paladin.DivineStorm.SpellID {
						onSpecialOrSwingProc.Cast(sim, spellEffect.Target)
					} else {
						onSpecialOrSwingProcCleave.Cast(sim, spellEffect.Target)
					}
				}
			}
		},
	})

	aura := paladin.SealOfCommandAura
	baseCost := paladin.BaseMana * 0.14
	paladin.SealOfCommand = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    auraActionID, // Seal of Command self buff.
		SpellSchool: core.SpellSchoolHoly,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(paladin.Talents.Benediction)),
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if paladin.CurrentSeal != nil {
				paladin.CurrentSeal.Deactivate(sim)
			}
			paladin.CurrentSeal = aura
			paladin.CurrentSeal.Activate(sim)
		},
	})
}
