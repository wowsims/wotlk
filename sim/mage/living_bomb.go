package mage

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerLivingBombSpell() {

	actionID := core.ActionID{SpellID: 55360}
	actionIDDot := core.ActionID{SpellID: 55359} // I want the dot to be separately trackable for metrics
	actionIDSpell := core.ActionID{SpellID: 44457}
	baseCost := .22 * mage.BaseMana
	bonusCrit := float64(mage.Talents.WorldInFlames+mage.Talents.CriticalMass) * 2 * core.CritRatingPerCritChance

	livingBombExplosionEffect := core.SpellEffect{
		BaseDamage:     core.BaseDamageConfigMagicNoRoll(690, 1.5/3.5),
		OutcomeApplier: mage.fireSpellOutcomeApplier(mage.bonusCritDamage),
	}

	livingBombExplosionSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | HotStreakSpells,

		BonusCritRating:  bonusCrit,
		DamageMultiplier: mage.spellDamageMultiplier,
		ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),
		ApplyEffects:     core.ApplyEffectFuncAOEDamageCapped(mage.Env, livingBombExplosionEffect),
		// ApplyEffects: core.ApplyEffectFuncDirectDamage(livingBombExplosionEffect),
	})

	target := mage.CurrentTarget

	lbOutcomeApplier := mage.OutcomeFuncTick()
	if mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfLivingBomb) {
		lbOutcomeApplier = mage.fireSpellOutcomeApplier(mage.bonusCritDamage)
	}

	mage.LivingBomb = mage.RegisterSpell(core.SpellConfig{
		ActionID:     actionIDSpell,
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskEmpty,
		Flags:        SpellFlagMage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,

				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			OutcomeApplier: mage.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				mage.LivingBombDots[mage.CurrentTarget.Index].Apply(sim)
			},
		}),
	})

	livingBombDotSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:         actionIDDot,
		SpellSchool:      core.SpellSchoolFire,
		ProcMask:         core.ProcMaskEmpty,
		Flags:            SpellFlagMage,
		Cast:             core.CastConfig{},
		BonusCritRating:  bonusCrit,
		DamageMultiplier: mage.spellDamageMultiplier,
		ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),
	})

	mage.LivingBombDots[target.Index] = core.NewDot(core.Dot{
		Spell: livingBombDotSpell,
		Aura: target.RegisterAura(core.Aura{
			Label:    "LivingBomb-" + strconv.Itoa(int(mage.Index)),
			ActionID: actionID,
			Tag:      "LivingBomb",
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				mage.LivingBombNotActive.Dequeue()
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				livingBombExplosionSpell.Cast(sim, target)
				mage.LivingBombNotActive.Enqueue(target)
			},
		}),

		NumberOfTicks:       4,
		TickLength:          time.Second * 3,
		AffectedByCastSpeed: false,

		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			BaseDamage:     core.BaseDamageConfigMagicNoRoll(345, .2),
			OutcomeApplier: lbOutcomeApplier,
			IsPeriodic:     true,
		}),
	})

}
