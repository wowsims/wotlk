package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var ScourgeStrikeActionID = core.ActionID{SpellID: 55271}

// this is just a simple spell because it has no rune costs and is really just a wrapper.
func (dk *Deathknight) registerScourgeStrikeShadowDamageSpell() *core.Spell {
	diseaseMulti := dk.dkDiseaseMultiplier(0.12)

	// This spell (70890) is marked as "Ignore Damage Taken Modifiers" and "Ignore Caster Damage Modifiers", but does neither.
	//  E.g. Ebon Plague affects it like a normal spell, but caster damage modifiers (Apply Aura: Mod Damage Done % (Shadow))
	//  are affecting it additively (e.g. Blood Presence (+15%), Desolation (+5%), and Black Ice (+10%) add only 30% instead of
	//  the regular ~32.9% damage).

	return dk.Unit.RegisterSpell(core.SpellConfig{
		ActionID:    ScourgeStrikeActionID.WithTag(2),
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage | core.ProcMaskProc,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreAttackerModifiers,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.LastScourgeStrikeDamage * diseaseMulti * dk.dkCountActiveDiseases(target) * dk.bonusCoeffs.scourgeStrikeShadowMultiplier
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
		},
	})
}

func (dk *Deathknight) registerScourgeStrikeSpell() {
	if !dk.Talents.ScourgeStrike {
		return
	}

	shadowDamageSpell := dk.registerScourgeStrikeShadowDamageSpell()
	bonusBaseDamage := dk.sigilOfAwarenessBonus() + dk.sigilOfArthriticBindingBonus()
	hasGlyph := dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfScourgeStrike)

	dk.ScourgeStrike = dk.RegisterSpell(core.SpellConfig{
		ActionID:    ScourgeStrikeActionID.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		RuneCost: core.RuneCostOptions{
			FrostRuneCost:  1,
			UnholyRuneCost: 1,
			RunicPowerGain: 15 + 2.5*float64(dk.Talents.Dirge) + dk.scourgeborneBattlegearRunicPowerBonus(),
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		BonusCritRating: (dk.subversionCritBonus() + dk.viciousStrikesCritChanceBonus() + dk.scourgeborneBattlegearCritBonus()) * core.CritRatingPerCritChance,

		DamageMultiplier: .7 *
			[]float64{1.0, 1.07, 1.13, 1.2}[dk.Talents.Outbreak] *
			dk.scourgelordsBattlegearDamageBonus(ScourgelordBonusSpellSS),

		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.ViciousStrikes),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 800 +
				bonusBaseDamage +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()
			baseDamage *= dk.RoRTSBonus(target)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			spell.SpendRefundableCost(sim, result)

			if result.Landed() && dk.DiseasesAreActive(target) {
				dk.LastScourgeStrikeDamage = result.Damage
				shadowDamageSpell.Cast(sim, target)

				if hasGlyph {
					// Extend FF by 3
					if dk.FrostFeverSpell.Dot(target).IsActive() && dk.FrostFeverExtended[target.Index] < 3 {
						dk.FrostFeverExtended[target.Index]++
						dk.FrostFeverSpell.Dot(target).UpdateExpires(dk.FrostFeverSpell.Dot(target).ExpiresAt() + 3*time.Second)
					}
					// Extend BP by 3
					if dk.BloodPlagueSpell.Dot(target).IsActive() && dk.BloodPlagueExtended[target.Index] < 3 {
						dk.BloodPlagueExtended[target.Index]++
						dk.BloodPlagueSpell.Dot(target).UpdateExpires(dk.BloodPlagueSpell.Dot(target).ExpiresAt() + 3*time.Second)
					}
				}
			}

			spell.DealDamage(sim, result)
		},
	})
}
