package mage

const (
	MageRuneChestBurnout         = 415460
	MageRuneChestEnlightenment   = 415729
	MageRuneChestFingersOfFrost  = 401741
	MageRuneChestRegeneration    = 401743
	MageRuneHandsArcaneBlast     = 401729
	MageRuneHandsIceLance        = 401732
	MageRuneHandsLivingBomb      = 401731
	MageRuneHandsRewindTime      = 401734
	MageRuneLegsArcaneSurge      = 425168
	MageRuneLegsIceVeins         = 425169
	MageRuneLegsLivingFlame      = 401744
	MageRuneLegsMassRegeneration = 415467
)

func (mage *Mage) ApplyRunes() {
	mage.registerArcaneBlastSpell()
	mage.registerEnlightenment()
	mage.registerFingersOfFrost()
	mage.registerIceLanceSpell()
	mage.registerIcyVeins()
	mage.registerLivingBombSpell()
	mage.registerLivingFlameSpell()
	mage.registerRuneBurnout()
}
