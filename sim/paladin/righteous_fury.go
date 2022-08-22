package paladin

func (paladin *Paladin) ActivateRighteousFury() {
	paladin.PseudoStats.DamageTakenMultiplier *= 1 - 0.02*float64(paladin.Talents.ImprovedRighteousFury)

	paladin.PseudoStats.HolySpellThreatMultiplier *= 1.8

	// Extra threat provided to all tanks on certain buff activation, for Paladins that is RF.
	paladin.PseudoStats.ThreatMultiplier *= 1.43
}
