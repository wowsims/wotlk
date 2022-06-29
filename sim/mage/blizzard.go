package mage

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (mage *Mage) registerBlizzardSpell() {
	actionID := core.ActionID{SpellID: 27085}
	baseCost := 1645.0

	blizzardDot := core.NewDot(core.Dot{
		Aura: mage.RegisterAura(core.Aura{
			Label:    "Blizzard",
			ActionID: actionID,
		}),
		NumberOfTicks:       8,
		TickLength:          time.Second * 1,
		AffectedByCastSpeed: true,
		TickEffects: core.TickFuncAOESnapshotCapped(mage.Env, 3620, core.SpellEffect{
			ProcMask: core.ProcMaskPeriodicDamage,
			DamageMultiplier: mage.spellDamageMultiplier *
				(1 + 0.02*float64(mage.Talents.PiercingIce)) *
				(1 + 0.01*float64(mage.Talents.ArcticWinds)),

			ThreatMultiplier: 1 - (0.1/3)*float64(mage.Talents.FrostChanneling),

			BaseDamage:     core.BaseDamageConfigMagicNoRoll(184, 0.119),
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
				Cost: baseCost *
					(1 - 0.05*float64(mage.Talents.FrostChanneling)) *
					(1 - 0.01*float64(mage.Talents.ElementalPrecision)),

				GCD:         core.GCDDefault,
				ChannelTime: time.Second * 8,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDot(blizzardDot),
	})
	blizzardDot.Spell = mage.Blizzard
}
