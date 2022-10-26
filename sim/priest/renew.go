package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerRenewSpell() {
	actionID := core.ActionID{SpellID: 48068}
	baseCost := 0.17 * priest.BaseMana
	spellCoeff := priest.renewSpellCoefficient()

	if priest.Talents.EmpoweredRenew > 0 {
		priest.EmpoweredRenew = priest.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 63543},
			SpellSchool: core.SpellSchoolHoly,
			ProcMask:    core.ProcMaskSpellHealing,
			Flags:       core.SpellFlagNoOnCastComplete,

			BonusCritRating: float64(priest.Talents.HolySpecialization) * 1 * core.CritRatingPerCritChance,
			DamageMultiplier: 1 *
				float64(priest.renewTicks()) *
				priest.renewHealingMultiplier() *
				.05 * float64(priest.Talents.EmpoweredRenew) *
				core.TernaryFloat64(priest.HasSetBonus(ItemSetZabrasRaiment, 4), 1.1, 1),
			CritMultiplier:   priest.DefaultHealingCritMultiplier(),
			ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseHealing := 280 + spellCoeff*spell.HealingPower()
				spell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)
			},
		})
	}

	priest.Renew = priest.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellHealing,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - []float64{0, .04, .07, .10}[priest.Talents.MentalAgility]),
				GCD:  core.GCDDefault,
			},
		},

		DamageMultiplier: priest.renewHealingMultiplier(),
		ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.SpellMetrics[target.UnitIndex].Hits++
			priest.RenewHots[target.UnitIndex].Apply(sim)

			if priest.EmpoweredRenew != nil {
				priest.EmpoweredRenew.Cast(sim, target)
			}
		},
	})

	priest.RenewHots = make([]*core.Dot, len(priest.Env.AllUnits))
	for _, unit := range priest.Env.AllUnits {
		if !priest.IsOpponent(unit) {
			priest.RenewHots[unit.UnitIndex] = priest.makeRenewHot(unit)
		}
	}
}

func (priest *Priest) makeRenewHot(target *core.Unit) *core.Dot {
	healingCoeff := priest.renewSpellCoefficient()
	return core.NewDot(core.Dot{
		Spell: priest.Renew,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Renew" + strconv.Itoa(int(priest.Index)),
			ActionID: priest.Renew.ActionID,
		}),
		NumberOfTicks: priest.renewTicks(),
		TickLength:    time.Second * 3,
		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			dot.SnapshotBaseDamage = 280 + healingCoeff*dot.Spell.HealingPower()
			dot.SnapshotAttackerMultiplier = dot.Spell.CasterHealingMultiplier()
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeTick)
		},
	})
}

func (priest *Priest) renewTicks() int {
	return 5 - core.TernaryInt(priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfRenew), 1, 0)
}

func (priest *Priest) renewHealingMultiplier() float64 {
	return 1 *
		(1 + 0.01*float64(priest.Talents.TwinDisciplines)) *
		(1 + 0.05*float64(priest.Talents.ImprovedRenew)) *
		core.TernaryFloat64(priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfRenew), 1.25, 1)
}

func (priest *Priest) renewSpellCoefficient() float64 {
	return (1.88 + .05*float64(priest.Talents.EmpoweredRenew)) / 5
}
