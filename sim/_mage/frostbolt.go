package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (mage *Mage) registerFrostboltSpell() {
	spellCoeff := (3.0 / 3.5) + 0.05*float64(mage.Talents.EmpoweredFrostbolt)

	replProcChance := float64(mage.Talents.EnduringWinter) / 3
	var replSrc core.ReplenishmentSource
	if replProcChance > 0 {
		replSrc = mage.Env.Raid.NewReplenishmentSource(core.ActionID{SpellID: 44561})
	}

	mage.Frostbolt = mage.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 42842},
		SpellSchool:  core.SpellSchoolFrost,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | BarrageSpells | core.SpellFlagAPL,
		MissileSpeed: 28,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.11,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second*3 - time.Millisecond*100*time.Duration(mage.Talents.ImprovedFrostbolt+mage.Talents.EmpoweredFrostbolt),
			},
		},

		BonusCritRating: 0 +
			core.TernaryFloat64(mage.HasSetBonus(ItemSetKhadgarsRegalia, 4), 5*core.CritRatingPerCritChance, 0),
		DamageMultiplierAdditive: 1 +
			.01*float64(mage.Talents.ChilledToTheBone) +
			core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfFrostbolt), .05, 0) +
			core.TernaryFloat64(mage.HasSetBonus(ItemSetTempestRegalia, 4), .05, 0),
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage+float64(mage.Talents.IceShards)/3),
		ThreatMultiplier: 1 - (0.1/3)*float64(mage.Talents.FrostChanneling),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(804, 866) + spellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if replProcChance == 1 || sim.RandomFloat("Enduring Winter") < replProcChance {
					mage.Env.Raid.ProcReplenishment(sim, replSrc)
				}
			})
		},
	})
}
