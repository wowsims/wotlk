package healing

import (
	"github.com/wowsims/classic/sim/common"
	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
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
	if !hpriest.Renew.CurHot().IsActive() {
		return hpriest.Renew
	} else if hpriest.PowerWordShield.CanCast(sim, hpriest.CurrentTarget) {
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
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				for _, unit := range hpriest.Env.Raid.AllUnits {
					renewHot := hpriest.Renew.Hot(unit)
					if renewHot != nil && !renewHot.IsActive() {
						success := hpriest.Renew.Cast(sim, unit)
						return success, hpriest.Renew.CurCast.Cost
					}
				}
				panic("No valid Renew target")
			},
			Condition: func(sim *core.Simulation) bool {
				for _, unit := range hpriest.Env.Raid.AllUnits {
					renewHot := hpriest.Renew.Hot(unit)
					if renewHot != nil && !renewHot.IsActive() {
						return true
					}
				}
				return false
			},
		},
		int32(proto.HealingPriest_Rotation_PowerWordShield): {
			Spell: hpriest.PowerWordShield,
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				for _, unit := range hpriest.Env.Raid.AllUnits {
					if hpriest.PowerWordShield.CanCast(sim, unit) {
						success := hpriest.PowerWordShield.Cast(sim, unit)
						return success, hpriest.PowerWordShield.CurCast.Cost
					}
				}
				panic("No valid PowerWordShield target")
			},
			Condition: func(sim *core.Simulation) bool {
				for _, unit := range hpriest.Env.Raid.AllUnits {
					if hpriest.PowerWordShield.CanCast(sim, unit) {
						return true
					}
				}
				return false
			},
		},
		int32(proto.HealingPriest_Rotation_CircleOfHealing): {
			Spell: hpriest.CircleOfHealing,
			Condition: func(sim *core.Simulation) bool {
				return hpriest.CircleOfHealing.IsReady(sim)
			},
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
			Spell: hpriest.PenanceHeal,
			Condition: func(sim *core.Simulation) bool {
				return hpriest.PenanceHeal.IsReady(sim)
			},
		},
		int32(proto.HealingPriest_Rotation_BindingHeal): {
			Spell: hpriest.BindingHeal,
		},
	})
}
