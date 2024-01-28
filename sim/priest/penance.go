package priest

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (priest *Priest) registerPenanceHealSpell() {
	priest.PenanceHeal = priest.makePenanceSpell(true)
}

func (priest *Priest) RegisterPenanceSpell() {
	if !priest.HasRune(proto.PriestRune_RuneHandsPenance) {
		return
	}
	priest.Penance = priest.makePenanceSpell(false)
}

// https://www.wowhead.com/classic/spell=402284/penance
// https://www.wowhead.com/classic/news/patch-1-15-build-52124-ptr-datamining-season-of-discovery-runes-336044
func (priest *Priest) makePenanceSpell(isHeal bool) *core.Spell {
	var procMask core.ProcMask
	// TODO: Classic verify numbers
	spellCoeff := 0.285
	flags := core.SpellFlagChanneled | core.SpellFlagAPL
	if isHeal {
		flags |= core.SpellFlagHelpful
		procMask = core.ProcMaskSpellHealing
	} else {
		procMask = core.ProcMaskSpellDamage
	}

	level := float64(priest.GetCharacter().Level)
	baseCalc := (9.456667 + 0.635108*level + 0.039063*level*level)
	baseDamage := baseCalc * 1.28

	return priest.RegisterSpell(core.SpellConfig{
		ActionID:      core.ActionID{SpellID: 402284},
		SpellSchool:   core.SpellSchoolHoly,
		ProcMask:      procMask,
		Flags:         flags,
		RequiredLevel: 1,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.16,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 12,
			},
		},

		BonusCritRating:  1 + 0.01*float64(priest.Talents.HolySpecialization),
		DamageMultiplier: 1,
		CritMultiplier: core.TernaryFloat64(
			isHeal,
			priest.DefaultHealingCritMultiplier(),
			priest.DefaultSpellCritMultiplier(),
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
