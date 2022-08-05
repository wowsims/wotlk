package mage

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerFrostfireBoltSpell() {
	actionID := core.ActionID{SpellID: 47610}
	baseCost := .14 * mage.BaseMana

	frostfireGlyphBonus := int32(0)
	if mage.HasGlyph(int32(proto.MageMajorGlyph_GlyphOfFrostfire)) {
		frostfireGlyphBonus = 1
	}

	bonusCrit := 0.0
	if mage.MageTier.t9_4 {
		bonusCrit += 5 * core.CritRatingPerCritChance
	}

	mage.FrostfireBolt = mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire | core.SpellSchoolFrost,
		Flags:       SpellFlagMage | HotStreakSpells | core.SpellFlagBinary,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:            core.ProcMaskSpellDamage,
			BonusSpellHitRating: 0,

			BonusSpellCritRating: bonusCrit +
				float64(mage.Talents.CriticalMass+frostfireGlyphBonus)*2*core.CritRatingPerCritChance +
				float64(mage.Talents.ImprovedScorch)*core.CritRatingPerCritChance,

			DamageMultiplier: mage.spellDamageMultiplier *
				(1 + 0.02*float64(mage.Talents.FirePower+frostfireGlyphBonus+mage.Talents.PiercingIce)) *
				(1 + .04*float64(mage.Talents.TormentTheWeak)),

			ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul) - .04*float64(mage.Talents.FrostChanneling),

			BaseDamage:     core.BaseDamageConfigMagic(722, 838, 3.0/3.5),
			OutcomeApplier: mage.OutcomeFuncMagicHitAndCrit(mage.SpellCritMultiplier(1, mage.bonusCritDamage+float64(mage.Talents.IceShards)/3)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					mage.FrostfireDot.Apply(sim)
				}
			},
		}),
	})

	target := mage.CurrentTarget
	mage.FrostfireDot = core.NewDot(core.Dot{
		Spell: mage.FrostfireBolt,
		Aura: target.RegisterAura(core.Aura{
			Label:    "FrostfireBolt-" + strconv.Itoa(int(mage.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 3,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask: core.ProcMaskPeriodicDamage,

			DamageMultiplier: mage.spellDamageMultiplier *
				(1 + 0.02*float64(mage.Talents.FirePower+frostfireGlyphBonus+mage.Talents.PiercingIce)) *
				(1 + .04*float64(mage.Talents.TormentTheWeak)),

			ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul) - .04*float64(mage.Talents.FrostChanneling),

			BaseDamage:     core.BaseDamageConfigFlat(90 / 3),
			OutcomeApplier: mage.OutcomeFuncTick(),
			IsPeriodic:     true,
		}),
	})
}
