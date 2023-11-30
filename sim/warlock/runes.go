package warlock

import (
	"github.com/wowsims/classic/sod/sim/core"
	"github.com/wowsims/classic/sod/sim/core/proto"
	"github.com/wowsims/classic/sod/sim/core/stats"
)

func (warlock *Warlock) ApplyRunes() {
	warlock.applyDemonicTactics()
}

func (warlock *Warlock) everlastingAfflictionRefresh(sim *core.Simulation, target *core.Unit) {
	if !warlock.HasRune(proto.WarlockRune_RuneLegsEverlastingAffliction) {
		return
	}

	if warlock.Corruption.Dot(target).IsActive() {
		warlock.Corruption.Dot(target).Rollover(sim)
	}
}

func (warlock *Warlock) applyDemonicTactics() {
	if !warlock.HasRune(proto.WarlockRune_RuneChestDemonicTactics) {
		return
	}

	warlock.AddStat(stats.SpellCrit, 10*core.SpellCritRatingPerCritChance)
}

// TODO: Classic warlock demo pact rune
// func (warlock *Warlock) updateDPASP(sim *core.Simulation) {
// 	if sim.CurrentTime < 0 {
// 		return
// 	}

// 	dpspCurrent := warlock.DemonicPactAura.ExclusiveEffects[0].Priority
// 	currentTimeJump := sim.CurrentTime.Seconds() - warlock.PreviousTime.Seconds()

// 	if currentTimeJump > 0 {
// 		warlock.DPSPAggregate += dpspCurrent * currentTimeJump
// 		warlock.Metrics.UpdateDpasp(dpspCurrent * currentTimeJump)

// 		if sim.Log != nil {
// 			warlock.Log(sim, "[Info] Demonic Pact spell power bonus average [%.0f]",
// 				warlock.DPSPAggregate/sim.CurrentTime.Seconds())
// 		}
// 	}

// 	warlock.PreviousTime = sim.CurrentTime
// }

// func (warlock *Warlock) setupDemonicPact() {
// 	if warlock.Talents.DemonicPact == 0 {
// 		return
// 	}

// 	dpMult := 0.02 * float64(warlock.Talents.DemonicPact)
// 	warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1. + dpMult
// 	warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1. + dpMult
// 	warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexArcane] *= 1. + dpMult
// 	warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] *= 1. + dpMult
// 	warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= 1. + dpMult

// 	if warlock.Options.Summon == proto.Warlock_Options_NoSummon {
// 		return
// 	}

// 	icd := core.Cooldown{
// 		Timer:    warlock.NewTimer(),
// 		Duration: 1 * time.Second,
// 	}

// 	var demonicPactAuras [25]*core.Aura
// 	for _, party := range warlock.Party.Raid.Parties {
// 		for _, player := range party.Players {
// 			demonicPactAuras[player.GetCharacter().Index] = core.DemonicPactAura(player.GetCharacter())
// 		}
// 	}
// 	warlock.DemonicPactAura = demonicPactAuras[warlock.Index]

// 	warlock.Pet.RegisterAura(core.Aura{
// 		Label:    "Demonic Pact Hidden Aura",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			warlock.PreviousTime = 0
// 			aura.Activate(sim)
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			warlock.updateDPASP(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if !result.DidCrit() || !icd.IsReady(sim) {
// 				return
// 			}

// 			icd.Use(sim)

// 			lastBonus := 0.0
// 			if warlock.DemonicPactAura.IsActive() {
// 				lastBonus = warlock.DemonicPactAura.ExclusiveEffects[0].Priority
// 			}
// 			newSPBonus := math.Round(dpMult * warlock.GetStat(stats.SpellPower))

// 			if warlock.DemonicPactAura.RemainingDuration(sim) < 10*time.Second || newSPBonus >= lastBonus {
// 				warlock.updateDPASP(sim)
// 				for _, dpAura := range demonicPactAuras {
// 					if dpAura != nil {
// 						dpAura.ExclusiveEffects[0].SetPriority(sim, newSPBonus)
// 						dpAura.Activate(sim)
// 					}
// 				}
// 			}
// 		},
// 	})
// }
