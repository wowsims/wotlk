package priest

// Pre-cache values instead of checking gear during sim runs
func (priest *Priest) registerSetBonuses() {
	priest.T7TwoSetBonus = priest.HasSetBonus(ItemSetValorous, 2)
	priest.T7FourSetBonus = priest.HasSetBonus(ItemSetValorous, 4)
	priest.T8TwoSetBonus = priest.HasSetBonus(ItemSetConquerorSanct, 2)
	priest.T8FourSetBonus = priest.HasSetBonus(ItemSetConquerorSanct, 4)
	priest.T9TwoSetBonus = priest.HasSetBonus(ItemSetZabras, 2)
	priest.T9FourSetBonus = priest.HasSetBonus(ItemSetZabras, 4)
	priest.T10TwoSetBonus = priest.HasSetBonus(ItemSetCrimsonAcolyte, 2)
	priest.T10FourSetBonus = priest.HasSetBonus(ItemSetCrimsonAcolyte, 4)
}
