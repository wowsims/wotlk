package core

import (
	"math"
	"testing"
)

func TestSplitMix64(t *testing.T) {
	x := SplitMix64{1234567}
	min, max := 1.0, 0.0
	distribution := make([]int, 500)
	n := 100_000_000
	for i := 0; i < n; i++ {
		f := x.NextFloat64()
		if f < min {
			min = f
		}
		if f > max {
			max = f
		}
		distribution[int(math.Trunc(f*500))]++
	}
	if min < 0 || max >= 1 {
		t.Fatalf("min = %f < 0 || max = %f >= 1", min, max)
	}
	e := float64(n) / 500
	var chiSquare float64
	for _, v := range distribution {
		chiSquare += (float64(v) - e) * (float64(v) - e) / e
	}
	if chiSquare > 540.93 {
		t.Fatalf("fails chi-square (k = 500) at a = 0.1 (%.1f >= 540.93)", chiSquare)
	}
	t.Logf("chiSquare = %.1f", chiSquare)
}

var result float64

func BenchmarkRnds(b *testing.B) {
	r := NewGoRand(444)
	b.Run("GoRand", func(b *testing.B) {
		b.ReportAllocs()
		var sum float64
		for i := 0; i < b.N; i++ {
			sum += r.NextFloat64()
			sum += r.NextFloat64()
			sum += r.NextFloat64()
			sum += r.NextFloat64()
			sum += r.NextFloat64()
		}
		result += sum
	})

	b.Run("GoRandSetup", func(b *testing.B) {
		b.ReportAllocs()
		var sum float64
		for i := 0; i < b.N; i++ {
			r := NewGoRand(444 + uint64(i))
			sum += r.NextFloat64()
			r = NewGoRand(555 + uint64(i))
			sum += r.NextFloat64()
			r = NewGoRand(666 + uint64(i))
			sum += r.NextFloat64()
			r = NewGoRand(777 + uint64(i))
			sum += r.NextFloat64()
			r = NewGoRand(888 + uint64(i))
			sum += r.NextFloat64()
		}
		result += sum
	})

	sm := NewSplitMix(444)
	b.Run("SplitMix64", func(b *testing.B) {
		b.ReportAllocs()
		var sum float64
		for i := 0; i < b.N; i++ {
			sum += sm.NextFloat64()
			sum += sm.NextFloat64()
			sum += sm.NextFloat64()
			sum += sm.NextFloat64()
			sum += sm.NextFloat64()
		}
		result += sum
	})

	b.Run("SplitMix64Setup", func(b *testing.B) {
		b.ReportAllocs()
		var sum float64
		for i := 0; i < b.N; i++ {
			sm := NewSplitMix(444 + uint64(i))
			sum += sm.NextFloat64()
			sm = NewSplitMix(555 + uint64(i))
			sum += sm.NextFloat64()
			sm = NewSplitMix(666 + uint64(i))
			sum += sm.NextFloat64()
			sm = NewSplitMix(777 + uint64(i))
			sum += sm.NextFloat64()
			sm = NewSplitMix(888 + uint64(i))
			sum += sm.NextFloat64()
		}
		result += sum
	})

	b.Run("Addition", func(b *testing.B) {
		var sum float64
		for i := 0; i < b.N; i++ {
			sum += 3.14159
			sum += 3.14159
			sum += 3.14159
			sum += 3.14159
			sum += 3.14159
		}
		result += sum
	})
}
