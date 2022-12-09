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
				cast.GCD = dk.GetModifiedGCD()
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
	}, func(sim *core.Simulation) bool {
		return dk.CastCostPossible(sim, 60.0, 0, 0, 0) && dk.SummonGargoyle.IsReady(sim)
	}, func(sim *core.Simulation) {
		dk.UpdateMajorCooldowns()
	})
}

type GargoylePet struct {
	core.Pet

	dkOwner *Deathknight

	GargoyleStrike *core.Spell
}

func (dk *Deathknight) NewGargoyle() *GargoylePet {
	// Remove any hit that would be given by NocS as it does not translate to pets
	nocsHit := 0.0
	if dk.nervesOfColdSteelActive() {
		nocsHit = (float64(dk.Talents.NervesOfColdSteel) / 8.0) * 17.0
	}
	gargoyle := &GargoylePet{
		Pet: core.NewPet(
			"Gargoyle",
			&dk.Character,
			stats.Stats{
				stats.Stamina:  1000,
				stats.SpellHit: -nocsHit * core.SpellHitRatingPerHitChance,
			},
			func(ownerStats stats.Stats) stats.Stats {
				// Convert dk melee hit to garg spell hit
				// We convert 8 melee hit to 17 spell hit (as thats how pets scale their hit/expertise)
				ownerHitChance := (ownerStats[stats.MeleeHit] / core.MeleeHitRatingPerHitChance)
				hitRatingFromOwner := ((ownerHitChance / 8.0) * 17.0) * core.SpellHitRatingPerHitChance

				return stats.Stats{
					stats.AttackPower: ownerStats[stats.AttackPower],
					stats.SpellHit:    hitRatingFromOwner,
					stats.SpellHaste:  (ownerStats[stats.MeleeHaste] / dk.PseudoStats.MeleeHasteRatingPerHastePercent) * core.HasteRatingPerHastePercent,
				}
			},
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

func (garg *GargoylePet) registerGargoyleStrikeSpell() {
	attackPowerModifier := (1.0 + 0.04*float64(garg.dkOwner.Talents.Impurity)) / 3.0
	var outcomeApplier core.OutcomeApplier

	garg.GargoyleStrike = garg.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 51963},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: time.Millisecond * 2000,
			},
			OnCastComplete: func(sim *core.Simulation, spell *core.Spell) {
				// Gargoyle doesn't use GCD, so we recast the spell over and over
				garg.GargoyleStrike.Cast(sim, garg.CurrentTarget)
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1.5,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 2.05*sim.Roll(51, 69) + attackPowerModifier*spell.MeleeAttackPower()
			result := spell.CalcDamage(sim, target, baseDamage, outcomeApplier)
			spell.DealDamage(sim, result)
		},
	})
	outcomeApplier = garg.GargoyleStrike.OutcomeMagicCritFixedChance(0.05)
}
