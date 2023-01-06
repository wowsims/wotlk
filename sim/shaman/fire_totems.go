package shaman

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (shaman *Shaman) registerSearingTotemSpell() {
	actionID := core.ActionID{SpellID: 58704}
	baseCost := 0.07 * shaman.BaseMana

	shaman.SearingTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskEmpty,
		Flags:        SpellFlagTotem,
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

		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1 + float64(shaman.Talents.CallOfFlame)*0.05,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			shaman.MagmaTotemDot.Cancel(sim)
			shaman.FireElemental.Disable(sim)
			shaman.SearingTotemDot.Apply(sim)
			if !shaman.Totems.UseFireMcd {
				// +1 needed because of rounding issues with totem tick time.
				shaman.NextTotemDrops[FireTotem] = sim.CurrentTime + time.Second*60 + 1
			}
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
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			baseDamage := sim.Roll(90, 120) + 0.167*dot.Spell.SpellPower()
			dot.Spell.CalcAndDealDamage(sim, target, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
		},
	})

	if shaman.Totems.Fire != proto.FireTotem_SearingTotem || !shaman.Totems.UseFireMcd {
		return
	}
	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell:    shaman.SearingTotem,
		Priority: core.CooldownPriorityDefault, // TODO needs to be altered due to snap shotting.
		Type:     core.CooldownTypeDPS,
		CanActivate: func(s *core.Simulation, c *core.Character) bool {
			if shaman.Totems.Fire != proto.FireTotem_SearingTotem {
				return false
			}
			if shaman.SearingTotemDot.IsActive() || shaman.FireElemental.IsEnabled() || shaman.FireElementalTotem.IsReady(s) {
				return false
			}
			return true
		},
	})
}

func (shaman *Shaman) registerMagmaTotemSpell() {
	actionID := core.ActionID{SpellID: 58734}
	baseCost := 0.27 * shaman.BaseMana

	shaman.MagmaTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskEmpty,
		Flags:        SpellFlagTotem,
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

		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1 + float64(shaman.Talents.CallOfFlame)*0.05,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			shaman.SearingTotemDot.Cancel(sim)
			shaman.FireElemental.Disable(sim)
			shaman.MagmaTotemDot.Apply(sim)
			if !shaman.Totems.UseFireMcd {
				// +1 needed because of rounding issues with totem tick time.
				shaman.NextTotemDrops[FireTotem] = sim.CurrentTime + time.Second*20 + 1
			}
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

		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			baseDamage := 371 + 0.1*dot.Spell.SpellPower()
			baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.Targets {
				dot.Spell.CalcAndDealDamage(sim, &aoeTarget.Unit, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	if shaman.Totems.Fire != proto.FireTotem_MagmaTotem || !shaman.Totems.UseFireMcd {
		return // don't add magma totem to the CDs
	}

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell:    shaman.MagmaTotem,
		Priority: core.CooldownPriorityDefault, // TODO needs to be altered due to snap shotting.
		Type:     core.CooldownTypeDPS,
		CanActivate: func(s *core.Simulation, c *core.Character) bool {
			if shaman.Totems.Fire != proto.FireTotem_MagmaTotem {
				return false
			}
			if shaman.MagmaTotemDot.IsActive() || shaman.FireElemental.IsEnabled() || shaman.FireElementalTotem.IsReady(s) {
				return false
			}
			return true
		},
	})
}
