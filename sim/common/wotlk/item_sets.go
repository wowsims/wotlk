package wotlk

import (
	"strconv"
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
			agent.GetCharacter().AddStats(stats.Stats{stats.SpellPower: 222, stats.HealingPower: 222})
			applyShardOfTheGods(agent.GetCharacter(), false)
		},
	},
})

var ItemSetShinyShardOfTheGods = core.NewItemSet(core.ItemSet{
	Name: "Shiny Shard of the Gods",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStats(stats.Stats{stats.SpellPower: 250, stats.HealingPower: 250})
			applyShardOfTheGods(agent.GetCharacter(), true)
		},
	},
})

func applyShardOfTheGods(character *core.Character, isHeroic bool) {
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
	})

	target := character.CurrentTarget
	dot := core.NewDot(core.Dot{
		Spell: dotSpell,
		Aura: target.RegisterAura(core.Aura{
			Label:    name + "-" + strconv.Itoa(int(character.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 6,
		TickLength:    time.Second * 2,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage:     core.BaseDamageConfigFlat(tickAmount),
			OutcomeApplier: character.OutcomeFuncTick(),
			IsPeriodic:     true,
		}),
	})

	MakeProcTriggerAura(&character.Unit, ProcTrigger{
		Name:       name + " Trigger",
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.25,
		ICD:        time.Second * 50,
		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellEffect) {
			dot.Apply(sim)
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
