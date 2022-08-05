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

	castTime := time.Millisecond*3500 - time.Millisecond*100*time.Duration(mage.Talents.ImprovedFireball)
	if mage.HasGlyph(int32(proto.MageMajorGlyph_GlyphOfFireball)) {
		castTime -= time.Millisecond * 150
	}

	bonusCrit := 0.0
	if mage.MageTier.t9_4 {
		bonusCrit += 5 * core.CritRatingPerCritChance
	}

	mage.Fireball = mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		Flags:       SpellFlagMage | BarrageSpells | HotStreakSpells,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
				GCD:      core.GCDDefault,
				CastTime: castTime,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:            core.ProcMaskSpellDamage,
			BonusSpellHitRating: 0,

			BonusSpellCritRating: bonusCrit +
				float64(mage.Talents.CriticalMass)*2*core.CritRatingPerCritChance +
				float64(mage.Talents.ImprovedScorch)*core.CritRatingPerCritChance,

			DamageMultiplier: mage.spellDamageMultiplier *
				(1 + 0.02*float64(mage.Talents.FirePower+mage.Talents.SpellImpact)) *
				(1 + .04*float64(mage.Talents.TormentTheWeak)),

			ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),

			BaseDamage:     core.BaseDamageConfigMagic(898, 1143, 1.0+0.05*float64(mage.Talents.EmpoweredFire)),
			OutcomeApplier: mage.OutcomeFuncMagicHitAndCrit(mage.SpellCritMultiplier(1, mage.bonusCritDamage)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() && !mage.HasGlyph(int32(proto.MageMajorGlyph_GlyphOfFireball)) {
					mage.FireballDot.Apply(sim)
				}
			},
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
