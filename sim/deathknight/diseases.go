package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (dk *Deathknight) drwCountActiveDiseases(target *core.Unit) float64 {
	count := 0
	if dk.Talents.DancingRuneWeapon {
		if dk.RuneWeapon.FrostFeverSpell.Dot(target).IsActive() {
			count++
		}
		if dk.RuneWeapon.BloodPlagueSpell.Dot(target).IsActive() {
			count++
		}
	}
	return float64(count)
}

func (dk *Deathknight) dkCountActiveDiseases(target *core.Unit) float64 {
	count := 0
	if dk.FrostFeverSpell.Dot(target).IsActive() {
		count++
	}
	if dk.BloodPlagueSpell.Dot(target).IsActive() {
		count++
	}
	if dk.EbonPlagueOrCryptFeverAura[target.Index].IsActive() {
		count++
	}
	return float64(count)
}

func (dk *Deathknight) dkCountActiveDiseasesBcb(target *core.Unit) float64 {
	count := 0
	if dk.FrostFeverSpell.Dot(target).IsActive() {
		count++
	}
	if dk.BloodPlagueSpell.Dot(target).IsActive() {
		count++
	}
	if target.HasActiveAuraWithTag("EbonPlaguebringer") {
		count++
	}
	return float64(count)
}

// diseaseMultiplier calculates the bonus based on if you have DarkrunedBattlegear 4p.
//
//	This function is slow so should only be used during initialization.
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
	dk.FrostFeverSpell = dk.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 55095},
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagDisease,

		DamageMultiplier: core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfIcyTouch), 1.2, 1.0),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "FrostFever",
				Tag:   "FrostFever",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					if dk.IcyTalonsAura != nil {
						dk.IcyTalonsAura.Activate(sim)
					}
					if dk.EbonPlagueOrCryptFeverAura[aura.Unit.Index] != nil {
						dk.EbonPlagueOrCryptFeverAura[aura.Unit.Index].Activate(sim)
					}
				},
			},
			NumberOfTicks: 5 + dk.Talents.Epidemic,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				// 80.0 * 0.32 * 1.15 base, 0.055 * 1.15
				dot.SnapshotBaseDamage = 29.44 + 0.06325*dk.getImpurityBonus(dot.Spell)

				if !isRollover {
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
					dot.SnapshotAttackerMultiplier *= dk.RoRTSBonus(target)
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
				dk.doWanderingPlague(sim, dot.Spell, result)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := spell.Dot(target)
			dot.Apply(sim)
		},
	})
	dk.FrostFeverExtended = make([]int, dk.Env.GetNumTargets())
}

func (dk *Deathknight) registerBloodPlague() {
	// Tier9 4Piece
	canCrit := dk.HasSetBonus(ItemSetThassariansBattlegear, 4)

	dk.BloodPlagueSpell = dk.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 55078},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagDisease,

		DamageMultiplier: 1,
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "BloodPlague",
				Tag:   "BloodPlague",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					if dk.EbonPlagueOrCryptFeverAura[aura.Unit.Index] != nil {
						dk.EbonPlagueOrCryptFeverAura[aura.Unit.Index].Activate(sim)
					}
				},
			},
			NumberOfTicks: 5 + dk.Talents.Epidemic,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				// 80.0 * 0.394 * 1.15 for base, 0.055 * 1.15 for ap coeff
				dot.SnapshotBaseDamage = 36.248 + 0.06325*dk.getImpurityBonus(dot.Spell)

				if !isRollover {
					dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
					dot.SnapshotAttackerMultiplier *= dk.RoRTSBonus(target)
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				var result *core.SpellResult
				if canCrit {
					result = dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				} else {
					result = dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
				}
				dk.doWanderingPlague(sim, dot.Spell, result)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := spell.Dot(target)
			dot.Apply(sim)
		},
	})
	dk.BloodPlagueExtended = make([]int, dk.Env.GetNumTargets())
}
func (dk *Deathknight) registerDrwDiseaseDots() {
	dk.registerDrwFrostFever()
	dk.registerDrwBloodPlague()
}

func (dk *Deathknight) registerDrwFrostFever() {
	dk.RuneWeapon.FrostFeverSpell = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 55095},
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagDisease,

		DamageMultiplier: core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfIcyTouch), 1.2, 1.0),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "DrwFrostFever",
			},
			NumberOfTicks: 5 + dk.Talents.Epidemic,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				// 80.0 * 0.32 * 1.15 base, 0.055 * 1.15
				dot.SnapshotBaseDamage = 29.44 + 0.06325*dk.getImpurityBonus(dot.Spell)

				if !isRollover {
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})

	if !dk.Inputs.NewDrw {
		dk.RuneWeapon.FrostFeverSpell.DamageMultiplier *= 0.5
		dk.RuneWeapon.FrostFeverSpell.Flags |= core.SpellFlagIgnoreAttackerModifiers
	}
}

func (dk *Deathknight) registerDrwBloodPlague() {
	// Tier9 4Piece
	canCrit := dk.HasSetBonus(ItemSetThassariansBattlegear, 4)

	dk.RuneWeapon.BloodPlagueSpell = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 55078},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagDisease,

		DamageMultiplier: 1,
		CritMultiplier:   dk.RuneWeapon.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "DrwBloodPlague",
			},
			NumberOfTicks: 5 + dk.Talents.Epidemic,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				// 80.0 * 0.394 * 1.15 for base, 0.055 * 1.15 for ap coeff
				dot.SnapshotBaseDamage = 36.248 + 0.06325*dk.getImpurityBonus(dot.Spell)

				if !isRollover {
					dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				var result *core.SpellResult
				if canCrit {
					result = dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				} else {
					result = dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
				}
				dk.doWanderingPlague(sim, dot.Spell, result)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})

	if !dk.Inputs.NewDrw {
		dk.RuneWeapon.BloodPlagueSpell.DamageMultiplier *= 0.5
		dk.RuneWeapon.BloodPlagueSpell.Flags |= core.SpellFlagIgnoreAttackerModifiers
	}
}

func (dk *Deathknight) doWanderingPlague(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	if dk.Talents.WanderingPlague == 0 {
		return
	}

	// 500ms ICD
	if sim.CurrentTime < dk.LastTickTime+500*time.Millisecond {
		return
	}

	attackTable := dk.AttackTables[result.Target.UnitIndex]
	physCritChance := spell.PhysicalCritChance(result.Target, attackTable)
	if sim.RandomFloat("Wandering Plague Roll") < physCritChance {
		dk.LastTickTime = sim.CurrentTime
		dk.LastDiseaseDamage = result.Damage / dk.WanderingPlague.TargetDamageMultiplier(attackTable, false)
		dk.WanderingPlague.Cast(sim, result.Target)
	}
}
