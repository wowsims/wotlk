package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerCorruptionSpell() {
	actionID := core.ActionID{SpellID: 47813}
	spellSchool := core.SpellSchoolShadow
	baseAdditiveMultiplier := warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, true)
	baseCost := 0.14 * warlock.BaseMana

	warlock.Corruption = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(warlock.Talents.Suppression)),
				GCD:  core.GCDDefault,
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:        core.ProcMaskSpellDamage,
			OutcomeApplier:  warlock.OutcomeFuncMagicHit(),
			OnSpellHitDealt: applyDotOnLanded(&warlock.CorruptionDot),
		}),
	})

	target := warlock.CurrentTarget
	spellCoefficient := 0.2 + 0.12*float64(warlock.Talents.EmpoweredCorruption)/6 + 0.01*float64(warlock.Talents.EverlastingAffliction)
	applier := warlock.OutcomeFuncTick()
	if warlock.Talents.Pandemic {
		applier = warlock.OutcomeFuncMagicCrit(warlock.SpellCritMultiplier(1, 1))
	}

	warlock.CorruptionDot = core.NewDot(core.Dot{
		Spell: warlock.Corruption,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Corruption-" + strconv.Itoa(int(warlock.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks:       6,
		TickLength:          time.Second * 3,
		AffectedByCastSpeed: warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfQuickDecay),
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: baseAdditiveMultiplier,
			ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(1080/6, spellCoefficient),
			BonusSpellCritRating: core.CritRatingPerCritChance * (3*float64(warlock.Talents.Malediction) +
				5*core.TernaryFloat64(warlock.HasSetBonus(ItemSetDarkCovensRegalia, 2), 1, 0)),
			OutcomeApplier: applier,
			IsPeriodic:     true,
		}),
	})
}

func applyDotOnLanded(dot **core.Dot) func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
	return func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
		if spellEffect.Landed() {
			(*dot).Apply(sim)
		}
	}
}
