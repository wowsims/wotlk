package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func (druid *Druid) ThickHideMultiplier() float64 {
	thickHideMulti := 1.0

	if druid.Talents.ThickHide > 0 {
		thickHideMulti += 0.04 + 0.03*float64(druid.Talents.ThickHide-1)
	}

	return thickHideMulti
}

func (druid *Druid) BearArmorMultiplier() float64 {
	sotfMulti := 1.0 + 0.33/3.0
	return 4.7 * sotfMulti
}

func (druid *Druid) ApplyTalents() {
	druid.ApplyEquipScaling(stats.Armor, druid.ThickHideMultiplier())

	if druid.Talents.HeartOfTheWild > 0 {
		bonus := 0.04 * float64(druid.Talents.HeartOfTheWild)
		druid.MultiplyStat(stats.Intellect, 1.0+bonus)
	}

	if druid.Talents.ImprovedMarkOfTheWild > 0 {
		bonus := 0.07 * float64(druid.Talents.ImprovedMarkOfTheWild)
		druid.MultiplyStat(stats.Stamina, 1.0+bonus)
		druid.MultiplyStat(stats.Strength, 1.0+bonus)
		druid.MultiplyStat(stats.Agility, 1.0+bonus)
		druid.MultiplyStat(stats.Intellect, 1.0+bonus)
		druid.MultiplyStat(stats.Spirit, 1.0+bonus)
	}

	druid.setupNaturesGrace()
	druid.applyMoonkinForm()
}

func (druid *Druid) setupNaturesGrace() {
	if !druid.Talents.NaturesGrace {
		return
	}

	druid.NaturesGraceProcAura = druid.RegisterAura(core.Aura{
		Label:    "Natures Grace Proc",
		ActionID: core.ActionID{SpellID: 16886},
		Duration: time.Second * 15,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if aura.TimeActive(sim) >= spell.CastTime() {
				aura.Deactivate(sim)
			}
		},
	})

	druid.RegisterAura(core.Aura{
		Label:    "Natures Grace",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() {
				return
			}
			druid.NaturesGraceProcAura.Activate(sim)
		},
	})
}

// func (druid *Druid) registerNaturesSwiftnessCD() {
// 	if !druid.Talents.NaturesSwiftness {
// 		return
// 	}
// 	actionID := core.ActionID{SpellID: 17116}

// 	var nsAura *core.Aura
// 	nsSpell := druid.RegisterSpell(Humanoid|Moonkin|Tree, core.SpellConfig{
// 		ActionID: actionID,
// 		Flags:    core.SpellFlagNoOnCastComplete,
// 		Cast: core.CastConfig{
// 			CD: core.Cooldown{
// 				Timer:    druid.NewTimer(),
// 				Duration: time.Minute * 3,
// 			},
// 		},
// 		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
// 			nsAura.Activate(sim)
// 		},
// 	})

// 	nsAura = druid.RegisterAura(core.Aura{
// 		Label:    "Natures Swiftness",
// 		ActionID: actionID,
// 		Duration: core.NeverExpires,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			if druid.Starfire != nil {
// 				druid.Starfire.CastTimeMultiplier -= 1
// 			}
// 			if druid.Wrath != nil {
// 				druid.Wrath.CastTimeMultiplier -= 1
// 			}
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			if druid.Starfire != nil {
// 				druid.Starfire.CastTimeMultiplier += 1
// 			}
// 			if druid.Wrath != nil {
// 				druid.Wrath.CastTimeMultiplier += 1
// 			}
// 		},
// 		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
// 			if !druid.Wrath.IsEqual(spell) && !druid.Starfire.IsEqual(spell) {
// 				return
// 			}

// 			// Remove the buff and put skill on CD
// 			aura.Deactivate(sim)
// 			nsSpell.CD.Use(sim)
// 			druid.UpdateMajorCooldowns()
// 		},
// 	})

// 	druid.AddMajorCooldown(core.MajorCooldown{
// 		Spell: nsSpell.Spell,
// 		Type:  core.CooldownTypeDPS,
// 		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
// 			// Don't use NS unless we're casting a full-length starfire or wrath.
// 			return !character.HasTemporarySpellCastSpeedIncrease()
// 		},
// 	})
// }

// TODO: Classic bear
// func (druid *Druid) applyPrimalFury() {
// 	if druid.Talents.PrimalFury == 0 {
// 		return
// 	}

// 	procChance := []float64{0, 0.5, 1}[druid.Talents.PrimalFury]
// 	actionID := core.ActionID{SpellID: 37117}
// 	rageMetrics := druid.NewRageMetrics(actionID)
// 	cpMetrics := druid.NewComboPointMetrics(actionID)

// 	druid.RegisterAura(core.Aura{
// 		Label:    "Primal Fury",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if druid.InForm(Bear) {
// 				if result.Outcome.Matches(core.OutcomeCrit) {
// 					if sim.Proc(procChance, "Primal Fury") {
// 						druid.AddRage(sim, 5, rageMetrics)
// 					}
// 				}
// 			} else if druid.InForm(Cat) {
// 				if druid.IsMangle(spell) || druid.Shred.IsEqual(spell) || druid.Rake.IsEqual(spell) {
// 					if result.Outcome.Matches(core.OutcomeCrit) {
// 						if sim.Proc(procChance, "Primal Fury") {
// 							druid.AddComboPoints(sim, 1, cpMetrics)
// 						}
// 					}
// 				}
// 			}
// 		},
// 	})
// }

// TODO: Class druid omen
// func (druid *Druid) applyOmenOfClarity() {
// 	if !druid.Talents.OmenOfClarity {
// 		return
// 	}

// 	var affectedSpells []*DruidSpell
// 	druid.ClearcastingAura = druid.RegisterAura(core.Aura{
// 		Label:    "Clearcasting",
// 		ActionID: core.ActionID{SpellID: 16870},
// 		Duration: time.Second * 15,
// 		OnInit: func(aura *core.Aura, sim *core.Simulation) {
// 			affectedSpells = core.FilterSlice([]*DruidSpell{
// 				druid.DemoralizingRoar,
// 				druid.FerociousBite,
// 				druid.Lacerate,
// 				druid.MangleBear,
// 				druid.MangleCat,
// 				druid.Maul,
// 				druid.Rake,
// 				druid.Rip,
// 				druid.Shred,
// 				druid.SwipeBear,
// 				druid.SwipeCat,
// 			}, func(spell *DruidSpell) bool { return spell != nil })
// 		},
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			for _, spell := range affectedSpells {
// 				spell.CostMultiplier -= 1
// 			}
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			for _, spell := range affectedSpells {
// 				spell.CostMultiplier += 1
// 			}
// 		},
// 		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
// 			if aura.RemainingDuration(sim) == aura.Duration {
// 				// OnCastComplete is called after OnSpellHitDealt / etc, so don't deactivate
// 				// if it was just activated.
// 				return
// 			}

// 			for _, as := range affectedSpells {
// 				if as.IsEqual(spell) {
// 					aura.Deactivate(sim)
// 					break
// 				}
// 			}
// 		},
// 	})

// 	if !druid.Talents.OmenOfClarity {
// 		return
// 	}

// 	druid.ProcOoc = func(sim *core.Simulation) {
// 		druid.ClearcastingAura.Activate(sim)
// 		if lasherweave2P != nil {
// 			lasherweave2P.Activate(sim)
// 		}
// 	}

// 	druid.RegisterAura(core.Aura{
// 		Label:    "Omen of Clarity",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if !result.Landed() {
// 				return
// 			}

// 			// https://github.com/JamminL/wotlk-classic-bugs/issues/66#issuecomment-1182017571
// 			if druid.HurricaneTickSpell.IsEqual(spell) {
// 				curCastTickSpeed := spell.CurCast.ChannelTime.Seconds() / 10
// 				hurricaneCoeff := 1.0 - (7.0 / 9.0)
// 				spellCoeff := hurricaneCoeff * curCastTickSpeed
// 				chanceToProc := ((1.5 / 60) * 3.5) * spellCoeff
// 				if sim.RandomFloat("Clearcasting") < chanceToProc {
// 					druid.ProcOoc(sim)
// 				}
// 			} else if druid.AutoAttacks.PPMProc(sim, 3.5, core.ProcMaskMeleeWhiteHit, "Omen of Clarity", spell) { // Melee
// 				druid.ProcOoc(sim)
// 			} else if spell.Flags.Matches(SpellFlagOmenTrigger) { // Spells
// 				// Heavily based on comment here
// 				// https://github.com/JamminL/wotlk-classic-bugs/issues/66#issuecomment-1182017571
// 				// Instants are treated as 1.5
// 				// Uses current cast time rather than default cast time (PPM is constant with haste)
// 				castTime := spell.CurCast.CastTime.Seconds()
// 				if castTime == 0 {
// 					castTime = 1.5
// 				}

// 				chanceToProc := (castTime / 60) * 3.5
// 				if druid.Typhoon.IsEqual(spell) { // Add Typhoon
// 					chanceToProc *= 0.25
// 				} else if druid.Moonfire.IsEqual(spell) { // Add Moonfire
// 					chanceToProc *= 0.076
// 				} else if druid.GiftOfTheWild.IsEqual(spell) { // Add Gift of the Wild
// 					// the above comment says it's 0.0875 * (1-0.924) which apparently is out-dated,
// 					// there is no longer an instant suppression factor
// 					// we assume 30 targets (25man + pets)
// 					chanceToProc = 1 - math.Pow(1-chanceToProc, 30)
// 				} else {
// 					chanceToProc *= 0.666
// 				}
// 				if sim.RandomFloat("Clearcasting") < chanceToProc {
// 					druid.ProcOoc(sim)
// 				}
// 			}
// 		},
// 	})
// }
