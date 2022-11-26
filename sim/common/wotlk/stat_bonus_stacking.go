package wotlk

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type StackingProcAura struct {
	Aura          core.Aura
	BonusPerStack stats.Stats
}

func MakeStackingAura(character *core.Character, config StackingProcAura) *core.Aura {
	var bonusPerStack stats.Stats
	config.Aura.OnInit = func(aura *core.Aura, sim *core.Simulation) {
		bonusPerStack = config.BonusPerStack
	}
	config.Aura.OnStacksChange = func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
		character.AddStatsDynamic(sim, bonusPerStack.Multiply(float64(newStacks-oldStacks)))
	}
	return character.RegisterAura(config.Aura)
}

type StackingStatBonusEffect struct {
	Name       string
	ID         int32
	Bonus      stats.Stats
	Duration   time.Duration
	MaxStacks  int32
	Callback   Callback
	ProcMask   core.ProcMask
	Outcome    core.HitOutcome
	Harmful    bool
	ProcChance float64
}

func newStackingStatBonusEffect(config StackingStatBonusEffect) {
	core.NewItemEffect(config.ID, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := MakeStackingAura(character, StackingProcAura{
			Aura: core.Aura{
				Label:     config.Name + " Proc",
				ActionID:  core.ActionID{ItemID: config.ID},
				Duration:  config.Duration,
				MaxStacks: config.MaxStacks,
			},
			BonusPerStack: config.Bonus,
		})

		MakeProcTriggerAura(&character.Unit, ProcTrigger{
			Name:       config.Name,
			Callback:   config.Callback,
			ProcMask:   config.ProcMask,
			Outcome:    config.Outcome,
			Harmful:    config.Harmful,
			ProcChance: config.ProcChance,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
				procAura.AddStack(sim)
			},
		})
	})
}

type StackingStatBonusCD struct {
	Name        string
	ID          int32
	Bonus       stats.Stats
	Duration    time.Duration
	MaxStacks   int32
	CD          time.Duration
	Callback    Callback
	ProcMask    core.ProcMask
	Outcome     core.HitOutcome
	Harmful     bool
	ProcChance  float64
	IsDefensive bool
}

func newStackingStatBonusCD(config StackingStatBonusCD) {
	core.NewItemEffect(config.ID, func(agent core.Agent) {
		character := agent.GetCharacter()

		buffAura := MakeStackingAura(character, StackingProcAura{
			Aura: core.Aura{
				Label:     config.Name + " Aura",
				ActionID:  core.ActionID{ItemID: config.ID},
				Duration:  config.Duration,
				MaxStacks: config.MaxStacks,
			},
			BonusPerStack: config.Bonus,
		})

		applyProcTriggerCallback(&character.Unit, buffAura, ProcTrigger{
			Name:       config.Name,
			Callback:   config.Callback,
			ProcMask:   config.ProcMask,
			Outcome:    config.Outcome,
			Harmful:    config.Harmful,
			ProcChance: config.ProcChance,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				buffAura.AddStack(sim)
			},
		})

		var sharedTimer *core.Timer
		if config.IsDefensive {
			sharedTimer = character.GetDefensiveTrinketCD()
		} else {
			sharedTimer = character.GetOffensiveTrinketCD()
		}

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: config.ID},
			Flags:    core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: config.CD,
				},
				SharedCD: core.Cooldown{
					Timer:    sharedTimer,
					Duration: config.Duration,
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
		})
	})
}

func init() {
	core.NewItemEffect(38212, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := MakeStackingAura(character, StackingProcAura{
			Aura: core.Aura{
				Label:     "Death Knight's Anguish Proc",
				ActionID:  core.ActionID{ItemID: 38212},
				Duration:  time.Second * 20,
				MaxStacks: 10,
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
						aura.AddStack(sim)
					}
				},
			},
			BonusPerStack: stats.Stats{stats.MeleeCrit: 15, stats.SpellCrit: 15},
		})

		MakeProcTriggerAura(&character.Unit, ProcTrigger{
			Name:       "Death Knight's Anguish",
			Callback:   OnSpellHitDealt,
			ProcMask:   core.ProcMaskMeleeOrRanged,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.1,
			ICD:        time.Second * 45,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			},
		})
	})

	newStackingStatBonusEffect(StackingStatBonusEffect{
		Name:      "Majestic Dragon Figurine",
		ID:        40430,
		Duration:  time.Second * 10,
		MaxStacks: 10,
		Bonus:     stats.Stats{stats.Spirit: 18},
		Callback:  OnCastComplete,
	})
	newStackingStatBonusEffect(StackingStatBonusEffect{
		Name:      "Fury of the Five Fights",
		ID:        40431,
		Duration:  time.Second * 10,
		MaxStacks: 20,
		Bonus:     stats.Stats{stats.AttackPower: 16, stats.RangedAttackPower: 16},
		Callback:  OnSpellHitDealt,
		ProcMask:  core.ProcMaskMeleeOrRanged,
		Harmful:   true,
	})
	newStackingStatBonusEffect(StackingStatBonusEffect{
		Name:      "Illustration of the Dragon Soul",
		ID:        40432,
		Duration:  time.Second * 10,
		MaxStacks: 10,
		Bonus:     stats.Stats{stats.SpellPower: 20},
		Callback:  OnCastComplete,
		Harmful:   true,
	})
	newStackingStatBonusEffect(StackingStatBonusEffect{
		Name:       "DMC Berserker",
		ID:         42989,
		Duration:   time.Second * 12,
		MaxStacks:  3,
		Bonus:      stats.Stats{stats.MeleeCrit: 35, stats.SpellCrit: 35},
		Callback:   OnSpellHitDealt | OnSpellHitTaken,
		Harmful:    true,
		ProcChance: 0.5,
	})
	newStackingStatBonusEffect(StackingStatBonusEffect{
		Name:      "Eye of the Broodmother",
		ID:        45308,
		Duration:  time.Second * 10,
		MaxStacks: 5,
		Bonus:     stats.Stats{stats.SpellPower: 25},
		Callback:  OnCastComplete,
	})

	core.AddEffectsToTest = false

	newStackingStatBonusEffect(StackingStatBonusEffect{
		Name:      "Solance of the Defeated",
		ID:        47041,
		Duration:  time.Second * 10,
		MaxStacks: 8,
		Bonus:     stats.Stats{stats.MP5: 16},
		Callback:  OnCastComplete,
	})
	newStackingStatBonusEffect(StackingStatBonusEffect{
		Name:      "Solance of the Defeated H",
		ID:        47059,
		Duration:  time.Second * 10,
		MaxStacks: 8,
		Bonus:     stats.Stats{stats.MP5: 18},
		Callback:  OnCastComplete,
	})
	newStackingStatBonusEffect(StackingStatBonusEffect{
		Name:      "Solance of the Fallen",
		ID:        47271,
		Duration:  time.Second * 10,
		MaxStacks: 8,
		Bonus:     stats.Stats{stats.MP5: 16},
		Callback:  OnCastComplete,
	})
	newStackingStatBonusEffect(StackingStatBonusEffect{
		Name:      "Solance of the Fallen H",
		ID:        47432,
		Duration:  time.Second * 10,
		MaxStacks: 8,
		Bonus:     stats.Stats{stats.MP5: 18},
		Callback:  OnCastComplete,
	})
	newStackingStatBonusEffect(StackingStatBonusEffect{
		Name:      "Muradin's Spyglass",
		ID:        50340,
		Duration:  time.Second * 10,
		MaxStacks: 10,
		Bonus:     stats.Stats{stats.SpellPower: 18},
		Callback:  OnSpellHitDealt,
		ProcMask:  core.ProcMaskSpellDamage,
		Harmful:   true,
	})
	newStackingStatBonusEffect(StackingStatBonusEffect{
		Name:       "Unidentifiable Organ",
		ID:         50341,
		Duration:   time.Second * 10,
		MaxStacks:  10,
		Bonus:      stats.Stats{stats.Stamina: 24},
		Callback:   OnSpellHitTaken,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.6,
	})
	newStackingStatBonusEffect(StackingStatBonusEffect{
		Name:       "Unidentifiable Organ H",
		ID:         50344,
		Duration:   time.Second * 10,
		MaxStacks:  10,
		Bonus:      stats.Stats{stats.Stamina: 27},
		Callback:   OnSpellHitTaken,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.6,
	})
	newStackingStatBonusEffect(StackingStatBonusEffect{
		Name:      "Muradin's Spyglass H",
		ID:        50345,
		Duration:  time.Second * 10,
		MaxStacks: 10,
		Bonus:     stats.Stats{stats.SpellPower: 20},
		Callback:  OnSpellHitDealt,
		ProcMask:  core.ProcMaskSpellDamage,
		Harmful:   true,
	})
	newStackingStatBonusEffect(StackingStatBonusEffect{
		Name:      "Herkuml War Token",
		ID:        50355,
		Duration:  time.Second * 10,
		MaxStacks: 20,
		Bonus:     stats.Stats{stats.AttackPower: 17, stats.RangedAttackPower: 17},
		Callback:  OnSpellHitDealt,
		ProcMask:  core.ProcMaskMeleeOrRanged,
		Harmful:   true,
	})

	// Stacking CD effects

	newStackingStatBonusCD(StackingStatBonusCD{
		Name:        "Meteorite Crystal",
		ID:          46051,
		Duration:    time.Second * 20,
		MaxStacks:   20,
		Bonus:       stats.Stats{stats.MP5: 60},
		CD:          time.Minute * 2,
		Callback:    OnCastComplete,
		IsDefensive: true,
	})
	newStackingStatBonusCD(StackingStatBonusCD{
		Name:      "Victor's Call",
		ID:        47725,
		Duration:  time.Second * 20,
		MaxStacks: 5,
		Bonus:     stats.Stats{stats.AttackPower: 215, stats.RangedAttackPower: 215},
		CD:        time.Minute * 2,
		Callback:  OnSpellHitDealt,
		ProcMask:  core.ProcMaskMelee,
		Outcome:   core.OutcomeLanded,
	})
	newStackingStatBonusCD(StackingStatBonusCD{
		Name:      "Talisman of Volatile Power",
		ID:        47726,
		Duration:  time.Second * 20,
		MaxStacks: 8,
		Bonus:     stats.Stats{stats.MeleeHaste: 57, stats.SpellHaste: 57},
		CD:        time.Minute * 2,
		Callback:  OnSpellHitDealt,
		ProcMask:  core.ProcMaskSpellDamage,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})
	newStackingStatBonusCD(StackingStatBonusCD{
		Name:        "Ferver of the Frostborn",
		ID:          47727,
		Duration:    time.Second * 20,
		MaxStacks:   5,
		Bonus:       stats.Stats{stats.Armor: 1265},
		CD:          time.Minute * 2,
		Callback:    OnSpellHitTaken,
		Outcome:     core.OutcomeLanded,
		IsDefensive: true,
	})
	newStackingStatBonusCD(StackingStatBonusCD{
		Name:      "Fetish of Volatile Power",
		ID:        47879,
		Duration:  time.Second * 20,
		MaxStacks: 8,
		Bonus:     stats.Stats{stats.MeleeHaste: 57, stats.SpellHaste: 57},
		CD:        time.Minute * 2,
		Callback:  OnSpellHitDealt,
		ProcMask:  core.ProcMaskSpellDamage,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})
	newStackingStatBonusCD(StackingStatBonusCD{
		Name:      "Vengeance of the Forsaken",
		ID:        47881,
		Duration:  time.Second * 20,
		MaxStacks: 5,
		Bonus:     stats.Stats{stats.AttackPower: 215, stats.RangedAttackPower: 215},
		CD:        time.Minute * 2,
		Callback:  OnSpellHitDealt,
		ProcMask:  core.ProcMaskMelee,
		Outcome:   core.OutcomeLanded,
	})
	newStackingStatBonusCD(StackingStatBonusCD{
		Name:        "Eitrigg's Oath",
		ID:          47882,
		Duration:    time.Second * 20,
		MaxStacks:   5,
		Bonus:       stats.Stats{stats.Armor: 1265},
		CD:          time.Minute * 2,
		Callback:    OnSpellHitTaken,
		Outcome:     core.OutcomeLanded,
		IsDefensive: true,
	})
	newStackingStatBonusCD(StackingStatBonusCD{
		Name:      "Talisman of Volatile Power H",
		ID:        47946,
		Duration:  time.Second * 20,
		MaxStacks: 8,
		Bonus:     stats.Stats{stats.MeleeHaste: 64, stats.SpellHaste: 64},
		CD:        time.Minute * 2,
		Callback:  OnSpellHitDealt,
		ProcMask:  core.ProcMaskSpellDamage,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})
	newStackingStatBonusCD(StackingStatBonusCD{
		Name:      "Victor's Call H",
		ID:        47948,
		Duration:  time.Second * 20,
		MaxStacks: 5,
		Bonus:     stats.Stats{stats.AttackPower: 250, stats.RangedAttackPower: 250},
		CD:        time.Minute * 2,
		Callback:  OnSpellHitDealt,
		ProcMask:  core.ProcMaskMelee,
		Outcome:   core.OutcomeLanded,
	})
	newStackingStatBonusCD(StackingStatBonusCD{
		Name:        "Ferver of the Frostborn H",
		ID:          47949,
		Duration:    time.Second * 20,
		MaxStacks:   5,
		Bonus:       stats.Stats{stats.Armor: 1422},
		CD:          time.Minute * 2,
		Callback:    OnSpellHitTaken,
		Outcome:     core.OutcomeLanded,
		IsDefensive: true,
	})
	newStackingStatBonusCD(StackingStatBonusCD{
		Name:      "Fetish of Volatile Power H",
		ID:        48018,
		Duration:  time.Second * 20,
		MaxStacks: 8,
		Bonus:     stats.Stats{stats.MeleeHaste: 64, stats.SpellHaste: 64},
		CD:        time.Minute * 2,
		Callback:  OnSpellHitDealt,
		ProcMask:  core.ProcMaskSpellDamage,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})
	newStackingStatBonusCD(StackingStatBonusCD{
		Name:      "Vengeance of the Forsaken H",
		ID:        48020,
		Duration:  time.Second * 20,
		MaxStacks: 5,
		Bonus:     stats.Stats{stats.AttackPower: 250, stats.RangedAttackPower: 250},
		CD:        time.Minute * 2,
		Callback:  OnSpellHitDealt,
		ProcMask:  core.ProcMaskMelee,
		Outcome:   core.OutcomeLanded,
	})
	newStackingStatBonusCD(StackingStatBonusCD{
		Name:        "Eitrigg's Oath H",
		ID:          48021,
		Duration:    time.Second * 20,
		MaxStacks:   5,
		Bonus:       stats.Stats{stats.Armor: 1422},
		CD:          time.Minute * 2,
		Callback:    OnSpellHitTaken,
		Outcome:     core.OutcomeLanded,
		IsDefensive: true,
	})

	core.AddEffectsToTest = true
}
