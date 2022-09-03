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
	baseCost := 439.0 // 10% of base

	baseEffect := core.SpellEffect{
		ProcMask:         core.ProcMaskMeleeMHSpecial,
		DamageMultiplier: 1.0,
		ThreatMultiplier: 1,
		BaseDamage:       core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 124, 1.0, 1.0, true),
		OutcomeApplier:   wp.OutcomeFuncMeleeSpecialHitAndCrit(2),
	}

	numHits := core.MinInt32(2, wp.Env.GetNumTargets())
	effects := make([]core.SpellEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		effects = append(effects, baseEffect)
		effects[i].Target = wp.Env.GetTargetUnit(i)
	}

	return wp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47994},
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
		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})
}

func (wp *WarlockPet) newLashOfPain() *core.Spell {
	baseCost := 250.0
	return wp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47992},
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
			DamageMultiplier: 1.0,
			ThreatMultiplier: 1,
			// TODO: the hidden 5% damage modifier succ currently gets also applies to this ...
			BaseDamage:     core.BaseDamageConfigMagic(237, 237, 0.429),
			OutcomeApplier: wp.OutcomeFuncMagicHitAndCrit(1.5),
		}),
	})
}

func (wp *WarlockPet) newShadowBite() *core.Spell {
	actionID := core.ActionID{SpellID: 54053}
	baseCost := 131.0 // TODO: should be 3% of BaseMana, but it's unclear what that actually refers to with pets

	var onSpellHitDealt func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect)
	if wp.owner.Talents.ImprovedFelhunter > 0 {
		petManaMetrics := wp.NewManaMetrics(actionID)
		maxManaMult := 0.04 * float64(wp.owner.Talents.ImprovedFelhunter)
		onSpellHitDealt = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				wp.AddMana(sim, wp.MaxMana()*maxManaMult, petManaMetrics, true)
			}
		}
	}

	return wp.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
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
			ProcMask:         core.ProcMaskSpellDamage,
			DamageMultiplier: 1.0 + 0.03*float64(wp.owner.Talents.ShadowMastery),
			ThreatMultiplier: 1,
			BaseDamage: core.WrapBaseDamageConfig(core.BaseDamageConfigMagic(97+1, 97+41, 0.429),
				func(oldCalc core.BaseDamageCalculator) core.BaseDamageCalculator {
					return func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
						w := wp.owner
						dots := []*core.Dot{
							w.UnstableAfflictionDot, w.ImmolateDot, w.CurseOfAgonyDot,
							w.CurseOfDoomDot, w.CorruptionDot, w.ConflagrateDot,
							w.SeedDots[hitEffect.Target.Index], w.DrainSoulDot,
							// missing: drain life, shadowflame
						}
						counter := 0
						for _, dot := range dots {
							if dot.IsActive() {
								counter++
							}
						}

						return oldCalc(sim, hitEffect, spell) * (1.0 + 0.15*float64(counter))
					}
				}),
			OutcomeApplier:  wp.OutcomeFuncMagicHitAndCritBinary(1.5 + 0.1*float64(wp.owner.Talents.Ruin)),
			OnSpellHitDealt: onSpellHitDealt,
		}),
	})
}
