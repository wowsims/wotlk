package paladin

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) registerSealOfRighteousnessSpellAndAura() {
	/*
	 * Seal of Righteousness is an Spell/Aura that when active makes the paladin capable of procing
	 * 2 different SpellIDs depending on a paladin's casted spell or melee swing.
	 * NOTE:
	 *   Seal of Righteousness is unique in that it is the only seal that can proc off its own judgements.
	 *
	 * SpellID 20187 (Judgement of Righteousness):
	 *   - Procs off of any "Primary" Judgement (JoL, JoW, JoJ).
	 *   - Cannot miss or be dodged/parried.
	 *   - Deals hybrid AP/SP damage.
	 *   - Crits off of a melee modifier.
	 *
	 * SpellID 20154 (Seal of Righteousness):
	 *   - Procs off of any melee special ability, or white hit.
	 *   - Cannot miss or be dodged/parried.
	 *   - Deals hybrid AP/SP damage * current weapon speed.
	 *   - CANNOT CRIT.
	 */

	baseModifiers := Multiplicative{
		Additive{
			paladin.getItemSetLightswornBattlegearBonus4(),
			paladin.getMajorGlyphSealOfRighteousnessBonus(),
			paladin.getTalentSealsOfThePureBonus(),
		},
		Additive{paladin.getTalentTwoHandedWeaponSpecializationBonus()},
	}
	baseMultiplier := baseModifiers.Get()

	judgementModifiers := append(baseModifiers.Clone(),
		Additive{paladin.getMajorGlyphOfJudgementBonus(), paladin.getTalentTheArtOfWarBonus()},
	)
	judgementMultiplier := judgementModifiers.Get()

	onJudgementProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20187}, // Judgement of Righteousness.
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagSecondaryJudgement,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeOrRangedSpecial,
			DamageMultiplier: judgementMultiplier,
			ThreatMultiplier: 1,

			BonusCritRating: (6 * float64(paladin.Talents.Fanaticism) * core.CritRatingPerCritChance) +
				(core.TernaryFloat64(paladin.HasSetBonus(ItemSetTuralyonsBattlegear, 4) || paladin.HasSetBonus(ItemSetLiadrinsBattlegear, 4), 5, 0) * core.CritRatingPerCritChance),
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					// i = 1 + 0.2 * AP + 0.32 * HolP
					scaling := hybridScaling{
						AP: 0.20,
						SP: 0.32,
					}

					damage := 1 + (scaling.AP * hitEffect.MeleeAttackPower(spell.Unit)) + (scaling.SP * hitEffect.SpellPower(spell.Unit, spell))

					return damage
				},
			},
			OutcomeApplier: paladin.OutcomeFuncMeleeSpecialCritOnly(paladin.MeleeCritMultiplier()), // Secondary Judgements cannot miss if the Primary Judgement hit, only roll for crit.
		}),
	})

	onSpecialOrSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20154}, // Seal of Righteousness damage bonus.
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: baseMultiplier,
			ThreatMultiplier: 1,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					// weapon_speed * (0.022* AP + 0.044*HolP)

					scaling := hybridScaling{
						AP: 0.022,
						SP: 0.044,
					}

					damage := paladin.GetMHWeapon().ToProto().WeaponSpeed * ((scaling.AP * hitEffect.MeleeAttackPower(spell.Unit)) + (scaling.SP * hitEffect.SpellPower(spell.Unit, spell)))
					return damage
				},
			},
			OutcomeApplier: paladin.OutcomeFuncAlwaysHit(), // can't miss if attack landed
		}),
	})

	// Seal of Righteousness aura.
	auraActionID := core.ActionID{SpellID: 21084}
	paladin.SealOfRighteousnessAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Righteousness",
		Tag:      "Seal",
		ActionID: auraActionID,
		Duration: SealDuration,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			// Don't proc on misses or our own procs.
			if !spellEffect.Landed() || spell.SpellID == onJudgementProc.SpellID || spell.SpellID == onSpecialOrSwingProc.SpellID {
				return
			}

			// Differ between judgements and other melee abilities.
			if spell.Flags.Matches(SpellFlagPrimaryJudgement) {
				// SoR is the only seal that can proc off its own judgement.
				onSpecialOrSwingProc.Cast(sim, spellEffect.Target)
				onJudgementProc.Cast(sim, spellEffect.Target)
				if paladin.Talents.JudgementsOfTheJust > 0 {
					// Special JoJ talent behavior, procs swing seal on judgements
					// Yes, for SoR this means it proces TWICE on one judgement.
					onSpecialOrSwingProc.Cast(sim, spellEffect.Target)
				}
			} else {
				if spellEffect.IsMelee() {
					onSpecialOrSwingProc.Cast(sim, spellEffect.Target)
				}
			}
		},
	})

	aura := paladin.SealOfRighteousnessAura
	baseCost := paladin.BaseMana * 0.14
	paladin.SealOfRighteousness = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    auraActionID, // Seal of Righteousness self buff.
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
