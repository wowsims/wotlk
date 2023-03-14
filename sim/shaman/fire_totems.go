package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (shaman *Shaman) registerSearingTotemSpell() {
	var extraCastCondition core.CanCastCondition
	if shaman.Totems.Fire == proto.FireTotem_SearingTotem && shaman.Totems.UseFireMcd {
		extraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
			if shaman.Totems.Fire != proto.FireTotem_SearingTotem {
				return false
			}
			if shaman.SearingTotem.AOEDot().IsActive() || shaman.FireElemental.IsEnabled() || shaman.FireElementalTotem.IsReady(sim) {
				return false
			}
			return true
		}
	}

	shaman.SearingTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 58704},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       SpellFlagTotem,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.07,
			Multiplier: 1 -
				0.05*float64(shaman.Talents.TotemicFocus) -
				0.02*float64(shaman.Talents.MentalQuickness),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},
		ExtraCastCondition: extraCastCondition,

		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1 + float64(shaman.Talents.CallOfFlame)*0.05,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),

		Dot: core.DotConfig{
			IsAOE: true, // Not really AOE, but we don't want separate DoTs for each enemy.
			Aura: core.Aura{
				Label: "SearingTotem",
			},
			// These are the real tick values, but searing totem doesn't start its next
			// cast until the previous missile hits the target. We don't have an option
			// for target distance yet so just pretend the tick rate is lower.
			// https://wotlk.wowhead.com/spell=25530/attack
			//NumberOfTicks:        30,
			//TickLength:           time.Second * 2.2,
			NumberOfTicks: 24,
			TickLength:    time.Second * 60 / 24,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDamage := sim.Roll(90, 120) + 0.167*dot.Spell.SpellPower()
				dot.Spell.CalcAndDealDamage(sim, target, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			shaman.MagmaTotem.AOEDot().Cancel(sim)
			shaman.FireElemental.Disable(sim)
			spell.AOEDot().Apply(sim)
			if !shaman.Totems.UseFireMcd {
				// +1 needed because of rounding issues with totem tick time.
				shaman.NextTotemDrops[FireTotem] = sim.CurrentTime + time.Second*60 + 1
			}
		},
	})

	if extraCastCondition == nil {
		return
	}
	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell:    shaman.SearingTotem,
		Priority: core.CooldownPriorityDefault, // TODO needs to be altered due to snap shotting.
		Type:     core.CooldownTypeDPS,
	})
}

func (shaman *Shaman) registerMagmaTotemSpell() {
	var extraCastCondition core.CanCastCondition
	if shaman.Totems.Fire == proto.FireTotem_MagmaTotem && shaman.Totems.UseFireMcd {
		extraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
			if shaman.Totems.Fire != proto.FireTotem_MagmaTotem {
				return false
			}
			if shaman.MagmaTotem.AOEDot().IsActive() || shaman.FireElemental.IsEnabled() || shaman.FireElementalTotem.IsReady(sim) {
				return false
			}
			return true
		}
	}

	shaman.MagmaTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 58734},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       SpellFlagTotem,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.27,
			Multiplier: 1 -
				0.05*float64(shaman.Talents.TotemicFocus) -
				0.02*float64(shaman.Talents.MentalQuickness),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},
		ExtraCastCondition: extraCastCondition,

		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1 + float64(shaman.Talents.CallOfFlame)*0.05,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "MagmaTotem",
			},
			NumberOfTicks: 10,
			TickLength:    time.Second * 2,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDamage := 371 + 0.1*dot.Spell.SpellPower()
				baseDamage *= sim.Encounter.AOECapMultiplier()
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.Spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			shaman.SearingTotem.AOEDot().Cancel(sim)
			shaman.FireElemental.Disable(sim)
			spell.AOEDot().Apply(sim)
			if !shaman.Totems.UseFireMcd {
				// +1 needed because of rounding issues with totem tick time.
				shaman.NextTotemDrops[FireTotem] = sim.CurrentTime + time.Second*20 + 1
			}
		},
	})

	if extraCastCondition == nil {
		return
	}
	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell:    shaman.MagmaTotem,
		Priority: core.CooldownPriorityDefault, // TODO needs to be altered due to snap shotting.
		Type:     core.CooldownTypeDPS,
	})
}
