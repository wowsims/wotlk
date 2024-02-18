package wotlk

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var valkyrStats = stats.Stats{
	stats.Stamina: 1260,
}

type ValkyrPet struct {
	core.Pet
	smite         *core.Spell
	healthMetrics *core.ResourceMetrics
}

func newValkyr(character *core.Character) *ValkyrPet {
	return &ValkyrPet{
		Pet: core.NewPet("Valkyr", character, valkyrStats, func(ownerStats stats.Stats) stats.Stats {
			return stats.Stats{}
		}, false, true),
	}
}

func getSmiteConfig(valkyr *ValkyrPet, spellId int32, damageMin float64, damageMax float64) core.SpellConfig {
	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellDamage,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Millisecond * 1950, // about 16 instant-casts per 30s with some time left-over
				CastTime: 0,
			},
		},
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(damageMin, damageMax)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
			spell.DealDamage(sim, result)
			valkyr.GainHealth(sim, valkyr.MaxHealth()*0.25, valkyr.healthMetrics)
		},
		CritMultiplier: valkyr.DefaultSpellCritMultiplier(),
	}

}

func (valkyr *ValkyrPet) registerSmite(isHeroic bool) {
	spellId := int32(71841)
	if isHeroic {
		spellId = 71842
	}

	if valkyr.healthMetrics == nil {
		valkyr.healthMetrics = valkyr.NewHealthMetrics(core.ActionID{SpellID: spellId})
	}

	if isHeroic {
		smite := getSmiteConfig(valkyr, spellId, 1804, 2022)
		valkyr.smite = valkyr.GetOrRegisterSpell(smite)
	} else {
		smite := getSmiteConfig(valkyr, spellId, 1591, 1785)
		valkyr.smite = valkyr.GetOrRegisterSpell(smite)
	}
}

func (valkyr *ValkyrPet) Initialize() {}

func (valkyr *ValkyrPet) Reset(_ *core.Simulation) {}

func (valkyr *ValkyrPet) ExecuteCustomRotation(sim *core.Simulation) {
	valkyr.smite.Cast(sim, valkyr.CurrentTarget)
}

func (valkyr *ValkyrPet) GetPet() *core.Pet {
	return &valkyr.Pet
}

func MakeNibelungTriggerAura(agent core.Agent, isHeroic bool) {
	var auraSpellId, procSpellId int32

	if isHeroic {
		auraSpellId = 71844
		procSpellId = 71846
	} else {
		auraSpellId = 71843
		procSpellId = 71845
	}

	character := agent.GetCharacter()
	valkyrAura := character.RegisterAura(core.Aura{
		Label:    "Summon Val'kyr",
		ActionID: core.ActionID{SpellID: auraSpellId},
		Duration: time.Second * 30,
	})

	core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:            "Nibelung Proc",
		Callback:        core.CallbackOnCastComplete,
		ProcMask:        core.ProcMaskSpellOrProc,
		ProcMaskExclude: core.ProcMaskNotInSpellbook,
		Harmful:         true,
		ProcChance:      0.02,
		ICD:             time.Millisecond * 250,
		ActionID:        core.ActionID{SpellID: procSpellId},
		Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			for _, petAgent := range character.PetAgents {
				if valkyr, ok := petAgent.(*ValkyrPet); ok && !valkyr.IsEnabled() {
					valkyr.registerSmite(isHeroic)

					averageCasts := character.NibelungAverageCasts
					duration := min(time.Duration(averageCasts/16*30)*time.Second, time.Second*30)
					valkyrAura.Duration = max(duration, time.Millisecond*250)

					if averageCasts > 0 {
						valkyr.EnableWithTimeout(sim, petAgent, valkyrAura.Duration)
					}
					break
				}
			}

			valkyrAura.Activate(sim)
		},
	})
}

func ConstructValkyrPets(character *core.Character) {
	for i := 0; i < 10; i++ {
		valkyr := newValkyr(character)
		character.AddPet(valkyr)
	}
}
