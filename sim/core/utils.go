package core

import (
	"hash/fnv"
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

func MinInt(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func MaxInt(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func MinInt32(a int32, b int32) int32 {
	if a < b {
		return a
	} else {
		return b
	}
}

func MaxInt32(a int32, b int32) int32 {
	if a > b {
		return a
	} else {
		return b
	}
}

func MinFloat(a float64, b float64) float64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func MaxFloat(a float64, b float64) float64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func MinDuration(a time.Duration, b time.Duration) time.Duration {
	if a < b {
		return a
	} else {
		return b
	}
}

func MaxDuration(a time.Duration, b time.Duration) time.Duration {
	if a > b {
		return a
	} else {
		return b
	}
}

func MinTristate(a proto.TristateEffect, b proto.TristateEffect) proto.TristateEffect {
	if a < b {
		return a
	} else {
		return b
	}
}

func MaxTristate(a proto.TristateEffect, b proto.TristateEffect) proto.TristateEffect {
	if a > b {
		return a
	} else {
		return b
	}
}

func DurationFromSeconds(numSeconds float64) time.Duration {
	return time.Duration(float64(time.Second) * numSeconds)
}

func DurationFromProto(durProto *proto.Duration) time.Duration {
	if durProto == nil {
		return 0
	} else {
		return time.Millisecond * time.Duration(durProto.Ms)
	}
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
	h.Write([]byte(s))
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
	dst := make(map[KO]VO)
	for ki, vi := range src {
		ko, vo := f(ki, vi)
		dst[ko] = vo
	}
	return dst
}

// Returns a new slice containing only the elements for which f returns true.
func FilterSlice[T any](src []T, f func(T) bool) []T {
	var dst []T
	for _, e := range src {
		if f(e) {
			dst = append(dst, e)
		}
	}
	return dst
}

// Returns a new map containing only the key/value pairs for which f returns true.
func FilterMap[K comparable, V any](src map[K]V, f func(K, V) bool) map[K]V {
	dst := make(map[K]V)
	for k, v := range src {
		if f(k, v) {
			dst[k] = v
		}
	}
	return dst
}

func calcMeanAndStdev(sample []float64) (float64, float64) {
	n := len(sample)
	sum := 0.0
	sumSq := 0.0
	for i := 0; i < n; i++ {
		sum += sample[i]
		sumSq += sample[i] * sample[i]
	}

	return calcMeanAndStdevFromSums(n, sum, sumSq)
}
func calcMeanAndStdevFromSums(n int, sum float64, sumSq float64) (float64, float64) {
	mean := sum / float64(n)
	stdev := math.Abs(math.Sqrt(sumSq/float64(n) - mean*mean))
	return mean, stdev
}
