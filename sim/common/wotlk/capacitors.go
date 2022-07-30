package wotlk

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	//"github.com/wowsims/wotlk/sim/core/stats"
)

type CapacitorHandler func(*core.Simulation)

type CapacitorAura struct {
	Aura    core.Aura
	Handler CapacitorHandler
}

// Creates an aura which activates a handler function upon reaching a certain number of stacks.
func makeCapacitorAura(unit *core.Unit, config CapacitorAura) *core.Aura {
	handler := config.Handler
	config.Aura.OnStacksChange = func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
		if newStacks == aura.MaxStacks {
			handler(sim)
			aura.SetStacks(sim, 0)
		}
	}
	return unit.RegisterAura(config.Aura)
}

type CapacitorDamageEffect struct {
	Name      string
	ID        int32
	MaxStacks int32
	Trigger   ProcTrigger

	School     core.SpellSchool
	BaseDamage core.BaseDamageConfig
}

func newCapacitorDamageEffect(config CapacitorDamageEffect) {
	core.NewItemEffect(config.ID, func(agent core.Agent) {
		character := agent.GetCharacter()

		damageSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: config.ID},
			SpellSchool: config.School,
			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskEmpty,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				BaseDamage:       config.BaseDamage,
				OutcomeApplier:   character.OutcomeFuncMagicHitAndCrit(character.DefaultSpellCritMultiplier()),
			}),
		})

		capacitorAura := makeCapacitorAura(&character.Unit, CapacitorAura{
			Aura: core.Aura{
				Label:     config.Name,
				ActionID:  core.ActionID{ItemID: config.ID},
				Duration:  core.NeverExpires,
				MaxStacks: config.MaxStacks,
			},
			Handler: func(sim *core.Simulation) {
				damageSpell.Cast(sim, character.CurrentTarget)
			},
		})

		config.Trigger.Name = config.Name + " Trigger"
		config.Trigger.Handler = func(sim *core.Simulation, _ *core.Spell, _ *core.SpellEffect) {
			capacitorAura.Activate(sim)
			capacitorAura.AddStack(sim)
		}
		MakeProcTriggerAura(&character.Unit, config.Trigger)
	})
}

func init() {
	core.AddEffectsToTest = false

	newCapacitorDamageEffect(CapacitorDamageEffect{
		Name:      "Thunder Capacitor",
		ID:        38072,
		MaxStacks: 4,
		Trigger: ProcTrigger{
			Callback: OnSpellHitDealt,
			ProcMask: core.ProcMaskSpellDamage,
			Outcome:  core.OutcomeCrit,
			ICD:      time.Millisecond * 2500,
		},
		School:     core.SpellSchoolNature,
		BaseDamage: core.BaseDamageConfigRoll(1181, 1371),
	})
	newCapacitorDamageEffect(CapacitorDamageEffect{
		Name:      "Reign of the Unliving",
		ID:        47182,
		MaxStacks: 3,
		Trigger: ProcTrigger{
			Callback: OnSpellHitDealt,
			ProcMask: core.ProcMaskDirect,
			Outcome:  core.OutcomeCrit,
			ICD:      time.Millisecond * 2000,
		},
		School:     core.SpellSchoolFire,
		BaseDamage: core.BaseDamageConfigRoll(1741, 2023),
	})
	newCapacitorDamageEffect(CapacitorDamageEffect{
		Name:      "Reign of the Unliving H",
		ID:        47188,
		MaxStacks: 3,
		Trigger: ProcTrigger{
			Callback: OnSpellHitDealt,
			ProcMask: core.ProcMaskDirect,
			Outcome:  core.OutcomeCrit,
			ICD:      time.Millisecond * 2000,
		},
		School:     core.SpellSchoolFire,
		BaseDamage: core.BaseDamageConfigRoll(1959, 2275),
	})

	core.AddEffectsToTest = true

	newCapacitorDamageEffect(CapacitorDamageEffect{
		Name:      "Reign of the Dead",
		ID:        47316,
		MaxStacks: 3,
		Trigger: ProcTrigger{
			Callback: OnSpellHitDealt,
			ProcMask: core.ProcMaskDirect,
			Outcome:  core.OutcomeCrit,
			ICD:      time.Millisecond * 2000,
		},
		School:     core.SpellSchoolFire,
		BaseDamage: core.BaseDamageConfigRoll(1741, 2023),
	})
	newCapacitorDamageEffect(CapacitorDamageEffect{
		Name:      "Reign of the Dead H",
		ID:        47477,
		MaxStacks: 3,
		Trigger: ProcTrigger{
			Callback: OnSpellHitDealt,
			ProcMask: core.ProcMaskDirect,
			Outcome:  core.OutcomeCrit,
			ICD:      time.Millisecond * 2000,
		},
		School:     core.SpellSchoolFire,
		BaseDamage: core.BaseDamageConfigRoll(1959, 2275),
	})

	NewItemEffectWithHeroic(func(isHeroic bool) {
		name := "Tiny Abomination in a Jar"
		itemID := int32(50351)
		maxStacks := int32(8)
		if isHeroic {
			name += " H"
			itemID = 50706
			maxStacks = 7
		}

		core.NewItemEffect(itemID, func(agent core.Agent) {
			character := agent.GetCharacter()
			if !character.AutoAttacks.IsEnabled() {
				return
			}

			var mhSpell *core.Spell
			var ohSpell *core.Spell
			initSpells := func() {
				mhEffect := character.AutoAttacks.MHEffect
				mhEffect.DamageMultiplier *= 0.5
				mhSpell = character.GetOrRegisterSpell(core.SpellConfig{
					ActionID:     core.ActionID{ItemID: itemID}.WithTag(1),
					SpellSchool:  core.SpellSchoolPhysical,
					Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,
					ApplyEffects: core.ApplyEffectFuncDirectDamage(mhEffect),
				})

				if character.AutoAttacks.IsDualWielding {
					ohEffect := character.AutoAttacks.OHEffect
					ohEffect.DamageMultiplier *= 0.5
					ohSpell = character.GetOrRegisterSpell(core.SpellConfig{
						ActionID:     core.ActionID{ItemID: itemID}.WithTag(2),
						SpellSchool:  core.SpellSchoolPhysical,
						Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,
						ApplyEffects: core.ApplyEffectFuncDirectDamage(ohEffect),
					})
				}
			}

			capacitorAura := makeCapacitorAura(&character.Unit, CapacitorAura{
				Aura: core.Aura{
					Label:     name,
					ActionID:  core.ActionID{ItemID: itemID},
					Duration:  core.NeverExpires,
					MaxStacks: maxStacks,
					OnInit: func(aura *core.Aura, sim *core.Simulation) {
						initSpells()
					},
				},
				Handler: func(sim *core.Simulation) {
					if character.AutoAttacks.IsDualWielding {
						if sim.RandomFloat("Tiny Abom") < 0.5 {
							mhSpell.Cast(sim, character.CurrentTarget)
						} else {
							ohSpell.Cast(sim, character.CurrentTarget)
						}
					} else {
						mhSpell.Cast(sim, character.CurrentTarget)
					}
				},
			})

			MakeProcTriggerAura(&character.Unit, ProcTrigger{
				Name:       name + " Trigger",
				Callback:   OnSpellHitDealt,
				ProcMask:   core.ProcMaskMelee,
				Outcome:    core.OutcomeLanded,
				ProcChance: 0.45,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellEffect) {
					capacitorAura.Activate(sim)
					capacitorAura.AddStack(sim)
				},
			})
		})
	})
}
