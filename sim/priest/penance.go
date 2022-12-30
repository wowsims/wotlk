package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerPenanceHealSpell() {
	priest.PenanceHeal = priest.makePenanceSpell(true)
}

func (priest *Priest) RegisterPenanceSpell() {
	priest.Penance = priest.makePenanceSpell(false)
}

func (priest *Priest) makePenanceSpell(isHeal bool) *core.Spell {
	if !priest.Talents.Penance {
		return nil
	}

	actionID := core.ActionID{SpellID: 53007}
	baseCost := priest.BaseMana * 0.16

	penanceDots := make([]*core.Dot, len(priest.Env.AllUnits))

	var procMask core.ProcMask
	flags := core.SpellFlagChanneled
	if isHeal {
		flags |= core.SpellFlagHelpful
		procMask = core.ProcMaskSpellHealing
	} else {
		procMask = core.ProcMaskSpellDamage
	}

	spell := priest.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolHoly,
		ProcMask:     procMask,
		Flags:        flags,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					(1 - 0.05*float64(priest.Talents.ImprovedHealing)) *
					(1 - []float64{0, .04, .07, .10}[priest.Talents.MentalAgility]),
				GCD:         core.GCDDefault,
				ChannelTime: time.Second * 2,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Duration(float64(time.Second*12-core.TernaryDuration(priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfPenance), time.Second*2, 0)) * (1 - .1*float64(priest.Talents.Aspiration))),
			},
		},

		DamageMultiplier: 1 +
			core.TernaryFloat64(isHeal,
				1*
					(1+.02*float64(priest.Talents.SpiritualHealing))*
					(1+.01*float64(priest.Talents.BlessedResilience))*
					(1+.02*float64(priest.Talents.FocusedPower)),
				.05*float64(priest.Talents.SearingLight)) +
			.01*float64(priest.Talents.TwinDisciplines),
		CritMultiplier:   core.TernaryFloat64(isHeal, priest.DefaultHealingCritMultiplier(), priest.DefaultSpellCritMultiplier()),
		ThreatMultiplier: 0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if isHeal {
				spell.SpellMetrics[target.UnitIndex].Hits--
				hot := penanceDots[target.UnitIndex]
				hot.Apply(sim)
				// Do immediate tick
				hot.TickOnce(sim)
			} else {
				result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
				if result.Landed() {
					spell.SpellMetrics[target.UnitIndex].Hits--
					dot := penanceDots[target.UnitIndex]
					dot.Apply(sim)
					// Do immediate tick
					dot.TickOnce(sim)
				}
				spell.DealOutcome(sim, result)
			}
		},
	})

	for _, unit := range priest.Env.AllUnits {
		penanceDots[unit.UnitIndex] = priest.makePenanceDotOrHot(unit, spell, isHeal)
	}

	return spell
}

func (priest *Priest) makePenanceDotOrHot(target *core.Unit, spell *core.Spell, isHeal bool) *core.Dot {
	if isHeal == priest.IsOpponent(target) {
		return nil
	}

	return core.NewDot(core.Dot{
		Spell: spell,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Penance-" + strconv.Itoa(int(priest.Index)),
			ActionID: spell.ActionID,
		}),

		NumberOfTicks:       2,
		TickLength:          time.Second,
		AffectedByCastSpeed: true,

		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			if priest.IsOpponent(target) {
				baseDamage := 375 + 0.2290*dot.Spell.SpellPower()
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
			} else {
				baseHealing := sim.Roll(1484, 1676) + 0.5362*dot.Spell.HealingPower(target)
				dot.Spell.CalcAndDealPeriodicHealing(sim, target, baseHealing, dot.Spell.OutcomeHealingCrit)
			}
		},
	})
}
