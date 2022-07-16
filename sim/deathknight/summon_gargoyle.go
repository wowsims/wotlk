package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// The numbers in this file are VERY rough approximations based on logs.

func (deathknight *DeathKnight) registerSummonGargoyleCD() {
	if !deathknight.Talents.SummonGargoyle {
		return
	}

	baseCost := 60.0
	deathknight.SummonGargoyle = deathknight.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 31687},

		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    deathknight.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			deathknight.Gargoyle.EnableWithTimeout(sim, deathknight.Gargoyle, time.Second*30)
		},
	})

	deathknight.AddMajorCooldown(core.MajorCooldown{
		Spell:    deathknight.SummonGargoyle,
		Priority: core.CooldownPriorityDrums - 1, // Always prefer to cast after drums or lust so the gargoyle gets their benefits.
		Type:     core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			if deathknight.Gargoyle.IsEnabled() {
				return false
			}
			if character.CurrentRunicPower() < deathknight.SummonGargoyle.DefaultCast.Cost {
				return false
			}
			return true
		},
	})
}

type GargoylePet struct {
	core.Pet

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
	}
	//gargoyle.EnableManaBar()

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
	garg.GargoyleStrike.Cast(sim, garg.CurrentTarget)
}

// These numbers are just rough guesses based on looking at some logs.
var gargoyleBaseStats = stats.Stats{
	stats.Intellect:  100,
	stats.SpellPower: 300,
	stats.SpellHit:   3 * core.SpellHitRatingPerHitChance,
	stats.SpellCrit:  8 * core.CritRatingPerCritChance,
}

var gargoyleStatInheritance = func(ownerStats stats.Stats) stats.Stats {
	return ownerStats.DotProduct(stats.Stats{
		stats.MeleeHaste: ownerStats[stats.MeleeHaste],
		stats.SpellHaste: ownerStats[stats.MeleeHaste],
	})
}

func (garg *GargoylePet) registerGargoyleStrikeSpell() {
	garg.GargoyleStrike = garg.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 69520},
		SpellSchool: core.SpellSchoolNature,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      0,
				CastTime: time.Millisecond * 1500,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagic(256, 328, 1),
			OutcomeApplier:   garg.OutcomeFuncAlwaysHit(),
		}),
	})
}
