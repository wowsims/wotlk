package sod

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = false

	// Proc effects. Keep these in order by item ID.

	//MCP
	core.NewItemEffect(9449, func(agent core.Agent) {
		character := agent.GetCharacter()

		// Assumes that the user will swap pummelers to have the buff for the whole fight.
		character.AddStat(stats.MeleeHaste, 500)
	})

	//Thunderfury
	core.NewItemEffect(19019, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(19019)
		ppmm := character.AutoAttacks.NewPPMManager(6.0, procMask)

		procActionID := core.ActionID{SpellID: 21992}

		singleTargetSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    procActionID.WithTag(1),
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 0.5,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 300, spell.OutcomeMagicHitAndCrit)
			},
		})

		makeDebuffAura := func(target *core.Unit) *core.Aura {
			return target.GetOrRegisterAura(core.Aura{
				Label:    "Thunderfury",
				ActionID: procActionID,
				Duration: time.Second * 12,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					target.AddStatDynamic(sim, stats.NatureResistance, -25)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					target.AddStatDynamic(sim, stats.NatureResistance, 25)
				},
			})
		}

		numHits := min(5, character.Env.GetNumTargets())
		debuffAuras := make([]*core.Aura, len(character.Env.Encounter.TargetUnits))
		for i, target := range character.Env.Encounter.TargetUnits {
			debuffAuras[i] = makeDebuffAura(target)
		}

		bounceSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    procActionID.WithTag(2),
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,

			ThreatMultiplier: 1,
			FlatThreatBonus:  63,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				curTarget := target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					result := spell.CalcDamage(sim, curTarget, 0, spell.OutcomeMagicHit)
					if result.Landed() {
						debuffAuras[target.Index].Activate(sim)
					}
					spell.DealDamage(sim, result)
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
			},
		})

		character.RegisterAura(core.Aura{
			Label:    "Thunderfury",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Thunderfury") {
					singleTargetSpell.Cast(sim, result.Target)
					bounceSpell.Cast(sim, result.Target)
				}
			},
		})
	})

	// Void-Touched Leather Gloves
	core.NewItemEffect(211423, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 429867}

		buffAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Void Madness",
			ActionID: actionID,
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1.1)
				character.PseudoStats.ThreatMultiplier *= 1.2
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1.0/1.1)
				character.PseudoStats.ThreatMultiplier /= 1.2
			},
		})

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Flags:    core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 10,
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    activationSpell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	core.AddEffectsToTest = true
}
