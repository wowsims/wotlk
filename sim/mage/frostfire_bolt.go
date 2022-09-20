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

	mage.FrostfireBolt = mage.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolFire | core.SpellSchoolFrost,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | HotStreakSpells,
		MissileSpeed: 25,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
		},

		BonusCritRating: 0 +
			core.TernaryFloat64(mage.MageTier.t9_4, 5*core.CritRatingPerCritChance, 0) +
			core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfFrostfire), 2*core.CritRatingPerCritChance, 0) +
			float64(mage.Talents.CriticalMass)*2*core.CritRatingPerCritChance +
			float64(mage.Talents.ImprovedScorch)*1*core.CritRatingPerCritChance,
		DamageMultiplier: mage.spellDamageMultiplier *
			(1 + .02*float64(mage.Talents.PiercingIce)) *
			(1 + core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfFrostfire), .02, 0)) *
			(1 + .04*float64(mage.Talents.TormentTheWeak)) *
			(1 + .01*float64(mage.Talents.ChilledToTheBone)),
		ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul) - .04*float64(mage.Talents.FrostChanneling),

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BaseDamage:     core.BaseDamageConfigMagicNoRoll((722+838)/2, 3.0/3.5+float64(mage.Talents.EmpoweredFire)*.05),
			OutcomeApplier: mage.fireSpellOutcomeApplier(mage.bonusCritDamage + float64(mage.Talents.IceShards)/3),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					mage.FrostfireDot.Apply(sim)
				}
			},
		}),
	})

	target := mage.CurrentTarget
	mage.FrostfireDot = core.NewDot(core.Dot{
		Spell: mage.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolFire | core.SpellSchoolFrost,
			Flags:       SpellFlagMage | HotStreakSpells,

			DamageMultiplier: mage.FrostfireBolt.DamageMultiplier /
				(1 + core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfFrostfire), .02, 0)),
			ThreatMultiplier: mage.FrostfireBolt.ThreatMultiplier,
		}),
		Aura: target.RegisterAura(core.Aura{
			Label:    "FrostfireBolt-" + strconv.Itoa(int(mage.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 3,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			BaseDamage:     core.BaseDamageConfigFlat(90 / 3),
			OutcomeApplier: mage.OutcomeFuncTick(),
			IsPeriodic:     true,
		}),
	})
}
