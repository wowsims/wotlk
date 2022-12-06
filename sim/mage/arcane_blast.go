package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const ArcaneBlastBaseCastTime = time.Millisecond * 2500

func (mage *Mage) registerArcaneBlastSpell() {
	ArcaneBlastBaseManaCost := .07 * mage.BaseMana

	abAuraMultiplierPerStack := core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfArcaneBlast), .18, .15)
	mage.ArcaneBlastAura = mage.GetOrRegisterAura(core.Aura{
		Label:     "Arcane Blast",
		ActionID:  core.ActionID{SpellID: 36032},
		Duration:  time.Second * 8,
		MaxStacks: 4,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			oldMultiplier := 1 + float64(oldStacks)*abAuraMultiplierPerStack
			newMultiplier := 1 + float64(newStacks)*abAuraMultiplierPerStack
			mage.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexArcane] *= newMultiplier / oldMultiplier
		},
	})

	actionID := core.ActionID{SpellID: 42897}
	totalDiscount := 1 - .01*float64(mage.Talents.ArcaneFocus+mage.Talents.Precision)
	spellCoeff := 2.5/3.5 + .03*float64(mage.Talents.ArcaneEmpowerment)

	mage.ArcaneBlast = mage.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolArcane,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | BarrageSpells,
		ResourceType: stats.Mana,
		BaseCost:     ArcaneBlastBaseManaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     ArcaneBlastBaseManaCost * totalDiscount,
				GCD:      core.GCDDefault,
				CastTime: ArcaneBlastBaseCastTime,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.Cost = ArcaneBlastBaseManaCost*totalDiscount*
					(1+1.75*float64(mage.ArcaneBlastAura.GetStacks())) +
					.01*float64(mage.Talents.Precision)*ArcaneBlastBaseManaCost
				//This is really hacky. In essence for only arcane blast we need precision to apply to the
				//original base cost of the spell instead of as a cost multiplier, so add extra mana cost equal
				//to the mana saved from having precision as a cost multiplier.
			},
			AfterCast: func(sim *core.Simulation, spell *core.Spell) {
				if mage.ArcaneBlastAura.GetStacks() >= 4 {
					mage.num4CostAB++
				}
				mage.ArcaneBlastAura.Activate(sim)
				mage.ArcaneBlastAura.AddStack(sim)
			},
		},

		BonusHitRating: float64(mage.Talents.ArcaneFocus) * core.SpellHitRatingPerHitChance, // maybe precision shouldnt be here
		BonusCritRating: 0 +
			float64(mage.Talents.Incineration)*2*core.CritRatingPerCritChance +
			core.TernaryFloat64(mage.HasSetBonus(ItemSetKhadgarsRegalia, 4), 5*core.CritRatingPerCritChance, 0),
		DamageMultiplier: mage.spellDamageMultiplier * (1 + .04*float64(mage.Talents.TormentTheWeak)) * (1 + .02*float64(mage.Talents.SpellImpact)),
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage),
		ThreatMultiplier: 1 - 0.2*float64(mage.Talents.ArcaneSubtlety),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(1185, 1377) + spellCoeff*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
