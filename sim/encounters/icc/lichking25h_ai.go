package icc

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func addLichKing25H(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        36597,
			Name:      "Lich King (Heroic)",
			Level:     83,
			MobType:   proto.MobType_MobTypeUndead,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      103_151_165,
				stats.Armor:       10643,
				stats.AttackPower: 805,
				stats.BlockValue:  76,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       1.50,
			MinBaseDamage:    146497, // Est 160154 minimum debuffed Unmit
			SuppressDodge:    true,
			ParryHaste:       false,
			DualWield:        false,
			DualWieldPenalty: false,
			DamageSpread:     0.1557,
			TargetInputs:     make([]*proto.TargetInput, 0),
		},
		AI: NewLichKing25HAI(),
	})
	core.AddPresetEncounter("Lich King (Heroic)", []string{
		bossPrefix + "/Lich King (Heroic)",
	})
}

type LichKing25HAI struct {
	Target *core.Target

	SoulReaper     *core.Spell
	SoulReaperAura *core.Aura
}

func NewLichKing25HAI() core.AIFactory {
	return func() core.TargetAI {
		return &LichKing25HAI{}
	}
}

func (ai *LichKing25HAI) Initialize(target *core.Target, _ *proto.Target) {
	ai.Target = target
	ai.registerSoulReaperSpell(target)
}

func (ai *LichKing25HAI) Reset(*core.Simulation) {
	ai.SoulReaper.CD.Set(ai.SoulReaper.CD.Duration)
}

func (ai *LichKing25HAI) registerSoulReaperSpell(target *core.Target) {
	ai.SoulReaperAura = target.GetOrRegisterAura(core.Aura{
		Label:    "Soul Reaper",
		ActionID: core.ActionID{SpellID: 69410},
		Duration: time.Second * 5,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, 2.0)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, 0.5)
		},
	})

	ai.SoulReaper = target.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 69409},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Second * 30,
			},
			DefaultCast: core.Cast{
				GCD: time.Millisecond * 1620,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Soul Reaper",
				Duration: time.Second * 5,
			},
			NumberOfTicks: 1,
			TickLength:    time.Second * 5,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// Soul Reaper ticks cannot be partially resisted even though the initial hit can, so temporarily change the spell flag accordingly.
				dot.Spell.Flags = core.SpellFlagIgnoreResists

				// Perform damage calculation
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, 70000., dot.Spell.OutcomeAlwaysHit)

				// Activate Haste aura on boss
				ai.SoulReaperAura.Activate(sim)

				// Reset spell flag for next application
				dot.Spell.Flags = core.SpellFlagNone
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// 50% weapon damage
			baseDamage := 0.5 * spell.Unit.AutoAttacks.MH().EnemyWeaponDamage(sim, spell.MeleeAttackPower(), 0.1557)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)

			dot := spell.Dot(target)
			dot.Apply(sim)

			// Soul Reaper application resets the boss melee swing timer based on log analysis
			spell.Unit.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime, false)
		},
	})
}

func (ai *LichKing25HAI) ExecuteCustomRotation(sim *core.Simulation) {
	if !ai.Target.GCD.IsReady(sim) {
		return
	}

	if ai.Target.CurrentTarget != nil {
		if ai.SoulReaper.IsReady(sim) {
			// Based on log analysis, Soul Reaper appears to have a ~75% chance to "proc" on every 1.62 second server tick once it is off cooldown.
			// Note that analysis based only on the cast intervals supported a ~40% proc chance fit. However, many of the apparent delays in Soul Reaper casts are
			// due to Defile and Infest casts that take priority when the cooldowns overlap. Once these CD conflicts are corrected for, the variance in Soul Reaper
			// cast times is a fair bit lower. A more fleshed out boss AI would directly model the Defile and Infest casts, with more sophisticated APLs for tank CD
			// usage in order to delay pre-emptive CDs when the player knows that a 2 second cast will take place before any potential Soul Reaper.
			procRoll := sim.RandomFloat("Soul Reaper AI")

			if procRoll < 0.75 {
				ai.SoulReaper.Cast(sim, ai.Target.CurrentTarget)
				return
			}
		}
	}

	// Lich King follows the standard Classic WoW boss AI behavior of evaluating actions on a 1.62 second server tick.
	ai.Target.WaitUntil(sim, sim.CurrentTime+time.Millisecond*1620)
}
