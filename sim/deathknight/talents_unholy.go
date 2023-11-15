package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) ApplyUnholyTalents() {
	dk.PseudoStats.BaseDodge += 0.01 * float64(dk.Talents.Anticipation)
	dk.AddStat(stats.SpellHit, core.SpellHitRatingPerHitChance*float64(dk.Talents.Virulence))

	if dk.Talents.RavenousDead > 0 {
		dk.MultiplyStat(stats.Strength, 1.0+0.01*float64(dk.Talents.RavenousDead))
	}

	dk.applyNecrosis()
	dk.applyBloodCakedBlade()
	dk.applyUnholyBlight()
	dk.applyImpurity()
	dk.applyDesolation()
	dk.applyWanderingPlague()
	dk.applyEbonPlaguebringerOrCryptFever()
	dk.applyRageOfRivendare()
}

func (dk *Deathknight) viciousStrikesCritChanceBonus() float64 {
	return 3 * float64(dk.Talents.ViciousStrikes)
}

func (dk *Deathknight) applyRageOfRivendare() {
	if dk.Talents.RageOfRivendare == 0 {
		return
	}

	dk.AddStat(stats.Expertise, float64(dk.Talents.RageOfRivendare)*core.ExpertisePerQuarterPercentReduction)

	bonus := 1.0 + 0.02*float64(dk.Talents.RageOfRivendare)
	dk.RoRTSBonus = func(target *core.Unit) float64 {
		// assume if external ebon plaguebringer is active, then another DK will always have Blood Plague up
		if dk.MakeTSRoRAssumptions && target.HasActiveAura("EbonPlaguebringer-1") {
			return bonus
		}

		return core.TernaryFloat64(target.HasActiveAuraWithTag("BloodPlague"), bonus, 1.0)
	}
}

func (dk *Deathknight) applyImpurity() {
	dk.bonusCoeffs.impurityBonusCoeff = 1.0 + float64(dk.Talents.Impurity)*0.04
}

func (dk *Deathknight) getImpurityBonus(spell *core.Spell) float64 {
	return spell.MeleeAttackPower() * dk.bonusCoeffs.impurityBonusCoeff
}

func (dk *Deathknight) applyWanderingPlague() {
	if dk.Talents.WanderingPlague == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 50526}

	// this scales with damage taken modifiers, and very slightly with (shadow) damage dealt modifiers (~10%):
	// e.g. in blood presence, a frost fever tick for 1130 hits debuffed targets for 1146 (+1.5%), and non debuffed
	// targets for 1015 (-13%, +1.5%)

	dk.WanderingPlague = dk.Unit.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreAttackerModifiers | core.SpellFlagNoOnDamageDealt,

		DamageMultiplier: []float64{0.0, 0.33, 0.66, 1.0}[dk.Talents.WanderingPlague],
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.LastDiseaseDamage * dk.bonusCoeffs.wanderingPlagueMultiplier
			baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeAlwaysHit)
			}
		},
	})
}

func (dk *Deathknight) applyNecrosis() {
	if dk.Talents.Necrosis == 0 {
		return
	}

	dk.NecrosisCoeff = 0.04 * float64(dk.Talents.Necrosis)
	dk.Necrosis = dk.Unit.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 51460},
		SpellSchool:      core.SpellSchoolShadow,
		ProcMask:         core.ProcMaskEmpty,
		Flags:            core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers | core.SpellFlagNoOnDamageDealt,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
	})

	// Replace normal melee swing applier with one that also applies necrosis damage.
	// Doing it this way means you don't see necrosis dmg on the timeline but is faster.
	dk.AutoAttacks.MHConfig().ApplyEffects = dk.necrosisMHAuto
	dk.AutoAttacks.OHConfig().ApplyEffects = dk.necrosisOHAuto
}

func (dk *Deathknight) necrosisDamage(damage float64, sim *core.Simulation, target *core.Unit) {
	dk.Necrosis.SpellMetrics[target.UnitIndex].Casts++
	dk.Necrosis.CalcAndDealDamage(sim, target, damage*dk.NecrosisCoeff, dk.Necrosis.OutcomeAlwaysHit)
}
func (dk *Deathknight) necrosisOHAuto(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	baseDamage := spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower()) +
		spell.BonusWeaponDamage()

	if result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWhite); result.Damage > 0 {
		dk.necrosisDamage(result.Damage, sim, target)
	}
}
func (dk *Deathknight) necrosisMHAuto(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
		spell.BonusWeaponDamage()

	if result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWhite); result.Damage > 0 {
		dk.necrosisDamage(result.Damage, sim, target)
	}
}

func (dk *Deathknight) applyBloodCakedBlade() {
	if dk.Talents.BloodCakedBlade == 0 {
		return
	}

	procChance := float64(dk.Talents.BloodCakedBlade) * 0.10
	bloodCakedBladeHitMh := dk.bloodCakedBladeHit(true)
	bloodCakedBladeHitOh := dk.bloodCakedBladeHit(false)

	dk.BloodCakedBladeAura = core.MakePermanent(dk.RegisterAura(core.Aura{
		Label: "Blood-Caked Blade",
		// ActionID: core.ActionID{SpellID: 49628}, // Hide from metrics
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Damage == 0 || !spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}

			if sim.RandomFloat("Blood-Caked Blade Roll") < procChance {
				isMh := spell.ProcMask.Matches(core.ProcMaskMeleeMHAuto)
				if isMh {
					bloodCakedBladeHitMh.Cast(sim, result.Target)
				} else {
					bloodCakedBladeHitOh.Cast(sim, result.Target)
				}
			}
		},
	}))
}

func (dk *Deathknight) bloodCakedBladeHit(isMh bool) *core.Spell {
	return dk.Unit.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 50463}.WithTag(core.TernaryInt32(isMh, 1, 2)),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskProc,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagMeleeMetrics,

		DamageMultiplier: 1 *
			core.TernaryFloat64(isMh, 1, dk.nervesOfColdSteelBonus()),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage float64
			if isMh {
				baseDamage = 0 +
					spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			} else {
				baseDamage = 0 +
					spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			}
			baseDamage *= 0.25 + 0.125*dk.dkCountActiveDiseasesBcb(target)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialNoCrit)
		},
	})
}

func (dk *Deathknight) applyEbonPlaguebringerOrCryptFever() {
	if dk.Talents.CryptFever == 0 {
		return
	}

	ebonPlaguebringerBonusCrit := core.CritRatingPerCritChance * float64(dk.Talents.EbonPlaguebringer)
	dk.AddStat(stats.MeleeCrit, ebonPlaguebringerBonusCrit)
	dk.AddStat(stats.SpellCrit, ebonPlaguebringerBonusCrit)

	dk.EbonPlagueOrCryptFeverAura = dk.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.EbonPlaguebringerOrCryptFeverAura(dk.GetCharacter(), target, dk.Talents.Epidemic, dk.Talents.CryptFever, dk.Talents.EbonPlaguebringer)
	})
	dk.Env.RegisterPreFinalizeEffect(func() {
		dk.FrostFeverSpell.RelatedAuras = append(dk.FrostFeverSpell.RelatedAuras, dk.EbonPlagueOrCryptFeverAura)
		dk.BloodPlagueSpell.RelatedAuras = append(dk.BloodPlagueSpell.RelatedAuras, dk.EbonPlagueOrCryptFeverAura)
	})
}

func (dk *Deathknight) applyDesolation() {
	if dk.Talents.Desolation == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 66803}
	bonusDamageCoeff := 0.01 * float64(dk.Talents.Desolation)

	dk.DesolationAura = dk.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Desolation",
		Duration: time.Second * 20.0,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.ModifyDamageModifier(bonusDamageCoeff)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.ModifyDamageModifier(-bonusDamageCoeff)
		},
	})
}

func (dk *Deathknight) procUnholyBlight(sim *core.Simulation, target *core.Unit, deathCoilDamage float64) {
	if !dk.Talents.UnholyBlight {
		return
	}

	dot := dk.UnholyBlightSpell.Dot(target)

	newDamage := deathCoilDamage * 0.10
	outstandingDamage := core.TernaryFloat64(dot.IsActive(), dot.SnapshotBaseDamage*float64(dot.NumberOfTicks-dot.TickCount), 0)

	dot.SnapshotAttackerMultiplier = dk.UnholyBlightSpell.DamageMultiplier
	dot.SnapshotBaseDamage = (outstandingDamage + newDamage) / float64(dot.NumberOfTicks)

	dk.UnholyBlightSpell.Cast(sim, target)
}

func (dk *Deathknight) applyUnholyBlight() {
	if !dk.Talents.UnholyBlight {
		return
	}

	dk.UnholyBlightSpell = dk.Unit.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 50536},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers | core.SpellFlagNoOnDamageDealt,

		DamageMultiplier: core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfUnholyBlight), 1.4, 1),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "UnholyBlight",
			},
			NumberOfTicks: 10,
			TickLength:    time.Second * 1,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).ApplyOrReset(sim)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
		},
	})
}
