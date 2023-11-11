package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (paladin *Paladin) registerCrusaderStrikeSpell() {
	bonusDmg := core.TernaryFloat64(paladin.Ranged().ID == 31033, 36, 0) + // Libram of Righteous Power
		core.TernaryFloat64(paladin.Ranged().ID == 40191, 79, 0) // Libram of Radiance

	jowAuras := paladin.NewEnemyAuraArray(core.JudgementOfWisdomAura)
	jolAuras := paladin.NewEnemyAuraArray(core.JudgementOfLightAura)

	paladin.CrusaderStrike = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 35395},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.05,
			Multiplier: 1 *
				(1 - 0.02*float64(paladin.Talents.Benediction)) *
				core.TernaryFloat64(paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfCrusaderStrike), 0.8, 1),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // cs is on phys gcd, which cannot be hasted
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 4, // the cd is 4 seconds in 3.3
			},
		},

		BonusCritRating: core.TernaryFloat64(paladin.HasSetBonus(ItemSetAegisBattlegear, 4), 10, 0) * core.CritRatingPerCritChance,
		DamageMultiplierAdditive: 1 +
			paladin.getTalentSanctityOfBattleBonus() +
			paladin.getTalentTheArtOfWarBonus() +
			paladin.getItemSetGladiatorsVindicationBonusGloves(),
		DamageMultiplier: 0.75,
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := bonusDmg +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			jowAura := jowAuras.Get(target)
			if jowAura.IsActive() {
				jowAura.Refresh(sim)
			}

			jolAura := jolAuras.Get(target)
			if jolAura.IsActive() {
				jolAura.Refresh(sim)
			}
		},
	})
}
