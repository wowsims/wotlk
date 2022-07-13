package deathknight

import (
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

func (deathKnight *DeathKnight) registerBloodPresenceAura() {
	threatMult := 0.8
	//TODO: Include hps bonus
	damageBonusCoeff := 0.15

	deathKnight.BloodPresenceAura = deathKnight.GetOrRegisterAura(core.Aura{
		Label:    "Blood Presence",
		Tag:      "Presence",
		Priority: 1,
		ActionID: core.ActionID{SpellID: 50689},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier *= threatMult
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.0 + damageBonusCoeff
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier /= threatMult
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.0 + damageBonusCoeff
		},
	})
}

func (deathKnight *DeathKnight) registerFrostPresenceAura() {
	threatMult := 2.0735
	staminaBonusCoeff := 0.08
	armorBonusCoeff := 0.6

	deathKnight.FrostPresenceAura = deathKnight.GetOrRegisterAura(core.Aura{
		Label:    "Frost Presence",
		Tag:      "Presence",
		Priority: 1,
		ActionID: core.ActionID{SpellID: 48263},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier *= threatMult
			aura.Unit.AddStatDynamic(sim, stats.Armor, aura.Unit.GetStat(stats.Armor)*armorBonusCoeff)
			aura.Unit.AddStatDynamic(sim, stats.Stamina, aura.Unit.GetStat(stats.Stamina)*staminaBonusCoeff)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier /= threatMult
			aura.Unit.AddStatDynamic(sim, stats.Armor, -aura.Unit.GetStat(stats.Armor)*armorBonusCoeff)
			aura.Unit.AddStatDynamic(sim, stats.Stamina, -aura.Unit.GetStat(stats.Stamina)*staminaBonusCoeff)
		},
	})
}

func (deathKnight *DeathKnight) registerUnholyPresenceAura() {
	//gcdReductionTime := 500 * time.Millisecond
	attackSpeedBonus := 0.15

	deathKnight.UnholyPresenceAura = deathKnight.GetOrRegisterAura(core.Aura{
		Label:    "Unholy Presence",
		Tag:      "Presence",
		Priority: 1,
		ActionID: core.ActionID{SpellID: 48265},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.MeleeSpeedMultiplier *= 1.0 + attackSpeedBonus
			// TODO: Minus GCD time
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier /= 1.0 + attackSpeedBonus
			// TODO: Plus GCD time
		},
	})
}

func (deathKnight *DeathKnight) registerPresences() {
	deathKnight.registerBloodPresenceAura()
	deathKnight.registerFrostPresenceAura()
	deathKnight.registerUnholyPresenceAura()
}
