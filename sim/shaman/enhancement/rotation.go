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

type Rotation interface {
	DoAction(*EnhancementShaman, *core.Simulation)
	Reset(*EnhancementShaman, *core.Simulation)
}

//adaptive rotation, shamelessly stolen from elemental shaman
type AdaptiveRotation struct {
}

func (rotation *AdaptiveRotation) DoAction(enh *EnhancementShaman, sim *core.Simulation) {
	target := sim.GetTargetUnit(0)

	// if enh.IsFireNovaCastable(sim) {
	// 	if !enh.FireNova.Cast(sim, target) {
	// 		enh.WaitForMana(sim, enh.FireNova.CurCast.Cost)
	// 	}
	// 	return
	// }

	if enh.Talents.Stormstrike {
		if (enh.StormstrikeDebuffAura(target).GetStacks() > 0) && enh.Stormstrike.IsReady(sim) {
			if !enh.Stormstrike.Cast(sim, target) {
				enh.WaitForMana(sim, enh.Stormstrike.CurCast.Cost)
			}
			return
		}
	}

	if enh.Talents.MaelstromWeapon > 0 {
		if enh.MaelstromWeaponAura.GetStacks() == 5 {
			if !enh.LightningBolt.Cast(sim, target) {
				enh.WaitForMana(sim, enh.LightningBolt.CurCast.Cost)
			}
			return
		}
	}

	if enh.Talents.Stormstrike {
		if enh.Stormstrike.IsReady(sim) {
			if !enh.Stormstrike.Cast(sim, target) {
				enh.WaitForMana(sim, enh.Stormstrike.CurCast.Cost)
			}
			return
		}
	}

	if !enh.FlameShockDot.IsActive() && enh.FlameShock.IsReady(sim) {
		if !enh.FlameShock.Cast(sim, target) {
			enh.WaitForMana(sim, enh.FlameShock.CurCast.Cost)
		}
		return
	}

	if enh.EarthShock.IsReady(sim) {
		if !enh.EarthShock.Cast(sim, target) {
			enh.WaitForMana(sim, enh.EarthShock.CurCast.Cost)
		}
		return
	}

	enh.LightningShield.Cast(sim, nil)

	enh.DoNothing()
	return
}

func (rotation *AdaptiveRotation) Reset(enh *EnhancementShaman, sim *core.Simulation) {

}

func NewAdaptiveRotation(talents *proto.ShamanTalents) *AdaptiveRotation {
	return &AdaptiveRotation{}
}

type AgentAction interface {
	GetActionID() core.ActionID

	GetManaCost() float64

	Cast(sim *core.Simulation) bool
}
