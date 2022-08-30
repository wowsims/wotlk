package mage

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerFireballSpell() {
	actionID := core.ActionID{SpellID: 42833}
	baseCost := .19 * mage.BaseMana

	hasGlyph := mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfFireball)

	mage.Fireball = mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		Flags:       SpellFlagMage | BarrageSpells | HotStreakSpells,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
				CastTime: time.Millisecond*3500 -
					time.Millisecond*100*time.Duration(mage.Talents.ImprovedFireball) -
					core.TernaryDuration(hasGlyph, time.Millisecond*150, 0),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:            core.ProcMaskSpellDamage,
			BonusSpellHitRating: 0,

			BonusCritRating: 0 +
				float64(mage.Talents.CriticalMass)*2*core.CritRatingPerCritChance +
				float64(mage.Talents.ImprovedScorch)*core.CritRatingPerCritChance +
				core.TernaryFloat64(mage.MageTier.t9_4, 5*core.CritRatingPerCritChance, 0),

			DamageMultiplier: mage.spellDamageMultiplier *
				(1 + 0.02*float64(mage.Talents.SpellImpact)) *
				(1 + .04*float64(mage.Talents.TormentTheWeak)),

			ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),

			BaseDamage: core.BaseDamageConfigMagic(898, 1143, 1.0+0.05*float64(mage.Talents.EmpoweredFire)),
			// BaseDamage:     core.BaseDamageConfigMagicNoRoll((898 + 1143)/2, 1.0+0.05*float64(mage.Talents.EmpoweredFire)),
			OutcomeApplier: mage.fireSpellOutcomeApplier(mage.bonusCritDamage),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() && !hasGlyph {
					mage.FireballDot.Apply(sim)
				}
			},

			MissileSpeed: 22,
		}),
	})

	target := mage.CurrentTarget
	mage.FireballDot = core.NewDot(core.Dot{
		Spell: mage.Fireball,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Fireball-" + strconv.Itoa(int(mage.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 4,
		TickLength:    time.Second * 2,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask: core.ProcMaskPeriodicDamage,

			DamageMultiplier: mage.spellDamageMultiplier *
				(1 + 0.02*float64(mage.Talents.FirePower)) * (1 + .04*float64(mage.Talents.TormentTheWeak)),

			ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),

			BaseDamage:     core.BaseDamageConfigFlat(116 / 4),
			OutcomeApplier: mage.OutcomeFuncTick(),
			IsPeriodic:     true,
		}),
	})
}
