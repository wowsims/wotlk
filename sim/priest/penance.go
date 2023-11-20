package priest

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (priest *Priest) registerPenanceHealSpell() {
	priest.PenanceHeal = priest.makePenanceSpell(true)
}

func (priest *Priest) RegisterPenanceSpell() {
	if !priest.HasRuneById(PriestRuneHandsPenance) {
		return
	}
	priest.Penance = priest.makePenanceSpell(false)
}

func (priest *Priest) makePenanceSpell(isHeal bool) *core.Spell {
	var procMask core.ProcMask
	// TODO: Classic verify numbers
	spellCoeff := 0.2290       // Using wrath default for now
	baseDamage := float64(128) // https://www.wowhead.com/classic/news/patch-1-15-build-52124-ptr-datamining-season-of-discovery-runes-336044
	flags := core.SpellFlagChanneled | core.SpellFlagAPL
	if isHeal {
		flags |= core.SpellFlagHelpful
		procMask = core.ProcMaskSpellHealing
	} else {
		procMask = core.ProcMaskSpellDamage
	}

	return priest.RegisterSpell(core.SpellConfig{
		ActionID:      core.ActionID{SpellID: 402174},
		SpellSchool:   core.SpellSchoolHoly,
		ProcMask:      procMask,
		Flags:         flags,
		RequiredLevel: 1,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.16,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:         core.GCDDefault,
				ChannelTime: time.Second * 2,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 12,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier: core.TernaryFloat64(
			isHeal,
			priest.DefaultHealingCritMultiplier(),
			priest.SpellCritMultiplier(1+0.01*float64(priest.Talents.HolySpecialization), 1),
		),
		ThreatMultiplier: 0,

		Dot: core.Ternary(!isHeal, core.DotConfig{
			Aura: core.Aura{
				Label: "Penance",
			},
			NumberOfTicks:       2,
			TickLength:          time.Second,
			AffectedByCastSpeed: true,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dmg := baseDamage + (spellCoeff * dot.Spell.SpellPower())
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, dmg, dot.OutcomeTick)
			},
		}, core.DotConfig{}),
		// TODO: Classic healing
		Hot: core.Ternary(isHeal, core.DotConfig{
			Aura: core.Aura{
				Label: "Penance",
			},
			NumberOfTicks:       2,
			TickLength:          time.Second,
			AffectedByCastSpeed: true,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseHealing := sim.Roll(1484, 1676) + 0.5362*dot.Spell.HealingPower(target)
				dot.Spell.CalcAndDealPeriodicHealing(sim, target, baseHealing, dot.Spell.OutcomeHealingCrit)
			},
		}, core.DotConfig{}),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if isHeal {
				spell.SpellMetrics[target.UnitIndex].Hits--
				hot := spell.Hot(target)
				hot.Apply(sim)
				// Do immediate tick
				hot.TickOnce(sim)
			} else {
				result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
				if result.Landed() {
					spell.SpellMetrics[target.UnitIndex].Hits--
					dot := spell.Dot(target)
					dot.Apply(sim)
					// Do immediate tick
					dot.TickOnce(sim)
				}
				spell.DealOutcome(sim, result)
			}
		},
	})
}
