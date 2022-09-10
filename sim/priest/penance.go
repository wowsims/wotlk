package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) RegisterPenanceSpell() {
	if !priest.Talents.Penance {
		return
	}

	actionID := core.ActionID{SpellID: 53007}
	baseCost := priest.BaseMana * 0.16

	penanceDots := make([]*core.Dot, len(priest.Env.AllUnits))

	damageEffect := core.ApplyEffectFuncDirectDamage(core.SpellEffect{
		ProcMask:         core.ProcMaskSpellDamage,
		ThreatMultiplier: 0,
		OutcomeApplier:   priest.OutcomeFuncMagicHit(),
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				dot := penanceDots[spellEffect.Target.UnitIndex]
				dot.Apply(sim)
				// Do immediate tick
				dot.TickOnce()
			}
		},
	})
	healingEffect := func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		hot := penanceDots[target.UnitIndex]
		hot.Apply(sim)
		// Do immediate tick
		hot.TickOnce()
	}

	priest.Penance = priest.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolHoly,
		Flags:        core.SpellFlagChanneled,
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

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if priest.IsOpponent(target) {
				damageEffect(sim, target, spell)
			} else {
				healingEffect(sim, target, spell)
			}
		},
	})

	for _, unit := range priest.Env.AllUnits {
		penanceDots[unit.UnitIndex] = priest.makePenanceDotOrHot(unit)
	}
}

func (priest *Priest) makePenanceDotOrHot(target *core.Unit) *core.Dot {
	var effect core.SpellEffect
	if priest.IsOpponent(target) {
		effect = core.SpellEffect{
			ProcMask:   core.ProcMaskPeriodicDamage,
			IsPeriodic: true,

			DamageMultiplier: 1 + 0.05*float64(priest.Talents.SearingLight),
			ThreatMultiplier: 0,

			BaseDamage:     core.BaseDamageConfigMagicNoRoll(375, .4286),
			OutcomeApplier: priest.OutcomeFuncMagicHit(),
		}
	} else {
		effect = core.SpellEffect{
			ProcMask:   core.ProcMaskPeriodicHealing,
			IsPeriodic: true,
			IsHealing:  true,

			BonusCritRating:  float64(priest.Talents.HolySpecialization) * 1 * core.CritRatingPerCritChance,
			DamageMultiplier: 1 + 0.05*float64(priest.Talents.SearingLight),
			ThreatMultiplier: 0,

			BaseDamage:     core.BaseDamageConfigHealing(1484, 1676, .5362),
			OutcomeApplier: priest.OutcomeFuncHealingCrit(priest.DefaultHealingCritMultiplier()),
		}
	}

	return core.NewDot(core.Dot{
		Spell: priest.Penance,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Penance-" + strconv.Itoa(int(priest.Index)),
			ActionID: priest.Penance.ActionID,
		}),

		NumberOfTicks:       2,
		TickLength:          time.Second,
		AffectedByCastSpeed: true,

		TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncDirectDamage(effect)),
	})
}
