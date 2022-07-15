package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warlock *Warlock) registerUnstableAffSpell() {
	baseCost := 0.15 * warlock.BaseMana

	warlock.UnstableAff = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 30405},
		SpellSchool:  core.SpellSchoolShadow,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - 0.02*float64(warlock.Talents.Suppression)),
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * (1500 - 200 * core.TernaryDuration(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfUnstableAffliction), 1, 0)),
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:        core.ProcMaskEmpty,
			OutcomeApplier:  warlock.OutcomeFuncMagicHit(),
			OnSpellHitDealt: applyDotOnLanded(&warlock.UnstableAffDot),
		}),
	})
}

func (warlock *Warlock) registerUnstableAffDot() {
	target := warlock.CurrentTarget
	ticksNumber := 6
	spellCoefficient := 1.2 / float64(ticksNumber) + 0.01 * float64(warlock.Talents.EverlastingAffliction)
	applier := warlock.OutcomeFuncTick()
	if warlock.Talents.Pandemic {
		applier = warlock.OutcomeFuncMagicCrit(warlock.SpellCritMultiplier(1, 1))
	}

	warlock.UnstableAffDot = core.NewDot(core.Dot{
		Spell: warlock.UnstableAff,
		Aura: target.RegisterAura(core.Aura{
			Label:    "unstableaff-" + strconv.Itoa(int(warlock.Index)),
			ActionID: core.ActionID{SpellID: 30405},
		}),
		NumberOfTicks: ticksNumber,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			DamageMultiplier: 1 * (1 + 0.05*core.TernaryFloat64(warlock.Talents.SiphonLife, 1, 0)),
			ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(1150/6, spellCoefficient),
			BonusCritRating:  3 * core.CritRatingPerCritChance * float64(warlock.Talents.Malediction),
			OutcomeApplier:   applier,
			IsPeriodic:       true,
			ProcMask:         core.ProcMaskPeriodicDamage,
		}),
	})
}
