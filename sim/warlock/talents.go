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

	if warlock.Talents.Emberstorm > 0 {
		warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= .02 * float64(warlock.Talents.Emberstorm)
	}
}

func (warlock *Warlock) applyWeaponImbue() {
	level := warlock.GetCharacter().Level
	// TODO: Classic warlock firestone + fire damage on hit
	if warlock.Options.WeaponImbue == proto.Warlock_Options_Firestone {
		if level >= 56 {
			warlock.AddStat(stats.FirePower, 21)
		} else if level >= 36 {
			warlock.AddStat(stats.FirePower, 14)
		} else if level >= 28 {
			warlock.AddStat(stats.FirePower, 10)
		}
	}
	if warlock.Options.WeaponImbue == proto.Warlock_Options_Spellstone {
		warlock.AddStat(stats.SpellCrit, 1*core.SpellCritRatingPerCritChance)
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

// 	nightfallProcChance := 0.02*float64(warlock.Talents.Nightfall) +
// 		0.04*core.TernaryFloat64(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfCorruption), 1, 0)

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
