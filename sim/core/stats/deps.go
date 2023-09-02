package stats

import (
	"fmt"
)

// This stat list is arranged such that evaluating dependencies in this order
// is always safe, given the set of dependencies actually used in the game.
//
// Note that many stats are omitted from this list, because they are not used
// in any dependencies.
var safeDepsOrder = []Stat{
	Strength,
	Agility,
	Stamina,
	Intellect,
	Spirit,
	BonusArmor,
	Armor,
	AttackPower,
	RangedAttackPower,
	SpellPower,
	Health,
	Mana,
	MP5,
	SpellCrit,
	MeleeCrit,
	Defense,
	Block,
	BlockValue,
	Dodge,
	Parry,
}

func isSafeDep(s Stat) bool {
	for _, v := range safeDepsOrder {
		if s == v {
			return true
		}
	}
	return false
}
func isValidDep(src Stat, dst Stat) bool {
	if !isSafeDep(src) || !isSafeDep(dst) {
		return false
	}

	// Check that src occurs before dst in the list
	for _, v := range safeDepsOrder {
		if v == src {
			return true
		} else if v == dst {
			return false
		}
	}
	return false
}
func validateDep(src Stat, dst Stat) {
	if !isValidDep(src, dst) {
		panic("Invalid stat dependency: " + src.StatName() + " --> " + dst.StatName())
	}
}

type StatDependency struct {
	dynamic bool
	enabled bool
	src     Stat
	dst     Stat

	// Note that amount is treated differently depending on whether src and dst
	// stats are the same.
	amount float64
}

func (sd StatDependency) String() string {
	if sd.src == sd.dst {
		return fmt.Sprintf("%s *= %0.2f", sd.src.StatName(), sd.amount)
	} else {
		return fmt.Sprintf("%s += %s * %0.2f", sd.dst.StatName(), sd.src.StatName(), sd.amount)
	}
}

// Manages dependencies between stats.
//
// Some examples:
// Increases your AP by 30% of your Int
// Increases agility by X%
// Reduces armor by 50%
type StatDependencyManager struct {
	deps      []*StatDependency
	finalized bool
}

func NewStatDependencyManager() StatDependencyManager {
	return StatDependencyManager{}
}

func (sdm *StatDependencyManager) AddStatDependency(src Stat, dst Stat, amount float64) {
	validateDep(src, dst)
	if sdm.IsFinalized() {
		panic("StatDependencyManager already finalized!")
	}
	if src == dst {
		panic("For same-stat dependencies, use MultiplyStat instead!")
	}

	sdm.deps = append(sdm.deps, &StatDependency{
		dynamic: false,
		enabled: true,
		src:     src,
		dst:     dst,
		amount:  amount,
	})
}

func (sdm *StatDependencyManager) MultiplyStat(s Stat, amount float64) {
	validateDep(s, s)
	if sdm.IsFinalized() {
		panic("StatDependencyManager already finalized!")
	}

	sdm.deps = append(sdm.deps, &StatDependency{
		dynamic: false,
		enabled: true,
		src:     s,
		dst:     s,
		amount:  amount,
	})
}

func (sdm *StatDependencyManager) NewDynamicStatDependency(src Stat, dst Stat, amount float64) *StatDependency {
	validateDep(src, dst)
	if sdm.IsFinalized() {
		panic("StatDependencyManager already finalized!")
	}
	if src == dst {
		panic("For same-stat dependencies, use NewDynamicMultiplyStat instead!")
	}

	dep := &StatDependency{
		dynamic: true,
		enabled: false,
		src:     src,
		dst:     dst,
		amount:  amount,
	}
	sdm.deps = append(sdm.deps, dep)
	return dep
}

func (sdm *StatDependencyManager) NewDynamicMultiplyStat(s Stat, amount float64) *StatDependency {
	validateDep(s, s)
	if sdm.IsFinalized() {
		panic("StatDependencyManager already finalized!")
	}

	dep := &StatDependency{
		dynamic: true,
		enabled: false,
		src:     s,
		dst:     s,
		amount:  amount,
	}
	sdm.deps = append(sdm.deps, dep)
	return dep
}

func (sdm *StatDependencyManager) sortDeps() {
	deps := make([]*StatDependency, 0, len(sdm.deps))

	// By looping through the stats in order of safeDeps, we guarantee proper
	// sorting of dependencies.
	for i, srcStat := range safeDepsOrder {
		for _, dstStat := range safeDepsOrder[i:] {
			// Combine all static deps into 1 for performance.
			startAmount := 0.0
			if srcStat == dstStat {
				startAmount = 1
			}

			amount := startAmount
			for _, dep := range sdm.deps {
				if dep.src != srcStat || dep.dst != dstStat {
					continue
				}

				if dep.dynamic {
					// Dynamic deps need to remain separate, so
					// they can be turned on/off.
					deps = append(deps, dep)
				} else {
					if srcStat == dstStat {
						amount *= dep.amount
					} else {
						amount += dep.amount
					}
				}
			}

			if amount != startAmount {
				deps = append(deps, &StatDependency{
					enabled: true,
					src:     srcStat,
					dst:     dstStat,
					amount:  amount,
				})
			}
		}
	}

	sdm.deps = deps
}

func (sdm *StatDependencyManager) FinalizeStatDeps() {
	if sdm.IsFinalized() {
		panic("StatDependencyManager already finalized!")
	}
	sdm.sortDeps()
	sdm.finalized = true
}

func (sdm *StatDependencyManager) ResetStatDeps() {
	for _, dep := range sdm.deps {
		if dep.dynamic {
			dep.enabled = false
		}
	}
}

func (sdm *StatDependencyManager) IsFinalized() bool {
	return sdm.finalized
}

func (sdm *StatDependencyManager) ApplyStatDependencies(s Stats) Stats {
	for _, dep := range sdm.deps {
		if dep.enabled {
			if dep.src == dep.dst {
				s[dep.dst] *= dep.amount
			} else {
				s[dep.dst] += s[dep.src] * dep.amount
			}
		}
	}
	return s
}

func (sdm *StatDependencyManager) SortAndApplyStatDependencies(s Stats) Stats {
	sdm.sortDeps()
	return sdm.ApplyStatDependencies(s)
}

// Returns whether the state changed.
func (sdm *StatDependencyManager) EnableDynamicStatDep(dep *StatDependency) bool {
	if !dep.enabled {
		dep.enabled = true
		return true
	}
	return false
}

// Returns whether the state changed.
func (sdm *StatDependencyManager) DisableDynamicStatDep(dep *StatDependency) bool {
	if dep.enabled {
		dep.enabled = false
		return true
	}
	return false
}
