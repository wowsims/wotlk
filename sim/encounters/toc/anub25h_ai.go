package toc

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func addAnub25H(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        34564,
			Name:      "Anub'arak",
			Level:     83,
			MobType:   proto.MobType_MobTypeUndead,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      27_192_750,
				stats.Armor:       10643,
				stats.AttackPower: 805,
				stats.BlockValue:  76,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       1.50,
			MinBaseDamage:    58411, // Est 63856 minimum debuffed Unmit
			SuppressDodge:    false,
			ParryHaste:       false,
			DualWield:        false,
			DualWieldPenalty: false,
			DamageSpread:     0.45,
			TargetInputs:     make([]*proto.TargetInput, 0),
		},
		AI: NewAnub25HAI(),
	})
	core.AddPresetEncounter("Anub'arak", []string{
		bossPrefix + "/Anub'arak",
	})
}

type Anub25HAI struct {
	Target *core.Target

	FreezingSlash     *core.Spell
	LeechingSwarm     *core.Spell
	LeechingSwarmHeal *core.Spell
}

func NewAnub25HAI() core.AIFactory {
	return func() core.TargetAI {
		return &Anub25HAI{}
	}
}

func (ai *Anub25HAI) Initialize(target *core.Target, _ *proto.Target) {
	ai.Target = target
	ai.registerFreezingSlashSpell(target)
	ai.registerLeechingSwarmSpell(target)
}

func (ai *Anub25HAI) Reset(*core.Simulation) {
}

func (ai *Anub25HAI) registerFreezingSlashSpell(target *core.Target) {
	actionID := core.ActionID{SpellID: 66012}

	ai.FreezingSlash = target.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagIgnoreResists,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Second * 20,
			},
			DefaultCast: core.Cast{
				GCD: time.Millisecond * 1620,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Freezing Slash",
				Duration: time.Second * 3,
			},
			NumberOfTicks: 1,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				target.PseudoStats.Stunned = true
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				target.PseudoStats.Stunned = false
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// 25% weapon damage
			baseDamage := 0.25 * spell.Unit.AutoAttacks.MH().EnemyWeaponDamage(sim, spell.MeleeAttackPower(), 0.45)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)

			dot := spell.Dot(target)
			dot.Apply(sim)
		},
	})
}

func (ai *Anub25HAI) registerLeechingSwarmSpell(target *core.Target) {
	actionID := core.ActionID{SpellID: 66118}

	// Add a dummy spell for the main tank to keep track of the effective raid DPS loss caused by the leech ticks
	if ai.Target.CurrentTarget != nil {
		ai.LeechingSwarmHeal = ai.Target.CurrentTarget.RegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Flags:    core.SpellFlagNoOnDamageDealt,
		})
	}

	ai.LeechingSwarm = target.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagIgnoreModifiers,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Millisecond * 1620,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Leeching Swarm",
			},
			NumberOfTicks: math.MaxInt32,
			TickLength:    time.Second,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDamage := max(0.3*target.CurrentHealth(), 250.)
				result := dot.Spell.CalcAndDealPeriodicDamage(sim, target, baseDamage, dot.Spell.OutcomeAlwaysHit)

				if ai.Target.CurrentTarget != nil && ai.Target.Env.Raid.Size() == 1 {
					healingResult := ai.LeechingSwarmHeal.NewResult(&ai.Target.Unit)
					healingResult.Outcome = result.Outcome
					healingResult.Damage = -result.Damage * 2.3 * 0.5
					ai.LeechingSwarmHeal.DealPeriodicDamage(sim, healingResult)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Raid.GetActiveUnits() {
				dot := spell.Dot(aoeTarget)
				dot.Apply(sim)
			}
		},
	})
}

func (ai *Anub25HAI) ExecuteCustomRotation(sim *core.Simulation) {
	if !ai.Target.GCD.IsReady(sim) {
		return
	}

	// Cast Leeching Swarm once at the start of the encounter (since this AI only models Phase 3)
	if sim.CurrentTime < time.Millisecond*1620 {
		ai.LeechingSwarm.Cast(sim, &ai.Target.Unit)
		return
	}

	if ai.Target.CurrentTarget != nil &&
		ai.FreezingSlash.IsReady(sim) &&
		sim.CurrentTime >= ai.FreezingSlash.CD.Duration {
		// Based on log analysis, Freezing Slash appears to have a ~30% chance to "proc" on every 1.62 second server tick once it is off cooldown.
		procRoll := sim.RandomFloat("Freezing Slash AI")

		if procRoll < 0.3 {
			ai.FreezingSlash.Cast(sim, ai.Target.CurrentTarget)
			return
		}
	}

	// Anub follows the standard Classic WoW boss AI behavior of evaluating actions on a 1.62 second server tick.
	ai.Target.WaitUntil(sim, sim.CurrentTime+time.Millisecond*1620)
}
