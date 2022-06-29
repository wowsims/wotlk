package paladin

func (paladin *Paladin) ActivateRighteousFury() {
	paladin.PseudoStats.DamageTakenMultiplier *= 1 - 0.02*float64(paladin.Talents.ImprovedRighteousFury)
	paladin.PseudoStats.HolySpellThreatMultiplier *= 1.6 + 0.1*float64(paladin.Talents.ImprovedRighteousFury)
}
