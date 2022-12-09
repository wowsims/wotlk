package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// TODO: Cleanup obliterate the same way we did for plague strike
var ObliterateActionID = core.ActionID{SpellID: 51425}

func (dk *Deathknight) newObliterateHitSpell(isMH bool) *RuneSpell {
	bonusBaseDamage := dk.sigilOfAwarenessBonus()
	diseaseMulti := dk.dkDiseaseMultiplier(0.125)
	diseaseConsumptionChance := []float64{1.0, 0.67, 0.34, 0.0}[dk.Talents.Annihilation]

	rs := &RuneSpell{}
	conf := core.SpellConfig{
		ActionID:    ObliterateActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    dk.threatOfThassarianProcMask(isMH),
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		BonusCritRating: (dk.rimeCritBonus() + dk.subversionCritBonus() + dk.annihilationCritBonus() + dk.scourgeborneBattlegearCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: .8 *
			core.TernaryFloat64(isMH, 1, dk.nervesOfColdSteelBonus()) *
			core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfObliterate), 1.25, 1.0) *
			dk.scourgelordsBattlegearDamageBonus(dk.Obliterate),
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.GuileOfGorefiend),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage float64
			if isMH {
				baseDamage = 584 +
					bonusBaseDamage +
					spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			} else {
				// SpellID 66974
				baseDamage = 292 +
					bonusBaseDamage +
					spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			}
			baseDamage *= dk.RoRTSBonus(target) *
				(1.0 + dk.dkCountActiveDiseases(target)*diseaseMulti) *
				dk.mercilessCombatBonus(sim)

			result := spell.CalcDamage(sim, target, baseDamage, dk.threatOfThassarianOutcomeApplier(spell))

			if isMH {
				rs.OnResult(sim, result)
				dk.LastOutcome = result.Outcome
				dk.threatOfThassarianProc(sim, result, dk.ObliterateOhHit)

				if sim.RandomFloat("Annihilation") < diseaseConsumptionChance {
					dk.FrostFeverDisease[target.Index].Deactivate(sim)
					dk.BloodPlagueDisease[target.Index].Deactivate(sim)
				}

				if sim.RandomFloat("Rime") < dk.rimeHbChanceProc() {
					dk.RimeAura.Activate(sim)
				}
			}

			spell.DealDamage(sim, result)
		},
	}

	if isMH {
		amountOfRunicPower := 15.0 + 2.5*float64(dk.Talents.ChillOfTheGrave) + dk.scourgeborneBattlegearRunicPowerBonus()
		conf.ResourceType = stats.RunicPower
		conf.BaseCost = float64(core.NewRuneCost(uint8(amountOfRunicPower), 0, 1, 1, 0))
		conf.Cast = core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: conf.BaseCost,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.GetModifiedGCD()
			},
			IgnoreHaste: true,
		}
		rs.Refundable = true
		rs.ConvertType = RuneTypeFrost | RuneTypeUnholy
		if dk.Talents.DeathRuneMastery == 3 {
			rs.DeathConvertChance = 1.0
		} else {
			rs.DeathConvertChance = float64(dk.Talents.DeathRuneMastery) * 0.33
		}
	}

	if isMH {
		return dk.RegisterSpell(rs, conf, func(sim *core.Simulation) bool {
			return dk.CastCostPossible(sim, 0.0, 0, 1, 1) && dk.Obliterate.IsReady(sim)
		}, nil)
	} else {
		return dk.RegisterSpell(rs, conf, nil, nil)
	}
}

func (dk *Deathknight) registerObliterateSpell() {
	dk.ObliterateMhHit = dk.newObliterateHitSpell(true)
	dk.ObliterateOhHit = dk.newObliterateHitSpell(false)
	dk.Obliterate = dk.ObliterateMhHit
}

