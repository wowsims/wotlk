package enhancement

import (
	"time"

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

type PriorityRotation struct {
	options *proto.EnhancementShaman_Rotation
}

// PRIORITY ROTATION (default)
func (rotation *PriorityRotation) DoAction(enh *EnhancementShaman, sim *core.Simulation) {
	target := enh.CurrentTarget

	//calculate cast times for weaving
	lbCastTime := enh.ApplyCastSpeed(enh.LightningBolt.DefaultCast.CastTime - (time.Millisecond * time.Duration(500*enh.MaelstromWeaponAura.GetStacks())))
	lvbCastTime := enh.ApplyCastSpeed(enh.LavaBurst.DefaultCast.CastTime)
	//calculate swing times for weaving
	timeUntilSwing := enh.AutoAttacks.NextAttackAt() - sim.CurrentTime
	if sim.CurrentTime > enh.AutoAttacks.NextAttackAt() { //just a little safeguard. possibly unnessecary
		timeUntilSwing = enh.AutoAttacks.MH.SwingDuration
	}

	//TODO: find a real prio for these, this is just feelcraft rn
	if enh.Talents.Stormstrike {
		if !enh.StormstrikeDebuffAura(target).IsActive() && enh.Stormstrike.IsReady(sim) {
			if !enh.Stormstrike.Cast(sim, target) {
				enh.WaitForMana(sim, enh.Stormstrike.CurCast.Cost)
			}
			return
		}
	}

	if enh.MaelstromWeaponAura.GetStacks() == 5 {
		if !enh.LightningBolt.Cast(sim, target) {
			enh.WaitForMana(sim, enh.LightningBolt.CurCast.Cost)
		}
		return
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
			enh.DoNothing()
		}
		return
	}

	if rotation.options.LavaburstWeave {
		if enh.MaelstromWeaponAura.GetStacks() >= 1 && enh.LavaBurst.IsReady(sim) {
			if lvbCastTime < timeUntilSwing {
				if !enh.LavaBurst.Cast(sim, target) {
					enh.DoNothing()
				}
				return
			}
		}
	}

	if rotation.options.LightningboltWeave {
		if enh.MaelstromWeaponAura.GetStacks() >= rotation.options.MaelstromweaponMinStack {
			if lbCastTime < timeUntilSwing {
				if !enh.LightningBolt.Cast(sim, target) {
					enh.DoNothing()
				}
				return
			}
		}
	}

	if enh.EarthShock.IsReady(sim) {
		if !enh.EarthShock.Cast(sim, target) {
			enh.DoNothing()
		}
		return
	}

	if !enh.LightningShieldAura.IsActive() && enh.LightningShieldAura != nil {
		enh.LightningShield.Cast(sim, nil)
		return
	}

	if enh.Totems.Fire != proto.FireTotem_NoFireTotem {
		if enh.FireNova.IsReady(sim) && enh.CurrentMana() > rotation.options.FirenovaManaThreshold {
			if !enh.FireNova.Cast(sim, target) {
				enh.DoNothing()
			}
			return
		}
	}

	if enh.Talents.LavaLash && enh.AutoAttacks.IsDualWielding { //TODO: potentially raise the prio when certain relics are equipped. TBD
		if enh.LavaLash.IsReady(sim) {
			if !enh.LavaLash.Cast(sim, target) {
				enh.WaitForMana(sim, enh.LavaLash.CurCast.Cost)
			}
			return
		}
	}

	// if nothing else,
	enh.DoNothing()
	return
}

func (rotation *PriorityRotation) Reset(enh *EnhancementShaman, sim *core.Simulation) {

}

func NewPriorityRotation(talents *proto.ShamanTalents, options *proto.EnhancementShaman_Rotation) *PriorityRotation {
	return &PriorityRotation{
		options: options,
	}
}

//	CUSTOM ROTATION (advanced) (also WIP).
//TODO: figure out how to do this (probably too complicated to copy hunters)

type AgentAction interface {
	GetActionID() core.ActionID

	GetManaCost() float64

	Cast(sim *core.Simulation) bool
}
