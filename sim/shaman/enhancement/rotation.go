package enhancement

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (enh *EnhancementShaman) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	if enh.GCD.IsReady(sim) {
		enh.tryUseGCD(sim)
	}
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

type AgentAction interface {
	GetActionID() core.ActionID

	GetManaCost() float64

	Cast(sim *core.Simulation) bool
}

type Rotation interface {
	DoAction(*EnhancementShaman, *core.Simulation)
	Reset(*EnhancementShaman, *core.Simulation)
}

///////////////////
// Base Rotation //
///////////////////
type BaseRotation struct {
}

func (rotation *BaseRotation) shouldCastStormstrikeNoDebuff(enh *EnhancementShaman, sim *core.Simulation, target *core.Unit) bool {
	if enh.Talents.Stormstrike {
		return (!enh.StormstrikeDebuffAura(target).IsActive() && enh.Stormstrike.IsReady(sim))
	}
	return false
}

func (rotation *BaseRotation) shouldCastStormstrike(enh *EnhancementShaman, sim *core.Simulation, target *core.Unit) bool {
	if enh.Talents.Stormstrike {
		return (enh.Stormstrike.IsReady(sim))
	}
	return false
}

func (rotation *BaseRotation) shouldCastLightningBoltInstant(enh *EnhancementShaman, sim *core.Simulation, target *core.Unit) bool {
	if enh.Talents.MaelstromWeapon > 0 {
		if enh.MaelstromWeaponAura.GetStacks() == 5 {
			return true
		}
	}
	return false
}

func (rotation *BaseRotation) shouldCastLightningBoltWeave(enh *EnhancementShaman, sim *core.Simulation, target *core.Unit) bool {
	if enh.Talents.MaelstromWeapon > 0 && enh.MaelstromWeaponAura.GetStacks() >= 3 {
		lbCastTime := enh.LightningBolt.CurCast.CastTime
		timeUntilSwing := enh.AutoAttacks.NextAttackAt() - sim.CurrentTime
		if sim.CurrentTime > enh.AutoAttacks.NextAttackAt() {
			timeUntilSwing = enh.AutoAttacks.MH.SwingDuration
		}
		if lbCastTime < timeUntilSwing {
			return true
		}
	}
	return false
}

func (rotation *BaseRotation) shouldCastFlameShock(enh *EnhancementShaman, sim *core.Simulation, target *core.Unit) bool {
	return (!enh.FlameShockDot.IsActive() && enh.FlameShock.IsReady(sim))
}

func (rotation *BaseRotation) shouldCastEarthShock(enh *EnhancementShaman, sim *core.Simulation, target *core.Unit) bool {
	return (enh.EarthShock.IsReady(sim))
}

func (rotation *BaseRotation) shouldCastLightningShield(enh *EnhancementShaman, sim *core.Simulation, target *core.Unit) bool {
	return (!enh.LightningShieldAura.IsActive())
}

func (rotation *BaseRotation) shouldCastFireNova(enh *EnhancementShaman, sim *core.Simulation, target *core.Unit) bool {
	if enh.Totems.Fire != proto.FireTotem_NoFireTotem {
		if enh.FireNova.IsReady(sim) && enh.CurrentMana() > 4000 {
			return true
		}
	}
	return false
}

func (rotation *BaseRotation) shouldCastLavaLash(enh *EnhancementShaman, sim *core.Simulation, target *core.Unit) bool {
	if enh.Talents.LavaLash {
		if enh.LavaLash.IsReady(sim) {
			return true
		}
	}
	return false
}

func (rotation *BaseRotation) DoAction(enh *EnhancementShaman, sim *core.Simulation) {
	target := sim.GetTargetUnit(0)

	if rotation.shouldCastStormstrikeNoDebuff(enh, sim, target) {
		if !enh.Stormstrike.Cast(sim, target) {
			enh.WaitForMana(sim, enh.Stormstrike.CurCast.Cost)
		}
		return
	}

	if rotation.shouldCastLightningBoltInstant(enh, sim, target) {
		if !enh.LightningBolt.Cast(sim, target) {
			enh.WaitForMana(sim, enh.LightningBolt.CurCast.Cost)
		}
		return
	}

	if rotation.shouldCastStormstrike(enh, sim, target) {
		if !enh.Stormstrike.Cast(sim, target) {
			enh.WaitForMana(sim, enh.Stormstrike.CurCast.Cost)
		}
		return
	}

	if rotation.shouldCastFlameShock(enh, sim, target) {
		if !enh.FlameShock.Cast(sim, target) {
			enh.WaitForMana(sim, enh.FlameShock.CurCast.Cost)
		}
		return
	}

	if rotation.shouldCastLightningBoltWeave(enh, sim, target) {
		if !enh.LightningBolt.Cast(sim, target) {
			enh.WaitForMana(sim, enh.LightningBolt.CurCast.Cost)
		}
		return
	}

	if rotation.shouldCastEarthShock(enh, sim, target) {
		if !enh.EarthShock.Cast(sim, target) {
			enh.WaitForMana(sim, enh.EarthShock.CurCast.Cost)
		}
		return
	}

	if rotation.shouldCastLightningShield(enh, sim, target) {
		enh.LightningShield.Cast(sim, nil)
		return
	}

	if rotation.shouldCastFireNova(enh, sim, target) {
		if !enh.FireNova.Cast(sim, target) {
			enh.WaitForMana(sim, enh.FireNova.CurCast.Cost)
		}
		return
	}

	if rotation.shouldCastLavaLash(enh, sim, target) {
		if !enh.LavaLash.Cast(sim, target) {
			enh.WaitForMana(sim, enh.LavaLash.CurCast.Cost)
		}
		return
	}

	enh.DoNothing()
	return
}

func (rotation *BaseRotation) Reset(enh *EnhancementShaman, sim *core.Simulation) {

}

func NewBaseRotation(talents *proto.ShamanTalents) *BaseRotation {
	return &BaseRotation{}
}

//////////////////////////////////////
// Priority Rotation - Configurable //
//////////////////////////////////////
type PriorityRotation struct {
	BaseRotation
}

// func (rotation *PriorityRotation) DoAction(enh *EnhancementShaman, sim *core.Simulation) {

// }

func NewPriorityRotation(talents *proto.ShamanTalents) *PriorityRotation {
	return &PriorityRotation{}
}
