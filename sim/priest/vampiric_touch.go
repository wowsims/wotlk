package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerVampiricTouchSpell() {
	actionID := core.ActionID{SpellID: 48160}
	baseCost := priest.BaseMana * 0.16

	priest.VampiricTouch = priest.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		BonusHitRating:   float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
		BonusCritRating:  float64(priest.Talents.MindMelt)*3*core.CritRatingPerCritChance + core.TernaryFloat64(priest.HasSetBonus(ItemSetCrimsonAcolyte, 2), 5, 0)*core.CritRatingPerCritChance,
		DamageMultiplier: 1 + float64(priest.Talents.Darkness)*0.02,
		CritMultiplier:   priest.SpellCritMultiplier(1, 1),
		ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				priest.AddShadowWeavingStack(sim)
				priest.VampiricTouchDot.Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
		ExpectedDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) *core.SpellResult {
			baseDamage := 850/5 + 0.4*spell.SpellPower()
			baseDamage *= float64(priest.VampiricTouchDot.NumberOfTicks)

			if priest.Talents.Shadowform {
				return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicCrit)
			} else {
				return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
			}
		},
	})

	target := priest.CurrentTarget

	priest.VampiricTouchDot = core.NewDot(core.Dot{
		Spell: priest.VampiricTouch,
		Aura: target.RegisterAura(core.Aura{
			Label:    "VampiricTouch-" + strconv.Itoa(int(priest.Index)),
			ActionID: actionID,
		}),

		NumberOfTicks:       5 + core.TernaryInt32(priest.HasSetBonus(ItemSetZabras, 2), 2, 0),
		TickLength:          time.Second * 3,
		AffectedByCastSpeed: priest.Talents.Shadowform,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			dot.SnapshotBaseDamage = 850/5 + 0.4*dot.Spell.SpellPower()
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
