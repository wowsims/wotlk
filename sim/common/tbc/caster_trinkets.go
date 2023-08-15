package tbc

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = false
	// Offensive trinkets. Keep these in order by item ID.
	core.NewSimpleStatOffensiveTrinketEffect(32483, stats.Stats{stats.SpellHaste: 175}, time.Second*20, time.Minute*2)  // Skull of Gul'dan
	core.NewSimpleStatOffensiveTrinketEffect(33829, stats.Stats{stats.SpellPower: 211}, time.Second*20, time.Minute*2)  // Hex Shrunken Head
	core.NewSimpleStatOffensiveTrinketEffect(34429, stats.Stats{stats.SpellPower: 320}, time.Second*15, time.Second*90) // Shifting Naaru Sliver

	// Even though these item effects are handled elsewhere, add them so they are
	// detected for automatic testing.
	for _, itemID := range core.AlchStoneItemIDs {
		core.NewItemEffect(itemID, func(core.Agent) {})
	}

	core.NewItemEffect(21625, func(agent core.Agent) { // Scarab Brooch
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 21625}

		shieldSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 26470},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskSpellHealing,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			Shield: core.ShieldConfig{
				Aura: core.Aura{
					Label:    "Scarab Brooch Shield",
					Duration: time.Second * 30,
				},
			},
		})

		activeAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Persistent Shield",
			ActionID: core.ActionID{SpellID: 26467},
			Callback: core.CallbackOnHealDealt,
			Duration: time.Second * 30,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				shieldSpell.Shield(result.Target).Apply(sim, result.Damage*0.15)
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				activeAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
		})
	})

	core.AddEffectsToTest = true
}
