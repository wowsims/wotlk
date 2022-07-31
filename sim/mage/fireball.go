package mage

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerFireballSpell() {
	actionID := core.ActionID{SpellID: 27070}
	baseCost := 425.0

	mage.Fireball = mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		Flags:       SpellFlagMage,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					(1 - 0.01*float64(mage.Talents.Pyromaniac)),
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*3500 - time.Millisecond*100*time.Duration(mage.Talents.ImprovedFireball),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:            core.ProcMaskSpellDamage,
			BonusSpellHitRating: 0,

			BonusSpellCritRating: 0 +
				float64(mage.Talents.CriticalMass)*2*core.CritRatingPerCritChance +
				float64(mage.Talents.Pyromaniac)*1*core.CritRatingPerCritChance,

			DamageMultiplier: mage.spellDamageMultiplier *
				(1 + 0.02*float64(mage.Talents.FirePower)) *
				core.TernaryFloat64(mage.HasSetBonus(ItemSetTempestRegalia, 4), 1.05, 1),

			ThreatMultiplier: 1 - 0.05*float64(mage.Talents.BurningSoul),

			BaseDamage:     core.BaseDamageConfigMagic(649, 821, 1.0),
			OutcomeApplier: mage.OutcomeFuncMagicHitAndCrit(mage.SpellCritMultiplier(1, 0.25*float64(mage.Talents.SpellPower))),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
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
				(1 + 0.02*float64(mage.Talents.FirePower)) *
				core.TernaryFloat64(mage.HasSetBonus(ItemSetTempestRegalia, 4), 1.05, 1),

			ThreatMultiplier: 1 - 0.05*float64(mage.Talents.BurningSoul),

			BaseDamage:     core.BaseDamageConfigFlat(84 / 4),
			OutcomeApplier: mage.OutcomeFuncTick(),
			IsPeriodic:     true,
		}),
	})
}
