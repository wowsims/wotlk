package naxxrammas

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func addPatchwerk25(bossPrefix string) {
	core.AddPresetTarget(core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: proto.Target{
			Id:        16028,
			Name:      "Patchwerk 25",
			Level:     83,
			MobType:   proto.MobType_MobTypeUndead,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      16_950_147,
				stats.Armor:       10643,
				stats.AttackPower: 640,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       1.6,
			MinBaseDamage:    14135,
			SuppressDodge:    false,
			ParryHaste:       false,
			DualWield:        true,
			DualWieldPenalty: false,
		},
		AI: NewPatchwerk25AI(),
	})
	core.AddPresetEncounter("Patchwerk 25", []string{
		bossPrefix + "/Patchwerk 25",
	})
}

type Patchwerk25AI struct {
	Target *core.Target

	HatefulStrike *core.Spell
	Frenzy        *core.Spell
}

func NewPatchwerk25AI() core.AIFactory {
	return func() core.TargetAI {
		return &Patchwerk25AI{}
	}
}

func (ai *Patchwerk25AI) Initialize(target *core.Target) {
	ai.Target = target

	ai.registerHatefulStrikeSpell(target)
	ai.registerFrenzySpell(target)
}

func (ai *Patchwerk25AI) registerHatefulStrikeSpell(target *core.Target) {
	actionID := core.ActionID{SpellID: 59192}

	ai.HatefulStrike = target.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Second * 3,
			},
		},

		DamageMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BaseDamage:     core.BaseDamageConfigRoll(79000, 81000),
			OutcomeApplier: target.OutcomeFuncEnemyMeleeWhite(),
		}),
	})

}

func (ai *Patchwerk25AI) registerFrenzySpell(target *core.Target) {
	actionID := core.ActionID{SpellID: 28131}
	frenzyAura := target.GetOrRegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Frenzy",
		Duration: 5 * time.Minute,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.PhysicalDamageDealtMultiplier *= 1.25
			aura.Unit.MultiplyMeleeSpeed(sim, 1.4)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.PhysicalDamageDealtMultiplier /= 1.25
			aura.Unit.MultiplyMeleeSpeed(sim, 1.0/1.4)
		},
	})

	ai.Frenzy = target.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Minute * 5,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			frenzyAura.Activate(sim)
		},
	})
}

func (ai *Patchwerk25AI) DoAction(sim *core.Simulation) {
	if ai.Frenzy.IsReady(sim) && sim.GetRemainingDurationPercent() < 0.05 {
		ai.Frenzy.Cast(sim, ai.Target.CurrentTarget)
	}

	if ai.HatefulStrike.IsReady(sim) {
		ai.HatefulStrike.Cast(sim, ai.Target.CurrentTarget)
	}

	if ai.Target.GCD.IsReady(sim) {
		waitUntil := ai.HatefulStrike.ReadyAt()
		ai.Target.WaitUntil(sim, waitUntil)
	}
}
