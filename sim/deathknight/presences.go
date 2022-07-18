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
	return deathKnight.CastCostPossible(sim, 0.0, 1, 0, 0) && deathKnight.BloodPressence.IsReady(sim)
}

func (deathKnight *DeathKnight) CanFrostPresence(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 0.0, 0, 1, 0) && deathKnight.FrostPressence.IsReady(sim)
}

func (deathKnight *DeathKnight) CanUnholyPresence(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 0.0, 0, 0, 1) && deathKnight.UnholyPressence.IsReady(sim)
}

func (deathKnight *DeathKnight) registerBloodPresenceAura() {
	threatMult := 0.8
	//TODO: Include hps bonus
	damageBonusCoeff := 0.15

	deathKnight.BloodPressence = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 50689},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    deathKnight.NewTimer(),
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

			deathKnight.ModifyAdditiveDamageModifier(sim, damageBonusCoeff)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier /= threatMult

			deathKnight.ModifyAdditiveDamageModifier(sim, -damageBonusCoeff)
		},
	})
}

func (deathKnight *DeathKnight) registerFrostPresenceAura() {
	threatMult := 2.0735
	staminaBonusCoeff := 0.08
	armorBonusCoeff := 0.6

	deathKnight.FrostPressence = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 48263},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    deathKnight.NewTimer(),
				Duration: time.Millisecond * 1500,
			},
		},
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dkSpellCost := deathKnight.DetermineOptimalCost(sim, 0, 1, 0)
			deathKnight.Spend(sim, spell, dkSpellCost)
			deathKnight.ChangePresence(sim, FrostPresence)
		},
	})

	var armorGained = 0.0
	var staminaGained = 0.0
	deathKnight.FrostPresenceAura = deathKnight.GetOrRegisterAura(core.Aura{
		Label:    "Frost Presence",
		Tag:      "Presence",
		Priority: 1,
		ActionID: core.ActionID{SpellID: 48263},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier *= threatMult

			// TODO: Refactor to a dynamic implementation
			armorGained = aura.Unit.GetStat(stats.Armor) * armorBonusCoeff
			staminaGained = aura.Unit.GetStat(stats.Stamina) * staminaBonusCoeff

			aura.Unit.AddStatDynamic(sim, stats.Armor, armorGained)
			aura.Unit.AddStatDynamic(sim, stats.Stamina, staminaGained)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier /= threatMult
			aura.Unit.AddStatDynamic(sim, stats.Armor, -armorGained)
			aura.Unit.AddStatDynamic(sim, stats.Stamina, -staminaGained)
		},
	})
}

func (deathKnight *DeathKnight) registerUnholyPresenceAura() {
	deathKnight.UnholyPressence = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 48265},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    deathKnight.NewTimer(),
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
			aura.Unit.PseudoStats.MeleeSpeedMultiplier *= 1.15
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.MeleeSpeedMultiplier /= 1.15
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
	deathKnight.registerBloodPresenceAura()
	deathKnight.registerUnholyPresenceAura()
	deathKnight.registerFrostPresenceAura()
}
