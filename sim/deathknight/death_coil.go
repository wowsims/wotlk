package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var DeathCoilActionID = core.ActionID{SpellID: 49895}

func (dk *Deathknight) registerDeathCoilSpell() {
	bonusFlatDamage := 443 + dk.sigilOfTheWildBuckBonus() + dk.sigilOfTheVengefulHeartDeathCoil()
	dk.DeathCoil = dk.RegisterSpell(core.SpellConfig{
		ActionID:    DeathCoilActionID,
		Flags:       core.SpellFlagAPL,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: 40,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		BonusCritRating: dk.darkrunedBattlegearCritBonus() * core.CritRatingPerCritChance,
		DamageMultiplier: (1 + float64(dk.Talents.Morbidity)*0.05) +
			core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDarkDeath), 0.15, 0.0),
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := (bonusFlatDamage + 0.15*dk.getImpurityBonus(spell)) * dk.RoRTSBonus(target)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() && dk.Talents.UnholyBlight {
				dk.procUnholyBlight(sim, target, result.Damage)
			}
			spell.DealDamage(sim, result)
		},
	})
}

func (dk *Deathknight) registerDrwDeathCoilSpell() {
	bonusFlatDamage := 443 + dk.sigilOfTheWildBuckBonus() + dk.sigilOfTheVengefulHeartDeathCoil()

	dk.RuneWeapon.DeathCoil = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    DeathCoilActionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		BonusCritRating: dk.darkrunedBattlegearCritBonus() * core.CritRatingPerCritChance,
		DamageMultiplier: (1.0 + float64(dk.Talents.Morbidity)*0.05) *
			core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDarkDeath), 1.15, 1.0),
		CritMultiplier:   dk.RuneWeapon.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := bonusFlatDamage + 0.15*dk.RuneWeapon.getImpurityBonus(spell)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
