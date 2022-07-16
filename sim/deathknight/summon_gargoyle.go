package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) registerSummonGargoyleCD() {
	if !deathKnight.Talents.SummonGargoyle {
		return
	}

	summonGargoyleAura := deathKnight.RegisterAura(core.Aura{
		Label:    "Summon Gargoyle",
		ActionID: core.ActionID{SpellID: 49206},
		Duration: time.Second * 30,
	})

	baseCost := 60.0
	deathKnight.SummonGargoyle = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 49206},

		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    deathKnight.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			deathKnight.Gargoyle.EnableWithTimeout(sim, deathKnight.Gargoyle, time.Second*30)

			// Add % atack speed modifiers
			deathKnight.Gargoyle.MultiplyCastSpeed(deathKnight.PseudoStats.MeleeSpeedMultiplier)

			// Add a dummy aura to show in metrics
			summonGargoyleAura.Activate(sim)

			// Start casting after a short 1 second delay to simulate the summon animation
			// Might need tweaking after testing the exact possible delay
			pa := core.PendingAction{
				NextActionAt: sim.CurrentTime + time.Second*1,
				Priority:     core.ActionPriorityAuto,
				OnAction: func(s *core.Simulation) {
					deathKnight.Gargoyle.GargoyleStrike.Cast(sim, deathKnight.CurrentTarget)
				},
			}
			sim.AddPendingAction(&pa)
		},
	})

	deathKnight.AddMajorCooldown(core.MajorCooldown{
		Spell:    deathKnight.SummonGargoyle,
		Priority: core.CooldownPriorityDrums - 1, // Always prefer to cast after drums or lust so the gargoyle gets their benefits.
		Type:     core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			if deathKnight.Gargoyle.IsEnabled() {
				return false
			}
			if character.CurrentRunicPower() < deathKnight.SummonGargoyle.DefaultCast.Cost {
				return false
			}
			return true
		},
	})
}

type GargoylePet struct {
	core.Pet

	dkOwner *DeathKnight

	GargoyleStrike *core.Spell
}

func (deathKnight *DeathKnight) NewGargoyle() *GargoylePet {
	gargoyle := &GargoylePet{
		Pet: core.NewPet(
			"Gargoyle",
			&deathKnight.Character,
			gargoyleBaseStats,
			gargoyleStatInheritance,
			false,
		),
		dkOwner: deathKnight,
	}

	// NightOfTheDead
	gargoyle.PseudoStats.DamageTakenMultiplier *= (1.0 - float64(deathKnight.Talents.NightOfTheDead)*0.45)

	deathKnight.AddPet(gargoyle)

	return gargoyle
}

func (garg *GargoylePet) GetPet() *core.Pet {
	return &garg.Pet
}

func (garg *GargoylePet) Initialize() {
	garg.registerGargoyleStrikeSpell()
}

func (garg *GargoylePet) Reset(sim *core.Simulation) {
}

func (garg *GargoylePet) OnGCDReady(sim *core.Simulation) {
	// Gargoyle has no GCD on his cast so just do nothing here
	// else we get the error that this unit is not using its gcd
	garg.DoNothing()
}

// These numbers are just rough guesses
var gargoyleBaseStats = stats.Stats{
	stats.Stamina: 1000,
}

var gargoyleStatInheritance = func(ownerStats stats.Stats) stats.Stats {
	return stats.Stats{
		stats.AttackPower: ownerStats[stats.AttackPower],
		stats.SpellHit:    ownerStats[stats.SpellHit],
		stats.SpellHaste:  ownerStats[stats.MeleeHaste],
	}
}

func (garg *GargoylePet) registerGargoyleStrikeSpell() {
	attackPowerModifier := 0.3333333333333333 * (1.0 + 0.04*float64(garg.dkOwner.Talents.Impurity))

	garg.GargoyleStrike = garg.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 69520},
		SpellSchool: core.SpellSchoolNature,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: time.Millisecond * 1500,
			},
			OnCastComplete: func(sim *core.Simulation, spell *core.Spell) {
				// Gargoyle doesnt use GCD so we recast the spell over and over
				garg.GargoyleStrike.Cast(sim, garg.CurrentTarget)
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return 120 + hitEffect.MeleeAttackPower(spell.Unit)*attackPowerModifier
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: garg.OutcomeFuncCritFixedChance(0.05, 1.5),
		}),
	})
}
