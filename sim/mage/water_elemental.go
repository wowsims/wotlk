package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// The numbers in this file are VERY rough approximations based on logs.

func (mage *Mage) registerSummonWaterElementalCD() {
	if !mage.Talents.SummonWaterElemental {
		return
	}

	if mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfEternalWater) {
		// Makes pet permanent, so doesn't use the CD.
		return
	}

	baseCost := mage.BaseMana * 0.16
	summonDuration := time.Second*45 + time.Second*5*time.Duration(mage.Talents.EnduringWinter)
	mage.SummonWaterElemental = mage.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 31687},

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer: mage.NewTimer(),
				Duration: time.Duration(float64(time.Minute*3)*(1-0.1*float64(mage.Talents.ColdAsIce))) -
					core.TernaryDuration(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfWaterElemental), time.Second*30, 0),
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.waterElemental.EnableWithTimeout(sim, mage.waterElemental, summonDuration)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell:    mage.SummonWaterElemental,
		Priority: core.CooldownPriorityDrums + 1, // Always prefer to cast before drums or lust so the ele gets their benefits.
		Type:     core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			if mage.waterElemental.IsEnabled() {
				return false
			}
			if character.CurrentMana() < mage.SummonWaterElemental.DefaultCast.Cost {
				return false
			}
			return true
		},
	})
}

type WaterElemental struct {
	core.Pet

	// Water Ele almost never just stands still and spams like we want, it sometimes
	// does its own thing. This controls how much it does that.
	disobeyChance float64

	Waterbolt *core.Spell
}

func (mage *Mage) NewWaterElemental(disobeyChance float64) *WaterElemental {
	waterElemental := &WaterElemental{
		Pet: core.NewPet(
			"Water Elemental",
			&mage.Character,
			waterElementalBaseStats,
			waterElementalStatInheritance,
			mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfEternalWater),
			true,
		),
		disobeyChance: disobeyChance,
	}
	waterElemental.EnableManaBarWithModifier(0.333)

	mage.AddPet(waterElemental)

	return waterElemental
}

func (we *WaterElemental) GetPet() *core.Pet {
	return &we.Pet
}

func (we *WaterElemental) Initialize() {
	we.registerWaterboltSpell()
}

func (we *WaterElemental) Reset(sim *core.Simulation) {
}

func (we *WaterElemental) OnGCDReady(sim *core.Simulation) {
	spell := we.Waterbolt

	if sim.RandomFloat("Water Elemental Disobey") < we.disobeyChance {
		// Water ele has decided not to cooperate, so just wait for the cast time
		// instead of casting.
		we.WaitUntil(sim, sim.CurrentTime+spell.DefaultCast.CastTime)
		return
	}

	if success := spell.Cast(sim, we.CurrentTarget); !success {
		we.WaitForMana(sim, spell.CurCast.Cost)
	}
}

// These numbers are just rough guesses based on looking at some logs.
var waterElementalBaseStats = stats.Stats{
	stats.Mana:      1082,
	stats.Intellect: 369,
}

var waterElementalStatInheritance = func(ownerStats stats.Stats) stats.Stats {
	// These numbers are just rough guesses based on looking at some logs.
	return ownerStats.DotProduct(stats.Stats{
		stats.Stamina:   0.2238,
		stats.Intellect: 0.3,

		stats.SpellPower: 0.333,
		stats.SpellHit:   0.01,
		stats.SpellCrit:  0.01,
	})
}

func (we *WaterElemental) registerWaterboltSpell() {
	baseCost := we.BaseMana * 0.01

	we.Waterbolt = we.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 31707},
		SpellSchool:  core.SpellSchoolFrost,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2500,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   we.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(601, 673) + (2.5/3.0)*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
