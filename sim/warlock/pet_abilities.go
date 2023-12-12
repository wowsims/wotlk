package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (wp *WarlockPet) registerFireboltSpell() {
	warlockLevel := wp.owner.GetCharacter().Level
	// assuming max rank available
	rank := map[int32]int{25: 3, 40: 5, 50: 6, 60: 7}[warlockLevel]

	if rank == 0 {
		rank = 1
	}

	spellCoeff := [8]float64{0, .164, .314, .529, .571, .571, .571, .571}[rank]
	baseDamage := [8][]float64{{0, 0}, {7, 10}, {14, 16}, {25, 29}, {36, 41}, {52, 59}, {72, 80}, {85, 96}}[rank]
	spellId := [8]int32{0, 3110, 7799, 7800, 7801, 7802, 11762, 11763}[rank]
	manaCost := [8]float64{0, 10, 20, 35, 50, 70, 95, 115}[rank]
	level := [8]int{0, 1, 8, 18, 28, 38, 48, 58}[rank]

	wp.primaryAbility = wp.RegisterSpell(core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolFire,
		ProcMask:      core.ProcMaskSpellDamage,
		Rank:          rank,
		RequiredLevel: level,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Millisecond * 1000,
				CastTime: time.Millisecond * (2000 - time.Duration(500*wp.owner.Talents.ImprovedFirebolt)),
			},
		},

		DamageMultiplier: (1 + 0.1*float64(wp.owner.Talents.ImprovedImp)),
		CritMultiplier:   wp.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamage[0], baseDamage[1]) + spellCoeff*spell.SpellPower()

			if wp.owner.LakeOfFireAuras != nil && wp.owner.LakeOfFireAuras.Get(target).IsActive() {
				baseDamage *= 1.4
			}

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

// func (wp *WarlockPet) registerCleaveSpell() {
// 	numHits := min(2, wp.Env.GetNumTargets())

// 	wp.primaryAbility = wp.RegisterSpell(core.SpellConfig{
// 		ActionID:    core.ActionID{SpellID: 47994},
// 		SpellSchool: core.SpellSchoolPhysical,
// 		ProcMask:    core.ProcMaskMeleeMHSpecial,
// 		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

// 		ManaCost: core.ManaCostOptions{
// 			FlatCost: 439,
// 		},
// 		Cast: core.CastConfig{
// 			DefaultCast: core.Cast{
// 				GCD: core.GCDDefault,
// 			},
// 			IgnoreHaste: true,
// 			CD: core.Cooldown{
// 				Timer:    wp.NewTimer(),
// 				Duration: time.Second * 6,
// 			},
// 		},

// 		DamageMultiplier: 1,
// 		CritMultiplier:   2,
// 		ThreatMultiplier: 1,

// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			constBaseDamage := 124 + spell.BonusWeaponDamage()

// 			curTarget := target
// 			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
// 				baseDamage := constBaseDamage + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
// 				spell.CalcAndDealDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
// 				curTarget = sim.Environment.NextTargetUnit(curTarget)
// 			}
// 		},
// 	})
// }

// func (wp *WarlockPet) registerInterceptSpell() {
// 	wp.secondaryAbility = nil // not implemented
// }

func (wp *WarlockPet) registerLashOfPainSpell() {
	warlockLevel := wp.owner.GetCharacter().Level
	// assuming max rank available
	rank := map[int32]int{25: 3, 40: 5, 50: 6, 60: 7}[warlockLevel]

	if rank == 0 {
		rank = 1
	}

	spellCoeff := [7]float64{0, .429, .429, .429, .429, .429, .429}[rank]
	baseDamage := [7]float64{0, 33, 44, 60, 73, 87, 99}[rank]
	spellId := [7]int32{0, 7814, 7815, 7816, 11778, 11779, 11780}[rank]
	manaCost := [7]float64{0, 65, 80, 105, 125, 145, 160}[rank]
	level := [7]int{0, 20, 28, 36, 44, 52, 60}[rank]

	wp.primaryAbility = wp.RegisterSpell(core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolShadow,
		ProcMask:      core.ProcMaskSpellDamage,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    wp.NewTimer(),
				Duration: time.Second * (12 - time.Duration(3*wp.owner.Talents.ImprovedLashOfPain)),
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1.5,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := baseDamage + spellCoeff*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

// func (wp *WarlockPet) registerShadowBiteSpell() {
// 	actionID := core.ActionID{SpellID: 54053}

// 	var petManaMetrics *core.ResourceMetrics
// 	maxManaMult := 0.04 * float64(wp.owner.Talents.ImprovedFelhunter)
// 	impFelhunter := wp.owner.Talents.ImprovedFelhunter > 0
// 	if impFelhunter {
// 		petManaMetrics = wp.NewManaMetrics(actionID)
// 	}

// 	wp.primaryAbility = wp.RegisterSpell(core.SpellConfig{
// 		ActionID:    actionID,
// 		SpellSchool: core.SpellSchoolShadow,
// 		ProcMask:    core.ProcMaskSpellDamage,

// 		ManaCost: core.ManaCostOptions{
// 			// TODO: should be 3% of BaseMana, but it's unclear what that actually refers to with pets
// 			FlatCost: 131,
// 		},
// 		Cast: core.CastConfig{
// 			DefaultCast: core.Cast{
// 				GCD: core.GCDDefault,
// 			},
// 			IgnoreHaste: true,
// 			CD: core.Cooldown{
// 				Timer:    wp.NewTimer(),
// 				Duration: time.Second * (6 - time.Duration(2*wp.owner.Talents.ImprovedFelhunter)),
// 			},
// 		},

// 		DamageMultiplier: 1 + 0.03*float64(wp.owner.Talents.ShadowMastery),
// 		CritMultiplier:   1.5 + 0.1*float64(wp.owner.Talents.Ruin),
// 		ThreatMultiplier: 1,

// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			baseDamage := sim.Roll(97+1, 97+41) + 0.429*spell.SpellPower()

// 			w := wp.owner
// 			spells := []*core.Spell{
// 				w.UnstableAffliction,
// 				w.Immolate,
// 				w.CurseOfAgony,
// 				w.CurseOfDoom,
// 				w.Corruption,
// 				w.Conflagrate,
// 				w.Seed,
// 				w.DrainSoul,
// 				// missing: drain life, shadowflame
// 			}
// 			counter := 0
// 			for _, spell := range spells {
// 				if spell != nil && spell.Dot(target).IsActive() {
// 					counter++
// 				}
// 			}

// 			baseDamage *= 1 + 0.15*float64(counter)

// 			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
// 			if impFelhunter && result.Landed() {
// 				wp.AddMana(sim, wp.MaxMana()*maxManaMult, petManaMetrics)
// 			}
// 			spell.DealDamage(sim, result)
// 		},
// 	})
// }
