package healing

import (
	"github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (hpriest *HealingPriest) OnGCDReady(sim *core.Simulation) {
	hpriest.tryUseGCD(sim)
}

func (hpriest *HealingPriest) tryUseGCD(sim *core.Simulation) {
	if hpriest.CustomRotation != nil {
		hpriest.CustomRotation.Cast(sim)
	} else {
		spell := hpriest.chooseSpell(sim)

		if success := spell.Cast(sim, hpriest.CurrentTarget); !success {
			hpriest.WaitForMana(sim, spell.CurCast.Cost)
		}
	}
}

func (hpriest *HealingPriest) chooseSpell(sim *core.Simulation) *core.Spell {
	if !hpriest.RenewHots[hpriest.CurrentTarget.UnitIndex].IsActive() {
		return hpriest.Renew
	} else if hpriest.CanCastPWS(sim, hpriest.CurrentTarget) {
		return hpriest.PowerWordShield
	} else {
		for !hpriest.spellCycle[hpriest.nextCycleIndex].IsReady(sim) {
			hpriest.nextCycleIndex = (hpriest.nextCycleIndex + 1) % len(hpriest.spellCycle)
		}
		spell := hpriest.spellCycle[hpriest.nextCycleIndex]
		hpriest.nextCycleIndex = (hpriest.nextCycleIndex + 1) % len(hpriest.spellCycle)
		return spell
	}
}

func (hpriest *HealingPriest) makeCustomRotation() *common.CustomRotation {
	return common.NewCustomRotation(hpriest.rotation.CustomRotation, hpriest.GetCharacter(), map[int32]common.CustomSpell{
		int32(proto.HealingPriest_Rotation_GreaterHeal): {
			Spell: hpriest.GreaterHeal,
		},
		int32(proto.HealingPriest_Rotation_FlashHeal): {
			Spell: hpriest.FlashHeal,
		},
		int32(proto.HealingPriest_Rotation_Renew): {
			Spell: hpriest.Renew,
			Condition: func(sim *core.Simulation) bool {
				return !hpriest.RenewHots[hpriest.CurrentTarget.UnitIndex].IsActive()
			},
		},
		int32(proto.HealingPriest_Rotation_PowerWordShield): {
			Spell: hpriest.PowerWordShield,
			Condition: func(sim *core.Simulation) bool {
				return hpriest.CanCastPWS(sim, hpriest.CurrentTarget)
			},
		},
		int32(proto.HealingPriest_Rotation_CircleOfHealing): {
			Spell: hpriest.CircleOfHealing,
		},
		int32(proto.HealingPriest_Rotation_PrayerOfHealing): {
			Spell: hpriest.PrayerOfHealing,
		},
		int32(proto.HealingPriest_Rotation_PrayerOfMending): {
			Spell: hpriest.PrayerOfMending,
			Condition: func(sim *core.Simulation) bool {
				return hpriest.PrayerOfMending.IsReady(sim)
			},
		},
		int32(proto.HealingPriest_Rotation_Penance): {
			Spell: hpriest.Penance,
			Condition: func(sim *core.Simulation) bool {
				return hpriest.Penance.IsReady(sim)
			},
		},
		int32(proto.HealingPriest_Rotation_BindingHeal): {
			Spell: hpriest.BindingHeal,
		},
	})
}
