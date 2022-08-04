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
	UnsetPresence
)

func (dk *Deathknight) PresenceMatches(other Presence) bool {
	return (dk.Presence & other) != 0
}

func (dk *Deathknight) ChangePresence(sim *core.Simulation, newPresence Presence) {
	if dk.PresenceMatches(newPresence) {
		return
	}

	dk.Presence = newPresence
	if dk.PresenceMatches(BloodPresence) {
		dk.BloodPresenceAura.Activate(sim)
		dk.FrostPresenceAura.Deactivate(sim)
		dk.UnholyPresenceAura.Deactivate(sim)
	} else if dk.PresenceMatches(FrostPresence) {
		dk.FrostPresenceAura.Activate(sim)
		dk.BloodPresenceAura.Deactivate(sim)
		dk.UnholyPresenceAura.Deactivate(sim)
	} else if dk.PresenceMatches(UnholyPresence) {
		dk.UnholyPresenceAura.Activate(sim)
		dk.BloodPresenceAura.Deactivate(sim)
		dk.FrostPresenceAura.Deactivate(sim)
	}
}

func (dk *Deathknight) CanBloodPresence(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 0.0, 1, 0, 0) && dk.BloodPresence.IsReady(sim)
}

func (dk *Deathknight) CastBloodPresence(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanBloodPresence(sim) {
		return dk.BloodPresence.Cast(sim, target)
	}
	return false
}

func (dk *Deathknight) CanFrostPresence(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 0.0, 0, 1, 0) && dk.FrostPresence.IsReady(sim)
}

func (dk *Deathknight) CastFrostPresence(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanFrostPresence(sim) {
		return dk.FrostPresence.Cast(sim, target)
	}
	return false
}

func (dk *Deathknight) CanUnholyPresence(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 0.0, 0, 0, 1) && dk.UnholyPresence.IsReady(sim)
}

func (dk *Deathknight) CastUnholyPresence(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanUnholyPresence(sim) {
		return dk.UnholyPresence.Cast(sim, target)
	}
	return false
}

func (dk *Deathknight) registerBloodPresenceAura(timer *core.Timer) {
	threatMult := 0.8
	threatMultSubversion := 1.0 - dk.subversionThreatBonus()
	//TODO: Include hps bonus
	damageBonusCoeff := 0.15
	staminaMult := 1.0 + 0.04*float64(dk.Talents.ImprovedFrostPresence)
	damageTakenMult := 1.0 - 0.01*float64(dk.Talents.ImprovedFrostPresence)

	baseCost := float64(core.NewRuneCost(0, 1, 0, 0, 0))
	dk.BloodPresence = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 50689},
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
			},
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second,
			},
		},
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dk.ChangePresence(sim, BloodPresence)
		},
	})

	dk.BloodPresenceAura = dk.GetOrRegisterAura(core.Aura{
		Label:    "Blood Presence",
		Tag:      "Presence",
		Priority: 1,
		ActionID: core.ActionID{SpellID: 50689},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier *= threatMult
			aura.Unit.PseudoStats.ThreatMultiplier *= threatMultSubversion
			aura.Unit.PseudoStats.DamageTakenMultiplier *= damageTakenMult

			dk.ModifyAdditiveDamageModifier(sim, damageBonusCoeff)
			aura.Unit.AddStatDependencyDynamic(sim, stats.Stamina, stats.Stamina, staminaMult)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier /= threatMult
			aura.Unit.PseudoStats.ThreatMultiplier /= threatMultSubversion
			aura.Unit.PseudoStats.DamageTakenMultiplier /= damageTakenMult

			dk.ModifyAdditiveDamageModifier(sim, -damageBonusCoeff)
			aura.Unit.AddStatDependencyDynamic(sim, stats.Stamina, stats.Stamina, 1.0/staminaMult)
		},
	})
}

func (dk *Deathknight) registerFrostPresenceAura(timer *core.Timer) {
	threatMult := 2.0735
	staminaMult := 1.08
	armorMult := 1.6

	baseCost := float64(core.NewRuneCost(0, 0, 1, 0, 0))
	dk.FrostPresence = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48263},
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
			},
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second,
			},
		},
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dk.ChangePresence(sim, FrostPresence)
		},
	})

	dk.FrostPresenceAura = dk.GetOrRegisterAura(core.Aura{
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

func (dk *Deathknight) registerUnholyPresenceAura(timer *core.Timer) {
	threatMultSubversion := 1.0 - dk.subversionThreatBonus()
	staminaMult := 1.0 + 0.04*float64(dk.Talents.ImprovedFrostPresence)

	baseCost := float64(core.NewRuneCost(0, 0, 0, 1, 0))
	dk.UnholyPresence = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48265},
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
			},
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second,
			},
		},
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dk.ChangePresence(sim, UnholyPresence)
		},
	})

	dk.UnholyPresenceAura = dk.GetOrRegisterAura(core.Aura{
		Label:    "Unholy Presence",
		Tag:      "Presence",
		Priority: 1,
		ActionID: core.ActionID{SpellID: 48265},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier *= threatMultSubversion
			aura.Unit.AddStatDependencyDynamic(sim, stats.Stamina, stats.Stamina, staminaMult)
			dk.MultiplyMeleeSpeed(sim, 1.15)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier /= threatMultSubversion
			aura.Unit.AddStatDependencyDynamic(sim, stats.Stamina, stats.Stamina, 1.0/staminaMult)
			dk.MultiplyMeleeSpeed(sim, 1/1.15)
		},
	})
}

func (dk *Deathknight) getModifiedGCD() time.Duration {
	if dk.UnholyPresenceAura.IsActive() {
		return time.Second
	} else {
		return core.GCDDefault
	}
}

func (dk *Deathknight) registerPresences() {
	presenceTimer := dk.NewTimer()
	dk.registerBloodPresenceAura(presenceTimer)
	dk.registerUnholyPresenceAura(presenceTimer)
	dk.registerFrostPresenceAura(presenceTimer)
}
