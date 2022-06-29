package warlock

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

type PetAbilityType byte

const (
	Unknown PetAbilityType = iota
	Cleave
	Intercept
	LashOfPain
	Firebolt
)

// Returns whether the ability was successfully cast.
func (wp *WarlockPet) TryCast(sim *core.Simulation, target *core.Unit, spell *core.Spell) bool {
	if wp.CurrentMana() < spell.DefaultCast.Cost {
		return false
	}
	if !spell.IsReady(sim) {
		return false
	}

	spell.Cast(sim, target)
	return true
}

func (wp *WarlockPet) NewPetAbility(abilityType PetAbilityType, isPrimary bool) *core.Spell {
	switch abilityType {
	case Cleave:
		return wp.newCleave()
	case Intercept:
		return wp.newIntercept()
	case LashOfPain:
		return wp.newLashOfPain()
	case Firebolt:
		return wp.newFirebolt()
	case Unknown:
		return nil
	default:
		panic("Invalid pet ability type")
	}
}

func (wp *WarlockPet) newIntercept() *core.Spell {
	return nil
}

func (wp *WarlockPet) newFirebolt() *core.Spell {
	baseCost := 190.0
	return wp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 27267},
		SpellSchool: core.SpellSchoolFire,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*2000 - (time.Millisecond * time.Duration(250*wp.owner.Talents.ImprovedFirebolt)),
			},
			IgnoreHaste: true,
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			DamageMultiplier: 1.0 + (0.1 * float64(wp.owner.Talents.ImprovedImp)),
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagic(112, 127, 0.571),
			OutcomeApplier:   wp.OutcomeFuncMagicHitAndCrit(2),
		}),
	})
}
func (wp *WarlockPet) newCleave() *core.Spell {
	baseCost := 295.0 // 10% of base
	return wp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 30223},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    wp.NewTimer(),
				Duration: time.Second * 6,
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			DamageMultiplier: 1.0,
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 78, 1.0, true),
			OutcomeApplier:   wp.OutcomeFuncMeleeSpecialHitAndCrit(2),
		}),
	})
}

func (wp *WarlockPet) newLashOfPain() *core.Spell {
	baseCost := 190.0
	return wp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 27274},
		SpellSchool: core.SpellSchoolShadow,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    wp.NewTimer(),
				Duration: time.Second*12 - (time.Second * time.Duration(3*wp.owner.Talents.ImprovedLashOfPain)),
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			DamageMultiplier: 1.0 * (1.0 + (0.1 * float64(wp.owner.Talents.ImprovedSayaad))),
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagic(123, 123, 0.429),
			OutcomeApplier:   wp.OutcomeFuncMagicHitAndCrit(2),
		}),
	})
}
