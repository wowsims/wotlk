package hunter

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
)

type PetAbilityType int

const (
	Unknown PetAbilityType = iota
	Bite
	Claw
	FireBreath
	Gore
	LightningBreath
	PoisonSpit
	Screech
)

type PetAbility struct {
	Type PetAbilityType

	// Focus cost
	Cost float64

	*core.Spell
}

// Returns whether the ability was successfully cast.
func (ability *PetAbility) TryCast(sim *core.Simulation, target *core.Unit, hp *HunterPet) bool {
	if hp.currentFocus < ability.Cost {
		return false
	}
	if !ability.IsReady(sim) {
		return false
	}

	hp.SpendFocus(sim, ability.Cost, ability.ActionID)
	ability.Cast(sim, target)
	return true
}

func (hp *HunterPet) NewPetAbility(abilityType PetAbilityType, isPrimary bool) PetAbility {
	switch abilityType {
	case Bite:
		return hp.newBite(isPrimary)
	case Claw:
		return hp.newClaw(isPrimary)
	case FireBreath:
		return PetAbility{}
	case Gore:
		return hp.newGore(isPrimary)
	case LightningBreath:
		return hp.newLightningBreath(isPrimary)
	case PoisonSpit:
		return PetAbility{}
	case Screech:
		return hp.newScreech(isPrimary)
	case Unknown:
		return PetAbility{}
	default:
		panic("Invalid pet ability type")
	}
	return PetAbility{}
}

func (hp *HunterPet) newBite(isPrimary bool) PetAbility {
	return PetAbility{
		Type: Bite,
		Cost: 35,

		Spell: hp.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 27050},
			SpellSchool: core.SpellSchoolPhysical,
			Flags:       core.SpellFlagMeleeMetrics,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
				IgnoreHaste: true,
				CD: core.Cooldown{
					Timer:    hp.NewTimer(),
					Duration: time.Second * 10,
				},
			},

			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskMeleeMHSpecial,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				BaseDamage:       core.BaseDamageConfigRoll(108, 132),
				OutcomeApplier:   hp.OutcomeFuncMeleeSpecialHitAndCrit(2),
			}),
		}),
	}
}

func (hp *HunterPet) newClaw(isPrimary bool) PetAbility {
	return PetAbility{
		Type: Claw,
		Cost: 25,

		Spell: hp.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 27049},
			SpellSchool: core.SpellSchoolPhysical,
			Flags:       core.SpellFlagMeleeMetrics,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
				IgnoreHaste: true,
			},

			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskMeleeMHSpecial,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				BaseDamage:       core.BaseDamageConfigRoll(54, 76),
				OutcomeApplier:   hp.OutcomeFuncMeleeSpecialHitAndCrit(2),
			}),
		}),
	}
}

func (hp *HunterPet) newGore(isPrimary bool) PetAbility {
	return PetAbility{
		Type: Gore,
		Cost: 25,

		Spell: hp.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 35298},
			SpellSchool: core.SpellSchoolPhysical,
			Flags:       core.SpellFlagMeleeMetrics,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
				IgnoreHaste: true,
			},

			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskMeleeMHSpecial,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				BaseDamage: core.WrapBaseDamageConfig(core.BaseDamageConfigRoll(37, 61), func(oldCalculator core.BaseDamageCalculator) core.BaseDamageCalculator {
					return func(sim *core.Simulation, spellEffect *core.SpellEffect, spell *core.Spell) float64 {
						damage := oldCalculator(sim, spellEffect, spell)
						if sim.RandomFloat("Gore") < 0.5 {
							damage *= 2
						}
						return damage
					}
				}),
				OutcomeApplier: hp.OutcomeFuncMeleeSpecialHitAndCrit(2),
			}),
		}),
	}
}

func (hp *HunterPet) newLightningBreath(isPrimary bool) PetAbility {
	return PetAbility{
		Type: LightningBreath,
		Cost: 50,

		Spell: hp.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 25011},
			SpellSchool: core.SpellSchoolNature,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
				IgnoreHaste: true,
			},

			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskSpellDamage,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				BaseDamage:       core.BaseDamageConfigMagic(80, 93, 0.05),
				OutcomeApplier:   hp.OutcomeFuncMagicHitAndCrit(1.5),
			}),
		}),
	}
}

func (hp *HunterPet) newScreech(isPrimary bool) PetAbility {
	var debuffs []*core.Aura
	for _, target := range hp.Env.Encounter.Targets {
		debuffs = append(debuffs, core.ScreechAura(&target.Unit))
	}

	return PetAbility{
		Type: Screech,
		Cost: 20,

		Spell: hp.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 27051},
			SpellSchool: core.SpellSchoolPhysical,
			Flags:       core.SpellFlagMeleeMetrics,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
				IgnoreHaste: true,
			},

			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskMeleeMHSpecial,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				BaseDamage:       core.BaseDamageConfigRoll(33, 61),
				OutcomeApplier:   hp.OutcomeFuncMeleeSpecialHitAndCrit(2),

				OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spellEffect.Landed() {
						for _, debuff := range debuffs {
							debuff.Activate(sim)
						}
					}
				},
			}),
		}),
	}
}
