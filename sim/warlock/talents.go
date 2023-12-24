package warlock

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

// TODO: Classic warlock talents
func (warlock *Warlock) ApplyTalents() {
	// Demonic Embrace
	if warlock.Talents.DemonicEmbrace > 0 {
		warlock.MultiplyStat(stats.Stamina, 1+.03*(float64(warlock.Talents.DemonicEmbrace)))
		warlock.MultiplyStat(stats.Spirit, 1-.01*(float64(warlock.Talents.DemonicEmbrace)))
	}

	if warlock.Talents.ImprovedShadowBolt > 0 {
		warlock.applyImprovedShadowBolt()
	}
}

func (warlock *Warlock) applyImprovedShadowBolt() {
	warlock.ImprovedShadowBoltAuras = warlock.NewEnemyAuraArray(func(unit *core.Unit, level int32) *core.Aura {
		return core.ImprovedShadowBoltAura(unit, warlock.Talents.ImprovedShadowBolt)
	})
}

func (warlock *Warlock) applyWeaponImbue() {
	level := warlock.GetCharacter().Level
	if warlock.Options.WeaponImbue == proto.Warlock_Options_Firestone {
		warlock.applyFirestone()
	}
	if warlock.Options.WeaponImbue == proto.Warlock_Options_Spellstone {
		if level >= 55 {
			warlock.AddStat(stats.SpellCrit, 1*core.SpellCritRatingPerCritChance)
		}
	}
}

func (warlock *Warlock) applyFirestone() {
	level := warlock.GetCharacter().Level

	damageMin := 0.0
	damageMax := 0.0

	// TODO: Test for spell scaling
	spellCoeff := 0.0
	spellId := int32(0)

	// TODO: Test PPM
	ppm := warlock.AutoAttacks.NewPPMManager(8, core.ProcMaskMelee)

	if level >= 56 {
		warlock.AddStat(stats.FirePower, 21)
		damageMin = 80.0
		damageMax = 120.0
		spellId = 17949
	} else if level >= 46 {
		warlock.AddStat(stats.FirePower, 17)
		damageMin = 60.0
		damageMax = 90.0
		spellId = 17947
	} else if level >= 36 {
		warlock.AddStat(stats.FirePower, 14)
		damageMin = 40.0
		damageMax = 60.0
		spellId = 17945
	} else if level >= 28 {
		warlock.AddStat(stats.FirePower, 10)
		damageMin = 25.0
		damageMax = 35.0
		spellId = 758
	}

	if level >= 28 {
		fireProcSpell := warlock.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: spellId},
			SpellSchool: core.SpellSchoolFire,
			ProcMask:    core.ProcMaskSpellDamage,

			CritMultiplier:           1.5,
			DamageMultiplier:         1,
			ThreatMultiplier:         1,
			DamageMultiplierAdditive: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(damageMin, damageMax) + spellCoeff*spell.SpellPower()

				// TODO: Test if LoF Buffs this
				//if warlock.LakeOfFireAuras != nil && warlock.LakeOfFireAuras.Get(target).IsActive() {
				//	baseDamage *= 1.4
				//}

				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
			},
		})

		core.MakePermanent(warlock.GetOrRegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: spellId},
			Label:    "Firestone Proc",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if !ppm.Proc(sim, spell.ProcMask, "Firestone Proc") {
					return
				}

				fireProcSpell.Cast(sim, result.Target)
			},
		}))
	}
}

// func (warlock *Warlock) setupPyroclasm() {
// 	if warlock.Talents.Pyroclasm <= 0 {
// 		return
// 	}

// 	pyroclasmDamageBonus := 1 + 0.02*float64(warlock.Talents.Pyroclasm)

// 	warlock.PyroclasmAura = warlock.RegisterAura(core.Aura{
// 		Label:    "Pyroclasm",
// 		ActionID: core.ActionID{SpellID: 63244},
// 		Duration: time.Second * 10,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= pyroclasmDamageBonus
// 			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= pyroclasmDamageBonus
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= pyroclasmDamageBonus
// 			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= pyroclasmDamageBonus
// 		},
// 	})

// 	warlock.RegisterAura(core.Aura{
// 		Label:    "Pyroclasm Talent Hidden Aura",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if (spell == warlock.Conflagrate || spell == warlock.SearingPain) && result.DidCrit() {
// 				warlock.PyroclasmAura.Activate(sim)
// 			}
// 		},
// 	})
// }

// func (warlock *Warlock) setupNightfall() {
// 	if warlock.Talents.Nightfall <= 0 {
// 		return
// 	}

// 	nightfallProcChance := 0.02*float64(warlock.Talents.Nightfall)

// 	warlock.NightfallProcAura = warlock.RegisterAura(core.Aura{
// 		Label:    "Nightfall Shadow Trance",
// 		ActionID: core.ActionID{SpellID: 17941},
// 		Duration: time.Second * 10,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			warlock.ShadowBolt.CastTimeMultiplier -= 1
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			warlock.ShadowBolt.CastTimeMultiplier += 1
// 		},
// 		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
// 			// Check if the shadowbolt was instant cast and not a normal one
// 			if spell == warlock.ShadowBolt && spell.CurCast.CastTime == 0 {
// 				aura.Deactivate(sim)
// 			}
// 		},
// 	})

// 	warlock.RegisterAura(core.Aura{
// 		Label:    "Nightfall Hidden Aura",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if spell == warlock.Corruption { // TODO: also works on drain life...
// 				if sim.Proc(nightfallProcChance, "Nightfall") {
// 					warlock.NightfallProcAura.Activate(sim)
// 				}
// 			}
// 		},
// 	})
// }
