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
	baseAdditiveMultiplier := warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, true)

	warlock.UnstableAffliction = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - 0.02*float64(warlock.Talents.Suppression)),
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * (1500 - 200*core.TernaryDuration(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfUnstableAffliction), 1, 0)),
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:        core.ProcMaskEmpty,
			OutcomeApplier:  warlock.OutcomeFuncMagicHit(),
			OnSpellHitDealt: applyDotOnLanded(&warlock.UnstableAfflictionDot),
		}),
	})

	target := warlock.CurrentTarget
	spellCoefficient := 0.2 + 0.01*float64(warlock.Talents.EverlastingAffliction)
	applier := warlock.OutcomeFuncTick()
	if warlock.Talents.Pandemic {
		applier = warlock.OutcomeFuncMagicCrit(warlock.SpellCritMultiplier(1, 1))
	}

	warlock.UnstableAfflictionDot = core.NewDot(core.Dot{
		Spell: warlock.UnstableAffliction,
		Aura: target.RegisterAura(core.Aura{
			Label:    "UnstableAffliction-" + strconv.Itoa(int(warlock.Index)),
			ActionID: core.ActionID{SpellID: 47843},
		}),
		NumberOfTicks: 5,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			DamageMultiplier:     baseAdditiveMultiplier,
			ThreatMultiplier:     1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
			BaseDamage:           core.BaseDamageConfigMagicNoRoll(1150/5, spellCoefficient),
			BonusSpellCritRating: 3 * core.CritRatingPerCritChance * float64(warlock.Talents.Malediction),
			OutcomeApplier:       applier,
			IsPeriodic:           true,
			ProcMask:             core.ProcMaskPeriodicDamage,
		}),
	})
}
