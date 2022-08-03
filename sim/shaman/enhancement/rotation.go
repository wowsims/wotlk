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

func (rotation *AdaptiveRotation) chooseSpell(enh *EnhancementShaman, sim *core.Simulation, target *core.Unit) *core.Spell {
	if enh.Talents.Stormstrike && !enh.StormstrikeDebuffAura(target).IsActive() && enh.Stormstrike.IsReady(sim) {
		return enh.Stormstrike
	}
	if enh.Talents.MaelstromWeapon > 0 && enh.MaelstromWeaponAura.GetStacks() == 5 {
		return enh.LightningBolt
	}
	if enh.Talents.Stormstrike && enh.Stormstrike.IsReady(sim) {
		return enh.Stormstrike
	}
	if !enh.FlameShockDot.IsActive() && enh.FlameShock.IsReady(sim) {
		return enh.FlameShock
	}
	if enh.Talents.MaelstromWeapon > 0 && enh.MaelstromWeaponAura.GetStacks() >= 1 {
		lbCastTime := enh.LightningBolt.DefaultCast.CastTime - (time.Millisecond * time.Duration(500*enh.MaelstromWeaponAura.GetStacks()))
		lbCastTime = enh.ApplyCastSpeed(lbCastTime)
		timeUntilSwing := enh.AutoAttacks.NextAttackAt() - sim.CurrentTime
		if sim.CurrentTime > enh.AutoAttacks.NextAttackAt() {
			timeUntilSwing = enh.AutoAttacks.MH.SwingDuration
		}
		if lbCastTime < timeUntilSwing {
			return enh.LightningBolt
		}
	}
	if enh.EarthShock.IsReady(sim) {
		return enh.EarthShock
	}
	if !enh.LightningShieldAura.IsActive() && enh.LightningShield.IsReady(sim) {
		return enh.LightningShield
	}
	if enh.Totems.Fire != proto.FireTotem_NoFireTotem {
		if enh.FireNova.IsReady(sim) && enh.CurrentMana() > 4000 {
			return enh.FireNova
		}
	}
	if enh.Talents.LavaLash && enh.LavaLash.IsReady(sim) {
		return enh.LavaLash
	}
	return nil
}

func (rotation *AdaptiveRotation) DoAction(enh *EnhancementShaman, sim *core.Simulation) {
	target := sim.GetTargetUnit(0)

	spell := rotation.chooseSpell(enh, sim, target)
	if spell == nil {
		enh.DoNothing()
		return
	}

	if !spell.Cast(sim, target) {
		enh.WaitForMana(sim, spell.CurCast.Cost)
	}
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
