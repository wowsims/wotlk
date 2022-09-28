package core

import (
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Function for calculating the base damage of a spell.
type BaseDamageCalculator func(*Simulation, *SpellEffect, *Spell) float64

type BaseDamageConfig struct {
	// Lambda for calculating the base damage.
	Calculator BaseDamageCalculator
}

func BuildBaseDamageConfig(calculator BaseDamageCalculator) BaseDamageConfig {
	return BaseDamageConfig{
		Calculator: calculator,
	}
}

func WrapBaseDamageConfig(config BaseDamageConfig, wrapper func(oldCalculator BaseDamageCalculator) BaseDamageCalculator) BaseDamageConfig {
	return BaseDamageConfig{
		Calculator: wrapper(config.Calculator),
	}
}

// Creates a BaseDamageCalculator function which returns a flat value.
func BaseDamageFuncFlat(damage float64) BaseDamageCalculator {
	return func(_ *Simulation, _ *SpellEffect, _ *Spell) float64 {
		return damage
	}
}
func BaseDamageConfigFlat(damage float64) BaseDamageConfig {
	return BuildBaseDamageConfig(BaseDamageFuncFlat(damage))
}

// Creates a BaseDamageCalculator function with a single damage roll.
func BaseDamageFuncRoll(minFlatDamage float64, maxFlatDamage float64) BaseDamageCalculator {
	if minFlatDamage == maxFlatDamage {
		return BaseDamageFuncFlat(minFlatDamage)
	} else {
		deltaDamage := maxFlatDamage - minFlatDamage
		return func(sim *Simulation, _ *SpellEffect, _ *Spell) float64 {
			return damageRollOptimized(sim, minFlatDamage, deltaDamage)
		}
	}
}
func BaseDamageConfigRoll(minFlatDamage float64, maxFlatDamage float64) BaseDamageConfig {
	return BuildBaseDamageConfig(BaseDamageFuncRoll(minFlatDamage, maxFlatDamage))
}

func BaseDamageFuncMagic(minFlatDamage float64, maxFlatDamage float64, spellCoefficient float64) BaseDamageCalculator {
	if spellCoefficient == 0 {
		return BaseDamageFuncRoll(minFlatDamage, maxFlatDamage)
	}

	if minFlatDamage == 0 && maxFlatDamage == 0 {
		return func(_ *Simulation, hitEffect *SpellEffect, spell *Spell) float64 {
			return spell.SpellPower() * spellCoefficient
		}
	} else if minFlatDamage == maxFlatDamage {
		return func(sim *Simulation, hitEffect *SpellEffect, spell *Spell) float64 {
			damage := spell.SpellPower() * spellCoefficient
			return damage + minFlatDamage
		}
	} else {
		deltaDamage := maxFlatDamage - minFlatDamage
		return func(sim *Simulation, hitEffect *SpellEffect, spell *Spell) float64 {
			damage := spell.SpellPower() * spellCoefficient
			damage += damageRollOptimized(sim, minFlatDamage, deltaDamage)
			return damage
		}
	}
}
func BaseDamageConfigMagic(minFlatDamage float64, maxFlatDamage float64, spellCoefficient float64) BaseDamageConfig {
	return BuildBaseDamageConfig(BaseDamageFuncMagic(minFlatDamage, maxFlatDamage, spellCoefficient))
}
func BaseDamageConfigMagicNoRoll(flatDamage float64, spellCoefficient float64) BaseDamageConfig {
	return BaseDamageConfigMagic(flatDamage, flatDamage, spellCoefficient)
}

func BaseDamageFuncHealing(minFlatHealing float64, maxFlatHealing float64, spellCoefficient float64) BaseDamageCalculator {
	if spellCoefficient == 0 {
		return BaseDamageFuncRoll(minFlatHealing, maxFlatHealing)
	}

	if minFlatHealing == 0 && maxFlatHealing == 0 {
		return func(_ *Simulation, hitEffect *SpellEffect, spell *Spell) float64 {
			return spell.HealingPower() * spellCoefficient
		}
	} else if minFlatHealing == maxFlatHealing {
		return func(sim *Simulation, hitEffect *SpellEffect, spell *Spell) float64 {
			damage := spell.HealingPower() * spellCoefficient
			return damage + minFlatHealing
		}
	} else {
		deltaHealing := maxFlatHealing - minFlatHealing
		return func(sim *Simulation, hitEffect *SpellEffect, spell *Spell) float64 {
			damage := spell.HealingPower() * spellCoefficient
			damage += damageRollOptimized(sim, minFlatHealing, deltaHealing)
			return damage
		}
	}
}
func BaseDamageConfigHealing(minFlatHealing float64, maxFlatHealing float64, spellCoefficient float64) BaseDamageConfig {
	return BuildBaseDamageConfig(BaseDamageFuncHealing(minFlatHealing, maxFlatHealing, spellCoefficient))
}
func BaseDamageConfigHealingNoRoll(flatHealing float64, spellCoefficient float64) BaseDamageConfig {
	return BaseDamageConfigHealing(flatHealing, flatHealing, spellCoefficient)
}

func MultiplyByStacks(config BaseDamageConfig, aura *Aura) BaseDamageConfig {
	return WrapBaseDamageConfig(config, func(oldCalculator BaseDamageCalculator) BaseDamageCalculator {
		return func(sim *Simulation, hitEffect *SpellEffect, spell *Spell) float64 {
			return oldCalculator(sim, hitEffect, spell) * float64(aura.GetStacks())
		}
	})
}

func BaseDamageFuncMelee(minFlatDamage float64, maxFlatDamage float64, spellCoefficient float64) BaseDamageCalculator {
	if spellCoefficient == 0 {
		return BaseDamageFuncRoll(minFlatDamage, maxFlatDamage)
	}

	if minFlatDamage == 0 && maxFlatDamage == 0 {
		return func(_ *Simulation, hitEffect *SpellEffect, spell *Spell) float64 {
			return spellCoefficient * spell.MeleeAttackPower()
		}
	} else if minFlatDamage == maxFlatDamage {
		return func(sim *Simulation, hitEffect *SpellEffect, spell *Spell) float64 {
			return minFlatDamage + spellCoefficient*spell.MeleeAttackPower()
		}
	} else {
		deltaDamage := maxFlatDamage - minFlatDamage
		return func(sim *Simulation, hitEffect *SpellEffect, spell *Spell) float64 {
			damage := spellCoefficient * spell.MeleeAttackPower()
			damage += damageRollOptimized(sim, minFlatDamage, deltaDamage)
			return damage
		}
	}
}
func BaseDamageConfigMelee(minFlatDamage float64, maxFlatDamage float64, spellCoefficient float64) BaseDamageConfig {
	return BuildBaseDamageConfig(BaseDamageFuncMelee(minFlatDamage, maxFlatDamage, spellCoefficient))
}

type Hand bool

const MainHand Hand = true
const OffHand Hand = false

func BaseDamageFuncMeleeWeapon(hand Hand, normalized bool, flatBonus float64, includeBonusWeaponDamage bool) BaseDamageCalculator {
	// Bonus weapon damage applies after OH penalty: https://www.youtube.com/watch?v=bwCIU87hqTs
	if normalized {
		if hand == MainHand {
			return func(sim *Simulation, hitEffect *SpellEffect, spell *Spell) float64 {
				damage := spell.Unit.AutoAttacks.MH.CalculateNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
				damage = damage + flatBonus
				if includeBonusWeaponDamage {
					damage += spell.BonusWeaponDamage()
				}
				return damage
			}
		} else {
			return func(sim *Simulation, hitEffect *SpellEffect, spell *Spell) float64 {
				damage := spell.Unit.AutoAttacks.OH.CalculateNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
				damage = damage*0.5 + flatBonus
				if includeBonusWeaponDamage {
					damage += spell.BonusWeaponDamage()
				}
				return damage
			}
		}
	} else {
		if hand == MainHand {
			return func(sim *Simulation, hitEffect *SpellEffect, spell *Spell) float64 {
				damage := spell.Unit.AutoAttacks.MH.CalculateWeaponDamage(sim, spell.MeleeAttackPower())
				damage = damage + flatBonus
				if includeBonusWeaponDamage {
					damage += spell.BonusWeaponDamage()
				}
				return damage
			}
		} else {
			return func(sim *Simulation, hitEffect *SpellEffect, spell *Spell) float64 {
				damage := spell.Unit.AutoAttacks.OH.CalculateWeaponDamage(sim, spell.MeleeAttackPower())
				damage = damage*0.5 + flatBonus
				if includeBonusWeaponDamage {
					damage += spell.BonusWeaponDamage()
				}
				return damage
			}
		}
	}
}
func BaseDamageConfigMeleeWeapon(hand Hand, normalized bool, flatBonus float64, includeBonusWeaponDamage bool) BaseDamageConfig {
	calculator := BaseDamageFuncMeleeWeapon(hand, normalized, flatBonus, includeBonusWeaponDamage)
	return BuildBaseDamageConfig(calculator)
}

func BaseDamageFuncRangedWeapon(flatBonus float64) BaseDamageCalculator {
	return func(sim *Simulation, hitEffect *SpellEffect, spell *Spell) float64 {
		return spell.Unit.AutoAttacks.Ranged.CalculateWeaponDamage(sim, spell.RangedAttackPower(hitEffect.Target)) +
			flatBonus +
			spell.BonusWeaponDamage()
	}
}
func BaseDamageConfigRangedWeapon(flatBonus float64) BaseDamageConfig {
	return BuildBaseDamageConfig(BaseDamageFuncRangedWeapon(flatBonus))
}

func BaseDamageFuncEnemyWeapon(hand Hand) BaseDamageCalculator {
	if hand == MainHand {
		return func(sim *Simulation, hitEffect *SpellEffect, spell *Spell) float64 {
			ap := MaxFloat(0, spell.Unit.stats[stats.AttackPower])
			return spell.Unit.AutoAttacks.MH.EnemyWeaponDamage(sim, ap)
		}
	} else {
		return func(sim *Simulation, hitEffect *SpellEffect, spell *Spell) float64 {
			ap := MaxFloat(0, spell.Unit.stats[stats.AttackPower])
			return spell.Unit.AutoAttacks.MH.EnemyWeaponDamage(sim, ap) * 0.5
		}
	}
}
func BaseDamageConfigEnemyWeapon(hand Hand) BaseDamageConfig {
	return BuildBaseDamageConfig(BaseDamageFuncEnemyWeapon(hand))
}

// Performs an actual damage roll. Keep this internal because the 2nd parameter
// is the delta rather than maxDamage, which is error-prone.
func damageRollOptimized(sim *Simulation, minDamage float64, deltaDamage float64) float64 {
	return minDamage + deltaDamage*sim.RandomFloat("Damage Roll")
}

// For convenience, but try to use damageRollOptimized in most cases.
func DamageRoll(sim *Simulation, minDamage float64, maxDamage float64) float64 {
	return damageRollOptimized(sim, minDamage, maxDamage-minDamage)
}

func DamageRollFunc(minDamage float64, maxDamage float64) func(*Simulation) float64 {
	deltaDamage := maxDamage - minDamage
	return func(sim *Simulation) float64 {
		return damageRollOptimized(sim, minDamage, deltaDamage)
	}
}
