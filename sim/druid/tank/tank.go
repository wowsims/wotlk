package tank

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
	"github.com/wowsims/wotlk/sim/druid"
)

func RegisterFeralTankDruid() {
	core.RegisterAgentFactory(
		proto.Player_FeralTankDruid{},
		proto.Spec_SpecFeralTankDruid,
		func(character core.Character, options proto.Player) core.Agent {
			return NewFeralTankDruid(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_FeralTankDruid)
			if !ok {
				panic("Invalid spec value for Feral Tank Druid!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewFeralTankDruid(character core.Character, options proto.Player) *FeralTankDruid {
	tankOptions := options.GetFeralTankDruid()

	selfBuffs := druid.SelfBuffs{}
	if tankOptions.Options.InnervateTarget != nil {
		selfBuffs.InnervateTarget = *tankOptions.Options.InnervateTarget
	} else {
		selfBuffs.InnervateTarget.TargetIndex = -1
	}

	bear := &FeralTankDruid{
		Druid:    druid.New(character, druid.Bear, selfBuffs, *tankOptions.Talents),
		Rotation: *tankOptions.Rotation,
		Options:  *tankOptions.Options,
	}

	bear.EnableRageBar(bear.Options.StartingRage, 1, func(sim *core.Simulation) {
		if bear.GCD.IsReady(sim) {
			bear.TryUseCooldowns(sim)
			if bear.GCD.IsReady(sim) {
				bear.doRotation(sim)
			}
		}
	})

	bear.EnableAutoAttacks(bear, core.AutoAttackOptions{
		// Base paw weapon.
		MainHand: core.Weapon{
			BaseDamageMin:        130 - 27,
			BaseDamageMax:        130 + 27,
			SwingSpeed:           2.5,
			NormalizedSwingSpeed: 2.5,
			SwingDuration:        time.Millisecond * 2500,
			CritMultiplier:       bear.MeleeCritMultiplier(),
		},
		AutoSwingMelee: true,
		ReplaceMHSwing: func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
			return bear.TryMaul(sim, mhSwingSpell)
		},
	})

	// Prevents Windfury from applying.
	bear.HasMHWeaponImbue = true

	bear.PseudoStats.ThreatMultiplier *= 1.3 + 0.05*float64(bear.Talents.FeralInstinct)

	// Bear Form adds 210 AP (3 * Level).
	bear.AddStat(stats.AttackPower, 3*float64(core.CharacterLevel))

	// Dire Bear Form bonuses.
	bear.AddStatDependency(stats.Stamina, stats.Stamina, 1.0+0.25)

	dps := (((bear.Equip[proto.ItemSlot_ItemSlotMainHand].WeaponDamageMax - bear.Equip[proto.ItemSlot_ItemSlotMainHand].WeaponDamageMin) / 2.0) + bear.Equip[proto.ItemSlot_ItemSlotMainHand].WeaponDamageMin) / bear.Equip[proto.ItemSlot_ItemSlotMainHand].SwingSpeed
	fap := math.Floor((dps - 54.8) * 14)
	if fap > 0 {
		bear.AddStat(stats.AttackPower, fap)
	}

	return bear
}

type FeralTankDruid struct {
	*druid.Druid

	Rotation proto.FeralTankDruid_Rotation
	Options  proto.FeralTankDruid_Options
}

func (bear *FeralTankDruid) GetDruid() *druid.Druid {
	return bear.Druid
}

func (bear *FeralTankDruid) Initialize() {
	bear.Druid.Initialize()
	bear.RegisterBearSpells(float64(bear.Rotation.MaulRageThreshold))
}

func (bear *FeralTankDruid) ApplyGearBonuses() {
	bear.AddStat(stats.Armor, bear.Equip.Stats()[stats.Armor]*4)
}

func (bear *FeralTankDruid) Reset(sim *core.Simulation) {
	bear.Druid.Reset(sim)
	//bear.BearFormAura.Activate(sim)
}
