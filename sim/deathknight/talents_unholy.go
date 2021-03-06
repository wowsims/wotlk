package deathknight

import (
	//"github.com/wowsims/wotlk/sim/core/proto"

	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	//"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) ApplyUnholyTalents() {
	// Virulence
	dk.AddStat(stats.SpellHit, core.SpellHitRatingPerHitChance*float64(dk.Talents.Virulence))

	// Ravenous Dead
	if dk.Talents.RavenousDead > 0 {
		strengthCoeff := 0.01 * float64(dk.Talents.RavenousDead)
		dk.AddStatDependency(stats.Strength, stats.Strength, 1.0+strengthCoeff)
	}

	// Necrosis
	dk.applyNecrosis()

	// Blood-Caked Blade
	dk.applyBloodCakedBlade()

	// Unholy Blight
	dk.applyUnholyBlight()

	// Reaping
	// TODO:

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
}

func (dk *Deathknight) viciousStrikesCritDamageBonus() float64 {
	return 0.15 * float64(dk.Talents.ViciousStrikes)
}

func (dk *Deathknight) viciousStrikesCritChanceBonus() float64 {
	return 3 * float64(dk.Talents.ViciousStrikes)
}

func (dk *Deathknight) rageOfRivendareBonus(target *core.Unit) float64 {
	return core.TernaryFloat64(dk.TargetHasDisease(BloodPlagueAuraLabel, target), 1.0+0.02*float64(dk.Talents.RageOfRivendare), 1.0)
}

func (dk *Deathknight) applyImpurity(hitEffect *core.SpellEffect, unit *core.Unit) float64 {
	return hitEffect.MeleeAttackPower(unit) * (1.0 + float64(dk.Talents.Impurity)*0.04)
}

func (dk *Deathknight) applyWanderingPlague() {
	if dk.Talents.WanderingPlague == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 49655}

	wanderingPlagueMultiplier := 0.0
	if dk.Talents.WanderingPlague == 1 {
		wanderingPlagueMultiplier = 0.33
	} else if dk.Talents.WanderingPlague == 2 {
		wanderingPlagueMultiplier = 0.66
	} else if dk.Talents.WanderingPlague == 3 {
		wanderingPlagueMultiplier = 1.0
	}

	dk.WanderingPlague = dk.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagNone,

		ApplyEffects: core.ApplyEffectFuncAOEDamage(dk.Env, core.SpellEffect{
			// No proc mask, so it won't proc itself.
			ProcMask: core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(_ *core.Simulation, _ *core.SpellEffect, _ *core.Spell) float64 {
					return dk.LastDiseaseDamage * wanderingPlagueMultiplier
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
	necrosisHit := dk.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 51460},
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagNone,

		ApplyEffects: core.ApplyEffectFuncDirectDamageTargetModifiersOnly(core.SpellEffect{
			// No proc mask, so it won't proc itself.
			ProcMask: core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(_ *core.Simulation, _ *core.SpellEffect, _ *core.Spell) float64 {
					return curDmg * 0.04 * float64(dk.Talents.Necrosis)
				},
			},
			OutcomeApplier: dk.OutcomeFuncAlwaysHit(),
		}),
	})

	dk.NecrosisAura = core.MakePermanent(dk.RegisterAura(core.Aura{
		Label:    "Necrosis",
		ActionID: core.ActionID{SpellID: 51465},
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

	bloodCakedBladeHitMh := dk.bloodCakedBladeHit(true)
	bloodCakedBladeHitOh := dk.bloodCakedBladeHit(false)

	dk.BloodCakedBladeAura = core.MakePermanent(dk.RegisterAura(core.Aura{
		Label:    "Blood-Caked Blade",
		ActionID: core.ActionID{SpellID: 49628},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Damage == 0 || !spellEffect.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}

			if sim.RandomFloat("Blood-Caked Blade Roll") < 0.30 {
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
	mhBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, false, 0, 1.0, true)
	ohBaseDamage := core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 0, 1.0*dk.nervesOfColdSteelBonus(), true)

	return dk.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 50463}.WithTag(core.TernaryInt32(isMh, 1, 2)),
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ApplyEffects: core.ApplyEffectFuncDirectDamageTargetModifiersOnly(core.SpellEffect{
			// No proc mask, so it won't proc itself.
			ProcMask: core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, spellEffect *core.SpellEffect, spell *core.Spell) float64 {
					diseaseMultiplier := (0.25 + float64(dk.countActiveDiseases(spellEffect.Target))*0.125)
					if isMh {
						return mhBaseDamage(sim, spellEffect, spell) * diseaseMultiplier
					} else {
						return ohBaseDamage(sim, spellEffect, spell) * diseaseMultiplier
					}
				},
			},
			OutcomeApplier: dk.OutcomeFuncMeleeWeaponSpecialNoHitNoCrit(),
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

func (dk *Deathknight) applyUnholyBlight() {
	actionID := core.ActionID{SpellID: 50536}

	var curDamage = make([]float64, dk.Env.GetNumTargets())
	dk.UnholyBlightSpell = dk.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			curDamage[dk.CurrentTarget.Index] = dk.LastDeathCoilDamage
			dk.UnholyBlightDot[dk.CurrentTarget.Index].Apply(sim)
		},
	})

	dk.UnholyBlightDot = make([]*core.Dot, dk.Env.GetNumTargets())
	for _, encounterTarget := range dk.Env.Encounter.Targets {
		target := &encounterTarget.Unit

		dk.UnholyBlightDot[target.Index] = core.NewDot(core.Dot{
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
					Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
						return (0.10 * curDamage[target.Index]) / 10
					},
					TargetSpellCoefficient: 1,
				},
				OutcomeApplier: dk.OutcomeFuncAlwaysHit(),
			}),
		})

		dk.UnholyBlightDot[target.Index].Spell = dk.UnholyBlightSpell
	}
}

func (dk *Deathknight) reapingChance() float64 {
	reapingChance := 0.0
	if dk.Talents.Reaping == 1 {
		reapingChance = 0.33
	} else if dk.Talents.Reaping == 2 {
		reapingChance = 0.66
	} else if dk.Talents.Reaping == 3 {
		reapingChance = 1.0
	}
	return reapingChance
}

func (dk *Deathknight) reapingWillProc(sim *core.Simulation, reapingChance float64) bool {
	ohWillCast := sim.RandomFloat("Reaping") <= reapingChance
	return ohWillCast
}

func (dk *Deathknight) reapingProc(sim *core.Simulation, spell *core.Spell, runeCost core.RuneAmount) bool {
	if dk.Talents.Reaping > 0 {
		if runeCost.Blood > 0 {
			reapingChance := dk.reapingChance()

			if dk.reapingWillProc(sim, reapingChance) {
				slot := dk.SpendBloodRune(sim, spell.BloodRuneMetrics())
				dk.SetRuneAtIdxSlotToState(0, slot, core.RuneState_DeathSpent, core.RuneKind_Death)
				dk.SetAsGeneratedByReapingOrBoTN(slot)
				return true
			}
		}
	}
	return false
}
