package core

import "strconv"

type ShieldConfig struct {
	SelfOnly bool // Set to true to only create the self-shield.

	Spell *Spell

	Aura
}

// Rerpresents an absorption effect, e.g. Power Word: Shield.
type Shield struct {
	Spell *Spell

	// Embed Aura so we can use IsActive/Refresh/etc directly.
	*Aura
}

func (shield *Shield) Apply(sim *Simulation, shieldAmount float64) {
	caster := shield.Spell.Unit
	target := shield.Aura.Unit
	//attackTable := caster.AttackTables[target.UnitIndex]

	// Shields are not affected by healing pseudostats the same way heals are.
	// So we only apply the spell-specific multiplier.
	shieldAmount *= shield.Spell.DamageMultiplier

	shield.Aura.Deactivate(sim)
	shield.Aura.Activate(sim)

	threat := 0.0 // TODO
	shield.Spell.SpellMetrics[target.UnitIndex].TotalThreat += threat
	shield.Spell.SpellMetrics[target.UnitIndex].TotalShielding += shieldAmount
	shield.Spell.SpellMetrics[target.UnitIndex].Hits++

	if sim.Log != nil {
		caster.Log(sim, "%s %s Hit for %0.3f shielding. (Threat: %0.3f)", target.LogLabel(), shield.Spell.ActionID, shieldAmount, threat)
	}
}

func newShield(config Shield) *Shield {
	shield := &Shield{}
	*shield = config

	return shield
}

type ShieldArray []*Shield

func (shields ShieldArray) Get(target *Unit) *Shield {
	return shields[target.UnitIndex]
}

func (spell *Spell) createShields(config ShieldConfig) {
	if config.Aura.Label == "" {
		return
	}

	if config.Spell == nil {
		config.Spell = spell
	}
	shield := Shield{
		Spell: config.Spell,
	}

	auraConfig := config.Aura
	if auraConfig.ActionID.IsEmptyAction() {
		auraConfig.ActionID = shield.Spell.ActionID
	}

	caster := shield.Spell.Unit
	if config.SelfOnly {
		shield.Aura = caster.GetOrRegisterAura(auraConfig)
		spell.selfShield = newShield(shield)
	} else {
		auraConfig.Label += "-" + strconv.Itoa(int(caster.UnitIndex))
		if spell.shields == nil {
			spell.shields = make([]*Shield, len(caster.Env.AllUnits))
		}
		for _, target := range caster.Env.AllUnits {
			if !caster.IsOpponent(target) {
				shield.Aura = target.GetOrRegisterAura(auraConfig)
				spell.shields[target.UnitIndex] = newShield(shield)
			}
		}
	}
}
