package feral

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/druid"
)

func RegisterFeralDruid() {
	core.RegisterAgentFactory(
		proto.Player_FeralDruid{},
		proto.Spec_SpecFeralDruid,
		func(character core.Character, options proto.Player) core.Agent {
			return NewFeralDruid(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_FeralDruid)
			if !ok {
				panic("Invalid spec value for Feral Druid!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewFeralDruid(character core.Character, options proto.Player) *FeralDruid {
	feralOptions := options.GetFeralDruid()

	selfBuffs := druid.SelfBuffs{}
	if feralOptions.Options.InnervateTarget != nil {
		selfBuffs.InnervateTarget = *feralOptions.Options.InnervateTarget
	} else {
		selfBuffs.InnervateTarget.TargetIndex = -1
	}

	cat := &FeralDruid{
		Druid:   druid.New(character, druid.Cat, selfBuffs, *feralOptions.Talents),
		latency: time.Duration(feralOptions.Options.LatencyMs) * time.Millisecond,
	}

	cat.maxRipTicks = cat.MaxRipTicks()
	cat.setupRotation(feralOptions.Rotation)

	// Passive Cat Form threat reduction
	cat.PseudoStats.ThreatMultiplier *= 0.71

	cat.EnableEnergyBar(100.0, cat.OnEnergyGain)

	cat.EnableRageBar(core.RageBarOptions{RageMultiplier: 1, MHSwingSpeed: 2.5}, func(sim *core.Simulation) {})

	cat.EnableAutoAttacks(cat, core.AutoAttackOptions{
		// Base paw weapon.
		MainHand: core.Weapon{
			BaseDamageMin:        43,
			BaseDamageMax:        66,
			SwingSpeed:           1.0,
			NormalizedSwingSpeed: 1.0,
			SwingDuration:        time.Second,
			CritMultiplier:       cat.MeleeCritMultiplier(),
		},
		AutoSwingMelee: true,
	})

	return cat
}

type FeralDruid struct {
	*druid.Druid

	Rotation FeralDruidRotation

	missChance     float64
	readyToShift   bool
	waitingForTick bool
	latency        time.Duration
	maxRipTicks    int
}

func (cat *FeralDruid) GetDruid() *druid.Druid {
	return cat.Druid
}

func (cat *FeralDruid) MissChance() float64 {
	speffect := core.SpellEffect{
		ProcMask:         core.ProcMaskMeleeMHSpecial,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		Target:           cat.CurrentTarget,
	}
	at := cat.AttackTables[cat.CurrentTarget.UnitIndex]
	miss := at.BaseMissChance - speffect.PhysicalHitChance(&cat.Druid.Unit, at)
	dodge := at.BaseDodgeChance - speffect.ExpertisePercentage(&cat.Druid.Unit) - cat.CurrentTarget.PseudoStats.DodgeReduction
	return miss + dodge
}

func (cat *FeralDruid) Initialize() {
	cat.Druid.Initialize()
	cat.RegisterFeralSpells(0)
	cat.DelayDPSCooldownsForArmorDebuffs()
}

func (cat *FeralDruid) Reset(sim *core.Simulation) {
	cat.Druid.Reset(sim)
	cat.Druid.ClearForm(sim)
	cat.CatFormAura.Activate(sim)
	cat.readyToShift = false
	cat.waitingForTick = false
}
