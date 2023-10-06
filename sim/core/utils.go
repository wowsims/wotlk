package core

import (
	"hash/fnv"
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

func DurationFromSeconds(numSeconds float64) time.Duration {
	return time.Duration(float64(time.Second) * numSeconds)
}

func GetTristateValueInt32(effect proto.TristateEffect, regularValue int32, impValue int32) int32 {
	if effect == proto.TristateEffect_TristateEffectRegular {
		return regularValue
	} else if effect == proto.TristateEffect_TristateEffectImproved {
		return impValue
	} else {
		return 0
	}
}

func GetTristateValueFloat(effect proto.TristateEffect, regularValue float64, impValue float64) float64 {
	if effect == proto.TristateEffect_TristateEffectRegular {
		return regularValue
	} else if effect == proto.TristateEffect_TristateEffectImproved {
		return impValue
	} else {
		return 0
	}
}

func MakeTristateValue(hasRegular bool, hasImproved bool) proto.TristateEffect {
	if !hasRegular {
		return proto.TristateEffect_TristateEffectMissing
	} else if !hasImproved {
		return proto.TristateEffect_TristateEffectRegular
	} else {
		return proto.TristateEffect_TristateEffectImproved
	}
}

func hash(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

func Ternary[T any](condition bool, val1 T, val2 T) T {
	if condition {
		return val1
	} else {
		return val2
	}
}

func TernaryInt(condition bool, val1 int, val2 int) int {
	if condition {
		return val1
	} else {
		return val2
	}
}

func TernaryInt32(condition bool, val1 int32, val2 int32) int32 {
	if condition {
		return val1
	} else {
		return val2
	}
}

func TernaryFloat64(condition bool, val1 float64, val2 float64) float64 {
	if condition {
		return val1
	} else {
		return val2
	}
}

func TernaryDuration(condition bool, val1 time.Duration, val2 time.Duration) time.Duration {
	if condition {
		return val1
	} else {
		return val2
	}
}

func UnitLevelFloat64(unitLevel int32, maxLevelPlus0Val float64, maxLevelPlus1Val float64, maxLevelPlus2Val float64, maxLevelPlus3Val float64) float64 {
	if unitLevel == CharacterLevel {
		return maxLevelPlus0Val
	} else if unitLevel == CharacterLevel+1 {
		return maxLevelPlus1Val
	} else if unitLevel == CharacterLevel+2 {
		return maxLevelPlus2Val
	} else {
		return maxLevelPlus3Val
	}
}

func WithinToleranceFloat64(expectedValue float64, actualValue float64, tolerance float64) bool {
	return actualValue >= (expectedValue-tolerance) && actualValue <= (expectedValue+tolerance)
}

// Returns a new slice by applying f to each element in src.
func MapSlice[I any, O any](src []I, f func(I) O) []O {
	dst := make([]O, len(src))
	for i, e := range src {
		dst[i] = f(e)
	}
	return dst
}

// Returns a new map by applying f to each key/value pair in src.
func MapMap[KI comparable, VI any, KO comparable, VO any](src map[KI]VI, f func(KI, VI) (KO, VO)) map[KO]VO {
	dst := make(map[KO]VO, len(src))
	for ki, vi := range src {
		ko, vo := f(ki, vi)
		dst[ko] = vo
	}
	return dst
}

// Returns a new slice containing only the elements for which f returns true.
func FilterSlice[T any](src []T, f func(T) bool) []T {
	dst := make([]T, 0, len(src))
	for _, e := range src {
		if f(e) {
			dst = append(dst, e)
		}
	}
	return dst
}

// Returns a new map containing only the key/value pairs for which f returns true.
func FilterMap[K comparable, V any](src map[K]V, f func(K, V) bool) map[K]V {
	dst := make(map[K]V, len(src))
	for k, v := range src {
		if f(k, v) {
			dst[k] = v
		}
	}
	return dst
}

// Flattens a 2D slice into a 1D slice.
func Flatten[T any](src [][]T) []T {
	var n int
	for _, sublist := range src {
		n += len(sublist)
	}
	dst := make([]T, 0, n)
	for _, sublist := range src {
		dst = append(dst, sublist...)
	}
	return dst
}

type aggregator struct {
	n     int
	sum   float64
	sumSq float64
}

func (x *aggregator) add(v float64) {
	x.n++
	x.sum += v
	x.sumSq += v * v
}

func (x *aggregator) scale(f float64) {
	x.sum *= f
	x.sumSq *= f * f
}

func (x *aggregator) merge(y *aggregator) *aggregator {
	return &aggregator{n: x.n + y.n, sum: x.sum + y.sum, sumSq: x.sumSq + y.sumSq}
}

func (x *aggregator) meanAndStdDev() (float64, float64) {
	mean := x.sum / float64(x.n)
	stdDev := math.Sqrt(x.sumSq/float64(x.n) - mean*mean)
	return mean, stdDev
}
