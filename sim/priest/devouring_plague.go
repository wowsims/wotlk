package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerDevouringPlagueSpell() {
	actionID := core.ActionID{SpellID: 48300}
	baseCost := priest.BaseMana * 0.25
	target := priest.CurrentTarget
	initialMultiplier := 8 * 0.1 * float64(priest.Talents.ImprovedDevouringPlague)

	priest.DevouringPlague = priest.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        core.SpellFlagDisease,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - []float64{0, .04, .07, .10}[priest.Talents.MentalAgility]),
				GCD:  core.GCDDefault,
			},
		},

		BonusHitRating: float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
		BonusCritRating: 0 +
			3*float64(priest.Talents.MindMelt)*core.CritRatingPerCritChance +
			core.TernaryFloat64(priest.HasSetBonus(ItemSetCrimsonAcolyte, 2), 5, 0)*core.CritRatingPerCritChance,
		DamageMultiplier: 1 +
			float64(priest.Talents.Darkness)*0.02 +
			float64(priest.Talents.TwinDisciplines)*0.01 +
			float64(priest.Talents.ImprovedDevouringPlague)*0.05 +
			core.TernaryFloat64(priest.HasSetBonus(ItemSetConquerorSanct, 2), 0.15, 0),
		CritMultiplier:   priest.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1 - 0.05*float64(priest.Talents.ShadowAffinity),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var result *core.SpellResult
			if initialMultiplier == 0 {
				result = spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			} else {
				baseDamage := (1376/8 + 0.1849*spell.SpellPower()) * initialMultiplier
				result = spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			if result.Landed() {
				priest.AddShadowWeavingStack(sim)
				priest.DevouringPlagueDot.Apply(sim)
			}
		},
	})

	priest.DevouringPlagueDot = core.NewDot(core.Dot{
		Spell: priest.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskSpellDamage,
			Flags:       core.SpellFlagDisease,

			BonusHitRating: float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
			BonusCritRating: 0 +
				3*float64(priest.Talents.MindMelt)*core.CritRatingPerCritChance +
				core.TernaryFloat64(priest.HasSetBonus(ItemSetCrimsonAcolyte, 2), 5, 0)*core.CritRatingPerCritChance,
			DamageMultiplier: 1 +
				float64(priest.Talents.Darkness)*0.02 +
				float64(priest.Talents.TwinDisciplines)*0.01 +
				float64(priest.Talents.ImprovedDevouringPlague)*0.05 +
				core.TernaryFloat64(priest.HasSetBonus(ItemSetConquerorSanct, 2), 0.15, 0),
			CritMultiplier:   priest.SpellCritMultiplier(1, 1),
			ThreatMultiplier: 1 - 0.05*float64(priest.Talents.ShadowAffinity),
		}),
		Aura: target.RegisterAura(core.Aura{
			Label:    "DevouringPlague-" + strconv.Itoa(int(priest.Index)),
			ActionID: actionID,
		}),

		NumberOfTicks:       8,
		TickLength:          time.Second * 3,
		AffectedByCastSpeed: priest.Talents.Shadowform,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			dot.SnapshotBaseDamage = 1376/8 + 0.1849*dot.Spell.SpellPower()
			dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			if priest.Talents.Shadowform {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			} else {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			}
		},
	})
}
