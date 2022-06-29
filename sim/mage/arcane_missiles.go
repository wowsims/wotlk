package mage

import (
	"strconv"
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

// Note: AM doesn't charge its mana up-front, instead it charges 1/5 of the mana on each tick.
// This is probably not worth simming since no other spell in the game does this and AM isn't
// even a popular choice for arcane mages.
func (mage *Mage) registerArcaneMissilesSpell() {
	actionID := core.ActionID{SpellID: 38699}
	baseCost := 740.0

	bonusCrit := float64(mage.Talents.ArcanePotency) * 10 * core.SpellCritRatingPerCritChance

	mage.ArcaneMissiles = mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolArcane,
		Flags:       SpellFlagMage | core.SpellFlagChanneled,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 + float64(mage.Talents.EmpoweredArcaneMissiles)*0.02),

				GCD:         core.GCDDefault,
				ChannelTime: time.Second * 5,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:            core.ProcMaskEmpty,
			BonusSpellHitRating: float64(mage.Talents.ArcaneFocus) * 2 * core.SpellHitRatingPerHitChance,

			ThreatMultiplier: 1 - 0.2*float64(mage.Talents.ArcaneSubtlety),

			OutcomeApplier: mage.OutcomeFuncMagicHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					// CC has a special interaction with AM, gets the benefit of CC crit bonus from
					// the previous cast along with its own.
					if mage.ClearcastingAura != nil && mage.ClearcastingAura.IsActive() && mage.bonusAMCCCrit == 0 {
						mage.AddStatDynamic(sim, stats.SpellCrit, bonusCrit)
						mage.bonusAMCCCrit = bonusCrit
					}

					mage.ArcaneMissilesDot.Apply(sim)
				}
			},
		}),
	})

	target := mage.CurrentTarget
	mage.ArcaneMissilesDot = core.NewDot(core.Dot{
		Spell: mage.ArcaneMissiles,
		Aura: target.RegisterAura(core.Aura{
			Label:    "ArcaneMissiles-" + strconv.Itoa(int(mage.Index)),
			ActionID: actionID,
		}),

		NumberOfTicks:       5,
		TickLength:          time.Second,
		AffectedByCastSpeed: true,

		TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:            core.ProcMaskSpellDamage,
			BonusSpellHitRating: float64(mage.Talents.ArcaneFocus) * 2 * core.SpellHitRatingPerHitChance,

			DamageMultiplier: mage.spellDamageMultiplier * core.TernaryFloat64(ItemSetTempestRegalia.CharacterHasSetBonus(&mage.Character, 4), 1.05, 1),
			ThreatMultiplier: 1 - 0.2*float64(mage.Talents.ArcaneSubtlety),

			BaseDamage:     core.BaseDamageConfigMagicNoRoll(265, 1/3.5+0.03*float64(mage.Talents.EmpoweredArcaneMissiles)),
			OutcomeApplier: mage.OutcomeFuncMagicHitAndCrit(mage.SpellCritMultiplier(1, 0.25*float64(mage.Talents.SpellPower))),
		})),
	})
}
