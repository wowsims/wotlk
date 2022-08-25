package druid

import (
	"github.com/wowsims/wotlk/sim/core/proto"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerMoonfireSpell() {
	actionID := core.ActionID{SpellID: 48463}
	baseCost := 0.21 * druid.BaseMana

	iffCritBonus := core.TernaryFloat64(druid.CurrentTarget.HasActiveAura("Improved Faerie Fire"), float64(druid.Talents.ImprovedFaerieFire)*1*core.CritRatingPerCritChance, 0)
	improvedMoonfireDamageMultiplier := 0.05 * float64(druid.Talents.ImprovedMoonfire)
	moonfuryDamageMultiplier := 0.02 * float64(druid.Talents.Moonfury)

	moonfireGlyphBaseDamageMultiplier := core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMoonfire), 0.9, 0)
	moonfireGlyphDotDamageMultiplier := core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMoonfire), 0.75, 0)

	manaMetrics := druid.NewManaMetrics(actionID)

	druid.Moonfire = druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48463},
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
			BonusSpellCritRating: (float64(druid.Talents.ImprovedMoonfire) * 5 * core.CritRatingPerCritChance) + iffCritBonus,
			DamageMultiplier:     1 * (1 + improvedMoonfireDamageMultiplier + moonfuryDamageMultiplier - moonfireGlyphBaseDamageMultiplier),
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
			Label:    "Moonfire Dot",
			ActionID: actionID,
		}),
		NumberOfTicks: 4 + core.TernaryInt(druid.SetBonuses.balance_t6_2, 1, 0) + core.TernaryInt(druid.Talents.NaturesSplendor, 1, 0),
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 * (1 + improvedMoonfireDamageMultiplier + moonfuryDamageMultiplier + moonfireGlyphDotDamageMultiplier) * (1 + 0.01*float64(druid.Talents.Genesis)),
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(200, 0.13),
			OutcomeApplier:   druid.OutcomeFuncTick(),
			IsPeriodic:       true,
		}),
	})
}

func (druid *Druid) maxMoonfireTicks() int {
	base := 4
	thunderhearthRegalia := core.TernaryInt(druid.SetBonuses.balance_t6_2, 1, 0)
	natureSplendor := core.TernaryInt(druid.Talents.NaturesSplendor, 1, 0)
	starfireGlyph := core.TernaryInt(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfStarfire), 3, 0)
	return base + thunderhearthRegalia + natureSplendor + starfireGlyph
}
