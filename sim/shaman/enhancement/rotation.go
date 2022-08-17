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

//adaptive rotation, shamelessly stolen from elemental shaman
type AdaptiveRotation struct {
}

func (rotation *AdaptiveRotation) DoAction(enh *EnhancementShaman, sim *core.Simulation) {
	target := sim.GetTargetUnit(0)

	//calculate swing times for weaving
	timeUntilSwing := enh.AutoAttacks.NextAttackAt() - sim.CurrentTime
	if sim.CurrentTime > enh.AutoAttacks.NextAttackAt() { //just a little safeguard. possibly unnessecary
		timeUntilSwing = enh.AutoAttacks.MH.SwingDuration
	}

	//Calculate Weaving latency
	latency := time.Duration(enh.WeaveLatency) * time.Millisecond
	previousAttack := sim.CurrentTime - enh.AutoAttacks.PreviousAttackAt
	latency = core.TernaryDuration(previousAttack > latency, 0, latency-previousAttack)

	//calculate cast times for weaving
	lbCastTime := enh.ApplyCastSpeed(enh.LightningBolt.DefaultCast.CastTime-(time.Millisecond*time.Duration(500*enh.MaelstromWeaponAura.GetStacks()))) + latency
	lvbCastTime := enh.ApplyCastSpeed(enh.LavaBurst.DefaultCast.CastTime) + latency

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

	if enh.LavaburstWeave {
		if enh.MaelstromWeaponAura.GetStacks() >= 1 && enh.LavaBurst.IsReady(sim) {
			if lvbCastTime < timeUntilSwing {
				//delay cast if we have latency
				if latency > 0 {
					enh.HardcastWaitUntil(sim, sim.CurrentTime+latency, func(sim *core.Simulation, _ *core.Unit) {
						enh.LavaBurst.Cast(sim, target)
					})
					enh.DoNothing()
				} else if !enh.LavaBurst.Cast(sim, target) {
					enh.DoNothing()
				}
				return
			}
		}
	}

	if enh.LightningboltWeave {
		if enh.MaelstromWeaponAura.GetStacks() >= enh.MaelstromweaponMinStack {
			if lbCastTime < timeUntilSwing {
				//delay cast if we have latency
				if latency > 0 {
					enh.HardcastWaitUntil(sim, sim.CurrentTime+latency, func(sim *core.Simulation, _ *core.Unit) {
						enh.LightningBolt.Cast(sim, target)
					})
					enh.DoNothing()
				} else if !enh.LightningBolt.Cast(sim, target) {
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
		if enh.FireNova.IsReady(sim) && enh.CurrentMana() > 4000 { //TODO: make this configurable
			if !enh.FireNova.Cast(sim, target) {
				enh.DoNothing()
			}
			return
		}
	}

	if enh.Talents.LavaLash { //TODO: potentially raise the prio when certain relics are equipped. tbd if its worth it though
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
