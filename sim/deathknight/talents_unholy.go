package deathknight

import (
	//"github.com/wowsims/wotlk/sim/core/proto"

	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	//"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) ApplyUnholyTalents() {
	// Virulence
	deathKnight.AddStat(stats.SpellHit, core.SpellHitRatingPerHitChance*float64(deathKnight.Talents.Virulence))

	// Ravenous Dead
	if deathKnight.Talents.RavenousDead > 0 {
		strengthCoeff := 0.01 * float64(deathKnight.Talents.RavenousDead)
		deathKnight.AddStatDependency(stats.Strength, stats.Strength, strengthCoeff)
	}

	// Necrosis
	deathKnight.applyNecrosis()

	// Blood-Caked Blade
	deathKnight.applyBloodCakedBlade()

	// Unholy Blight
	deathKnight.applyUnholyBlight()

	// Reaping
	// TODO:

	// Desolation
	deathKnight.applyDesolation()

	// Wandering Plague
	deathKnight.applyWanderingPlague()

	// Crypt Fever
	// Ebon Plaguebringer
	deathKnight.applyEbonPlaguebringer()

	// Rage of Rivendare
	deathKnight.AddStat(stats.Expertise, float64(deathKnight.Talents.RageOfRivendare)*core.ExpertisePerQuarterPercentReduction)
}

func (deathKnight *DeathKnight) viciousStrikesBonus() float64 {
	return 0.15 * float64(deathKnight.Talents.ViciousStrikes)
}

func (deathKnight *DeathKnight) rageOfRivendareBonus(target *core.Unit) float64 {
	return core.TernaryFloat64(deathKnight.TargetHasDisease(BloodPlagueAuraLabel, target), 1.0+0.02*float64(deathKnight.Talents.RageOfRivendare), 1.0)
}

func (deathKnight *DeathKnight) applyImpurity(hitEffect *core.SpellEffect, unit *core.Unit) float64 {
	return hitEffect.MeleeAttackPower(unit) * (1.0 + float64(deathKnight.Talents.Impurity)*0.04)
}

func (deathKnight *DeathKnight) applyWanderingPlague() {
	if deathKnight.Talents.WanderingPlague == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 49655}

	wanderingPlagueMultiplier := 0.0
	if deathKnight.Talents.WanderingPlague == 1 {
		wanderingPlagueMultiplier = 0.33
	} else if deathKnight.Talents.WanderingPlague == 2 {
		wanderingPlagueMultiplier = 0.66
	} else if deathKnight.Talents.WanderingPlague == 3 {
		wanderingPlagueMultiplier = 1.0
	}

	deathKnight.WanderingPlague = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagNone,

		ApplyEffects: core.ApplyEffectFuncAOEDamage(deathKnight.Env, core.SpellEffect{
			// No proc mask, so it won't proc itself.
			ProcMask: core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(_ *core.Simulation, _ *core.SpellEffect, _ *core.Spell) float64 {
					return deathKnight.LastDiseaseDamage * wanderingPlagueMultiplier
				},
			},
			OutcomeApplier: deathKnight.OutcomeFuncAlwaysHit(),
		}),
	})
}

func (deathKnight *DeathKnight) applyNecrosis() {
	if deathKnight.Talents.Necrosis == 0 {
		return
	}

	target := deathKnight.CurrentTarget

	var curDmg float64
	necrosisHit := deathKnight.RegisterSpell(core.SpellConfig{
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
					return curDmg * 0.04 * float64(deathKnight.Talents.Necrosis)
				},
			},
			OutcomeApplier: deathKnight.OutcomeFuncAlwaysHit(),
		}),
	})

	deathKnight.NecrosisAura = deathKnight.RegisterAura(core.Aura{
		Label:    "Necrosis",
		ActionID: core.ActionID{SpellID: 51465},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			deathKnight.NecrosisAura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Damage == 0 || !spellEffect.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}

			curDmg = spellEffect.Damage
			necrosisHit.Cast(sim, target)
		},
	})
}

func (deathKnight *DeathKnight) applyBloodCakedBlade() {
	if deathKnight.Talents.BloodCakedBlade == 0 {
		return
	}

	target := deathKnight.CurrentTarget

	mhBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, false, 0, 1.0, true)
	ohBaseDamage := core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 0, 1.0, true)

	var isMH = false
	bloodCakedBladeHit := deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 50463},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ApplyEffects: core.ApplyEffectFuncDirectDamageTargetModifiersOnly(core.SpellEffect{
			// No proc mask, so it won't proc itself.
			ProcMask: core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, spellEffect *core.SpellEffect, spell *core.Spell) float64 {
					diseaseMultiplier := (0.25 + float64(deathKnight.countActiveDiseases(spellEffect.Target))*0.125)
					if isMH {
						return mhBaseDamage(sim, spellEffect, spell) * diseaseMultiplier
					} else {
						return ohBaseDamage(sim, spellEffect, spell) * diseaseMultiplier
					}
				},
			},
			OutcomeApplier: deathKnight.OutcomeFuncMeleeWeaponSpecialNoHitNoCrit(),
		}),
	})

	deathKnight.BloodCakedBladeAura = deathKnight.RegisterAura(core.Aura{
		Label:    "Blood-Caked Blade",
		ActionID: core.ActionID{SpellID: 49628},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			deathKnight.BloodCakedBladeAura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Damage == 0 || !spellEffect.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}

			if sim.RandomFloat("Blood-Caked Blade Roll") < 0.30 {
				isMH = spellEffect.ProcMask.Matches(core.ProcMaskMeleeMHAuto)
				bloodCakedBladeHit.Cast(sim, target)
			}
		},
	})
}

func (deathKnight *DeathKnight) applyEbonPlaguebringer() {
	if deathKnight.Talents.EbonPlaguebringer == 0 {
		return
	}

	ebonPlaguebringerBonusCrit := core.CritRatingPerCritChance * float64(deathKnight.Talents.EbonPlaguebringer)
	deathKnight.AddStat(stats.MeleeCrit, ebonPlaguebringerBonusCrit)
	deathKnight.AddStat(stats.SpellCrit, ebonPlaguebringerBonusCrit)

	target := deathKnight.CurrentTarget

	epAura := core.EbonPlaguebringerAura(target)
	epAura.Duration = time.Second * (15 + 3*time.Duration(deathKnight.Talents.Epidemic))

	deathKnight.EbonPlagueAura = epAura
}

func (deathKnight *DeathKnight) applyDesolation() {
	if deathKnight.Talents.Desolation == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 66803}

	deathKnight.DesolationAura = deathKnight.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Desolation",
		Duration: time.Second * 20.0,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.01*float64(deathKnight.Talents.Desolation)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.0 + 0.01*float64(deathKnight.Talents.Desolation)
		},
	})
}

func (deathKnight *DeathKnight) applyUnholyBlight() {
	actionID := core.ActionID{SpellID: 50536}
	target := deathKnight.CurrentTarget

	var curDamage = 0.0
	deathKnight.UnholyBlightDot = core.NewDot(core.Dot{
		Aura: target.RegisterAura(core.Aura{
			Label:    "UnholyBlight-" + strconv.Itoa(int(deathKnight.Index)),
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
					return (0.10 * curDamage) / 10
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: deathKnight.OutcomeFuncAlwaysHit(),
		}),
	})

	deathKnight.UnholyBlightSpell = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			curDamage = deathKnight.LastDeathCoilDamage
			deathKnight.UnholyBlightDot.Apply(sim)
		},
	})

	deathKnight.UnholyBlightDot.Spell = deathKnight.UnholyBlightSpell
}

func (deathKnight *DeathKnight) reapingChance() float64 {
	reapingChance := 0.0
	if deathKnight.Talents.Reaping == 1 {
		reapingChance = 0.33
	} else if deathKnight.Talents.Reaping == 2 {
		reapingChance = 0.66
	} else if deathKnight.Talents.Reaping == 3 {
		reapingChance = 1.0
	}
	return reapingChance
}

func (deathKnight *DeathKnight) reapingWillProc(sim *core.Simulation, reapingChance float64) bool {
	ohWillCast := sim.RandomFloat("Reaping") <= reapingChance
	return ohWillCast
}

func (deathKnight *DeathKnight) reapingProc(sim *core.Simulation, spell *core.Spell, runeCost core.DKRuneCost) bool {
	if deathKnight.Talents.Reaping > 0 {
		if runeCost.Blood > 0 {
			reapingChance := deathKnight.reapingChance()

			if deathKnight.reapingWillProc(sim, reapingChance) {
				slot := deathKnight.SpendBloodRune(sim, spell.BloodRuneMetrics())
				deathKnight.SetRuneAtSlotToState(0, slot, core.RuneState_DeathSpent, core.RuneKind_Death)
				deathKnight.SetAsGeneratedByReapingOrBoTN(slot)
				return true
			}
		}
	}
	return false
}
