package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerBlizzardSpell() {
	actionID := core.ActionID{SpellID: 42939}
	baseCost := .74 * mage.BaseMana

	blizzardDot := core.NewDot(core.Dot{
		Aura: mage.RegisterAura(core.Aura{
			Label:    "Blizzard",
			ActionID: actionID,
		}),
		NumberOfTicks:       8,
		TickLength:          time.Second * 1,
		AffectedByCastSpeed: true,
		TickEffects: core.TickFuncAOESnapshotCapped(mage.Env, core.SpellEffect{
			ProcMask:       core.ProcMaskPeriodicDamage,
			BaseDamage:     core.BaseDamageConfigMagicNoRoll(352, 0.119),
			OutcomeApplier: mage.OutcomeFuncTick(),
			IsPeriodic:     true,
		}),
	})

	mage.Blizzard = mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFrost,
		Flags:       SpellFlagMage | core.SpellFlagChanneled,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,

				GCD:         core.GCDDefault,
				ChannelTime: time.Second * 8,
			},
		},

		DamageMultiplier: mage.spellDamageMultiplier,
		ThreatMultiplier: 1 - (0.1/3)*float64(mage.Talents.FrostChanneling),

		ApplyEffects: core.ApplyEffectFuncDot(blizzardDot),
	})
	blizzardDot.Spell = mage.Blizzard
}
