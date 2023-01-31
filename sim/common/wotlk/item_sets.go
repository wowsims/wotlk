package wotlk

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Keep these in alphabetical order.

var ItemSetPurifiedShardOfTheGods = core.NewItemSet(core.ItemSet{
	Name: "Purified Shard of the Gods",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStats(stats.Stats{stats.SpellPower: 222})
			applyShardOfTheGodsDamageProc(agent.GetCharacter(), false)
			applyShardOfTheGodsHealingProc(agent.GetCharacter(), false)
		},
	},
})

var ItemSetShinyShardOfTheGods = core.NewItemSet(core.ItemSet{
	Name: "Shiny Shard of the Gods",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStats(stats.Stats{stats.SpellPower: 250})
			applyShardOfTheGodsDamageProc(agent.GetCharacter(), true)
			applyShardOfTheGodsHealingProc(agent.GetCharacter(), true)
		},
	},
})

func applyShardOfTheGodsDamageProc(character *core.Character, isHeroic bool) {
	name := "Searing Flames"
	actionID := core.ActionID{SpellID: 69729}
	tickAmount := 477.0
	if isHeroic {
		name += " H"
		actionID = core.ActionID{SpellID: 69730}
		tickAmount = 532.0
	}

	dotSpell := character.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: name,
			},
			NumberOfTicks: 6,
			TickLength:    time.Second * 2,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = tickAmount
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},
	})

	core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:       name + " Trigger",
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.25,
		ICD:        time.Second * 50,
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			dotSpell.Dot(result.Target).Apply(sim)
		},
	})
}

func applyShardOfTheGodsHealingProc(character *core.Character, isHeroic bool) {
	name := "Cauterizing Heal"
	actionID := core.ActionID{SpellID: 69733}
	minHeal := 2269.0
	maxHeal := 2773.0
	if isHeroic {
		name += " H"
		actionID = core.ActionID{SpellID: 69734}
		minHeal = 2530.0
		maxHeal = 3092.0
	}

	spell := character.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   character.DefaultHealingCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseHealing := sim.Roll(minHeal, maxHeal)
			spell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)
		},
	})

	core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:       name + " Trigger",
		Callback:   core.CallbackOnHealDealt,
		ProcChance: 0.25,
		ICD:        time.Second * 50,
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			spell.Cast(sim, result.Target)
		},
	})
}

func makeUndeadSet(setName string) *core.ItemSet {
	return core.NewItemSet(core.ItemSet{
		Name: setName,
		Bonuses: map[int32]core.ApplyEffect{
			2: func(agent core.Agent) {
				character := agent.GetCharacter()
				if character.CurrentTarget.MobType == proto.MobType_MobTypeUndead {
					character.PseudoStats.DamageDealtMultiplier *= 1.01
				}
			},
			3: func(agent core.Agent) {
				character := agent.GetCharacter()
				if character.CurrentTarget.MobType == proto.MobType_MobTypeUndead {
					character.PseudoStats.DamageDealtMultiplier *= 1.02 / 1.01
				}
			},
			4: func(agent core.Agent) {
				character := agent.GetCharacter()
				if character.CurrentTarget.MobType == proto.MobType_MobTypeUndead {
					character.PseudoStats.DamageDealtMultiplier *= 1.03 / 1.02
				}
			},
		},
	})
}

var ItemSetBlessedBattlegearOfUndeadSlaying = makeUndeadSet("Blessed Battlegear of Undead Slaying")
var ItemSetBlessedRegaliaOfUndeadCleansing = makeUndeadSet("Blessed Regalia of Undead Cleansing")
var ItemSetBlessedGarbOfTheUndeadSlayer = makeUndeadSet("Blessed Garb of the Undead Slayer")
var ItemSetUndeadSlayersBlessedArmor = makeUndeadSet("Undead Slayer's Blessed Armor")
