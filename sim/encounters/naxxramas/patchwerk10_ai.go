package naxxramas

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func addPatchwerk10(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        16028,
			Name:      "Patchwerk",
			Level:     83,
			MobType:   proto.MobType_MobTypeUndead,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      5_691_835,
				stats.Armor:       10643,
				stats.AttackPower: 574,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       1.6,
			MinBaseDamage:    14135,
			DamageSpread:     0.3333,
			SuppressDodge:    false,
			ParryHaste:       false,
			DualWield:        true,
			DualWieldPenalty: false,
			TargetInputs:     make([]*proto.TargetInput, 0),
		},
		AI: NewPatchwerk10AI(),
	})
	core.AddPresetEncounter("Patchwerk", []string{
		bossPrefix + "/Patchwerk",
	})
}

type Patchwerk10AI struct {
	Target *core.Target

	HatefulStrike *core.Spell
	Frenzy        *core.Spell
}

func NewPatchwerk10AI() core.AIFactory {
	return func() core.TargetAI {
		return &Patchwerk10AI{}
	}
}

func (ai *Patchwerk10AI) Initialize(target *core.Target, config *proto.Target) {
	ai.Target = target

	ai.registerHatefulStrikeSpell(target)
	ai.registerFrenzySpell(target)
}

func (ai *Patchwerk10AI) Reset(*core.Simulation) {
}

func (ai *Patchwerk10AI) registerHatefulStrikeSpell(target *core.Target) {
	actionID := core.ActionID{SpellID: 59192}

	ai.HatefulStrike = target.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Second * 2,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(27750, 32250)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeEnemyMeleeWhite)
		},
	})
}

func (ai *Patchwerk10AI) registerFrenzySpell(target *core.Target) {
	actionID := core.ActionID{SpellID: 28131}
	frenzyAura := target.GetOrRegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Frenzy",
		Duration: 5 * time.Minute,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.25
			aura.Unit.MultiplyMeleeSpeed(sim, 1.4)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.25
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

func (ai *Patchwerk10AI) ExecuteCustomRotation(sim *core.Simulation) {
	if ai.Target.CurrentTarget == nil {
		return
	}

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
