package enhancement

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (enh *EnhancementShaman) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	// if enh.GCD.IsReady(sim) {
	// 	enh.tryUseGCD(sim)
	// }
}

func (enh *EnhancementShaman) OnGCDReady(sim *core.Simulation) {
	enh.tryUseGCD(sim)
}

func (enh *EnhancementShaman) tryUseGCD(sim *core.Simulation) {
	if enh.TryDropTotems(sim) {
		return
	}
	enh.rotation.DoAction(enh, sim)
}

type Rotation interface {
	DoAction(*EnhancementShaman, *core.Simulation)
	Reset(*EnhancementShaman, *core.Simulation)
}

const prioritySize = 9
const (
	StormstrikeApplyDebuff = iota
	LightningBolt
	Stormstrike
	FlameShock
	Weave
	EarthShock
	LightningShield
	FireNova
	LavaLash
)

type PriorityRotation struct {
	options       *proto.EnhancementShaman_Rotation
	spellPriority [prioritySize]Spell
}

type Cast func(sim *core.Simulation, target *core.Unit) bool
type Condition func(sim *core.Simulation, target *core.Unit) bool

type Spell struct {
	readyAt   func() time.Duration
	cast      Cast
	condition Condition
}

// PRIORITY ROTATION (default)
func (rotation *PriorityRotation) DoAction(enh *EnhancementShaman, sim *core.Simulation) {
	target := enh.CurrentTarget

	upcomingCD := enh.AutoAttacks.NextAttackAt()
	var cast Cast
	for _, spell := range rotation.spellPriority {
		if spell.condition(sim, target) && spell.cast(sim, target) {
			return
		}

		readyAt := spell.readyAt()
		if readyAt > 0 && upcomingCD > readyAt {
			upcomingCD = readyAt
			cast = spell.cast
		}
	}

	enh.WaitUntil(sim, upcomingCD)

	if cast != nil {
		enh.HardcastWaitUntil(sim, upcomingCD, func(sim *core.Simulation, target *core.Unit) {
			enh.GCD.Reset()
			cast(sim, target)
		})
	}
}

func (rotation *PriorityRotation) Reset(enh *EnhancementShaman, sim *core.Simulation) {

}

func NewPriorityRotation(enh *EnhancementShaman, options *proto.EnhancementShaman_Rotation) *PriorityRotation {
	rotation := PriorityRotation{
		options: options,
	}

	rotation.buildPriority(enh)

	return &rotation
}

func (rotation *PriorityRotation) buildPriority(enh *EnhancementShaman) {
	stormstrikeApplyDebuff := Spell{
		condition: func(sim *core.Simulation, target *core.Unit) bool {
			return !enh.StormstrikeDebuffAura(target).IsActive() && enh.Stormstrike.IsReady(sim)
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.Stormstrike.Cast(sim, target)
		},
		readyAt: func() time.Duration {
			return enh.Stormstrike.ReadyAt()
		},
	}

	instantLightningBolt := Spell{
		condition: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.MaelstromWeaponAura.GetStacks() == 5
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.LightningBolt.Cast(sim, target)
		},
		readyAt: func() time.Duration {
			return 0
		},
	}

	stormstrike := Spell{
		condition: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.Stormstrike.IsReady(sim)
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {

			return enh.Stormstrike.Cast(sim, target)
		},
		readyAt: func() time.Duration {
			return enh.Stormstrike.ReadyAt()
		},
	}

	weave := Spell{
		condition: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.MaelstromWeaponAura.GetStacks() >= rotation.options.MaelstromweaponMinStack
		},

		cast: func(sim *core.Simulation, target *core.Unit) bool {
			if rotation.options.LavaburstWeave && enh.CastLavaBurstWeave(sim, target) {
				return true
			}

			if rotation.options.LightningboltWeave && enh.CastLightningBoltWeave(sim, target) {
				return true
			}

			return false
		},
		readyAt: func() time.Duration {
			return 0
		},
	}

	flameShock := Spell{
		condition: func(sim *core.Simulation, target *core.Unit) bool {
			return !enh.FlameShockDot.IsActive() && enh.FlameShock.IsReady(sim)
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.FlameShock.Cast(sim, target)
		},
		readyAt: func() time.Duration {
			return enh.FlameShock.ReadyAt()
		},
	}

	earthShock := Spell{
		condition: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.EarthShock.IsReady(sim)
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.EarthShock.Cast(sim, target)
		},
		readyAt: func() time.Duration {
			return enh.EarthShock.ReadyAt()
		},
	}

	lightningShield := Spell{
		condition: func(sim *core.Simulation, target *core.Unit) bool {
			return !enh.LightningShieldAura.IsActive() && enh.LightningShieldAura != nil
		},
		cast: func(sim *core.Simulation, _ *core.Unit) bool {
			return enh.LightningShield.Cast(sim, nil)
		},
		readyAt: func() time.Duration {
			return 0
		},
	}

	fireNova := Spell{
		condition: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.Totems.Fire != proto.FireTotem_NoFireTotem && enh.FireNova.IsReady(sim) && enh.CurrentMana() > rotation.options.FirenovaManaThreshold
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.FireNova.Cast(sim, target)
		},
		readyAt: func() time.Duration {
			return enh.FireNova.ReadyAt()
		},
	}

	lavaLash := Spell{
		condition: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.LavaLash.IsReady(sim)
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.LavaLash.Cast(sim, target)
		},
		readyAt: func() time.Duration {
			return enh.LavaLash.ReadyAt()
		},
	}

	//This can allow for a custom prio rotation, using a ENUM for default rotation order for now.
	var spellPriority [prioritySize]Spell
	spellPriority[StormstrikeApplyDebuff] = stormstrikeApplyDebuff
	spellPriority[LightningBolt] = instantLightningBolt
	spellPriority[Stormstrike] = stormstrike
	spellPriority[FlameShock] = flameShock
	spellPriority[EarthShock] = earthShock
	spellPriority[LightningShield] = lightningShield
	spellPriority[FireNova] = fireNova
	spellPriority[LavaLash] = lavaLash
	spellPriority[Weave] = weave

	rotation.spellPriority = spellPriority
}

//	CUSTOM ROTATION (advanced) (also WIP).
//TODO: figure out how to do this (probably too complicated to copy hunters)

type AgentAction interface {
	GetActionID() core.ActionID

	GetManaCost() float64

	Cast(sim *core.Simulation) bool
}
