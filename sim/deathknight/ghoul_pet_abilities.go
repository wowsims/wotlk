package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"time"
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

	ghoul.SpendFocus(sim, ability.Cost, ability.ActionID)
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
			ProcMask:    core.ProcMaskMeleeMHSpecial,
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: time.Second,
				},
				IgnoreHaste: true,
			},

			DamageMultiplier: 1.5,
			CritMultiplier:   2,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := 0 +
					spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()

				result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				if !result.Landed() {
					ghoulPet.AddFocus(sim, 32, spell.ActionID)
				}
			},
		}),
	}
}
