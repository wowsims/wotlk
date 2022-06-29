package druid

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (druid *Druid) registerHurricaneSpell() {
	actionID := core.ActionID{SpellID: 27012}
	baseCost := 1905.0

	hurricaneDot := core.NewDot(core.Dot{
		Aura: druid.RegisterAura(core.Aura{
			Label:    "Hurricane",
			ActionID: actionID,
		}),
		NumberOfTicks:       10,
		TickLength:          time.Second * 1,
		AffectedByCastSpeed: true,
		TickEffects: core.TickFuncAOESnapshot(druid.Env, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(206, 0.107),
			OutcomeApplier:   druid.OutcomeFuncTick(),
			IsPeriodic:       true,
		}),
	})

	druid.Hurricane = druid.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		Flags:       core.SpellFlagChanneled,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:        baseCost,
				GCD:         core.GCDDefault,
				ChannelTime: time.Second * 10,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 60,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDot(hurricaneDot),
	})
	hurricaneDot.Spell = druid.Hurricane
}

func (druid *Druid) ShouldCastHurricane(sim *core.Simulation, rotation proto.BalanceDruid_Rotation) bool {
	return rotation.Hurricane && druid.Hurricane.IsReady(sim)
}
