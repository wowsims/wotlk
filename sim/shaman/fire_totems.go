package shaman

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (shaman *Shaman) registerSearingTotemSpell() {
	actionID := core.ActionID{SpellID: 58704}
	baseCost := baseMana * 0.07

	shaman.SearingTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		Flags:       SpellFlagTotem,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost -
					baseCost*float64(shaman.Talents.TotemicFocus)*0.05 -
					baseCost*float64(shaman.Talents.MentalQuickness)*0.02,
				GCD: time.Second,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			shaman.SearingTotemDot.Apply(sim)
			// +1 needed because of rounding issues with Searing totem tick time.
			shaman.NextTotemDrops[FireTotem] = sim.CurrentTime + time.Second*60 + 1
		},
	})

	target := shaman.CurrentTarget
	shaman.SearingTotemDot = core.NewDot(core.Dot{
		Spell: shaman.SearingTotem,
		Aura: target.RegisterAura(core.Aura{
			Label:    "SearingTotem-" + strconv.Itoa(int(shaman.Index)),
			ActionID: actionID,
		}),
		// These are the real tick values, but searing totem doesn't start its next
		// cast until the previous missile hits the target. We don't have an option
		// for target distance yet so just pretend the tick rate is lower.
		// https://wotlk.wowhead.com/spell=25530/attack
		//NumberOfTicks:        30,
		//TickLength:           time.Second * 2.2,
		NumberOfTicks: 24,
		TickLength:    time.Second * 60 / 24,
		TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:            core.ProcMaskEmpty,
			BonusSpellHitRating: float64(shaman.Talents.ElementalPrecision) * 2 * core.SpellHitRatingPerHitChance,
			DamageMultiplier:    1 + float64(shaman.Talents.CallOfFlame)*0.05,
			BaseDamage:          core.BaseDamageConfigMagic(50, 66, 0.167),
			OutcomeApplier:      shaman.OutcomeFuncMagicHitAndCrit(shaman.ElementalCritMultiplier(0)),
		})),
	})
}

func (shaman *Shaman) registerMagmaTotemSpell() {
	actionID := core.ActionID{SpellID: 58734}
	baseCost := baseMana * 0.27

	shaman.MagmaTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		Flags:       SpellFlagTotem,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost -
					baseCost*float64(shaman.Talents.TotemicFocus)*0.05 -
					baseCost*float64(shaman.Talents.MentalQuickness)*0.02,
				GCD: time.Second,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			shaman.MagmaTotemDot.Apply(sim)
			shaman.NextTotemDrops[FireTotem] = sim.CurrentTime + time.Second*20 + 1
		},
	})

	target := shaman.CurrentTarget
	shaman.MagmaTotemDot = core.NewDot(core.Dot{
		Spell: shaman.MagmaTotem,
		Aura: target.RegisterAura(core.Aura{
			Label:    "MagmaTotem-" + strconv.Itoa(int(shaman.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 10,
		TickLength:    time.Second * 2,
		TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncAOEDamageCapped(shaman.Env, 1550, core.SpellEffect{
			ProcMask:            core.ProcMaskEmpty,
			BonusSpellHitRating: float64(shaman.Talents.ElementalPrecision) * 2 * core.SpellHitRatingPerHitChance,
			DamageMultiplier:    1 + float64(shaman.Talents.CallOfFlame)*0.05,
			// TODO: find magma totem sp coeff
			BaseDamage:     core.BaseDamageConfigMagicNoRoll(371, 0.067),
			OutcomeApplier: shaman.OutcomeFuncMagicHitAndCrit(shaman.ElementalCritMultiplier(0)),
		})),
	})
}

// func (shaman *Shaman) FireNovaTickLength() time.Duration {
// 	return time.Second * time.Duration(4-shaman.Talents.ImprovedFireTotems)
// }

// func (shaman *Shaman) registerNovaTotemSpell() {
// 	actionID := core.ActionID{SpellID: 25537}
// 	baseCost := 765.0

// 	tickLength := shaman.FireNovaTickLength()
// 	shaman.FireNovaTotem = shaman.RegisterSpell(core.SpellConfig{
// 		ActionID:    actionID,
// 		SpellSchool: core.SpellSchoolFire,
// 		Flags:       SpellFlagTotem,

// 		ResourceType: stats.Mana,
// 		BaseCost:     baseCost,

// 		Cast: core.CastConfig{
// 			DefaultCast: core.Cast{
// 				Cost: baseCost -
// 					baseCost*float64(shaman.Talents.TotemicFocus)*0.05 -
// 					baseCost*float64(shaman.Talents.MentalQuickness)*0.02,
// 				GCD: time.Second,
// 			},
// 			CD: core.Cooldown{
// 				Timer:    shaman.NewTimer(),
// 				Duration: time.Second * 15,
// 			},
// 		},

// 		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
// 			shaman.MagmaTotemDot.Cancel(sim)
// 			shaman.SearingTotemDot.Cancel(sim)
// 			shaman.FireNovaTotemDot.Apply(sim)
// 			shaman.NextTotemDrops[FireTotem] = sim.CurrentTime + tickLength + 1
// 			shaman.tryTwistFireNova(sim)
// 		},
// 	})

// 	target := shaman.CurrentTarget
// 	shaman.FireNovaTotemDot = core.NewDot(core.Dot{
// 		Spell: shaman.FireNovaTotem,
// 		Aura: target.RegisterAura(core.Aura{
// 			Label:    "FireNovaTotem-" + strconv.Itoa(int(shaman.Index)),
// 			ActionID: actionID,
// 		}),
// 		NumberOfTicks: 1,
// 		TickLength:    tickLength,
// 		TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncAOEDamageCapped(shaman.Env, 9975, core.SpellEffect{
// 			ProcMask:            core.ProcMaskEmpty,
// 			BonusSpellHitRating: float64(shaman.Talents.ElementalPrecision) * 2 * core.SpellHitRatingPerHitChance,
// 			DamageMultiplier:    1 + float64(shaman.Talents.CallOfFlame)*0.05,
// 			BaseDamage:          core.BaseDamageConfigMagic(654, 730, 0.214),
// 			OutcomeApplier:      shaman.OutcomeFuncMagicHitAndCrit(shaman.ElementalCritMultiplier()),
// 		})),
// 	})
// }
