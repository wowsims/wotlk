package warlock

// Pre-cache values instead of checking gear during sim runs
func (warlock *Warlock) registerSetBonuses() {
	warlock.T7TwoSetBonus = warlock.HasSetBonus(ItemSetPlagueheartGarb, 2)
	warlock.T7FourSetBonus = warlock.HasSetBonus(ItemSetPlagueheartGarb, 4)
	warlock.T8TwoSetBonus = warlock.HasSetBonus(ItemSetDeathbringerGarb, 2)
	warlock.T8FourSetBonus = warlock.HasSetBonus(ItemSetDeathbringerGarb, 4)
	warlock.T9TwoSetBonus = warlock.HasSetBonus(ItemSetGuldansRegalia, 2)
	warlock.T9FourSetBonus = warlock.HasSetBonus(ItemSetGuldansRegalia, 4)
	warlock.T10TwoSetBonus = warlock.HasSetBonus(ItemSetDarkCovensRegalia, 2)
	warlock.T10FourSetBonus = warlock.HasSetBonus(ItemSetDarkCovensRegalia, 4)
}
