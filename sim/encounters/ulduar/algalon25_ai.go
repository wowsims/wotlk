package ulduar

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func addAlgalon25(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        32871,
			Name:      "Algalon",
			Level:     83,
			MobType:   proto.MobType_MobTypeElemental,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      41_834_998,
				stats.Armor:       10643,
				stats.AttackPower: 805,
				stats.BlockValue:  76,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       1.00,
			MinBaseDamage:    63649,
			SuppressDodge:    false,
			ParryHaste:       false,
			DualWield:        true,
			DualWieldPenalty: false,
			DamageSpread:     0.1,
			TargetInputs:     make([]*proto.TargetInput, 0),
		},
		AI: NewAlgalon25AI(),
	})
	core.AddPresetEncounter("Algalon", []string{
		bossPrefix + "/Algalon",
	})
}

type Algalon25AI struct {
	Target *core.Target

	QuantumStrike      *core.Spell
	PhasePunch         *core.Spell
	BlackHoleExplosion *core.Spell
	CosmicSmash        *core.Spell
}

func NewAlgalon25AI() core.AIFactory {
	return func() core.TargetAI {
		return &Algalon25AI{}
	}
}

func (ai *Algalon25AI) Initialize(target *core.Target, _ *proto.Target) {
	ai.Target = target

	ai.registerQuantumStrikeSpell(target)
	ai.registerPhasePunchSpell(target)
	ai.registerBlackHoleExplosionSpell(target)
	ai.registerCosmicSmashSpell(target)

}

func (ai *Algalon25AI) Reset(*core.Simulation) {
}

func (ai *Algalon25AI) registerQuantumStrikeSpell(target *core.Target) {
	actionID := core.ActionID{SpellID: 64592}

	ai.QuantumStrike = target.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Millisecond * 3200,
			},
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(34125, 35875)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
		},
	})
}

func (ai *Algalon25AI) registerPhasePunchSpell(target *core.Target) {
	actionID := core.ActionID{SpellID: 64412}

	ai.PhasePunch = target.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Millisecond * 16000,
			},
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(8788, 10212)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
		},
	})
}

func (ai *Algalon25AI) registerBlackHoleExplosionSpell(target *core.Target) {
	actionID := core.ActionID{SpellID: 65108}

	ai.BlackHoleExplosion = target.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNone,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Millisecond * 30000,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(20475, 21525)
			for _, aoeTarget := range sim.Raid.GetActiveUnits() {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeAlwaysHit)
			}
		},
	})
}

// This is a distance based spell which we obviously cannot model accurately, however it is
// apparent from logs that you are able to reduce the damage to pretty low levels. Therefore
// the assumption for the sim is that you do a pretty good job at it but there is always a
// closest one and a far one.
func (ai *Algalon25AI) registerCosmicSmashSpell(target *core.Target) {
	actionID := core.ActionID{SpellID: 64596}

	ai.CosmicSmash = target.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagIgnoreResists,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Millisecond * 25000,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Raid.GetActiveUnits() {
				// There are always 3 damage events at different distances
				spell.CalcAndDealDamage(sim, aoeTarget, sim.Roll(200, 800), spell.OutcomeAlwaysHit)
				spell.CalcAndDealDamage(sim, aoeTarget, sim.Roll(500, 2500), spell.OutcomeAlwaysHit)
				spell.CalcAndDealDamage(sim, aoeTarget, sim.Roll(800, 5000), spell.OutcomeAlwaysHit)
			}
		},
	})
}

func (ai *Algalon25AI) ExecuteCustomRotation(sim *core.Simulation) {
	if ai.Target.CurrentTarget != nil {
		if ai.BlackHoleExplosion.IsReady(sim) && sim.CurrentTime >= ai.BlackHoleExplosion.CD.Duration {
			ai.BlackHoleExplosion.Cast(sim, ai.Target.CurrentTarget)
		}

		if ai.CosmicSmash.IsReady(sim) && sim.CurrentTime >= ai.CosmicSmash.CD.Duration {
			ai.CosmicSmash.Cast(sim, ai.Target.CurrentTarget)
		}

		if ai.PhasePunch.IsReady(sim) && sim.CurrentTime >= ai.PhasePunch.CD.Duration {
			ai.PhasePunch.Cast(sim, ai.Target.CurrentTarget)
			return
		}

		if ai.QuantumStrike.IsReady(sim) && sim.CurrentTime >= ai.QuantumStrike.CD.Duration {
			ai.QuantumStrike.Cast(sim, ai.Target.CurrentTarget)
			return
		}
	}

	if ai.Target.GCD.IsReady(sim) {
		nextEventAt := sim.CurrentTime + time.Minute

		// All possible next events
		events := []time.Duration{
			max(ai.BlackHoleExplosion.ReadyAt(), ai.BlackHoleExplosion.CD.Duration),
			max(ai.CosmicSmash.ReadyAt(), ai.CosmicSmash.CD.Duration),
		}

		if ai.Target.CurrentTarget != nil {
			events = append(events, max(ai.PhasePunch.ReadyAt(), ai.PhasePunch.CD.Duration))
			events = append(events, max(ai.QuantumStrike.ReadyAt(), ai.QuantumStrike.CD.Duration))
		}

		for _, elem := range events {
			if elem > sim.CurrentTime && elem < nextEventAt {
				nextEventAt = elem
			}
		}

		if nextEventAt == 0 {
			nextEventAt = time.Millisecond * 100
		}

		ai.Target.WaitUntil(sim, nextEventAt)
	}
}
