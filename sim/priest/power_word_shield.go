package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerPowerWordShieldSpell() {
	actionID := core.ActionID{SpellID: 48066}
	baseCost := 0.23 * priest.BaseMana
	coeff := 0.8057 + 0.08*float64(priest.Talents.BorrowedTime)

	priest.PowerWordShield = priest.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 -
					[]float64{0, .04, .07, .10}[priest.Talents.MentalAgility] -
					core.TernaryFloat64(priest.Talents.SoulWarding, .15, 0)),
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			weakenedSoul := priest.WeakenedSouls[target.UnitIndex]
			if weakenedSoul.IsActive() {
				panic("Cannot cast PWS on target with Weakened Soul!")
			}

			shieldAmount := 2230.0 + priest.GetStat(stats.HealingPower)*coeff
			shield := priest.PWSShields[target.UnitIndex]
			shield.Apply(sim, shieldAmount)
			weakenedSoul.Activate(sim)
		},
	})

	priest.PWSShields = make([]*core.Shield, len(priest.Env.AllUnits))
	priest.WeakenedSouls = make([]*core.Aura, len(priest.Env.AllUnits))
	for _, unit := range priest.Env.AllUnits {
		if !priest.IsOpponent(unit) {
			priest.PWSShields[unit.UnitIndex] = priest.makePWSShield(unit)
			priest.WeakenedSouls[unit.UnitIndex] = priest.makeWeakenedSoul(unit)
		}
	}
}

func (priest *Priest) makePWSShield(target *core.Unit) *core.Shield {
	return core.NewShield(core.Shield{
		Spell: priest.PowerWordShield,
		Aura: target.GetOrRegisterAura(core.Aura{
			Label:    "Power Word Shield",
			ActionID: priest.PowerWordShield.ActionID,
			Duration: time.Second * 30,
		}),
	})
}

func (priest *Priest) makeWeakenedSoul(target *core.Unit) *core.Aura {
	return target.GetOrRegisterAura(core.Aura{
		Label:    "Weakened Soul",
		ActionID: core.ActionID{SpellID: 6788},
		Duration: time.Second * 15,
	})
}

func (priest *Priest) CanCastPWS(target *core.Unit) bool {
	return !priest.WeakenedSouls[target.UnitIndex].IsActive()
}
