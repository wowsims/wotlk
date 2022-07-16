package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) registerDivineStormSpell() {
	baseCost := paladin.BaseMana * 0.12

	baseMultiplier := 1.0
	// Additive bonuses
	baseMultiplier += 0.05 * float64(paladin.Talents.TheArtOfWar)

	baseEffectMH := core.SpellEffect{ // wait how will this work, something like whirlwind
		ProcMask: core.ProcMaskMeleeMHSpecial,

		DamageMultiplier: baseMultiplier,
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfigMeleeWeapon(
			core.MainHand,
			false, // ds is not subject to normalisation
			core.TernaryFloat64(paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 45510, 235, 0), // libram of discord adds 235 to ds' damage
			// (much akin to what stuff like hs or ms have intrinsically)
			(1.1), // base 1.1 multiplier, can be further improved by 10% via taow for a grand total of 1.21. NOTE: Unlike cs, ds tooltip IS NOT updated to reflect this.
			true,
		),
		OutcomeApplier: paladin.OutcomeFuncMeleeSpecialHitAndCrit(paladin.MeleeCritMultiplier()),
	}

	numHits := core.MinInt32(4, paladin.Env.GetNumTargets())
	effects := make([]core.SpellEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		mhEffect := baseEffectMH
		mhEffect.Target = paladin.Env.GetTargetUnit(i)
		effects = append(effects, mhEffect)
	}

	paladin.DivineStorm = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53385},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(paladin.Talents.Benediction)),
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true, // ds is on phys gcd, which cannot be hasted
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})
}
