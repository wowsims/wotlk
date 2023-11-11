package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (hunter *Hunter) registerChimeraShotSpell() {
	if !hunter.Talents.ChimeraShot {
		return
	}

	ssProcSpell := hunter.chimeraShotSerpentStingSpell()

	hunter.ChimeraShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53209},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.12,
			Multiplier: 1 -
				0.03*float64(hunter.Talents.Efficiency) -
				0.05*float64(hunter.Talents.MasterMarksman),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second*10 - core.TernaryDuration(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfChimeraShot), time.Second*1, 0),
			},
		},

		DamageMultiplier: 1 * hunter.markedForDeathMultiplier(),
		CritMultiplier:   hunter.critMultiplier(true, true, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.2*spell.RangedAttackPower(target) +
				hunter.AutoAttacks.Ranged().BaseDamage(sim) +
				hunter.AmmoDamageBonus +
				spell.BonusWeaponDamage()
			baseDamage *= 1.25

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
			if result.Landed() {
				if hunter.SerpentSting.Dot(target).IsActive() {
					hunter.SerpentSting.Dot(target).Rollover(sim)
					ssProcSpell.Cast(sim, target)
				} else if hunter.ScorpidStingAuras.Get(target).IsActive() {
					hunter.ScorpidStingAuras.Get(target).Refresh(sim)
				}
			}
			spell.DealDamage(sim, result)
		},
	})
}

func (hunter *Hunter) chimeraShotSerpentStingSpell() *core.Spell {
	return hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53353},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		DamageMultiplierAdditive: 1 +
			0.1*float64(hunter.Talents.ImprovedStings) +
			core.TernaryFloat64(hunter.HasSetBonus(ItemSetScourgestalkerBattlegear, 2), .1, 0),
		DamageMultiplier: 1 *
			(2.0 + core.TernaryFloat64(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfSerpentSting), 0.8, 0)) *
			hunter.markedForDeathMultiplier(),
		CritMultiplier:   hunter.critMultiplier(true, false, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 242 + 0.04*spell.RangedAttackPower(target)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedCritOnly)
		},
	})
}
