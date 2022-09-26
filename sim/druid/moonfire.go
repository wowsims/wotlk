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

	improvedMoonfireDamageMultiplier := 0.05 * float64(druid.Talents.ImprovedMoonfire)

	moonfireGlyphBaseDamageMultiplier := core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMoonfire), 0.9, 0)
	moonfireGlyphDotDamageMultiplier := core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMoonfire), 0.75, 0)

	// T9-2P
	dotOutcomeApplier := druid.OutcomeFuncTick()
	if druid.SetBonuses.balance_t9_2 {
		dotOutcomeApplier = druid.OutcomeFuncMagicHitAndCrit()
	}

	manaMetrics := druid.NewManaMetrics(core.ActionID{SpellID: 24858})

	druid.Moonfire = druid.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48463},
		SpellSchool:  core.SpellSchoolArcane,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * druid.TalentsBonuses.moonglowMultiplier,
				GCD:  core.GCDDefault,
			},
		},

		BonusCritRating: float64(druid.Talents.ImprovedMoonfire) * 5 * core.CritRatingPerCritChance,
		DamageMultiplierAdditive: 1 +
			improvedMoonfireDamageMultiplier +
			druid.TalentsBonuses.moonfuryMultiplier -
			moonfireGlyphBaseDamageMultiplier,
		DamageMultiplier: 1,
		CritMultiplier:   druid.SpellCritMultiplier(1, druid.TalentsBonuses.vengeanceModifier),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(406, 476) + 0.15*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				druid.MoonfireDot.Apply(sim)
				if result.DidCrit() && druid.Talents.MoonkinForm {
					druid.AddMana(sim, 0.02*druid.MaxMana(), manaMetrics, true)
				}
			}
			spell.DealDamage(sim, &result)
		},
	})

	target := druid.CurrentTarget
	druid.MoonfireDot = core.NewDot(core.Dot{
		Spell: druid.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 48463},
			SpellSchool: core.SpellSchoolArcane,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplierAdditive: 1 +
				improvedMoonfireDamageMultiplier +
				druid.TalentsBonuses.moonfuryMultiplier +
				moonfireGlyphDotDamageMultiplier,
			DamageMultiplier: 1 *
				druid.TalentsBonuses.genesisMultiplier,
			CritMultiplier:   druid.SpellCritMultiplier(1, druid.TalentsBonuses.vengeanceModifier),
			ThreatMultiplier: 1,
		}),
		Aura: target.RegisterAura(core.Aura{
			Label:    "Moonfire Dot",
			ActionID: actionID,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				druid.Starfire.BonusCritRating += core.CritRatingPerCritChance * float64(druid.Talents.ImprovedInsectSwarm)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				druid.Starfire.BonusCritRating -= core.CritRatingPerCritChance * float64(druid.Talents.ImprovedInsectSwarm)
			},
		}),
		NumberOfTicks: 4 + core.TernaryInt(druid.SetBonuses.balance_t6_2, 1, 0) + druid.TalentsBonuses.naturesSplendorTick,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			BaseDamage:     core.BaseDamageConfigMagicNoRoll(200, 0.13),
			OutcomeApplier: dotOutcomeApplier,
			IsPeriodic:     true,
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
