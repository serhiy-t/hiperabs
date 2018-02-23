package benchmarks_test

import (
	"testing"

	"github.com/stykhanskyy/hiperabs/reference"

	"github.com/stykhanskyy/hiperabs/reference/benchmarks"
)

func benchmarkNativeMap(size int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		hashSet := make(map[benchmarks.Point]struct{})
		var s struct{}
		for x := 0; x < size; x++ {
			hashSet[benchmarks.Point{X: x, Y: x}] = s
		}
	}
}

func benchmarkInterfaceMap(size int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		hashSet := make(map[interface{}]bool)
		for x := 0; x < size; x++ {
			hashSet[benchmarks.Point{X: x, Y: x}] = true
		}
	}
}

func benchmarkHashSet(size int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		hashSet := benchmarks.NewHashSet(benchmarks.Point{})
		point := benchmarks.Point{}
		pointRef := reference.ToPointer(&point)
		for x := 0; x < size; x++ {
			point = benchmarks.Point{X: x, Y: x}
			hashSet.Add(pointRef)
		}
	}
}

func benchmarkInterfaceHashSet(size int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		pointPtr := &benchmarks.Point{}
		hashSet := benchmarks.NewHashSet(pointPtr)
		pointRef := reference.ToPointer(&pointPtr)
		for x := 0; x < size; x++ {
			pointPtr = &benchmarks.Point{X: x, Y: x}
			hashSet.Add(pointRef)
		}
	}
}

func BenchmarkNativeMap10000000(b *testing.B) { benchmarkNativeMap(10000000, b) }

func BenchmarkInterfaceMap10000000(b *testing.B) { benchmarkInterfaceMap(10000000, b) }

func BenchmarkHashSet10000000(b *testing.B) { benchmarkHashSet(10000000, b) }

func BenchmarkInterfaceHashSet10000000(b *testing.B) { benchmarkInterfaceHashSet(10000000, b) }
