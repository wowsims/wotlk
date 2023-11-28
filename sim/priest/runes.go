package priest

func (priest *Priest) ApplyRunes() {
	priest.registerVoidPlagueSpell()
	priest.RegisterPenanceSpell()
	priest.registerShadowWordDeathSpell()
}
