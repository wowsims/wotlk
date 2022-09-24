package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// TODO: Cleanup obliterate the same way we did for plague strike
var ObliterateActionID = core.ActionID{SpellID: 51425}

func (dk *Deathknight) newObliterateHitSpell(isMH bool, onhit func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect)) *RuneSpell {
	bonusBaseDamage := dk.sigilOfAwarenessBonus()
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 584.0+bonusBaseDamage, true)
	if !isMH {
		// SpellID 66974
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, true, 292.0+bonusBaseDamage, true)
	}

	diseaseMulti := dk.dkDiseaseMultiplier(0.125)

	effect := core.SpellEffect{
		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) *
					(1.0 + dk.dkCountActiveDiseases(hitEffect.Target)*diseaseMulti) *
					dk.RoRTSBonus(hitEffect.Target) *
					dk.mercilessCombatBonus(sim)
			},
		},

		OnSpellHitDealt: onhit,
	}

	procMask := dk.threatOfThassarianProcMasks(isMH, &effect)

	conf := core.SpellConfig{
		ActionID:    ObliterateActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		BonusCritRating: (dk.rimeCritBonus() + dk.subversionCritBonus() + dk.annihilationCritBonus() + dk.scourgeborneBattlegearCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: .8 *
			core.TernaryFloat64(isMH, 1, dk.nervesOfColdSteelBonus()) *
			core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfObliterate), 1.25, 1.0) *
			dk.scourgelordsBattlegearDamageBonus(dk.Obliterate),
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.GuileOfGorefiend),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	}

	rs := &RuneSpell{}
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
				cast.GCD = dk.getModifiedGCD()
			},
			IgnoreHaste: true,
		}
		conf.ApplyEffects = dk.withRuneRefund(rs, effect, false)
		if dk.Talents.DeathRuneMastery == 3 {
			rs.DeathConvertChance = 1.0
		} else {
			rs.DeathConvertChance = float64(dk.Talents.DeathRuneMastery) * 0.33
		}
		rs.ConvertType = RuneTypeFrost | RuneTypeUnholy
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
	diseaseConsumptionChance := []float64{1.0, 0.67, 0.34, 0.0}[dk.Talents.Annihilation]

	dk.ObliterateMhHit = dk.newObliterateHitSpell(true, func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
		dk.LastOutcome = spellEffect.Outcome
		dk.threatOfThassarianProc(sim, spellEffect, dk.ObliterateOhHit)

		if sim.RandomFloat("Annihilation") < diseaseConsumptionChance {
			dk.FrostFeverDisease[spellEffect.Target.Index].Deactivate(sim)
			dk.BloodPlagueDisease[spellEffect.Target.Index].Deactivate(sim)
		}

		if sim.RandomFloat("Rime") < dk.rimeHbChanceProc() {
			dk.RimeAura.Activate(sim)
		}
	})
	dk.ObliterateOhHit = dk.newObliterateHitSpell(false, nil)
	dk.Obliterate = dk.ObliterateMhHit
}
