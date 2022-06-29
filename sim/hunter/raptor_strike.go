package hunter

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (hunter *Hunter) registerRaptorStrikeSpell() {
	baseCost := 120.0

	hunter.RaptorStrike = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 27014},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.2*float64(hunter.Talents.Resourcefulness)),
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeMHSpecial,

			BonusCritRating:  float64(hunter.Talents.SavageStrikes) * 10 * core.MeleeCritRatingPerCritChance,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 170, 1, true),
			OutcomeApplier: hunter.OutcomeFuncMeleeSpecialHitAndCrit(hunter.critMultiplier(false, hunter.CurrentTarget)),
		}),
	})
}

// Returns true if the regular melee swing should be used, false otherwise.
func (hunter *Hunter) TryRaptorStrike(sim *core.Simulation) *core.Spell {
	if hunter.Rotation.Weave == proto.Hunter_Rotation_WeaveAutosOnly || !hunter.RaptorStrike.IsReady(sim) || hunter.CurrentMana() < hunter.RaptorStrike.DefaultCast.Cost {
		return nil
	}

	return hunter.RaptorStrike
}
