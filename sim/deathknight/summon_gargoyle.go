package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerSummonGargoyleCD() {
	if !dk.Talents.SummonGargoyle {
		return
	}

	summonGargoyleAura := dk.RegisterAura(core.Aura{
		Label:    "Summon Gargoyle",
		ActionID: core.ActionID{SpellID: 49206},
		Duration: time.Second * 30,
	})

	baseCost := float64(core.NewRuneCost(60.0, 0, 0, 0, 0))
	dk.SummonGargoyle = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID: core.ActionID{SpellID: 49206},

		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: baseCost,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			dk.Gargoyle.EnableWithTimeout(sim, dk.Gargoyle, time.Second*30)
			dk.Gargoyle.CancelGCDTimer(sim)

			// Add % atack speed modifiers
			dk.Gargoyle.PseudoStats.CastSpeedMultiplier = 1.0
			dk.Gargoyle.MultiplyCastSpeed(dk.PseudoStats.MeleeSpeedMultiplier)

			// Add a dummy aura to show in metrics
			summonGargoyleAura.Activate(sim)

			// Start casting after a short 1 second delay to simulate the summon animation
			// Might need tweaking after testing the exact possible delay
			pa := core.PendingAction{
				NextActionAt: sim.CurrentTime + time.Second*1,
				Priority:     core.ActionPriorityAuto,
				OnAction: func(s *core.Simulation) {
					dk.Gargoyle.GargoyleStrike.Cast(sim, dk.CurrentTarget)
				},
			}
			sim.AddPendingAction(&pa)
		},
	})
}

func (dk *Deathknight) CanSummonGargoyle(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 60.0, 0, 0, 0) && dk.SummonGargoyle.IsReady(sim)
}

func (dk *Deathknight) CastSummonGargoyle(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanSummonGargoyle(sim) {
		res := dk.SummonGargoyle.Cast(sim, target)
		if res {
			dk.UpdateMajorCooldowns()
		}
		return res
	}
	return false
}

type GargoylePet struct {
	core.Pet

	dkOwner *Deathknight

	GargoyleStrike *core.Spell
}

func (dk *Deathknight) NewGargoyle() *GargoylePet {
	gargoyle := &GargoylePet{
		Pet: core.NewPet(
			"Gargoyle",
			&dk.Character,
			gargoyleBaseStats,
			gargoyleStatInheritance,
			false,
			true,
		),
		dkOwner: dk,
	}

	// NightOfTheDead
	gargoyle.PseudoStats.DamageTakenMultiplier *= (1.0 - float64(dk.Talents.NightOfTheDead)*0.45)

	dk.AddPet(gargoyle)

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
					return 130 + hitEffect.MeleeAttackPower(spell.Unit)*attackPowerModifier
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: garg.OutcomeFuncCritFixedChance(0.05, 1.5),
		}),
	})
}
