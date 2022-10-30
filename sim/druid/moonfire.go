package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"

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
	dotCanCrit := druid.setBonuses.balance_t9_2

	manaMetrics := druid.NewManaMetrics(core.ActionID{SpellID: 24858})

	druid.Moonfire = druid.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48463},
		SpellSchool:  core.SpellSchoolArcane,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * druid.talentBonuses.moonglow,
				GCD:  core.GCDDefault,
			},
		},

		BonusCritRating: float64(druid.Talents.ImprovedMoonfire) * 5 * core.CritRatingPerCritChance,
		DamageMultiplier: 1 +
			improvedMoonfireDamageMultiplier +
			druid.talentBonuses.moonfury -
			moonfireGlyphBaseDamageMultiplier,
		CritMultiplier:   druid.SpellCritMultiplier(1, druid.talentBonuses.vengeance),
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
			spell.DealDamage(sim, result)
		},
	})

	druid.MoonfireDot = core.NewDot(core.Dot{
		Spell: druid.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 48463},
			SpellSchool: core.SpellSchoolArcane,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplier: 1 +
				druid.talentBonuses.genesis +
				improvedMoonfireDamageMultiplier +
				druid.talentBonuses.moonfury +
				moonfireGlyphDotDamageMultiplier,

			CritMultiplier:   druid.SpellCritMultiplier(1, druid.talentBonuses.vengeance),
			ThreatMultiplier: 1,
		}),
		Aura: druid.CurrentTarget.RegisterAura(core.Aura{
			Label:    "Moonfire Dot",
			ActionID: actionID,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				druid.Starfire.BonusCritRating += core.CritRatingPerCritChance * float64(druid.Talents.ImprovedInsectSwarm)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				druid.Starfire.BonusCritRating -= core.CritRatingPerCritChance * float64(druid.Talents.ImprovedInsectSwarm)
			},
		}),
		NumberOfTicks: 4 + core.TernaryInt(druid.setBonuses.balance_t6_2, 1, 0) + druid.talentBonuses.naturesSplendor,
		TickLength:    time.Second * 3,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
			dot.SnapshotBaseDamage = 200 + 0.13*dot.Spell.SpellPower()
			attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
			dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			if dotCanCrit {
				// TODO: This allows misses... probably a bug.
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeMagicHitAndSnapshotCrit)
			} else {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			}
		},
	})
}

func (druid *Druid) maxMoonfireTicks() int {
	base := 4
	thunderhearthRegalia := core.TernaryInt(druid.setBonuses.balance_t6_2, 1, 0)
	natureSplendor := druid.talentBonuses.naturesSplendor
	starfireGlyph := core.TernaryInt(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfStarfire), 3, 0)
	return base + thunderhearthRegalia + natureSplendor + starfireGlyph
}
