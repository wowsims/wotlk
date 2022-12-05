package wotlk

func NewItemEffectWithHeroic(f func(isHeroic bool)) {
	f(true)
	f(false)
}
