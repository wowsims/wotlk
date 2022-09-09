package enhancement

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (enh *EnhancementShaman) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
}

func (enh *EnhancementShaman) OnGCDReady(sim *core.Simulation) {
	enh.tryUseGCD(sim)
}

func (enh *EnhancementShaman) tryUseGCD(sim *core.Simulation) {
	// TODO move this into the rotation, also this uses waitForMana if it was unable to cast the totem
	// that will need to be pulled out so we are not waiting for a magma totem mana cost.
	if enh.TryDropTotems(sim) {
		return
	}
	enh.rotation.DoAction(enh, sim)
}

type Rotation interface {
	DoAction(*EnhancementShaman, *core.Simulation)
	Reset(*EnhancementShaman, *core.Simulation)
}

// PRIORITY ROTATION (default)
func (rotation *PriorityRotation) DoAction(enh *EnhancementShaman, sim *core.Simulation) {
	target := enh.CurrentTarget

	//TODO Add in a mana guard

	// We could choose to not wait for auto attacks if we don't have any MW stacks,
	// this would reduce the amount of DoAction calls by a little bit; might not be a issue though.
	upcomingCD := enh.AutoAttacks.NextAttackAt()
	var cast Cast
	for _, spell := range rotation.spellPriority {
		if !spell.condition(sim, target) {
			continue
		}

		if spell.cast(sim, target) {
			return
		}

		readyAt := spell.readyAt()
		if readyAt > sim.CurrentTime && upcomingCD > readyAt {
			upcomingCD = readyAt
			cast = spell.cast
		}
	}

	//Lets wait on a upcoming CD or AutoAttack
	enh.WaitUntil(sim, upcomingCD)

	//Incase the next auto is our best CD then there are no spells to cast.
	if cast != nil {
		//We have a upcoming CD and we know what to cast lets just do that.
		enh.HardcastWaitUntil(sim, upcomingCD, func(sim *core.Simulation, target *core.Unit) {
			enh.GCD.Reset()
			cast(sim, target)
		})
	}
}

func (rotation *PriorityRotation) Reset(enh *EnhancementShaman, sim *core.Simulation) {

}

// The amont of spells in the priority array
const prioritySize = 9

//Default Priority Order
const (
	LightningBolt = iota
	StormstrikeApplyDebuff
	Weave
	Stormstrike
	FlameShock
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
type ReadyAt func() time.Duration

//Holds all the spell info we need to make decisions
type Spell struct {
	readyAt ReadyAt
	cast    Cast
	// Must pass this check to cast or use readyAt, a special condition to be met
	condition Condition
}

func NewPriorityRotation(enh *EnhancementShaman, options *proto.EnhancementShaman_Rotation) *PriorityRotation {
	rotation := PriorityRotation{
		options: options,
	}

	rotation.buildPriorityRotation(enh)

	return &rotation
}

func (rotation *PriorityRotation) buildPriorityRotation(enh *EnhancementShaman) {
	stormstrikeApplyDebuff := Spell{
		condition: func(sim *core.Simulation, target *core.Unit) bool {
			return !enh.StormstrikeDebuffAura(target).IsActive()
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.Stormstrike.IsReady(sim) && enh.Stormstrike.Cast(sim, target)
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
			//Checking if we learned the spell, ie untalented
			return enh.Stormstrike != nil
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {
			//TODO add in SS delay so we don't loose flametongues, if Last attack = current time
			return enh.Stormstrike.IsReady(sim) && enh.Stormstrike.Cast(sim, target)
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
			//TODO add in check for how much time we have left, IE we don't want to cast a FS with 4 seconds left.
			return !enh.FlameShockDot.IsActive()
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.FlameShock.IsReady(sim) && enh.FlameShock.Cast(sim, target)
		},
		readyAt: func() time.Duration {
			return enh.FlameShock.ReadyAt()
		},
	}

	earthShock := Spell{
		condition: func(sim *core.Simulation, target *core.Unit) bool {
			// No special condition needed
			return true
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.EarthShock.IsReady(sim) && enh.EarthShock.Cast(sim, target)
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
			return enh.Totems.Fire != proto.FireTotem_NoFireTotem && enh.CurrentMana() > rotation.options.FirenovaManaThreshold
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.FireNova.IsReady(sim) && enh.FireNova.Cast(sim, target)
		},
		readyAt: func() time.Duration {
			return enh.FireNova.ReadyAt()
		},
	}

	lavaLash := Spell{
		condition: func(sim *core.Simulation, target *core.Unit) bool {
			//Checking if we learned the spell, ie untalented
			return enh.LavaLash != nil
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {
			//TODO add in LL delay so we don't loose flametongues, if Last attack = current time
			return enh.LavaLash.IsReady(sim) && enh.LavaLash.Cast(sim, target)
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
