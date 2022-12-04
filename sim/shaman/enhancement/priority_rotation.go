package enhancement

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

// Default Priority Order
const (
	LightningBolt = iota
	StormstrikeApplyDebuff
	WeaveLavaBurst
	WeaveLightningBolt
	MagmaTotem
	Stormstrike
	FlameShock
	EarthShock
	FrostShock
	FireNova
	DelayedWeave
	LightningShield
	DropAllTotems
	LavaLash
	NumberSpells // Used to get the max number of spells in the prio list, keep at bottoom
)

type PriorityRotation struct {
	options       *proto.EnhancementShaman_Rotation
	spellPriority []Spell
}

type Cast func(sim *core.Simulation, target *core.Unit) bool
type Condition func(sim *core.Simulation, target *core.Unit) bool
type ReadyAt func() time.Duration

// Holds all the spell info we need to make decisions
type Spell struct {
	readyAt ReadyAt
	cast    Cast
	// Must pass this check to cast or use readyAt, a special condition to be met
	condition Condition
}

func NewPriorityRotation(enh *EnhancementShaman, options *proto.EnhancementShaman_Rotation) *PriorityRotation {
	rotation := &PriorityRotation{
		options: options,
	}

	rotation.buildPriorityRotation(enh)

	return rotation
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
			//TODO add in SS delay so we don't lose flametongues, if Last attack = current time
			return enh.Stormstrike.IsReady(sim) && enh.Stormstrike.Cast(sim, target)
		},
		readyAt: func() time.Duration {
			return enh.Stormstrike.ReadyAt()
		},
	}

	weaveLightningBolt := Spell{
		condition: func(sim *core.Simulation, target *core.Unit) bool {
			return rotation.options.LightningboltWeave && enh.MaelstromWeaponAura.GetStacks() >= rotation.options.MaelstromweaponMinStack && enh.CurrentMana() >= enh.LightningBolt.CurCast.Cost
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {
			reactionTime := time.Millisecond * time.Duration(rotation.options.AutoWeaveDelay)
			return enh.CastLightningBoltWeave(sim, reactionTime)
		},
		readyAt: func() time.Duration {
			return 0
		},
	}

	weaveLavaBurst := Spell{
		condition: func(sim *core.Simulation, target *core.Unit) bool {
			if enh.CurrentMana() < enh.LavaBurst.CurCast.Cost {
				return false
			}

			return rotation.options.LavaburstWeave && enh.MaelstromWeaponAura.GetStacks() >= rotation.options.MaelstromweaponMinStack
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {
			reactionTime := time.Millisecond * time.Duration(rotation.options.AutoWeaveDelay)
			return rotation.options.LavaburstWeave && enh.LavaBurst.IsReady(sim) && enh.CastLavaBurstWeave(sim, reactionTime)
		},
		readyAt: func() time.Duration {
			return 0
		},
	}

	flameShock := Spell{
		condition: func(sim *core.Simulation, target *core.Unit) bool {
			return rotation.options.WeaveFlameShock && enh.FlameShockDot.RemainingDuration(sim) <= time.Duration(rotation.options.FlameShockClipTicks)*enh.FlameShockDot.TickLength
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
			return rotation.options.PrimaryShock == proto.EnhancementShaman_Rotation_Earth
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.EarthShock.IsReady(sim) && enh.EarthShock.Cast(sim, target)
		},
		readyAt: func() time.Duration {
			return enh.EarthShock.ReadyAt()
		},
	}

	frostShock := Spell{
		condition: func(sim *core.Simulation, target *core.Unit) bool {
			return rotation.options.PrimaryShock == proto.EnhancementShaman_Rotation_Frost
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.FrostShock.IsReady(sim) && enh.FrostShock.Cast(sim, target)
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
			if enh.Totems.Fire == proto.FireTotem_NoFireTotem {
				return false
			}

			return enh.NextTotemDrops[2] > sim.CurrentTime && enh.CurrentMana() > rotation.options.FirenovaManaThreshold
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
			//TODO add in LL delay so we don't lose flametongues, if Last attack = current time
			return enh.LavaLash.IsReady(sim) && enh.LavaLash.Cast(sim, target)
		},
		readyAt: func() time.Duration {
			return enh.LavaLash.ReadyAt()
		},
	}

	magmaTotem := Spell{
		condition: func(sim *core.Simulation, target *core.Unit) bool {
			//TODO : rework totems to make them easier to work with.
			return enh.Totems.Fire == proto.FireTotem_MagmaTotem && enh.NextTotemDrops[2] <= sim.CurrentTime
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.MagmaTotem.Cast(sim, target)
		},
		readyAt: func() time.Duration {
			return 0
		},
	}

	dropAllTotems := Spell{
		condition: func(sim *core.Simulation, target *core.Unit) bool {
			return true
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {
			// TODO : totems need to be dropped all at once.
			return enh.TryDropTotems(sim)
		},
		readyAt: func() time.Duration {
			return 0
		},
	}

	delayedWeave := Spell{
		condition: func(sim *core.Simulation, target *core.Unit) bool {
			if enh.CurrentMana() < enh.LightningBolt.CurCast.Cost {
				return false
			}

			return rotation.options.LightningboltWeave && rotation.options.DelayGcdWeave > 0 && enh.MaelstromWeaponAura.GetStacks() >= rotation.options.MaelstromweaponMinStack
		},
		cast: func(sim *core.Simulation, target *core.Unit) bool {
			timeUntilSwing := enh.AutoAttacks.NextAttackAt() - sim.CurrentTime
			if timeUntilSwing <= time.Millisecond*time.Duration(rotation.options.DelayGcdWeave) && timeUntilSwing != 0 {
				delay := enh.AutoAttacks.NextAttackAt() + time.Millisecond*100
				if delay < sim.CurrentTime {
					return false
				}

				enh.HardcastWaitUntil(sim, delay, func(_ *core.Simulation, _ *core.Unit) {
					enh.GCD.Reset()
					enh.CastLightningBoltWeave(sim, 0)
				})

				enh.WaitUntil(sim, delay)
				return true
			}

			return false
		},
		readyAt: func() time.Duration {
			return 0
		},
	}

	//Normal Priority Rotation
	var spellPriority []Spell
	if rotation.options.RotationType == proto.EnhancementShaman_Rotation_Priority {
		spellPriority = make([]Spell, NumberSpells)
		spellPriority[StormstrikeApplyDebuff] = stormstrikeApplyDebuff
		spellPriority[LightningBolt] = instantLightningBolt
		spellPriority[Stormstrike] = stormstrike
		spellPriority[FlameShock] = flameShock
		spellPriority[EarthShock] = earthShock
		spellPriority[LightningShield] = lightningShield
		spellPriority[FireNova] = fireNova
		spellPriority[LavaLash] = lavaLash
		spellPriority[WeaveLightningBolt] = weaveLightningBolt
		spellPriority[FrostShock] = frostShock
		spellPriority[WeaveLavaBurst] = weaveLavaBurst
		spellPriority[MagmaTotem] = magmaTotem
		spellPriority[DropAllTotems] = dropAllTotems
		spellPriority[DelayedWeave] = delayedWeave
	}

	//Custom Priority Rotation
	if rotation.options.RotationType == proto.EnhancementShaman_Rotation_Custom && rotation.options.CustomRotation != nil {
		spellPriority := make([]Spell, 0, len(rotation.options.CustomRotation.Spells))

		// Turn weaving off, will enable them if they have been added.
		rotation.options.LightningboltWeave = false
		rotation.options.LavaburstWeave = false
		rotation.options.WeaveFlameShock = false
		rotation.options.PrimaryShock = proto.EnhancementShaman_Rotation_None
		for _, customSpellProto := range rotation.options.CustomRotation.Spells {
			switch customSpellProto.Spell {
			case int32(proto.EnhancementShaman_Rotation_StormstrikeDebuffMissing):
				spellPriority = append(spellPriority, stormstrikeApplyDebuff)
			case int32(proto.EnhancementShaman_Rotation_LightningBolt):
				spellPriority = append(spellPriority, instantLightningBolt)
			case int32(proto.EnhancementShaman_Rotation_LightningBoltWeave):
				rotation.options.LightningboltWeave = true
				spellPriority = append(spellPriority, weaveLightningBolt)
			case int32(proto.EnhancementShaman_Rotation_LightningBoltDelayedWeave):
				rotation.options.LightningboltWeave = true
				spellPriority = append(spellPriority, delayedWeave)
			case int32(proto.EnhancementShaman_Rotation_Stormstrike):
				spellPriority = append(spellPriority, stormstrike)
			case int32(proto.EnhancementShaman_Rotation_FlameShock):
				rotation.options.WeaveFlameShock = true
				spellPriority = append(spellPriority, flameShock)
			case int32(proto.EnhancementShaman_Rotation_FireNova):
				spellPriority = append(spellPriority, fireNova)
			case int32(proto.EnhancementShaman_Rotation_LavaLash):
				spellPriority = append(spellPriority, lavaLash)
			case int32(proto.EnhancementShaman_Rotation_EarthShock):
				if rotation.options.PrimaryShock == proto.EnhancementShaman_Rotation_None {
					rotation.options.PrimaryShock = proto.EnhancementShaman_Rotation_Earth
					spellPriority = append(spellPriority, earthShock)
				}
			case int32(proto.EnhancementShaman_Rotation_LightningShield):
				spellPriority = append(spellPriority, lightningShield)
			case int32(proto.EnhancementShaman_Rotation_FrostShock):
				if rotation.options.PrimaryShock == proto.EnhancementShaman_Rotation_None {
					rotation.options.PrimaryShock = proto.EnhancementShaman_Rotation_Frost
					spellPriority = append(spellPriority, frostShock)
				}
			case int32(proto.EnhancementShaman_Rotation_LavaBurst):
				rotation.options.LavaburstWeave = true
				spellPriority = append(spellPriority, weaveLavaBurst)
			case int32(proto.EnhancementShaman_Rotation_MagmaTotem):
				spellPriority = append(spellPriority, magmaTotem)
			}
		}

		rotation.spellPriority = spellPriority
		return
	}

	rotation.spellPriority = spellPriority
}

func (rotation *PriorityRotation) DoAction(enh *EnhancementShaman, sim *core.Simulation) {
	target := enh.CurrentTarget

	cheapestSpell := core.TernaryFloat64(enh.LavaLash == nil, enh.LightningBolt.CurCast.Cost, enh.LavaBurst.CurCast.Cost)

	// Incase the array is empty
	if len(rotation.spellPriority) == 0 {
		enh.DoNothing()
	}

	//Mana guard, our cheapest spell.
	if enh.CurrentMana() < cheapestSpell {
		// Lets top off lightning shield stacks before waiting for mana.
		if enh.LightningShieldAura.GetStacks() < 3 {
			enh.LightningShield.Cast(sim, nil)
		}
		enh.WaitForMana(sim, cheapestSpell)
		return
	}

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
		if readyAt > sim.CurrentTime && (upcomingCD > readyAt || upcomingCD < sim.CurrentTime) {
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
