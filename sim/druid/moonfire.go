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

	iffCritBonus := core.TernaryFloat64(druid.CurrentTarget.HasActiveAura("Improved Faerie Fire"), druid.TalentsBonuses.iffBonusCrit, 0)
	improvedMoonfireDamageMultiplier := 0.05 * float64(druid.Talents.ImprovedMoonfire)

	moonfireGlyphBaseDamageMultiplier := core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMoonfire), 0.9, 0)
	moonfireGlyphDotDamageMultiplier := core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMoonfire), 0.75, 0)

	// T9-2P
	dotOutcomeApplier := druid.OutcomeFuncTick()
	if druid.SetBonuses.balance_t9_2 {
		dotOutcomeApplier = druid.OutcomeFuncMagicHitAndCrit(druid.SpellCritMultiplier(1, druid.TalentsBonuses.vengeanceModifier))
	}

	manaMetrics := druid.NewManaMetrics(core.ActionID{SpellID: 24858})

	druid.Moonfire = druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48463},
		SpellSchool: core.SpellSchoolArcane,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * druid.TalentsBonuses.moonglowMultiplier,
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: (float64(druid.Talents.ImprovedMoonfire) * 5 * core.CritRatingPerCritChance) + iffCritBonus,
			DamageMultiplier:     1 * (1 + improvedMoonfireDamageMultiplier + druid.TalentsBonuses.moonfuryMultiplier - moonfireGlyphBaseDamageMultiplier),
			ThreatMultiplier:     1,
			BaseDamage:           core.BaseDamageConfigMagic(305, 357, 0.15),
			OutcomeApplier:       druid.OutcomeFuncMagicHitAndCrit(druid.SpellCritMultiplier(1, druid.TalentsBonuses.vengeanceModifier)),
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
		NumberOfTicks: 4 + core.TernaryInt(druid.SetBonuses.balance_t6_2, 1, 0) + druid.TalentsBonuses.naturesSplendorTick,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 * (1 + improvedMoonfireDamageMultiplier + druid.TalentsBonuses.moonfuryMultiplier + moonfireGlyphDotDamageMultiplier) * druid.TalentsBonuses.genesisMultiplier,
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(200, 0.13),
			OutcomeApplier:   dotOutcomeApplier,
			IsPeriodic:       true,
		}),
	})
}

func (druid *Druid) maxMoonfireTicks() int {
	base := 4
	thunderhearthRegalia := core.TernaryInt(druid.SetBonuses.balance_t6_2, 1, 0)
	natureSplendor := druid.TalentsBonuses.naturesSplendorTick
	starfireGlyph := core.TernaryInt(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfStarfire), 3, 0)
	return base + thunderhearthRegalia + natureSplendor + starfireGlyph
}
