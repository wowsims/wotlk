package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerMoonfireSpell() {
	actionID := core.ActionID{SpellID: 26988}
	baseCost := 0.21 * druid.BaseMana
	iffCritBonus := core.TernaryFloat64(druid.CurrentTarget.HasActiveAura("Improved Faerie Fire"), float64(druid.Talents.ImprovedFaerieFire)*1*core.CritRatingPerCritChance, 0)
	manaMetrics := druid.NewManaMetrics(actionID)

	druid.Moonfire = druid.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolArcane,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.03*float64(druid.Talents.Moonglow)),
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: float64(druid.Talents.ImprovedMoonfire)*5*core.CritRatingPerCritChance + iffCritBonus,
			DamageMultiplier:     1 * (1 + 0.05*float64(druid.Talents.ImprovedMoonfire)) * (1 + 0.02*float64(druid.Talents.Moonfury)),
			ThreatMultiplier:     1,
			BaseDamage:           core.BaseDamageConfigMagic(305, 357, 0.15),
			OutcomeApplier:       druid.OutcomeFuncMagicHitAndCrit(druid.SpellCritMultiplier(1, 0.2*float64(druid.Talents.Vengeance))),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					druid.MoonfireDot.Apply(sim)
					if spellEffect.Outcome.Matches(core.OutcomeCrit) {
						hasMoonkinForm := core.TernaryFloat64(druid.Talents.MoonkinForm, 1, 0)
						druid.AddMana(sim, druid.MaxMana()*0.02*hasMoonkinForm, manaMetrics, true)
					}
				}
			},
		}),
	})

	target := druid.CurrentTarget
	druid.MoonfireDot = core.NewDot(core.Dot{
		Spell: druid.Moonfire,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Moonfire",
			ActionID: actionID,
		}),
		NumberOfTicks: 4 + core.TernaryInt(druid.HasSetBonus(ItemSetThunderheartRegalia, 2), 1, 0) + core.TernaryInt(druid.Talents.NaturesSplendor, 1, 0),
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 * (1 + 0.05*float64(druid.Talents.ImprovedMoonfire)) * (1 + 0.02*float64(druid.Talents.Moonfury)) * (1 + 0.01*float64(druid.Talents.Genesis)),
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(200, 0.13),
			OutcomeApplier:   druid.OutcomeFuncTick(),
			IsPeriodic:       true,
		}),
	})
}
