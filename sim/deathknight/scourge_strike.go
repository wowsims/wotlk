package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var ScourgeStrikeActionID = core.ActionID{SpellID: 55271}

// this is just a simple spell because it has no rune costs and is really just a wrapper.
func (dk *Deathknight) registerScourgeStrikeShadowDamageSpell() *core.Spell {
	diseaseMulti := dk.dkDiseaseMultiplier(0.12)

	return dk.Unit.RegisterSpell(core.SpellConfig{
		ActionID:    ScourgeStrikeActionID.WithTag(2),
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagIgnoreResists | core.SpellFlagMeleeMetrics,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			BonusCritRating:  -100 * core.CritRatingPerCritChance, // Disable criticals for shadow hit
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			OutcomeApplier: dk.CurrentTarget.OutcomeFuncAlwaysHit(),

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return dk.LastScourgeStrikeDamage * (diseaseMulti * dk.dkCountActiveDiseases(hitEffect.Target))
				},
			},
		}),
	})
}

func (dk *Deathknight) registerScourgeStrikeSpell() {

	shadowDamageSpell := dk.registerScourgeStrikeShadowDamageSpell()
	bonusBaseDamage := dk.sigilOfAwarenessBonus(dk.ScourgeStrike) + dk.sigilOfArthriticBindingBonus()
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 800.0+bonusBaseDamage, 1.0, 0.7, true)
	outbreakBonus := []float64{1.0, 1.07, 1.13, 1.2}[dk.Talents.Outbreak]
	rpGain := 15.0 + 2.5*float64(dk.Talents.Dirge) + dk.scourgeborneBattlegearRunicPowerBonus()
	hasGlyph := dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfScourgeStrike)

	baseCost := float64(core.NewRuneCost(uint8(rpGain), 0, 1, 1, 0))
	rs := &RuneSpell{}
	dk.ScourgeStrike = dk.RegisterSpell(rs, core.SpellConfig{
		ActionID:     ScourgeStrikeActionID.WithTag(1),
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
			IgnoreHaste: true,
		},

		ApplyEffects: dk.withRuneRefund(rs, core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			BonusCritRating:  (dk.subversionCritBonus() + dk.viciousStrikesCritChanceBonus() + dk.scourgeborneBattlegearCritBonus()) * core.CritRatingPerCritChance,
			DamageMultiplier: outbreakBonus * dk.scourgelordsBattlegearDamageBonus(dk.ScourgeStrike),
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return weaponBaseDamage(sim, hitEffect, spell) * dk.RoRTSBonus(hitEffect.Target)
				},
				TargetSpellCoefficient: 1,
			},

			OutcomeApplier: dk.OutcomeFuncMeleeSpecialHitAndCrit(dk.MeleeCritMultiplier(1.0, dk.viciousStrikesCritDamageBonus())),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				dk.LastOutcome = spellEffect.Outcome
				if spellEffect.Landed() && dk.DiseasesAreActive(spellEffect.Target) {
					dk.LastScourgeStrikeDamage = spellEffect.Damage
					shadowDamageSpell.Cast(sim, spellEffect.Target)

					if hasGlyph {
						// Extend FF by 3
						if dk.FrostFeverDisease[spellEffect.Target.Index].IsActive() && dk.FrostFeverExtended[spellEffect.Target.Index] < 3 {
							dk.FrostFeverExtended[spellEffect.Target.Index]++
							dk.FrostFeverDisease[spellEffect.Target.Index].UpdateExpires(dk.FrostFeverDisease[spellEffect.Target.Index].ExpiresAt() + 3*time.Second)
						}
						// Extend BP by 3
						if dk.BloodPlagueDisease[spellEffect.Target.Index].IsActive() && dk.BloodPlagueExtended[spellEffect.Target.Index] < 3 {
							dk.BloodPlagueExtended[spellEffect.Target.Index]++
							dk.BloodPlagueDisease[spellEffect.Target.Index].UpdateExpires(dk.BloodPlagueDisease[spellEffect.Target.Index].ExpiresAt() + 3*time.Second)
						}
					}
				}
			},
		}, false),
	}, func(sim *core.Simulation) bool {
		return dk.Talents.ScourgeStrike && dk.CastCostPossible(sim, 0.0, 0, 1, 1) && dk.ScourgeStrike.IsReady(sim)
	}, nil)
}
