package paladin

import (
	"github.com/wowsims/wotlk/sim/core"
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
	 *   - Cannot miss or be dodged/parried.
	 *   - Deals hybrid AP/SP damage * current weapon speed.
	 *   - Crits off of a melee modifier.
	 */

	baseMultiplier := 1.0
	// Additive bonuses
	baseMultiplier += core.TernaryFloat64(paladin.HasSetBonus(ItemSetLightswornBattlegear, 4), .1, 0)
	baseMultiplier *= paladin.WeaponSpecializationMultiplier()

	onJudgementProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20467}, // Judgement of Command
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagJudgement,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeOrRangedSpecial,
			DamageMultiplier: baseMultiplier,
			ThreatMultiplier: 1,

			BonusCritRating: 6 * float64(paladin.Talents.Fanaticism) * core.CritRatingPerCritChance,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					scaling := hybridScaling{
						AP: 0.08,
						SP: 0.13,
					}

					minimum := (0.19 * paladin.GetMHWeapon().WeaponDamageMin) + (scaling.AP * hitEffect.MeleeAttackPower(spell.Unit)) + (scaling.SP * hitEffect.SpellPower(spell.Unit, spell))
					maximum := (0.19 * paladin.GetMHWeapon().WeaponDamageMax) + (scaling.AP * hitEffect.MeleeAttackPower(spell.Unit)) + (scaling.SP * hitEffect.SpellPower(spell.Unit, spell))

					deltaDamage := maximum - minimum
					damage := minimum + deltaDamage*sim.RandomFloat("Damage Roll")

					return damage
				},
			},
			OutcomeApplier: paladin.OutcomeFuncMeleeSpecialCritOnly(paladin.MeleeCritMultiplier()), // Secondary Judgements cannot miss if the Primary Judgement hit, only roll for crit.
		}),
	})

	baseEffect := core.SpellEffect{
		ProcMask:         core.ProcMaskEmpty,
		DamageMultiplier: baseMultiplier,
		ThreatMultiplier: 1,
		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				minimum := (0.36 * paladin.GetMHWeapon().WeaponDamageMin)
				maximum := (0.36 * paladin.GetMHWeapon().WeaponDamageMax)

				deltaDamage := maximum - minimum
				damage := minimum + deltaDamage*sim.RandomFloat("Damage Roll")

				return damage
			},
		},
		OutcomeApplier: paladin.OutcomeFuncMeleeSpecialCritOnly(paladin.MeleeCritMultiplier()), // can't miss if melee swing landed, but can crit
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

	// Seal of Command aura.
	auraActionID := core.ActionID{SpellID: 20375}
	paladin.SealOfCommandAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Command",
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
