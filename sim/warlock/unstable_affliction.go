package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerUnstableAfflictionSpell() {
	baseCost := 0.15 * warlock.BaseMana
	actionID := core.ActionID{SpellID: 47843}
	spellSchool := core.SpellSchoolShadow
	spellCoeff := 0.2 + 0.01*float64(warlock.Talents.EverlastingAffliction)
	canCrit := warlock.Talents.Pandemic

	warlock.UnstableAffliction = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ProcMask:     core.ProcMaskEmpty,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - 0.02*float64(warlock.Talents.Suppression)),
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * (1500 - 200*core.TernaryDuration(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfUnstableAffliction), 1, 0)),
			},
		},

		BonusCritRating: 0 +
			warlock.masterDemonologistShadowCrit() +
			3*core.CritRatingPerCritChance*float64(warlock.Talents.Malediction),
		DamageMultiplierAdditive: warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, true),
		CritMultiplier:           warlock.SpellCritMultiplier(1, 1),
		ThreatMultiplier:         1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			OutcomeApplier:  warlock.OutcomeFuncMagicHit(),
			OnSpellHitDealt: applyDotOnLanded(&warlock.UnstableAfflictionDot),
		}),
	})

	warlock.UnstableAfflictionDot = core.NewDot(core.Dot{
		Spell: warlock.UnstableAffliction,
		Aura: warlock.CurrentTarget.RegisterAura(core.Aura{
			Label:    "UnstableAffliction-" + strconv.Itoa(int(warlock.Index)),
			ActionID: core.ActionID{SpellID: 47843},
		}),
		NumberOfTicks: 5,
		TickLength:    time.Second * 3,
		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			dot.SnapshotBaseDamage = 1150/5 + spellCoeff*dot.Spell.SpellPower()
			attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
			dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			if canCrit {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			} else {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			}
		},
	})
}
