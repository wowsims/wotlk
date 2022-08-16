package mage

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerPyroblastSpell() {
	actionID := core.ActionID{SpellID: 42891}
	baseCost := .22 * mage.BaseMana

	mage.Pyroblast = mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		Flags:       SpellFlagMage | HotStreakSpells,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,

				GCD:      core.GCDDefault,
				CastTime: time.Second * 6,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if mage.HotStreakAura.IsActive() {
					bloodmageHasteAura.Activate(sim)
					cast.CastTime = 0
					// cast.AfterCastDelay could be used for CQS to avoid ignite munching. Going to wait to implement for now though
					if !mage.MageTier.t8_4 || sim.RandomFloat("MageT84PC") > .1 {
						mage.HotStreakAura.Deactivate(sim)
					}
				}
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskSpellDamage,

			BonusSpellCritRating: 0 +
				float64(mage.Talents.CriticalMass+mage.Talents.WorldInFlames)*2*core.CritRatingPerCritChance,

			DamageMultiplier: mage.spellDamageMultiplier * (1 + .04*float64(mage.Talents.TormentTheWeak)),

			ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),

			BaseDamage:     core.BaseDamageConfigMagic(1210, 1531, 1.15),
			OutcomeApplier: mage.fireSpellOutcomeApplier(mage.bonusCritDamage),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					mage.PyroblastDot.Apply(sim)
				}
			},

			MissileSpeed: 25,
		}),
	})

	target := mage.CurrentTarget
	mage.PyroblastDot = core.NewDot(core.Dot{
		Spell: mage.Pyroblast,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Pyroblast-" + strconv.Itoa(int(mage.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 4,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask: core.ProcMaskPeriodicDamage,
			DamageMultiplier: mage.spellDamageMultiplier * (1 + 0.02*float64(mage.Talents.FirePower)) *
				(1 + .04*float64(mage.Talents.TormentTheWeak)),

			ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),

			BaseDamage:     core.BaseDamageConfigMagicNoRoll(113, .02),
			OutcomeApplier: mage.OutcomeFuncTick(),
			IsPeriodic:     true,
		}),
	})
}
