package deathknight

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (dk *Deathknight) drwCountActiveDiseases(target *core.Unit) float64 {
	count := 0
	if dk.RuneWeapon.FrostFeverDisease[target.Index].IsActive() {
		count++
	}
	if dk.RuneWeapon.BloodPlagueDisease[target.Index].IsActive() {
		count++
	}
	return float64(count)
}

func (dk *Deathknight) dkCountActiveDiseases(target *core.Unit) float64 {
	count := 0
	if dk.FrostFeverDisease[target.Index].IsActive() {
		count++
	}
	if dk.BloodPlagueDisease[target.Index].IsActive() {
		count++
	}
	if dk.Talents.CryptFever > 0 && dk.CryptFeverAura[target.Index].IsActive() {
		count++
	} else if dk.Talents.EbonPlaguebringer > 0 && dk.EbonPlagueAura[target.Index].IsActive() {
		count++
	}
	return float64(count)
}

// diseaseMultiplier calculates the bonus based on if you have DarkrunedBattlegear 4p.
//  This function is slow so should only be used during initialization.
func (dk *Deathknight) dkDiseaseMultiplier(multiplier float64) float64 {
	if dk.Env.IsFinalized() {
		panic("dont call dk.diseaseMultiplier function during runtime, cache result during initialization")
	}
	if dk.HasSetBonus(ItemSetDarkrunedBattlegear, 4) {
		return multiplier * 1.2
	}
	return multiplier
}

func (dk *Deathknight) registerDiseaseDots() {
	dk.registerFrostFever()
	dk.registerBloodPlague()
}

func (dk *Deathknight) registerFrostFever() {
	actionID := core.ActionID{SpellID: 55095}

	flagTs := make([]bool, dk.Env.GetNumTargets())
	isRefreshing := make([]bool, dk.Env.GetNumTargets())

	dk.FrostFeverSpell = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFrost,
		Flags:       core.SpellFlagDisease,
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			if dk.FrostFeverDisease[unit.Index].IsActive() {
				isRefreshing[unit.Index] = true
			}
			dk.FrostFeverDisease[unit.Index].Apply(sim)
			isRefreshing[unit.Index] = false
			dk.FrostFeverDebuffAura[unit.Index].Activate(sim)
		},
	}, nil, nil)

	dk.FrostFeverDisease = make([]*core.Dot, dk.Env.GetNumTargets())
	dk.FrostFeverExtended = make([]int, dk.Env.GetNumTargets())

	var wpWrapper func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect)
	if dk.Talents.WanderingPlague > 0 {
		wpWrapper = dk.wpWrapper
	}
	for _, encounterTarget := range dk.Env.Encounter.Targets {
		target := &encounterTarget.Unit
		aura := core.Aura{
			Label:    "FrostFever-" + strconv.Itoa(int(dk.Index)),
			ActionID: actionID,
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if !isRefreshing[aura.Unit.Index] {
					flagTs[aura.Unit.Index] = false
				}
			},
		}
		if dk.Talents.IcyTalons > 0 {
			aura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
				dk.IcyTalonsAura.Activate(sim)
			}
		}
		dk.FrostFeverDisease[target.Index] = core.NewDot(core.Dot{
			Aura:          target.RegisterAura(aura),
			NumberOfTicks: 5 + int(dk.Talents.Epidemic),
			TickLength:    time.Second * 3,
			TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
				ProcMask:              core.ProcMaskPeriodicDamage,
				DamageMultiplier:      core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfIcyTouch), 1.2, 1.0),
				ThreatMultiplier:      1,
				IsPeriodic:            true,
				OnPeriodicDamageDealt: wpWrapper,
				BaseDamage: core.BaseDamageConfig{
					Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
						firstTsApply := !flagTs[hitEffect.Target.Index]
						flagTs[hitEffect.Target.Index] = true
						// 80.0 * 0.32 * 1.15 base, 0.055 * 1.15
						return (29.44 + dk.getImpurityBonus(hitEffect, spell.Unit)*0.06325) *
							core.TernaryFloat64(firstTsApply, 1.0, dk.RoRTSBonus(hitEffect.Target))
					},
					TargetSpellCoefficient: 1,
				},
				OutcomeApplier: dk.OutcomeFuncAlwaysHit(),
			}),
		})

		dk.FrostFeverDisease[target.Index].Spell = dk.FrostFeverSpell.Spell
	}
}

func (dk *Deathknight) registerBloodPlague() {
	actionID := core.ActionID{SpellID: 55078}

	flagRor := make([]bool, dk.Env.GetNumTargets())
	isRefreshing := make([]bool, dk.Env.GetNumTargets())

	dk.BloodPlagueSpell = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagDisease,
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			if dk.BloodPlagueDisease[unit.Index].IsActive() {
				isRefreshing[unit.Index] = true
			}
			dk.BloodPlagueDisease[unit.Index].Apply(sim)
			isRefreshing[unit.Index] = false
		},
	}, nil, nil)

	dk.BloodPlagueDisease = make([]*core.Dot, dk.Env.GetNumTargets())
	dk.BloodPlagueExtended = make([]int, dk.Env.GetNumTargets())

	// Tier9 4Piece
	outcomeApplier := dk.OutcomeFuncAlwaysHit()
	if dk.HasSetBonus(ItemSetThassariansBattlegear, 4) {
		outcomeApplier = dk.OutcomeFuncMagicCrit(dk.spellCritMultiplier())
	}

	var wpWrapper func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect)
	if dk.Talents.WanderingPlague > 0 {
		wpWrapper = dk.wpWrapper
	}
	for _, encounterTarget := range dk.Env.Encounter.Targets {
		target := &encounterTarget.Unit

		dk.BloodPlagueDisease[target.Index] = core.NewDot(core.Dot{
			Aura: target.RegisterAura(core.Aura{
				Label:    "BloodPlague-" + strconv.Itoa(int(dk.Index)),
				ActionID: actionID,
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if !isRefreshing[aura.Unit.Index] {
						flagRor[aura.Unit.Index] = false
					}
				},
			}),
			NumberOfTicks: 5 + int(dk.Talents.Epidemic),
			TickLength:    time.Second * 3,

			TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
				ProcMask:              core.ProcMaskPeriodicDamage,
				DamageMultiplier:      1,
				ThreatMultiplier:      1,
				IsPeriodic:            true,
				OnPeriodicDamageDealt: wpWrapper,
				BaseDamage: core.BaseDamageConfig{
					Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
						firstRorApply := !flagRor[hitEffect.Target.Index]
						flagRor[hitEffect.Target.Index] = true
						// 80.0 * 0.394 * 1.15 for base, 0.055 * 1.15 for ap coeff
						return (36.248 + dk.getImpurityBonus(hitEffect, spell.Unit)*0.06325) *
							core.TernaryFloat64(firstRorApply, 1.0, dk.RoRTSBonus(hitEffect.Target))
					},
					TargetSpellCoefficient: 1,
				},
				OutcomeApplier: outcomeApplier,
			}),
		})

		dk.BloodPlagueDisease[target.Index].Spell = dk.BloodPlagueSpell.Spell
	}
}
func (dk *Deathknight) registerDrwDiseaseDots() {
	dk.registerDrwFrostFever()
	dk.registerDrwBloodPlague()
}

func (dk *Deathknight) registerDrwFrostFever() {
	actionID := core.ActionID{SpellID: 55095}

	dk.RuneWeapon.FrostFeverSpell = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFrost,
		Flags:       core.SpellFlagDisease,
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dk.RuneWeapon.FrostFeverDisease[unit.Index].Apply(sim)
		},
	})

	dk.RuneWeapon.FrostFeverDisease = make([]*core.Dot, dk.Env.GetNumTargets())

	for _, encounterTarget := range dk.Env.Encounter.Targets {
		target := &encounterTarget.Unit

		dk.RuneWeapon.FrostFeverDisease[target.Index] = core.NewDot(core.Dot{
			Aura: target.RegisterAura(core.Aura{
				Label:    "DrwFrostFever-" + strconv.Itoa(int(dk.RuneWeapon.Index)),
				ActionID: actionID,
			}),
			NumberOfTicks: 5 + int(dk.Talents.Epidemic),
			TickLength:    time.Second * 3,
			TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
				ProcMask:         core.ProcMaskPeriodicDamage,
				DamageMultiplier: core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfIcyTouch), 1.2, 1.0),
				ThreatMultiplier: 1,
				IsPeriodic:       true,
				BaseDamage: core.BaseDamageConfig{
					Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
						// 80.0 * 0.32 * 1.15 base, 0.055 * 1.15
						return (29.44 + dk.RuneWeapon.getImpurityBonus(hitEffect, spell.Unit)*0.06325)
					},
					TargetSpellCoefficient: 1,
				},
				OutcomeApplier: dk.RuneWeapon.OutcomeFuncAlwaysHit(),
			}),
		})

		dk.RuneWeapon.FrostFeverDisease[target.Index].Spell = dk.RuneWeapon.FrostFeverSpell
	}
}

func (dk *Deathknight) registerDrwBloodPlague() {
	actionID := core.ActionID{SpellID: 55078}

	dk.RuneWeapon.BloodPlagueSpell = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagDisease,
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dk.RuneWeapon.BloodPlagueDisease[unit.Index].Apply(sim)
		},
	})

	dk.RuneWeapon.BloodPlagueDisease = make([]*core.Dot, dk.Env.GetNumTargets())

	// Tier9 4Piece
	outcomeApplier := dk.RuneWeapon.OutcomeFuncAlwaysHit()
	if dk.HasSetBonus(ItemSetThassariansBattlegear, 4) {
		outcomeApplier = dk.RuneWeapon.OutcomeFuncMagicCrit(dk.spellCritMultiplier())
	}

	for _, encounterTarget := range dk.Env.Encounter.Targets {
		target := &encounterTarget.Unit

		dk.RuneWeapon.BloodPlagueDisease[target.Index] = core.NewDot(core.Dot{
			Aura: target.RegisterAura(core.Aura{
				Label:    "DrwBloodPlague-" + strconv.Itoa(int(dk.RuneWeapon.Index)),
				ActionID: actionID,
			}),
			NumberOfTicks: 5 + int(dk.Talents.Epidemic),
			TickLength:    time.Second * 3,

			TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
				ProcMask:         core.ProcMaskPeriodicDamage,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				IsPeriodic:       true,
				BaseDamage: core.BaseDamageConfig{
					Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
						// 80.0 * 0.394 * 1.15 for base, 0.055 * 1.15 for ap coeff
						return (36.248 + dk.RuneWeapon.getImpurityBonus(hitEffect, spell.Unit)*0.06325)
					},
					TargetSpellCoefficient: 1,
				},
				OutcomeApplier: outcomeApplier,
			}),
		})

		dk.RuneWeapon.BloodPlagueDisease[target.Index].Spell = dk.RuneWeapon.BloodPlagueSpell
	}
}

func (dk *Deathknight) wpWrapper(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
	dk.doWanderingPlague(sim, spell, spellEffect)
}

func (dk *Deathknight) doWanderingPlague(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
	if dk.Talents.WanderingPlague == 0 {
		return
	}

	if dk.LastTickTime == sim.CurrentTime {
		return
	}

	physCritChance := spellEffect.PhysicalCritChance(spell.Unit, spell, dk.AttackTables[spellEffect.Target.UnitIndex])
	if sim.RandomFloat("Wandering Plague Roll") < physCritChance {
		dk.LastTickTime = sim.CurrentTime
		dk.LastDiseaseDamage = spellEffect.Damage
		dk.WanderingPlague.Cast(sim, spellEffect.Target)
	}
}
