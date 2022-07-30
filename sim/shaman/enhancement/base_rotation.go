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
	enh.DoNothing()
	return
}

func (rotation *BaseRotation) Reset(enh *EnhancementShaman, sim *core.Simulation) {

}

func NewBaseRotation(talents *proto.ShamanTalents) *BaseRotation {
	return &BaseRotation{}
}
