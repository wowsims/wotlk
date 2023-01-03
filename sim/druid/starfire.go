package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Idol IDs
const IvoryMoongoddess int32 = 27518
const ShootingStar int32 = 40321

func (druid *Druid) registerStarfireSpell() {
	actionID := core.ActionID{SpellID: 48465}
	baseCost := 0.16 * druid.BaseMana
	spellCoeff := 1.0
	bonusCoeff := 0.04 * float64(druid.Talents.WrathOfCenarius)

	idolSpellPower := core.TernaryFloat64(druid.Equip[core.ItemSlotRanged].ID == IvoryMoongoddess, 55, 0) +
		core.TernaryFloat64(druid.Equip[core.ItemSlotRanged].ID == ShootingStar, 165, 0)

	hasGlyph := druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfStarfire)
	maxMoonfireTicks := druid.moonfireTicks() + core.TernaryInt32(hasGlyph, 3, 0)

	druid.Starfire = druid.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolArcane,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagNaturesGrace | SpellFlagOmenTrigger,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - 0.03*float64(druid.Talents.Moonglow)),
				GCD:      core.GCDDefault,
				CastTime: druid.starfireCastTime(),
			},
		},

		BonusCritRating: 0 +
			2*float64(druid.Talents.NaturesMajesty)*core.CritRatingPerCritChance +
			core.TernaryFloat64(druid.HasSetBonus(ItemSetThunderheartRegalia, 4), 5*core.CritRatingPerCritChance, 0) +
			core.TernaryFloat64(druid.HasSetBonus(ItemSetDreamwalkerGarb, 4), 5*core.CritRatingPerCritChance, 0),
		DamageMultiplier: (1 + []float64{0.0, 0.03, 0.06, 0.1}[druid.Talents.Moonfury]) *
			core.TernaryFloat64(druid.HasSetBonus(ItemSetMalfurionsRegalia, 4), 1.04, 1),
		CritMultiplier:   druid.BalanceCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(1038, 1222) + ((spell.SpellPower() + idolSpellPower) * spellCoeff) + (spell.SpellPower() * bonusCoeff)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				if druid.EarthAndMoonAura != nil {
					druid.EarthAndMoonAura.Activate(sim)
				}
				if hasGlyph && druid.MoonfireDot.IsActive() && druid.MoonfireDot.NumberOfTicks < maxMoonfireTicks {
					druid.MoonfireDot.NumberOfTicks += 1
					druid.MoonfireDot.UpdateExpires(druid.MoonfireDot.ExpiresAt() + time.Second*3)
				}
			}
			spell.DealDamage(sim, result)
		},
	})
}

func (druid *Druid) starfireCastTime() time.Duration {
	return time.Millisecond*3500 - time.Millisecond*100*time.Duration(druid.Talents.StarlightWrath)
}
