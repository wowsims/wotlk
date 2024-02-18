package toc

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func addGormok25H(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        34796,
			Name:      "Gormok",
			Level:     83,
			MobType:   proto.MobType_MobTypeHumanoid,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      11_853_250,
				stats.Armor:       10643,
				stats.AttackPower: 805,
				stats.BlockValue:  76,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       1.50,
			MinBaseDamage:    39600, // Est 43K minimum debuffed Unmit
			SuppressDodge:    false,
			ParryHaste:       false,
			DualWield:        false,
			DualWieldPenalty: false,
			DamageSpread:     0.3333,
			TargetInputs:     make([]*proto.TargetInput, 0),
		},
		AI: NewGormok25HAI(),
	})
	core.AddPresetEncounter("Gormok", []string{
		bossPrefix + "/Gormok",
	})
}

type Gormok25HAI struct {
	Target *core.Target

	Impale          *core.Spell
	StaggeringStomp *core.Spell
	RisingAnger     *core.Spell
	RisingAngerAura *core.Aura

	//ValidStompTarget   bool
}

//func GormokTargetInputs() []*proto.TargetInput {
//	return []*proto.TargetInput{
//		{
//			Label:     "Getting Stomped On",
//			Tooltip:   "Keep this checked if you are melee",
//			InputType: proto.InputType_Bool,
//			BoolValue: true,
//		},
//	}
//}

func NewGormok25HAI() core.AIFactory {
	return func() core.TargetAI {
		return &Gormok25HAI{}
	}
}

func (ai *Gormok25HAI) Initialize(target *core.Target, _ *proto.Target) {
	ai.Target = target

	//ai.ValidStompTarget = config.TargetInputs[0].BoolValue

	ai.registerImpaleSpell(target)
	ai.registerStaggeringStompSpell(target)
	ai.registerRisingAngerSpell(target)

}

func (ai *Gormok25HAI) Reset(*core.Simulation) {
}

func (ai *Gormok25HAI) registerImpaleSpell(target *core.Target) {
	actionID := core.ActionID{SpellID: 66331}

	// TODO - Allegedly he can be Disarmed to suppress this ability?

	ai.Impale = target.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Millisecond * 10000,
			},
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Impale Bleed",
				MaxStacks: 99,
				Duration:  time.Second * 40,
			},
			NumberOfTicks: 20,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = sim.Roll(3938, 5062)
				dot.SnapshotBaseDamage *= float64(dot.Aura.GetStacks())

				if !isRollover {
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
					dot.Spell.DamageMultiplier = 1
					dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(attackTable)
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				}
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// 150% weapon damage
			baseDamage := 1.50 * spell.Unit.AutoAttacks.MH().EnemyWeaponDamage(sim, spell.MeleeAttackPower(), 0.3333)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)

			dot := spell.Dot(target)
			if dot.IsActive() {
				dot.Refresh(sim)
				dot.AddStack(sim)
				dot.TakeSnapshot(sim, true)
			} else {
				dot.Apply(sim)
				dot.SetStacks(sim, 1)
				dot.TakeSnapshot(sim, true)
			}
		},
	})
}

func (ai *Gormok25HAI) registerStaggeringStompSpell(target *core.Target) {
	actionID := core.ActionID{SpellID: 66330}

	ai.StaggeringStomp = target.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNone,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Second * 20,
			},
			DefaultCast: core.Cast{
				CastTime: time.Millisecond * 500,
				GCD:      core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {

			for _, aoeTarget := range sim.Raid.GetActiveUnits() {

				// TODO - Filter targets to melee only, right now it just hits everyone
				// TODO - Should this ignore armor? Damage in logs seems inconsistent
				baseDamage := sim.Roll(11700, 12300)
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeAlwaysHit)
				// TODO - Interrupts spellcasting for 8 seconds. Does NOT stun or knockdown
			}
		},
	})
}

func (ai *Gormok25HAI) registerRisingAngerSpell(target *core.Target) {

	actionID := core.ActionID{SpellID: 66636}

	ai.RisingAngerAura = target.GetOrRegisterAura(core.Aura{
		Label:     "Rising Anger",
		ActionID:  actionID.WithTag(1),
		MaxStacks: 99,
		Duration:  time.Second * 120,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1 + (.15 * float64(oldStacks))
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1 + (.15 * float64(newStacks))
		},
	})

	ai.RisingAnger = target.RegisterSpell(core.SpellConfig{
		ActionID:    actionID.WithTag(2),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNone,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Second * 20,
			},
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			ai.RisingAngerAura.Activate(sim)
			ai.RisingAngerAura.AddStack(sim)
		},
	})

}

func (ai *Gormok25HAI) ExecuteCustomRotation(sim *core.Simulation) {
	if ai.RisingAnger.IsReady(sim) && sim.CurrentTime >= ai.RisingAnger.CD.Duration && ai.Target.GCD.IsReady(sim) {
		ai.RisingAnger.Cast(sim, &ai.Target.Unit)
		return
	}

	if ai.StaggeringStomp.IsReady(sim) && sim.CurrentTime >= ai.StaggeringStomp.CD.Duration && ai.Target.GCD.IsReady(sim) {
		ai.StaggeringStomp.Cast(sim, &ai.Target.Unit)
		return
	}

	if ai.Target.CurrentTarget != nil {
		if ai.Impale.IsReady(sim) && sim.CurrentTime >= ai.Impale.CD.Duration && ai.Target.GCD.IsReady(sim) {
			ai.Impale.Cast(sim, ai.Target.CurrentTarget)
			return
		}
	}

	if ai.Target.GCD.IsReady(sim) {
		nextEventAt := sim.CurrentTime + time.Minute

		// All possible next events
		events := []time.Duration{
			max(ai.StaggeringStomp.ReadyAt(), ai.StaggeringStomp.CD.Duration),
			max(ai.RisingAnger.ReadyAt(), ai.RisingAnger.CD.Duration),
		}

		if ai.Target.CurrentTarget != nil {
			events = append(events, max(ai.Impale.ReadyAt(), ai.Impale.CD.Duration))
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
