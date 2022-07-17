package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) registerCrusaderStrikeSpell() {
	baseCost := paladin.BaseMana * 0.05

	baseMultiplier := 1.0
	// Additive bonuses
	baseMultiplier += 0.05 * float64(paladin.Talents.SanctityOfBattle)
	baseMultiplier += 0.05 * float64(paladin.Talents.TheArtOfWar)

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

			BaseDamage: core.BaseDamageConfigMeleeWeapon(
				core.MainHand,
				true, // cs is subject to normalisation
				core.TernaryFloat64(paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 31033, 36, 0), // theres librams that can improve cs damage (this is a 2.x one - 3.x wip)
				(0.75), // base multiplier's .75, can be improved by sanctity (15%), taow (10%) & pvp gloves (5%), stacking additively
				// for a grand total of .9375 / .975 multiplier, respectively, which is also UPDATED LIVE on the TOOLTIP.
				true,
			),
			OutcomeApplier: paladin.OutcomeFuncMeleeSpecialHitAndCrit(paladin.MeleeCritMultiplier()),
		}),
	})
}
