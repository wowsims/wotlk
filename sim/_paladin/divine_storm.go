package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (paladin *Paladin) registerDivineStormSpell() {
	bonusDmg := core.TernaryFloat64(paladin.Ranged().ID == 45510, 235, 0) + // Libram of Discord
		core.TernaryFloat64(paladin.Ranged().ID == 38362, 81, 0) // Venture Co. Libram of Retribution
	numHits := min(4, paladin.Env.GetNumTargets())
	results := make([]*core.SpellResult, numHits)

	paladin.DivineStorm = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53385},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.12,
			Multiplier: 1 - 0.02*float64(paladin.Talents.Benediction),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // ds is on phys gcd, which cannot be hasted
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		BonusCritRating: core.TernaryFloat64(paladin.HasSetBonus(ItemSetAegisBattlegear, 4), 10, 0) * core.CritRatingPerCritChance,
		// base 1.1 multiplier, can be further improved by 10% via taow for a grand total of 1.21. NOTE: Unlike cs, ds tooltip IS NOT updated to reflect this.
		DamageMultiplierAdditive: 1 +
			paladin.getTalentTheArtOfWarBonus() +
			paladin.getItemSetRedemptionBattlegearBonus2(),
		DamageMultiplier: 1.1,
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := bonusDmg +
					spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()

				results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			curTarget = target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				spell.DealDamage(sim, results[hitIndex])
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
		},
	})
}
