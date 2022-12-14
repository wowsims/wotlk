package core

// Rerpresents an absorption effect, e.g. Power Word: Shield.
type Shield struct {
	Spell *Spell

	// Embed Aura so we can use IsActive/Refresh/etc directly.
	*Aura
}

func (shield *Shield) Apply(sim *Simulation, shieldAmount float64) {
	caster := shield.Spell.Unit
	target := shield.Aura.Unit
	attackTable := caster.AttackTables[target.UnitIndex]

	shieldAmount *= shield.Spell.DamageMultiplier * caster.PseudoStats.HealingDealtMultiplier
	shieldAmount *= target.PseudoStats.HealingTakenMultiplier * attackTable.HealingDealtMultiplier

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

func NewShield(config Shield) *Shield {
	shield := &Shield{}
	*shield = config

	return shield
}

// Creates Shields for all allied units.
func NewAllyShieldArray(caster *Unit, config Shield, auraConfig Aura) []*Shield {
	shields := make([]*Shield, len(caster.Env.AllUnits))
	for _, target := range caster.Env.AllUnits {
		if !caster.IsOpponent(target) {
			config.Aura = target.GetOrRegisterAura(auraConfig)
			shields[target.UnitIndex] = NewShield(config)
		}
	}
	return shields
}
