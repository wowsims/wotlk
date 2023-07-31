package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (wp *WarlockPet) registerCleaveSpell() {
	numHits := core.MinInt32(2, wp.Env.GetNumTargets())

	wp.primaryAbility = wp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47994},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		ManaCost: core.ManaCostOptions{
			FlatCost: 439,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    wp.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			constBaseDamage := 124 + spell.BonusWeaponDamage()

			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := constBaseDamage + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
				spell.CalcAndDealDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
		},
	})
}

func (wp *WarlockPet) registerInterceptSpell() {
	wp.secondaryAbility = nil // not implemented
}

func (wp *WarlockPet) registerLashOfPainSpell() {
	wp.primaryAbility = wp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47992},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		ManaCost: core.ManaCostOptions{
			FlatCost: 250,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    wp.NewTimer(),
				Duration: time.Second * (12 - time.Duration(3*wp.owner.Talents.DemonicPower)),
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1.5,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO: the hidden 5% damage modifier succ currently gets also applies to this ...
			baseDamage := 237 + 0.429*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (wp *WarlockPet) registerShadowBiteSpell() {
	actionID := core.ActionID{SpellID: 54053}

	var petManaMetrics *core.ResourceMetrics
	maxManaMult := 0.04 * float64(wp.owner.Talents.ImprovedFelhunter)
	impFelhunter := wp.owner.Talents.ImprovedFelhunter > 0
	if impFelhunter {
		petManaMetrics = wp.NewManaMetrics(actionID)
	}

	wp.primaryAbility = wp.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		ManaCost: core.ManaCostOptions{
			// TODO: should be 3% of BaseMana, but it's unclear what that actually refers to with pets
			FlatCost: 131,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    wp.NewTimer(),
				Duration: time.Second * (6 - time.Duration(2*wp.owner.Talents.ImprovedFelhunter)),
			},
		},

		DamageMultiplier: 1 + 0.03*float64(wp.owner.Talents.ShadowMastery),
		CritMultiplier:   1.5 + 0.1*float64(wp.owner.Talents.Ruin),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(97+1, 97+41) + 0.429*spell.SpellPower()

			w := wp.owner
			spells := []*core.Spell{
				w.UnstableAffliction,
				w.Immolate,
				w.CurseOfAgony,
				w.CurseOfDoom,
				w.Corruption,
				w.Conflagrate,
				w.Seed,
				w.DrainSoul,
				// missing: drain life, shadowflame
			}
			counter := 0
			for _, spell := range spells {
				if spell != nil && spell.Dot(target).IsActive() {
					counter++
				}
			}

			baseDamage *= 1 + 0.15*float64(counter)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if impFelhunter && result.Landed() {
				wp.AddMana(sim, wp.MaxMana()*maxManaMult, petManaMetrics)
			}
			spell.DealDamage(sim, result)
		},
	})
}

func (wp *WarlockPet) registerFireboltSpell() {
	wp.primaryAbility = wp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47964},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,

		ManaCost: core.ManaCostOptions{
			FlatCost: 180,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * (2500 - time.Duration(250*wp.owner.Talents.DemonicPower)),
			},
		},

		DamageMultiplier: (1 + 0.1*float64(wp.owner.Talents.ImprovedImp)) *
			(1 + 0.2*core.TernaryFloat64(wp.owner.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfImp), 1, 0)),
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(203, 227) + 0.571*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
