package core

import (
	"time"
)

// An Exclusive effect is one which may not be active at the same time as other
// effects in the same category. For example, the armor reduction effects from
// Sunder Armor and Expose Armor are exclusive with each other.
//
// Within each ExclusiveCategory, the ExclusiveEffect with the highest Priority
// AND isEnabled == true is the one whose effect is applied.
type ExclusiveEffect struct {
	Aura     *Aura
	Priority float64

	OnGain   func(*ExclusiveEffect, *Simulation)
	OnExpire func(*ExclusiveEffect, *Simulation)

	Category  *ExclusiveCategory
	isEnabled bool
}

func (ee *ExclusiveEffect) IsActive() bool {
	return ee.Category.activeEffect == ee
}

type ExclusiveCategory struct {
	Name       string
	SingleAura bool // If true, only 1 aura in this category may be active at a time.
	effects    []*ExclusiveEffect

	activeEffect *ExclusiveEffect
}

func (ec *ExclusiveCategory) AnyActive() bool {
	return ec.activeEffect != nil
}

func (ec *ExclusiveCategory) GetActiveEffect() *ExclusiveEffect {
	return ec.activeEffect
}
func (ec *ExclusiveCategory) GetActiveAura() *Aura {
	if ec.activeEffect == nil {
		return nil
	} else {
		return ec.activeEffect.Aura
	}
}

func (ec *ExclusiveCategory) GetHighestPrioActiveEffect() *ExclusiveEffect {
	var effect *ExclusiveEffect
	for _, curEffect := range ec.effects {
		if curEffect.isEnabled && (effect == nil || curEffect.Priority > effect.Priority) {
			effect = curEffect
		}
	}
	return effect
}

func (ec *ExclusiveCategory) SetActive(sim *Simulation, newActiveEffect *ExclusiveEffect) {
	if newActiveEffect == ec.activeEffect {
		return
	}

	if ec.activeEffect != nil {
		ec.activeEffect.OnExpire(ec.activeEffect, sim)
	}
	ec.activeEffect = newActiveEffect
	if newActiveEffect != nil {
		newActiveEffect.OnGain(newActiveEffect, sim)
	}
}

type ExclusiveEffectManager struct {
	categories []*ExclusiveCategory
}

// Returns a category with the given name. Creates a new category if one doesn't already exist.
func (eem *ExclusiveEffectManager) GetExclusiveEffectCategory(categoryName string) *ExclusiveCategory {
	for i := 0; i < len(eem.categories); i++ {
		if eem.categories[i].Name == categoryName {
			return eem.categories[i]
		}
	}

	newCategory := &ExclusiveCategory{
		Name: categoryName,
	}
	eem.categories = append(eem.categories, newCategory)
	return newCategory
}

func (aura *Aura) NewExclusiveEffect(categoryName string, singleAura bool, config ExclusiveEffect) *ExclusiveEffect {
	if config.Aura != nil {
		panic("Don't specify aura in NewExclusiveEffect!")
	}

	if config.OnGain == nil {
		config.OnGain = func(*ExclusiveEffect, *Simulation) {}
	}
	if config.OnExpire == nil {
		config.OnExpire = func(*ExclusiveEffect, *Simulation) {}
	}

	eem := aura.Unit.ExclusiveEffectManager
	category := eem.GetExclusiveEffectCategory(categoryName)
	category.SingleAura = singleAura

	// If there is already an effect in this category with the same aura, use that instead.
	for _, effect := range category.effects {
		if effect.Aura == aura {
			return effect
		}
	}

	newEffect := &ExclusiveEffect{}
	*newEffect = config
	newEffect.Aura = aura
	newEffect.Category = category

	category.effects = append(category.effects, newEffect)
	aura.ExclusiveEffects = append(aura.ExclusiveEffects, newEffect)
	return newEffect
}

// Returns whether the effect is active.
func (ee *ExclusiveEffect) Activate(sim *Simulation) bool {
	if ee.isEnabled {
		return true
	}

	if ee.Category.SingleAura && ee.Category.activeEffect != nil && ee.Category.activeEffect != ee && (ee.Category.activeEffect.Priority > ee.Priority || (ee.Priority == ee.Category.activeEffect.Priority && ee.Category.activeEffect.Aura.RemainingDuration(sim) > ee.Aura.Duration)) {
		return false
	}

	ee.isEnabled = true

	if ee.Category.activeEffect == nil {
		ee.Category.SetActive(sim, ee)
	} else if ee.Priority >= ee.Category.activeEffect.Priority {
		if ee.Category.SingleAura && ee.Category.activeEffect != ee {
			ee.Category.activeEffect.Aura.Deactivate(sim)
		}
		ee.Category.SetActive(sim, ee)
	}

	return true
}
func (ee *ExclusiveEffect) Deactivate(sim *Simulation) {
	if !ee.isEnabled {
		return
	}
	ee.isEnabled = false

	if ee.Category.activeEffect == ee {
		ee.Category.SetActive(sim, ee.Category.GetHighestPrioActiveEffect())
	}
}

func (ee *ExclusiveEffect) SetPriority(sim *Simulation, newPrio float64) {
	if !ee.isEnabled {
		ee.Priority = newPrio
		return
	}

	curActiveEffect := ee.Category.activeEffect

	oldPrio := ee.Priority
	ee.Priority = newPrio
	newActiveEffect := ee.Category.GetHighestPrioActiveEffect()
	ee.Priority = oldPrio

	if curActiveEffect == ee && newActiveEffect == ee {
		ee.OnExpire(ee, sim)
		ee.Priority = newPrio
		ee.OnGain(ee, sim)
	} else if curActiveEffect != ee && newActiveEffect != ee {
		ee.Priority = newPrio
	} else if curActiveEffect == ee {
		// This effect is currently active but will be deactivated now.
		// Update prio after transition so it deactivates with the correct value.
		ee.Category.SetActive(sim, newActiveEffect)
		ee.Priority = newPrio
	} else {
		// This effect is currently not active but will be activated now.
		// Update prio before transition so it activates with the correct value.
		ee.Priority = newPrio
		ee.Category.SetActive(sim, newActiveEffect)
	}
}

// Returns if an aura should be refreshed, i.e. the aura is inactive/about to expire
// AND there are no other active effects of equal or greater strength.
func (aura *Aura) ShouldRefreshExclusiveEffects(sim *Simulation, refreshWindow time.Duration) bool {
	for _, ee := range aura.ExclusiveEffects {
		activeEffect := ee.Category.activeEffect
		if activeEffect == nil || ee.Priority > activeEffect.Priority {
			return true
		}

		if aura.MaxStacks > 0 {
			maxPriority := ee.Priority / float64(aura.GetStacks()) * float64(aura.MaxStacks)
			if maxPriority > activeEffect.Priority {
				return true
			}
		}

		if ee.Priority == activeEffect.Priority {
			anyWithLongRemainingDuration := false
			for _, effect := range ee.Category.effects {
				if effect.Priority == ee.Priority && effect.Aura.RemainingDuration(sim) > refreshWindow {
					anyWithLongRemainingDuration = true
					break
				}
			}
			if !anyWithLongRemainingDuration {
				return true
			}
		}
	}
	return false
}
func (spell *Spell) ShouldRefreshExclusiveEffects(sim *Simulation, target *Unit, refreshWindow time.Duration) bool {
	for _, auraArray := range spell.RelatedAuras {
		aura := auraArray.Get(target)
		if aura != nil && aura.ShouldRefreshExclusiveEffects(sim, refreshWindow) {
			return true
		}
	}
	return false
}

type ExclusiveCategoryArray []*ExclusiveCategory

func (categories ExclusiveCategoryArray) Get(target *Unit) *ExclusiveCategory {
	return categories[target.UnitIndex]
}

func (caster *Unit) NewEnemyExclusiveCategoryArray(makeExclusiveCategory func(*Unit) *ExclusiveCategory) ExclusiveCategoryArray {
	categories := make([]*ExclusiveCategory, len(caster.Env.AllUnits))
	for _, target := range caster.Env.AllUnits {
		if target.Type == EnemyUnit {
			categories[target.UnitIndex] = makeExclusiveCategory(target)
		}
	}
	return categories
}
func (caster *Unit) GetEnemyExclusiveCategories(category string) ExclusiveCategoryArray {
	return caster.NewEnemyExclusiveCategoryArray(func(target *Unit) *ExclusiveCategory {
		return target.GetExclusiveEffectCategory(category)
	})
}
