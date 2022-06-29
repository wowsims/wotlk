package feral

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
	"github.com/wowsims/tbc/sim/druid"
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

	cat.setupRotation(feralOptions.Rotation)

	// Passive Cat Form threat reduction
	cat.PseudoStats.ThreatMultiplier *= 0.71

	// Prevents Windfury from applying.
	cat.HasMHWeaponImbue = true

	cat.EnableEnergyBar(100.0, func(sim *core.Simulation) {
		if cat.GCD.IsReady(sim) {
			cat.doRotation(sim)
		}
	})

	cat.EnableAutoAttacks(cat, core.AutoAttackOptions{
		// Base paw weapon.
		MainHand: core.Weapon{
			BaseDamageMin:        43.5,
			BaseDamageMax:        66.5,
			SwingSpeed:           1.0,
			NormalizedSwingSpeed: 1.0,
			SwingDuration:        time.Second,
			CritMultiplier:       cat.MeleeCritMultiplier(),
		},
		AutoSwingMelee: true,
	})

	// Cat Form adds (2 x Level) AP + 1 AP per Agi
	cat.AddStat(stats.AttackPower, 140)
	cat.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Agility,
		ModifiedStat: stats.AttackPower,
		Modifier: func(agility float64, attackPower float64) float64 {
			return attackPower + agility*1
		},
	})

	cat.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.FeralAttackPower,
		ModifiedStat: stats.AttackPower,
		Modifier: func(feralAttackPower float64, attackPower float64) float64 {
			return attackPower + feralAttackPower*1
		},
	})

	return cat
}

type FeralDruid struct {
	*druid.Druid

	Rotation FeralDruidRotation

	readyToShift   bool
	waitingForTick bool
	latency        time.Duration
}

func (cat *FeralDruid) GetDruid() *druid.Druid {
	return cat.Druid
}

func (cat *FeralDruid) Initialize() {
	cat.Druid.Initialize()
	cat.RegisterCatSpells()
	cat.DelayDPSCooldownsForArmorDebuffs()
}

func (cat *FeralDruid) Reset(sim *core.Simulation) {
	cat.Druid.Reset(sim)
	cat.CatFormAura.Activate(sim)
	cat.readyToShift = false
	cat.waitingForTick = false
}
