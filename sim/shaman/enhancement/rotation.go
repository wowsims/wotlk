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

type CastType int32

const (
	SpellCast CastType = iota
	SpellWeave
)

func (rotation *AdaptiveRotation) chooseSpell(enh *EnhancementShaman, sim *core.Simulation, target *core.Unit) (*core.Spell, CastType) {
	if enh.Talents.Stormstrike && enh.Stormstrike.IsReady(sim) && !enh.StormstrikeDebuffAura(target).IsActive() {
		return enh.Stormstrike, SpellCast
	}

	if enh.Talents.MaelstromWeapon > 0 && enh.MaelstromWeaponAura.GetStacks() == 5 {
		return enh.LightningBolt, SpellCast
	}

	if enh.Talents.Stormstrike && enh.Stormstrike.IsReady(sim) {
		return enh.Stormstrike, SpellCast
	}

	if !enh.FlameShockDot.IsActive() && enh.FlameShock.IsReady(sim) {
		return enh.FlameShock, SpellCast
	}

	if enh.LavaburstWeave && enh.Talents.MaelstromWeapon > 0 && enh.MaelstromWeaponAura.GetStacks() >= enh.MaelstromWeaveThreshold && enh.LavaBurst.IsReady(sim) {
		return enh.LavaBurst, SpellWeave
	}

	if enh.Talents.MaelstromWeapon > 0 && enh.MaelstromWeaponAura.GetStacks() >= enh.MaelstromWeaveThreshold {
		return enh.LightningBolt, SpellWeave
	}

	if enh.EarthShock.IsReady(sim) {
		return enh.EarthShock, SpellCast
	}

	if !enh.LightningShieldAura.IsActive() && enh.LightningShield.IsReady(sim) {
		return enh.LightningShield, SpellCast
	}

	if enh.Totems.Fire != proto.FireTotem_NoFireTotem && enh.FireNova.IsReady(sim) && enh.CurrentMana() > enh.FireNovaManaThreshold {
		return enh.FireNova, SpellCast
	}

	if enh.Talents.LavaLash && enh.LavaLash.IsReady(sim) {
		return enh.LavaLash, SpellCast
	}

	// if enh.LavaburstWeave {
	// 	if enh.Talents.MaelstromWeapon > 0 && enh.MaelstromWeaponAura.GetStacks() >= 1 && enh.LavaBurst.IsReady(sim) {
	// 		lvbCastTime := enh.ApplyCastSpeed(enh.LavaBurst.DefaultCast.CastTime)
	// 		timeUntilSwing := enh.AutoAttacks.NextAttackAt() - sim.CurrentTime
	// 		if lvbCastTime < timeUntilSwing {
	// 			if !enh.LavaBurst.Cast(sim, target) {
	// 				enh.DoNothing()
	// 			}
	// 			return
	// 		}
	// 	}
	// }

	// if enh.Talents.MaelstromWeapon > 0 && enh.MaelstromWeaponAura.GetStacks() >= 1 {
	// 	lbCastTime := enh.LightningBolt.DefaultCast.CastTime - (time.Millisecond * time.Duration(500*enh.MaelstromWeaponAura.GetStacks()))
	// 	lbCastTime = enh.ApplyCastSpeed(lbCastTime)
	// 	timeUntilSwing := enh.AutoAttacks.NextAttackAt() - sim.CurrentTime
	// 	if sim.CurrentTime > enh.AutoAttacks.NextAttackAt() {
	// 		timeUntilSwing = enh.AutoAttacks.MH.SwingDuration
	// 	}
	// 	if lbCastTime < timeUntilSwing {
	// 		if !enh.LightningBolt.Cast(sim, target) {
	// 			enh.WaitForMana(sim, enh.LightningBolt.CurCast.Cost)
	// 		}
	// 		return
	// 	}
	// }

	//enh.LightningShield.Cast(sim, nil) // if nothing else, refresh lightning shield
	return nil, SpellCast
}

func (rotation *AdaptiveRotation) DoAction(enh *EnhancementShaman, sim *core.Simulation) {
	target := sim.GetTargetUnit(0)

	spell, castType := rotation.chooseSpell(enh, sim, target)

	if spell == nil {
		enh.DoNothing()
		return
	}

	switch castType {
	case SpellCast:
		if !spell.Cast(sim, target) {
			enh.WaitForMana(sim, spell.CurCast.Cost)
		}
		return
	case SpellWeave:
		spellCastTime := enh.ApplyCastSpeed(spell.DefaultCast.CastTime - core.TernaryDuration(enh.IsAffectedByMaelstromStacks(spell), time.Microsecond*time.Duration(500*enh.MaelstromWeaponAura.GetStacks()), 0))
		timeUntilSwing := enh.AutoAttacks.NextAttackAt() - sim.CurrentTime
		if sim.CurrentTime > enh.AutoAttacks.NextAttackAt() {
			timeUntilSwing = enh.AutoAttacks.MH.SwingDuration
		}
		if spellCastTime < timeUntilSwing {
			if !spell.Cast(sim, target) {
				enh.WaitForMana(sim, spell.CurCast.Cost)
			}
		}
		return
	}

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
