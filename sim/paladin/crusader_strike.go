package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) registerCrusaderStrikeSpell() {
	baseCost := paladin.BaseMana * 0.05

	baseModifiers := Modifiers{
		{0.05 * float64(paladin.Talents.SanctityOfBattle), 0.05 * float64(paladin.Talents.TheArtOfWar)},
	}
	baseMultiplier := baseModifiers.Get()

	paladin.CrusaderStrike = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 35395},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(paladin.Talents.Benediction)),
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true, // cs is on phys gcd, which cannot be hasted
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 4, // the cd is 4 seconds in 3.3
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeMHSpecial,

			DamageMultiplier: baseMultiplier,
			ThreatMultiplier: 1,
			BonusCritRating:  core.TernaryFloat64(paladin.HasSetBonus(ItemSetAegisBattlegear, 4), 10, 0) * core.CritRatingPerCritChance,

			BaseDamage: core.BaseDamageConfigMeleeWeapon(
				core.MainHand,
				true, // cs is subject to normalisation
				core.TernaryFloat64(paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 31033, 36, 0)+ // Libram of Righteous Power
					core.TernaryFloat64(paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 40191, 79, 0), // Libram of Radiance
				(0.75), // base multiplier's .75, can be improved by sanctity (15%), taow (10%) & pvp gloves (5%), stacking additively
				true,
			),
			OutcomeApplier: paladin.OutcomeFuncMeleeSpecialHitAndCrit(paladin.MeleeCritMultiplier()),
		}),
	})
}
