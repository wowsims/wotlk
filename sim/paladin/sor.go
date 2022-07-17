package paladin

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) setupSealOfRighteousness() {
	/*
	 * Seal of Righteousness is an Spell/Aura that when active makes the paladin capable of procing
	 * 2 different SpellIDs depending on a paladin's casted spell or melee swing.
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
	 *   - Crits off of a melee modifier.
	 */

	baseMultiplier := 1.0
	// Additive bonuses
	baseMultiplier += core.TernaryFloat64(paladin.HasSetBonus(ItemSetLightswornBattlegear, 4), .1, 0)
	baseMultiplier += core.TernaryFloat64(paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfSealOfRighteousness), .1, 0)
	baseMultiplier += 0.03 * float64(paladin.Talents.SealsOfThePure)
	baseMultiplier *= paladin.WeaponSpecializationMultiplier()

	onJudgementProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20187}, // Judgement of Righteousness.
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeOrRangedSpecial,
			DamageMultiplier: baseMultiplier,
			ThreatMultiplier: 1,

			BonusCritRating: 6 * float64(paladin.Talents.Fanaticism) * core.CritRatingPerCritChance,
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
			OutcomeApplier: paladin.OutcomeFuncMeleeSpecialCritOnly(paladin.MeleeCritMultiplier()), // can't miss if melee swing landed, but can crit
		}),
	})

	// Seal of Righteousness aura.
	auraActionID := core.ActionID{SpellID: 21084}
	aura := paladin.RegisterAura(core.Aura{
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
			if spell.Flags.Matches(SpellFlagJudgement) {
				onJudgementProc.Cast(sim, spellEffect.Target)
			} else {
				if spellEffect.IsMelee() {
					onSpecialOrSwingProc.Cast(sim, spellEffect.Target)
				}
			}
		},
	})

	baseCost := paladin.BaseMana * 0.14
	paladin.SealOfRighteousness = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    auraActionID, // Seal of Righteousness self buff.
		SpellSchool: core.SpellSchoolHoly,
		Flags:       SpellFlagSeal,

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
