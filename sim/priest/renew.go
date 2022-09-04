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

	priest.Renew = priest.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - []float64{0, .04, .07, .10}[priest.Talents.MentalAgility]),
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			priest.RenewHots[target.UnitIndex].Apply(sim)
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
	return core.NewDot(core.Dot{
		Spell: priest.Renew,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Renew" + strconv.Itoa(int(priest.Index)),
			ActionID: priest.Renew.ActionID,
		}),
		NumberOfTicks: 5 - core.TernaryInt(priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfRenew), 1, 0),
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:   core.ProcMaskPeriodicDamage,
			IsPeriodic: true,
			IsHealing:  true,

			DamageMultiplier: 1 *
				(1 + 0.01*float64(priest.Talents.TwinDisciplines)) *
				(1 + 0.05*float64(priest.Talents.ImprovedRenew)) *
				core.TernaryFloat64(priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfRenew), 1.25, 1),
			ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

			BaseDamage:     core.BaseDamageConfigHealingNoRoll(280, (1.88+.05*float64(priest.Talents.EmpoweredRenew))/5),
			OutcomeApplier: priest.OutcomeFuncTick(),
			//OnPeriodicDamageDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			//},
		}),
	})
}
