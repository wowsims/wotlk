package paladin

import (
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
	// TODO: Test whether T8 Prot 2pc also affects Judgement, once available
	// TODO: Verify whether these bonuses should indeed be additive with similar

	dotSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 31803, Tag: 2},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagMeleeMetrics,

		DamageMultiplier: 1 *
			(1 + paladin.getItemSetLightswornBattlegearBonus4() + paladin.getItemSetAegisPlateBonus2() + paladin.getTalentSealsOfThePureBonus()),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Holy Vengeance (DoT)",
				MaxStacks: 5,
			},
			NumberOfTicks: 5,
			TickLength:    time.Second * 3, // ticking every three seconds for a grand total of 15s of duration
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				tickValue := 0 +
					.013*dot.Spell.SpellPower() +
					.025*dot.Spell.MeleeAttackPower()
				dot.SnapshotBaseDamage = tickValue * float64(dot.GetStacks())

				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
			},
		},
	})
	paladin.SovDotSpell = dotSpell

	dotApplicationSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 31803, Tag: 1},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskProc,

		DamageMultiplier: 1 *
			(1 + paladin.getItemSetLightswornBattlegearBonus4() + paladin.getItemSetAegisPlateBonus2() + paladin.getTalentSealsOfThePureBonus()),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Does no damage, just applies dot and rolls.
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)

			if result.Landed() {
				dot := dotSpell.Dot(target)
				if !dot.IsActive() {
					dot.Apply(sim)
				}
				dot.AddStack(sim)
				dot.TakeSnapshot(sim, false)
				dot.Activate(sim)
			}
		},
	})

	onJudgementProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 31804}, // Judgement of Vengeance.
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskMeleeSpecial,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagSecondaryJudgement,

		BonusCritRating: (6 * float64(paladin.Talents.Fanaticism) * core.CritRatingPerCritChance) +
			(core.TernaryFloat64(paladin.HasSetBonus(ItemSetTuralyonsBattlegear, 4), 5, 0) * core.CritRatingPerCritChance),
		DamageMultiplier: 1 *
			(1 + paladin.getItemSetLightswornBattlegearBonus4() +
				paladin.getTalentSealsOfThePureBonus() + paladin.getMajorGlyphOfJudgementBonus() + paladin.getTalentTheArtOfWarBonus()) *
			(1 + paladin.getTalentTwoHandedWeaponSpecializationBonus()),
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// i = 1 + 0.22 * HolP + 0.14 * AP
			baseDamage := 1 +
				.22*spell.SpellPower() +
				.14*spell.MeleeAttackPower()

			// i = i * (1 + (0.10 * stacks))
			baseDamage *= 1 + .1*float64(dotSpell.Dot(target).GetStacks())

			// Secondary Judgements cannot miss if the Primary Judgement hit, only roll for crit.
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	onSpecialOrSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 42463}, // Seal of Vengeance damage bonus.
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskProc, // does proc certain spell damage-based items, e.g. Black Magic, Pendulum of Telluric Currents
		Flags:       core.SpellFlagMeleeMetrics,

		// (mult * weaponScaling / stacks)
		DamageMultiplier: 1 *
			(1 + paladin.getItemSetLightswornBattlegearBonus4() + paladin.getItemSetAegisPlateBonus2() + paladin.getTalentSealsOfThePureBonus()) *
			(1 + paladin.getTalentTwoHandedWeaponSpecializationBonus()) * .33 / 5,
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := paladin.MHWeaponDamage(sim, spell.MeleeAttackPower()) *
				float64(dotSpell.Dot(target).GetStacks())

			// can't miss if melee swing landed, but can crit
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
			paladin.maybeProcVengeance(sim, result)
		},
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

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Don't proc on misses or our own procs.
			if !result.Landed() || spell == dotSpell || spell == onJudgementProc || spell == onSpecialOrSwingProc {
				return
			}

			// Differ between judgements and other melee abilities.
			if spell.Flags.Matches(SpellFlagPrimaryJudgement) {
				onJudgementProc.Cast(sim, result.Target)
				if paladin.Talents.JudgementsOfTheJust > 0 {
					// Special JoJ talent behavior, procs swing seal on judgements
					if dotSpell.Dot(result.Target).GetStacks() > 0 {
						onSpecialOrSwingProc.Cast(sim, result.Target)
					}
				}
			} else {
				if spell.IsMelee() {
					if dotSpell.Dot(result.Target).GetStacks() > 0 {
						onSpecialOrSwingProc.Cast(sim, result.Target)
					}
				}
			}

			dotApplicableSpells := []*core.Spell{
				paladin.HammerOfTheRighteous,
				paladin.CrusaderStrike,
				paladin.DivineStorm,
				paladin.HammerOfWrath,
				paladin.ShieldOfRighteousness,
			}
			isApplicableSpell := false
			for _, validSpell := range dotApplicableSpells {
				if spell == validSpell {
					isApplicableSpell = true
				}
			}

			if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) || isApplicableSpell {
				dotApplicationSpell.Cast(sim, result.Target)
			}
		},
	})

	aura := paladin.SealOfVengeanceAura
	paladin.SealOfVengeance = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    auraActionID, // Seal of Vengeance self buff.
		SpellSchool: core.SpellSchoolHoly,
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

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if paladin.CurrentSeal != nil {
				paladin.CurrentSeal.Deactivate(sim)
			}
			paladin.CurrentSeal = aura
			paladin.CurrentSeal.Activate(sim)
		},
	})
}
