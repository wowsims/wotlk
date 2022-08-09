package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

type PetAbilityType int

const (
	Unknown PetAbilityType = iota
	Claw
)

// These IDs are needed for certain talents.
const ClawSpellID = 47468

type PetAbility struct {
	Type PetAbilityType

	// Focus cost
	Cost float64

	*core.Spell
}

// Returns whether the ability was successfully cast.
func (ability *PetAbility) TryCast(sim *core.Simulation, target *core.Unit, ghoul *GhoulPet) bool {
	if ghoul.currentFocus < ability.Cost {
		return false
	}
	if !ability.IsReady(sim) {
		return false
	}

	if !ghoul.PseudoStats.NoCost {
		ghoul.SpendFocus(sim, ability.Cost, ability.ActionID)
	}
	ability.Cast(sim, target)
	return true
}

func (ghoulPet *GhoulPet) NewPetAbility(abilityType PetAbilityType) PetAbility {
	switch abilityType {
	case Claw:
		return ghoulPet.newClaw()
	case Unknown:
		return PetAbility{}
	default:
		panic("Invalid pet ability type")
	}
}

func (ghoulPet *GhoulPet) newClaw() PetAbility {
	return PetAbility{
		Type: Claw,
		Cost: 40,

		Spell: ghoulPet.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: ClawSpellID},
			SpellSchool: core.SpellSchoolPhysical,
			Flags:       core.SpellFlagMeleeMetrics,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
			},

			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskMeleeMHSpecial,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				BaseDamage:       core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 0, 1.0, 1.5, true),
				OutcomeApplier:   ghoulPet.OutcomeFuncMeleeSpecialHitAndCrit(2),
			}),
		}),
	}
}
