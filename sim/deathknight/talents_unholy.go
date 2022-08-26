package deathknight

import (
	//"github.com/wowsims/wotlk/sim/core/proto"

	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	//"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) ApplyUnholyTalents() {
	// Virulence
	dk.AddStat(stats.SpellHit, core.SpellHitRatingPerHitChance*float64(dk.Talents.Virulence))

	// Ravenous Dead
	if dk.Talents.RavenousDead > 0 {
		strengthCoeff := 0.01 * float64(dk.Talents.RavenousDead)
		dk.MultiplyStat(stats.Strength, 1.0+strengthCoeff)
	}

	// Necrosis
	dk.applyNecrosis()

	// Blood-Caked Blade
	dk.applyBloodCakedBlade()

	// Unholy Blight
	dk.applyUnholyBlight()

	// Impurity
	dk.applyImpurity()

	// Desolation
	dk.applyDesolation()

	// Wandering Plague
	dk.applyWanderingPlague()

	// Crypt Fever
	dk.applyCryptFever()
	// Ebon Plaguebringer
	dk.applyEbonPlaguebringer()

	// Rage of Rivendare
	dk.AddStat(stats.Expertise, float64(dk.Talents.RageOfRivendare)*core.ExpertisePerQuarterPercentReduction)
	if dk.Talents.RageOfRivendare > 0 {
		dk.applyRageOfRivendare()
	}
}

func (dk *Deathknight) viciousStrikesCritDamageBonus() float64 {
	return 0.15 * float64(dk.Talents.ViciousStrikes)
}

func (dk *Deathknight) viciousStrikesCritChanceBonus() float64 {
	return 3 * float64(dk.Talents.ViciousStrikes)
}

func (dk *Deathknight) applyRageOfRivendare() {
	bonus := 1.0 + 0.02*float64(dk.Talents.RageOfRivendare)
	dk.RoRTSBonus = func(target *core.Unit) float64 {
		return core.TernaryFloat64(dk.BloodPlagueDisease[target.Index].IsActive(), bonus, 1.0)
	}
}

func (dk *Deathknight) applyImpurity() {
	dk.bonusCoeffs.impurityBonusCoeff = 1.0 + float64(dk.Talents.Impurity)*0.04
}

func (dk *Deathknight) getImpurityBonus(hitEffect *core.SpellEffect, unit *core.Unit) float64 {
	return (hitEffect.MeleeAttackPower(unit) + hitEffect.MeleeAttackPowerOnTarget()) * dk.bonusCoeffs.impurityBonusCoeff
}

func (dk *Deathknight) applyWanderingPlague() {
	if dk.Talents.WanderingPlague == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 49655}

	wanderingPlagueMultiplier := []float64{0.0, 0.33, 0.66, 1.0}[dk.Talents.WanderingPlague]

	dk.WanderingPlague = dk.Unit.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagIgnoreAttackerModifiers | core.SpellFlagIgnoreTargetModifiers,

		ApplyEffects: core.ApplyEffectFuncAOEDamageCapped(dk.Env, core.SpellEffect{
			ProcMask: core.ProcMaskSpellDamage,

			DamageMultiplier: wanderingPlagueMultiplier,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(_ *core.Simulation, _ *core.SpellEffect, _ *core.Spell) float64 {
					return dk.LastDiseaseDamage
				},
			},
			OutcomeApplier: dk.OutcomeFuncAlwaysHit(),
		}),
	})
}

func (dk *Deathknight) applyNecrosis() {
	if dk.Talents.Necrosis == 0 {
		return
	}

	var curDmg float64
	necrosisCoeff := 0.04 * float64(dk.Talents.Necrosis)
	necrosisHit := dk.Unit.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 51460},
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagIgnoreAttackerModifiers | core.SpellFlagIgnoreTargetModifiers,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskSpellDamage,

			DamageMultiplier: necrosisCoeff,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(_ *core.Simulation, _ *core.SpellEffect, _ *core.Spell) float64 {
					return curDmg
				},
			},
			OutcomeApplier: dk.OutcomeFuncAlwaysHit(),
		}),
	})

	dk.NecrosisAura = core.MakePermanent(dk.RegisterAura(core.Aura{
		Label: "Necrosis",
		// ActionID: core.ActionID{SpellID: 51465}, // hide from metrics
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Damage == 0 || !spellEffect.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}

			curDmg = spellEffect.Damage
			necrosisHit.Cast(sim, spellEffect.Target)
		},
	}))
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
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Damage == 0 || !spellEffect.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}

			if sim.RandomFloat("Blood-Caked Blade Roll") < procChance {
				isMh := spellEffect.ProcMask.Matches(core.ProcMaskMeleeMHAuto)
				if isMh {
					bloodCakedBladeHitMh.Cast(sim, spellEffect.Target)
				} else {
					bloodCakedBladeHitOh.Cast(sim, spellEffect.Target)
				}
			}
		},
	}))
}

func (dk *Deathknight) bloodCakedBladeHit(isMh bool) *core.Spell {
	mhBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, false, 0, 1.0, 1.0, true)
	ohBaseDamage := core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 0, dk.nervesOfColdSteelBonus(), 1.0, true)

	procMask := core.ProcMaskMeleeOHSpecial
	if isMh {
		procMask = core.ProcMaskMeleeMHSpecial
	}

	return dk.Unit.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 50463}.WithTag(core.TernaryInt32(isMh, 1, 2)),
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: procMask,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, spellEffect *core.SpellEffect, spell *core.Spell) float64 {
					diseaseMultiplier := (0.25 + dk.dkCountActiveDiseases(spellEffect.Target)*0.125)
					if isMh {
						return mhBaseDamage(sim, spellEffect, spell) * diseaseMultiplier
					} else {
						return ohBaseDamage(sim, spellEffect, spell) * diseaseMultiplier
					}
				},
			},
			OutcomeApplier: dk.OutcomeFuncMeleeWeaponSpecialNoCrit(),
		}),
	})
}

func (dk *Deathknight) applyCryptFever() {
	if dk.Talents.CryptFever == 0 {
		return
	}

	dk.CryptFeverAura = make([]*core.Aura, dk.Env.GetNumTargets())
	for _, encounterTarget := range dk.Env.Encounter.Targets {
		target := &encounterTarget.Unit

		cfAura := core.CryptFeverAura(target, int(dk.Index))
		cfAura.Duration = time.Second * (15 + 3*time.Duration(dk.Talents.Epidemic))

		dk.CryptFeverAura[target.Index] = cfAura
	}
}

func (dk *Deathknight) applyEbonPlaguebringer() {
	if dk.Talents.EbonPlaguebringer == 0 {
		return
	}

	ebonPlaguebringerBonusCrit := core.CritRatingPerCritChance * float64(dk.Talents.EbonPlaguebringer)
	dk.AddStat(stats.MeleeCrit, ebonPlaguebringerBonusCrit)
	dk.AddStat(stats.SpellCrit, ebonPlaguebringerBonusCrit)

	dk.EbonPlagueAura = make([]*core.Aura, dk.Env.GetNumTargets())
	for _, encounterTarget := range dk.Env.Encounter.Targets {
		target := &encounterTarget.Unit

		epAura := core.EbonPlaguebringerAura(target, int(dk.Index))
		epAura.Duration = time.Second * (15 + 3*time.Duration(dk.Talents.Epidemic))

		dk.EbonPlagueAura[target.Index] = epAura
	}
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
			dk.ModifyAdditiveDamageModifier(sim, bonusDamageCoeff)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.ModifyAdditiveDamageModifier(sim, -bonusDamageCoeff)
		},
	})
}

func (dk *Deathknight) procUnholyBlight(sim *core.Simulation, target *core.Unit, deathCoilDamage float64) {
	if !dk.Talents.UnholyBlight {
		return
	}

	unholyBlightDot := dk.UnholyBlightDot[target.Index]

	newUnholyBlightDamage := deathCoilDamage * 0.10
	if unholyBlightDot.IsActive() {
		newUnholyBlightDamage += dk.UnholyBlightTickDamage[target.Index] * float64(10-unholyBlightDot.TickCount)
	}
	dk.UnholyBlightTickDamage[target.Index] = newUnholyBlightDamage / 10

	// resets the length
	if unholyBlightDot.IsActive() {
		unholyBlightDot.Apply(sim)
	} else {
		dk.UnholyBlightSpell.Cast(sim, target)
	}
}

func (dk *Deathknight) applyUnholyBlight() {
	if !dk.Talents.UnholyBlight {
		return
	}

	glyphDmgBonus := core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfUnholyBlight), 1.4, 1.0)

	actionID := core.ActionID{SpellID: 50536}

	dk.UnholyBlightSpell = dk.Unit.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagIgnoreAttackerModifiers | core.SpellFlagIgnoreTargetModifiers,
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dk.UnholyBlightDot[dk.CurrentTarget.Index].Apply(sim)
		},
	})

	dk.UnholyBlightDot = make([]*core.Dot, dk.Env.GetNumTargets())
	dk.UnholyBlightTickDamage = make([]float64, dk.Env.GetNumTargets())
	for _, encounterTarget := range dk.Env.Encounter.Targets {
		target := &encounterTarget.Unit

		dk.UnholyBlightTickDamage[target.Index] = 0
		dk.UnholyBlightDot[target.Index] = core.NewDot(core.Dot{
			Spell: dk.UnholyBlightSpell,
			Aura: target.RegisterAura(core.Aura{
				Label:    "UnholyBlight-" + strconv.Itoa(int(dk.Index)),
				ActionID: actionID,
			}),
			NumberOfTicks: 10,
			TickLength:    time.Second * 1,
			TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
				ProcMask:         core.ProcMaskPeriodicDamage,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				IsPeriodic:       true,
				BaseDamage: core.BaseDamageConfig{
					Calculator: func(_ *core.Simulation, se *core.SpellEffect, _ *core.Spell) float64 {
						return dk.UnholyBlightTickDamage[se.Target.Index] * glyphDmgBonus
					},
				},
				OutcomeApplier: dk.OutcomeFuncTick(),
			}),
		})
	}
}
