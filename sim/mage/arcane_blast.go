package mage

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const ArcaneBlastBaseCastTime = time.Millisecond * 2500

func (mage *Mage) registerArcaneBlastSpell() {
	ArcaneBlastBaseManaCost := .07 * mage.BaseMana

	mage.ArcaneBlastAura = mage.GetOrRegisterAura(core.Aura{
		Label:     "Arcane Blast",
		ActionID:  core.ActionID{SpellID: 36032},
		Duration:  time.Second * 8,
		MaxStacks: 4,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			if !aura.IsActive() {
				return
			}
			mage.PseudoStats.ArcaneDamageDealtMultiplier = 1 + float64(newStacks)*.15
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.PseudoStats.ArcaneDamageDealtMultiplier = 1
		},
	})

	actionID := core.ActionID{SpellID: 42897}

	mage.ArcaneBlast = mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolArcane,
		Flags:       SpellFlagMage | BarrageSpells,

		ResourceType: stats.Mana,
		BaseCost:     ArcaneBlastBaseManaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: ArcaneBlastBaseManaCost,

				GCD:      core.GCDDefault,
				CastTime: ArcaneBlastBaseCastTime,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.Cost = spell.BaseCost * math.Pow(1.75, float64(mage.ArcaneBlastAura.GetStacks()))
			},
			OnCastComplete: func(sim *core.Simulation, _ *core.Spell) {
				mage.ArcaneBlastAura.Activate(sim)
				mage.ArcaneBlastAura.AddStack(sim)
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellHitRating:  float64(mage.Talents.ArcaneFocus+mage.Talents.Precision) * core.SpellHitRatingPerHitChance, // maybe precision shouldnt be here
			BonusSpellCritRating: float64(mage.Talents.Incineration) * 2 * core.CritRatingPerCritChance,

			DamageMultiplier: mage.spellDamageMultiplier + float64(mage.Talents.SpellImpact)*.02,
			ThreatMultiplier: 1 - 0.2*float64(mage.Talents.ArcaneSubtlety),

			BaseDamage:     core.BaseDamageConfigMagic(1185, 1377, (2.5/3.5)+.03*float64(mage.Talents.ArcaneEmpowerment)),
			OutcomeApplier: mage.OutcomeFuncMagicHitAndCrit(mage.SpellCritMultiplier(1, 0.25*float64(mage.Talents.SpellPower))),
		}),
	})
}
