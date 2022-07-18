package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type PetAbilityType byte

const (
	Unknown PetAbilityType = iota
	Cleave
	Intercept
	LashOfPain
	ShadowBite
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
	case ShadowBite:
		return wp.newShadowBite()
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
	baseCost := 180.0
	return wp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47964},
		SpellSchool: core.SpellSchoolFire,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * (2500 - time.Duration(250*wp.owner.Talents.DemonicPower)),
			},
			IgnoreHaste: true,
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskSpellDamage,
			DamageMultiplier: (1.0 + 0.1*float64(wp.owner.Talents.ImprovedImp)) *
				(1.0 + 0.2*core.TernaryFloat64(wp.owner.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfImp), 1, 0)),
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagic(203, 227, 0.571),
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
				Duration: time.Second * (12 - time.Duration(3*wp.owner.Talents.DemonicPower)),
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

func (wp *WarlockPet) newShadowBite() *core.Spell {
	baseCost := wp.BaseMana * 0.03

	return wp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 54053},
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
				Duration: time.Second * (6 - time.Duration(2*wp.owner.Talents.ImprovedFelhunter)),
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskSpellDamage,
			DamageMultiplier: (1.0 + 0.03*float64(wp.owner.Talents.ShadowMastery)) * (1 + 0.15*(core.TernaryFloat64(wp.owner.DrainSoulDot.IsActive(), 1, 0) + //core.TernaryFloat64(wp.owner.ConflagrateDot.IsActive(), 1, 0) +
				core.TernaryFloat64(wp.owner.CorruptionDot.IsActive(), 1, 0) + //core.TernaryFloat64(wp.owner.SeedDots.IsActive(), 1, 0) +
				core.TernaryFloat64(wp.owner.CurseOfDoomDot.IsActive(), 1, 0) + core.TernaryFloat64(wp.owner.CurseOfAgonyDot.IsActive(), 1, 0)+
				core.TernaryFloat64(wp.owner.UnstableAffDot.IsActive(), 1, 0) + core.TernaryFloat64(wp.owner.ImmolateDot.IsActive(), 1, 0))),
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagic(98, 138, 0.429), //TODO : change spellpower coefficient
			OutcomeApplier:   wp.OutcomeFuncMagicHitAndCrit(2),
		}),
	})
}
