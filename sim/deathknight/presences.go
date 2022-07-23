package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type Presence uint8

const (
	BloodPresence Presence = 1 << iota
	FrostPresence
	UnholyPresence
)

func (deathKnight *DeathKnight) PresenceMatches(other Presence) bool {
	return (deathKnight.Presence & other) != 0
}

func (deathKnight *DeathKnight) ChangePresence(sim *core.Simulation, newPresence Presence) {
	if deathKnight.PresenceMatches(newPresence) {
		return
	}

	deathKnight.Presence = newPresence
	if deathKnight.PresenceMatches(BloodPresence) {
		deathKnight.BloodPresenceAura.Activate(sim)
		deathKnight.FrostPresenceAura.Deactivate(sim)
		deathKnight.UnholyPresenceAura.Deactivate(sim)
	} else if deathKnight.PresenceMatches(FrostPresence) {
		deathKnight.FrostPresenceAura.Activate(sim)
		deathKnight.BloodPresenceAura.Deactivate(sim)
		deathKnight.UnholyPresenceAura.Deactivate(sim)
	} else if deathKnight.PresenceMatches(UnholyPresence) {
		deathKnight.UnholyPresenceAura.Activate(sim)
		deathKnight.BloodPresenceAura.Deactivate(sim)
		deathKnight.FrostPresenceAura.Deactivate(sim)
	}
}

func (deathKnight *DeathKnight) CanBloodPresence(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 0.0, 1, 0, 0) && deathKnight.BloodPresence.IsReady(sim)
}

func (deathKnight *DeathKnight) CastBloodPresence(sim *core.Simulation, target *core.Unit) bool {
	if deathKnight.CanBloodPresence(sim) {
		deathKnight.BloodPresence.Cast(sim, target)
		return true
	}
	return false
}

func (deathKnight *DeathKnight) CanFrostPresence(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 0.0, 0, 1, 0) && deathKnight.FrostPresence.IsReady(sim)
}

func (deathKnight *DeathKnight) CastFrostPresence(sim *core.Simulation, target *core.Unit) bool {
	if deathKnight.CanFrostPresence(sim) {
		deathKnight.FrostPresence.Cast(sim, target)
		return true
	}
	return false
}

func (deathKnight *DeathKnight) CanUnholyPresence(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 0.0, 0, 0, 1) && deathKnight.UnholyPresence.IsReady(sim)
}

func (deathKnight *DeathKnight) CastUnholyPresence(sim *core.Simulation, target *core.Unit) bool {
	if deathKnight.CanUnholyPresence(sim) {
		deathKnight.UnholyPresence.Cast(sim, target)
		return true
	}
	return false
}

func (deathKnight *DeathKnight) registerBloodPresenceAura(timer *core.Timer) {
	threatMult := 0.8
	threatMultSubversion := 0.75
	//TODO: Include hps bonus
	damageBonusCoeff := 0.15
	staminaMult := 1.0 + 0.04*float64(deathKnight.Talents.ImprovedFrostPresence)
	damageTakenMult := 1.0 - 0.01*float64(deathKnight.Talents.ImprovedFrostPresence)

	deathKnight.BloodPresence = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 50689},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Millisecond * 1500,
			},
		},
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dkSpellCost := deathKnight.DetermineOptimalCost(sim, 1, 0, 0)
			deathKnight.Spend(sim, spell, dkSpellCost)
			deathKnight.ChangePresence(sim, BloodPresence)
		},
	})

	deathKnight.BloodPresenceAura = deathKnight.GetOrRegisterAura(core.Aura{
		Label:    "Blood Presence",
		Tag:      "Presence",
		Priority: 1,
		ActionID: core.ActionID{SpellID: 50689},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier *= threatMult
			aura.Unit.PseudoStats.ThreatMultiplier *= threatMultSubversion
			aura.Unit.PseudoStats.DamageTakenMultiplier *= damageTakenMult

			deathKnight.ModifyAdditiveDamageModifier(sim, damageBonusCoeff)
			aura.Unit.AddStatDependencyDynamic(sim, stats.Stamina, stats.Stamina, staminaMult)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier /= threatMult
			aura.Unit.PseudoStats.ThreatMultiplier /= threatMultSubversion
			aura.Unit.PseudoStats.DamageTakenMultiplier /= damageTakenMult

			deathKnight.ModifyAdditiveDamageModifier(sim, -damageBonusCoeff)
			aura.Unit.AddStatDependencyDynamic(sim, stats.Stamina, stats.Stamina, 1.0/staminaMult)
		},
	})
}

func (deathKnight *DeathKnight) registerFrostPresenceAura(timer *core.Timer) {
	threatMult := 2.0735
	staminaMult := 1.08
	armorMult := 1.6

	deathKnight.FrostPresence = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 48263},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Millisecond * 1500,
			},
		},
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dkSpellCost := deathKnight.DetermineOptimalCost(sim, 0, 1, 0)
			deathKnight.Spend(sim, spell, dkSpellCost)
			deathKnight.ChangePresence(sim, FrostPresence)
		},
	})

	deathKnight.FrostPresenceAura = deathKnight.GetOrRegisterAura(core.Aura{
		Label:    "Frost Presence",
		Tag:      "Presence",
		Priority: 1,
		ActionID: core.ActionID{SpellID: 48263},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier *= threatMult

			aura.Unit.AddStatDependencyDynamic(sim, stats.Armor, stats.Armor, armorMult)
			aura.Unit.AddStatDependencyDynamic(sim, stats.Stamina, stats.Stamina, staminaMult)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier /= threatMult

			aura.Unit.AddStatDependencyDynamic(sim, stats.Armor, stats.Armor, 1.0/armorMult)
			aura.Unit.AddStatDependencyDynamic(sim, stats.Stamina, stats.Stamina, 1.0/staminaMult)
		},
	})
}

func (deathKnight *DeathKnight) registerUnholyPresenceAura(timer *core.Timer) {
	threatMultiplierSubversion := 0.75
	staminaMult := 1.0 + 0.04*float64(deathKnight.Talents.ImprovedFrostPresence)

	deathKnight.UnholyPresence = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 48265},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Millisecond * 1500,
			},
		},
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dkSpellCost := deathKnight.DetermineOptimalCost(sim, 0, 0, 1)
			deathKnight.Spend(sim, spell, dkSpellCost)
			deathKnight.ChangePresence(sim, UnholyPresence)
		},
	})

	deathKnight.UnholyPresenceAura = deathKnight.GetOrRegisterAura(core.Aura{
		Label:    "Unholy Presence",
		Tag:      "Presence",
		Priority: 1,
		ActionID: core.ActionID{SpellID: 48265},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier *= threatMultiplierSubversion
			aura.Unit.PseudoStats.MeleeSpeedMultiplier *= 1.15
			aura.Unit.AddStatDependencyDynamic(sim, stats.Stamina, stats.Stamina, staminaMult)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier /= threatMultiplierSubversion
			aura.Unit.PseudoStats.MeleeSpeedMultiplier /= 1.15
			aura.Unit.AddStatDependencyDynamic(sim, stats.Stamina, stats.Stamina, 1.0/staminaMult)
		},
	})
}

func (deathKnight *DeathKnight) getModifiedGCD() time.Duration {
	if deathKnight.UnholyPresenceAura.IsActive() {
		return time.Second
	} else {
		return core.GCDDefault
	}
}

func (deathKnight *DeathKnight) registerPresences() {
	presenceTimer := deathKnight.NewTimer()
	deathKnight.registerBloodPresenceAura(presenceTimer)
	deathKnight.registerUnholyPresenceAura(presenceTimer)
	deathKnight.registerFrostPresenceAura(presenceTimer)
}
