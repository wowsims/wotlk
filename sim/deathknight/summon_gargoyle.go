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

	dk.SummonGargoyle = dk.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 49206},

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: 60,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			dk.Gargoyle.EnableWithTimeout(sim, dk.Gargoyle, time.Second*30)
			dk.Gargoyle.CancelGCDTimer(sim)

			snapshottedMeleeSpeedMultipler := dk.PseudoStats.MeleeSpeedMultiplier
			dk.Gargoyle.meleeSpeedMultiplier = func() float64 {
				if dk.Gargoyle.isNerfedGargoyle {
					return dk.PseudoStats.MeleeSpeedMultiplier
				}
				return snapshottedMeleeSpeedMultipler
			}
			dk.Gargoyle.updateCastSpeed()

			// Add a dummy aura to show in metrics
			summonGargoyleAura.Activate(sim)

			// Start casting after a 2.5s delay to simulate the summon animation
			pa := core.PendingAction{
				NextActionAt: sim.CurrentTime + dk.GargoyleSummonDelay,
				Priority:     core.ActionPriorityAuto,
				OnAction: func(s *core.Simulation) {
					dk.OnGargoyleStartFirstCast()
					dk.Gargoyle.GargoyleStrike.Cast(sim, dk.CurrentTarget)
				},
			}
			sim.AddPendingAction(&pa)
		},
	})

	dk.AddMajorCooldown(core.MajorCooldown{
		Spell: dk.SummonGargoyle,
		Type:  core.CooldownTypeDPS,
	})
	if dk.Inputs.IsDps {
		// We use this for defining the min cast time of gargoyle
		// but we dont cast it with the MCD system in the dps sim
		dk.GetMajorCooldown(dk.SummonGargoyle.ActionID).Disable()
	}
}

type GargoylePet struct {
	core.Pet

	dkOwner *Deathknight

	GargoyleStrike *core.Spell

	ownerMeleeMultiplier float64
	meleeSpeedMultiplier func() float64
	isNerfedGargoyle     bool
}

func (dk *Deathknight) NewGargoyle(nerfedGargoyle bool) *GargoylePet {
	// Remove any hit that would be given by NocS as it does not translate to pets
	nocsHit := 0.0
	if dk.nervesOfColdSteelActive() {
		nocsHit = (float64(dk.Talents.NervesOfColdSteel) / 8.0) * 17.0
	}
	if dk.HasDraeneiHitAura {
		nocsHit = nocsHit + 1.0
	}

	var gargoyleDynamicStatInheritance core.PetStatInheritance = nil
	if nerfedGargoyle {
		gargoyleDynamicStatInheritance = func(ownerStats stats.Stats) stats.Stats {
			return stats.Stats{
				stats.SpellHaste: (ownerStats[stats.MeleeHaste] / dk.PseudoStats.MeleeHasteRatingPerHastePercent) * core.HasteRatingPerHastePercent,
			}
		}
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
			gargoyleDynamicStatInheritance,
			false,
			true,
		),
		dkOwner:              dk,
		isNerfedGargoyle:     nerfedGargoyle,
		ownerMeleeMultiplier: 1.0,
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
	garg.ownerMeleeMultiplier = 1.0
}

func (garg *GargoylePet) OnGCDReady(sim *core.Simulation) {
	// Gargoyle has no GCD on his cast so just do nothing here
	// else we get the error that this unit is not using its gcd
	garg.DoNothing()
}

func (garg *GargoylePet) updateCastSpeed() {
	garg.MultiplyCastSpeed(1.0 / garg.ownerMeleeMultiplier)
	garg.ownerMeleeMultiplier = garg.meleeSpeedMultiplier()
	garg.MultiplyCastSpeed(garg.ownerMeleeMultiplier)
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
				garg.updateCastSpeed()
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
