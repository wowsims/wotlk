package icc

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func addSindragosa25H(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        36853,
			Name:      "Sindragosa (Heroic)",
			Level:     83,
			MobType:   proto.MobType_MobTypeUndead,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      46_018_500,
				stats.Armor:       10643,
				stats.AttackPower: 805,
				stats.BlockValue:  76,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       1.50,
			MinBaseDamage:    88072, // Est 96282 minimum debuffed Unmit
			SuppressDodge:    true,
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
			DamageSpread:     0.5,
			TargetInputs:     SindragosaTargetInputs(),
		},
		AI: NewSindragosa25HAI(),
	})
	core.AddPresetEncounter("Sindragosa (Heroic)", []string{
		bossPrefix + "/Sindragosa (Heroic)",
	})
}

type Sindragosa25HAI struct {
	Target *core.Target

	ChilledToTheBone  *core.Spell
	FrostAura         *core.Spell
	FrostBreath       *core.Spell
	FrostBreathDebuff *core.Aura

	IncludeMysticBuffet bool
	MysticBuffetAuras   []*core.Aura
}

func SindragosaTargetInputs() []*proto.TargetInput {
	return []*proto.TargetInput{
		{
			Label:     "Include Mystic Buffet",
			Tooltip:   "Model the ramping magic damage taken debuff applied during Phase 3 of the encounter, in addition to the normal Phase 1 mechanics.",
			InputType: proto.InputType_Bool,
			BoolValue: false,
		},
	}
}

func NewSindragosa25HAI() core.AIFactory {
	return func() core.TargetAI {
		return &Sindragosa25HAI{}
	}
}

func (ai *Sindragosa25HAI) Initialize(target *core.Target, config *proto.Target) {
	ai.Target = target
	ai.IncludeMysticBuffet = config.TargetInputs[0].BoolValue

	ai.registerFrostAuraSpell(target)
	ai.registerFrostBreathSpell(target)
	ai.registerPermeatingChillAura(target)
	ai.registerMysticBuffetAuras()
}

func (ai *Sindragosa25HAI) Reset(sim *core.Simulation) {
	// Randomize time of first Frost Breath under the constraint of preserving the maximum number of possible breaths
	breathPeriod := time.Millisecond * 22680
	maxBreathsPossible := (sim.Duration - time.Millisecond*1500) / breathPeriod
	latestAllowedBreath := sim.Duration - time.Millisecond*1500 - breathPeriod*maxBreathsPossible - time.Millisecond*1620
	firstBreath := core.DurationFromSeconds(sim.RandomFloat("Frost Breath Timing") * latestAllowedBreath.Seconds())

	ai.FrostBreath.CD.Set(firstBreath)
}

func (ai *Sindragosa25HAI) registerPermeatingChillAura(target *core.Target) {
	ai.ChilledToTheBone = target.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 70106},
		SpellSchool:      core.SpellSchoolFrost,
		ProcMask:         core.ProcMaskSpellDamage,
		Flags:            core.SpellFlagNone,
		DamageMultiplier: 1,
		CritMultiplier:   1,
		ThreatMultiplier: 0,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Chilled to the Bone",
				MaxStacks: math.MaxInt32,
				Duration:  time.Second * 8,
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = 1000. * float64(dot.Aura.GetStacks())

				if !isRollover {
					// The player is technically dealing the damage to themselves on each Chilled to the Bone tick, so the ticks use the player's damage modifiers even though we're modeling it as a boss spell for cleaner metrics.
					dot.SnapshotAttackerMultiplier = target.PseudoStats.DamageDealtMultiplier * target.PseudoStats.SchoolDamageDealtMultiplier[dot.Spell.SchoolIndex]
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
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

	for _, party := range ai.Target.Env.Raid.Parties {
		for _, player := range party.PlayersAndPets {
			character := player.GetCharacter()
			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Permeating Chill",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeOrRanged,
				Harmful:    true,
				ProcChance: 0.2,
				ICD:        time.Second * 2,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					ai.ChilledToTheBone.Cast(sim, &character.Unit)
				},
			})
		}
	}
}

func (ai *Sindragosa25HAI) registerMysticBuffetAuras() {
	if !ai.IncludeMysticBuffet {
		return
	}

	ai.MysticBuffetAuras = make([]*core.Aura, 0)
	pendingActions := make([]*core.PendingAction, len(ai.Target.Env.AllUnits))

	for _, raidUnit := range ai.Target.Env.Raid.AllUnits {
		ai.MysticBuffetAuras = append(ai.MysticBuffetAuras, raidUnit.GetOrRegisterAura(core.Aura{
			Label:     "Mystic Buffet",
			ActionID:  core.ActionID{SpellID: 70127},
			MaxStacks: math.MaxInt32,
			Duration:  time.Second * 8,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] /= 1.0 + 0.2*float64(oldStacks)
				aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= 1.0 + 0.2*float64(newStacks)
			},
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				period := time.Second * 6
				numTicks := int(sim.GetRemainingDuration() / period)

				if pendingActions[aura.Unit.UnitIndex] != nil {
					pendingActions[aura.Unit.UnitIndex].Cancel(sim)
				}

				pendingActions[aura.Unit.UnitIndex] = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
					NumTicks: numTicks,
					Period:   period,
					OnAction: func(sim *core.Simulation) {
						aura.Refresh(sim)
						aura.AddStack(sim)
					},
				})
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				pendingActions[aura.Unit.UnitIndex].Cancel(sim)
			},
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				if pendingActions[aura.Unit.UnitIndex] != nil {
					pendingActions[aura.Unit.UnitIndex].Cancel(sim)
				}
			},
		}))
	}
}

func (ai *Sindragosa25HAI) registerFrostAuraSpell(target *core.Target) {
	ai.FrostAura = target.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 70084},
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagNone,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Millisecond * 1620,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Frost Aura",
			},
			NumberOfTicks: math.MaxInt32,
			TickLength:    time.Second * 3,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, 6000, dot.Spell.OutcomeAlwaysHit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Raid.GetActiveUnits() {
				dot := spell.Dot(aoeTarget)
				dot.Apply(sim)
			}

			// Bundle Mystic Buffet application with Frost Aura if requested
			if ai.IncludeMysticBuffet {
				for _, mysticBuffetAura := range ai.MysticBuffetAuras {
					mysticBuffetAura.Activate(sim)
				}
			}
		},
	})
}

func (ai *Sindragosa25HAI) registerFrostBreathSpell(target *core.Target) {
	// Phase 3 uses an intentionally weaker version of the Frost Breath spell, so set up two variants depending on whether Mystic Buffet is being modeled or not
	var spellID int32
	var minRoll float64
	var maxRoll float64

	if ai.IncludeMysticBuffet {
		spellID = 73061
		minRoll = 46250
		maxRoll = 53750
	} else {
		spellID = 69649
		minRoll = 55500
		maxRoll = 64500
	}

	actionID := core.ActionID{SpellID: spellID}

	if ai.Target.CurrentTarget != nil {
		ai.FrostBreathDebuff = ai.Target.CurrentTarget.GetOrRegisterAura(core.Aura{
			Label:     "Frost Breath",
			ActionID:  actionID,
			MaxStacks: math.MaxInt32,
			Duration:  time.Minute,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				aura.Unit.MultiplyAttackSpeed(sim, 1.0+0.5*float64(oldStacks))
				aura.Unit.MultiplyAttackSpeed(sim, 1.0/(1.0+0.5*float64(newStacks)))
			},
		})
	}

	ai.FrostBreath = target.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Second * 20,
			},
			DefaultCast: core.Cast{
				GCD:      time.Millisecond * 1620,
				CastTime: time.Millisecond * 1500,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(minRoll, maxRoll)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
			ai.FrostBreathDebuff.Activate(sim)
			ai.FrostBreathDebuff.AddStack(sim)
		},
	})
}

func (ai *Sindragosa25HAI) ExecuteCustomRotation(sim *core.Simulation) {
	if !ai.Target.GCD.IsReady(sim) {
		return
	}

	// Cast Frost Aura once at the start of the encounter.
	if sim.CurrentTime < time.Millisecond*1620 {
		ai.FrostAura.Cast(sim, &ai.Target.Unit)
		return
	}

	if ai.Target.CurrentTarget != nil && ai.FrostBreath.IsReady(sim) {
		ai.Target.Unit.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+time.Millisecond*1500, false)
		ai.FrostBreath.Cast(sim, ai.Target.CurrentTarget)
		return
	}

	// Sindragosa follows the standard Classic WoW boss AI behavior of evaluating actions on a 1.62 second server tick.
	ai.Target.WaitUntil(sim, sim.CurrentTime+time.Millisecond*1620)
}
