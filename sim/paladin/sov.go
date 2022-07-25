package paladin

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) registerSealOfVengeanceSpellAndAura() {
	/*
	 * Seal of Vengeance is an Spell/Aura that when active makes the paladin capable of procing
	 * 3 different SpellIDs depending on a paladin's casted spell or melee swing.
	 *
	 * SpellID 31803 (Holy Vengeance):
	 * 	 - "Hidden" proc that does a second melee roll on white hit to apply a DoT of
	 *     the same SpellID.
	 *   - Since this is a second roll, it can miss or be dodged/parried.
	 *   - Does no damage on its own, only the DoT does damage, DoT scales based on AP/SP.
	 *   - The DoT applied by this modifies all other procs.
	 *   - Cannot crit by default.
	 *
	 * SpellID 31804 (Judgement of Vengeance):
	 *   - Procs off of any "Primary" Judgement (JoL, JoW, JoJ).
	 *   - Cannot miss or be dodged/parried.
	 *   - Deals hybrid AP/SP damage, increased by 10% per stack of Holy Vengeance.
	 *   - Crits off of a melee modifier.
	 *
	 * SpellID 42463 (Seal of Vengeance):
	 *   - Procs off of any melee special ability, or white hit.
	 *   - Cannot miss or be dodged/parried.
	 *   - Deals % weapon damage, only after reaching 1 stack, increased by ~7% per stack of Holy Vengeance for a total of ~33%.
	 *   - Crits off of a melee modifier.
	 *
	 * TODO:
	 *  - Add set bonus and talent related modifiers.
	 *  - Fix expertise rating on glyph application.
	 */

	baseModifiers := Multiplicative{
		Additive{
			paladin.getItemSetLightswornBattlegearBonus4(),
			paladin.getTalentSealsOfThePureBonus(),
		},
		Additive{paladin.getTalentTwoHandedWeaponSpecializationBonus()},
	}
	baseMultiplier := baseModifiers.Get()

	judgementModifiers := append(baseModifiers.Clone(),
		Additive{paladin.getMajorGlyphOfJudgementBonus()},
	)
	judgementMultiplier := judgementModifiers.Get()

	dot := paladin.createSealOfVengeanceDot(baseMultiplier)
	paladin.SealOfVengeanceDot = dot

	onSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 31803, Tag: 1}, // Holy Vengeance.
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskEmpty, // Might need to be changed later if SOV secondary rolls can proc other things.
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					if !dot.IsActive() {
						dot.Apply(sim)
					}
					dot.AddStack(sim)
					dot.Reapply(sim)
					return 0 // Does no damage, just applies dot and rolls.
				},
			},
			OutcomeApplier: paladin.OutcomeFuncMeleeSpecialHit(),
		}),
	})

	onJudgementProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 31804}, // Judgement of Vengeance.
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
					// i = 1 + 0.22 * HolP + 0.14 * AP
					scaling := hybridScaling{
						AP: 0.14,
						SP: 0.22,
					}

					damage := 1 + (scaling.SP * hitEffect.SpellPower(spell.Unit, spell)) + (scaling.AP * hitEffect.MeleeAttackPower(spell.Unit))

					// i = i * (1 + (0.10 * stacks))
					stackMultiplier := 0.10
					stacksBonus := 1 + (float64(stackMultiplier) * float64(dot.GetStacks()))
					damage *= stacksBonus

					return damage
				},
			},
			OutcomeApplier: paladin.OutcomeFuncMeleeSpecialCritOnly(paladin.MeleeCritMultiplier()), // Secondary Judgements cannot miss if the Primary Judgement hit, only roll for crit.
		}),
	})

	weaponScaling := 0.33
	damagePerStack := weaponScaling / 5
	onSpecialOrSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 42463}, // Seal of Vengeance damage bonus.
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: baseMultiplier,
			ThreatMultiplier: 1,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					adjustedWeaponDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, false, 0, damagePerStack, false)
					damage := adjustedWeaponDamage(sim, hitEffect, spell) * float64(dot.GetStacks())
					return damage
				},
			},
			OutcomeApplier: paladin.OutcomeFuncMeleeSpecialCritOnly(paladin.MeleeCritMultiplier()), // can't miss if melee swing landed, but can crit
		}),
	})

	// Seal of Vengeance aura.
	auraActionID := core.ActionID{SpellID: 31801}
	paladin.SealOfVengeanceAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Vengeance",
		Tag:      "Seal",
		ActionID: auraActionID,
		Duration: SealDuration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfSealOfVengeance) {
				expertise := core.ExpertisePerQuarterPercentReduction * 10
				paladin.AddStatDynamic(sim, stats.Expertise, expertise)
			}
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfSealOfVengeance) {
				expertise := core.ExpertisePerQuarterPercentReduction * 10
				paladin.AddStatDynamic(sim, stats.Expertise, -expertise)
			}
		},

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			// Don't proc on misses or our own procs.
			if !spellEffect.Landed() || spell.SpellID == onSwingProc.SpellID || spell.SpellID == onJudgementProc.SpellID || spell.SpellID == onSpecialOrSwingProc.SpellID {
				return
			}

			// Only white hits can trigger this. (SoV dot)
			if spellEffect.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				onSwingProc.Cast(sim, spellEffect.Target)
			}

			// Differ between judgements and other melee abilities.
			if spell.Flags.Matches(SpellFlagPrimaryJudgement) {
				onJudgementProc.Cast(sim, spellEffect.Target)
			} else {
				if spellEffect.IsMelee() {
					if dot.GetStacks() > 0 && spellEffect.Target == paladin.CurrentTarget {
						onSpecialOrSwingProc.Cast(sim, spellEffect.Target)
					}
				}
			}
		},
	})

	aura := paladin.SealOfRighteousnessAura
	baseCost := paladin.BaseMana * 0.14
	paladin.SealOfVengeance = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    auraActionID, // Seal of Vengeance self buff.
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

func (paladin *Paladin) createSealOfVengeanceDot(multiplier float64) *core.Dot {
	target := paladin.CurrentTarget

	dotActionID := core.ActionID{SpellID: 31803, Tag: 2} // Holy Vengeance
	dotSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    dotActionID,
		SpellSchool: core.SpellSchoolHoly,
	})
	dotAura := target.RegisterAura(core.Aura{
		Label:     "Holy Vengeance (DoT) -" + strconv.Itoa(int(paladin.Index)),
		ActionID:  dotActionID,
		MaxStacks: 5,
	})
	dotScaling := hybridScaling{
		AP: 0.025,
		SP: 0.013,
	}
	return core.NewDot(core.Dot{
		Spell: dotSpell,
		Aura:  dotAura,

		NumberOfTicks: 5,
		TickLength:    time.Second * 3, // ticking every three seconds for a grand total of 15s of duration

		TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			IsPeriodic:       true,
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: multiplier,
			BaseDamage: core.MultiplyByStacks(core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					// (0.013 * HolP + 0.025 * AP)
					tickValue := (dotScaling.SP * hitEffect.SpellPower(spell.Unit, spell)) + (dotScaling.AP * hitEffect.MeleeAttackPower(spell.Unit))
					return tickValue
				},
				TargetSpellCoefficient: 1,
			}, dotAura),
			OutcomeApplier: paladin.OutcomeFuncAlwaysHit(),
			// OutcomeApplier: paladin.OutcomeFuncTick(),
		})),
	})
}
