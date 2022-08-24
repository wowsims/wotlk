package wotlk

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

type Callback uint16

func (c Callback) Matches(other Callback) bool {
	return (c & other) != 0
}

const (
	CallbackEmpty Callback = 0

	OnSpellHitDealt Callback = 1 << iota
	OnSpellHitTaken
	OnPeriodicDamageDealt
	OnCastComplete
)

type ProcHandler func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect)

type ProcTrigger struct {
	Name       string
	ActionID   core.ActionID
	Callback   Callback
	ProcMask   core.ProcMask
	Outcome    core.HitOutcome
	Harmful    bool
	ProcChance float64
	PPM        float64
	ICD        time.Duration
	Handler    ProcHandler
}

func applyProcTriggerCallback(unit *core.Unit, aura *core.Aura, config ProcTrigger) {
	var icd core.Cooldown
	if config.ICD != 0 {
		icd = core.Cooldown{
			Timer:    unit.NewTimer(),
			Duration: config.ICD,
		}
	}

	var ppmm core.PPMManager
	if config.PPM > 0 {
		ppmm = unit.AutoAttacks.NewPPMManager(config.PPM, config.ProcMask)
	}

	handler := config.Handler
	callback := func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
		if config.ProcMask != core.ProcMaskUnknown && !spellEffect.ProcMask.Matches(config.ProcMask) {
			return
		}
		if config.Outcome != core.OutcomeEmpty && !spellEffect.Outcome.Matches(config.Outcome) {
			return
		}
		if config.Harmful && spellEffect.Damage == 0 {
			return
		}
		if icd.Duration != 0 && !icd.IsReady(sim) {
			return
		}
		if config.ProcChance != 1 && sim.RandomFloat(config.Name) > config.ProcChance {
			return
		} else if config.PPM != 0 && !ppmm.Proc(sim, spellEffect.ProcMask, config.Name) {
			return
		}

		if icd.Duration != 0 {
			icd.Use(sim)
		}
		handler(sim, spell, spellEffect)
	}

	if config.Callback.Matches(OnSpellHitDealt) {
		aura.OnSpellHitDealt = callback
	}
	if config.Callback.Matches(OnSpellHitTaken) {
		aura.OnSpellHitTaken = callback
	}
	if config.Callback.Matches(OnPeriodicDamageDealt) {
		aura.OnPeriodicDamageDealt = callback
	}
	if config.ProcChance == 0 {
		config.ProcChance = 1
	}
	if config.Callback.Matches(OnCastComplete) {
		aura.OnCastComplete = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if icd.Duration != 0 && !icd.IsReady(sim) {
				return
			}
			if config.ProcChance != 1 && sim.RandomFloat(config.Name) > config.ProcChance {
				return
			}

			if icd.Duration != 0 {
				icd.Use(sim)
			}
			handler(sim, spell, nil)
		}
	}
}

func MakeProcTriggerAura(unit *core.Unit, config ProcTrigger) *core.Aura {
	aura := core.Aura{
		Label:    config.Name,
		ActionID: config.ActionID,
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	}

	applyProcTriggerCallback(unit, &aura, config)

	return unit.GetOrRegisterAura(aura)
}

func NewItemEffectWithHeroic(f func(isHeroic bool)) {
	f(true)
	f(false)
}
