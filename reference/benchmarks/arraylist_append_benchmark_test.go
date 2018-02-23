package benchmarks_test

import (
	"reflect"
	"testing"

	"github.com/stykhanskyy/hiperabs/reference/benchmarks"

	"github.com/stykhanskyy/hiperabs/reference"
)

func verifySlice(slice interface{}) {
	size := reflect.ValueOf(slice).Len()

	if size < 1 {
		panic("Empty slice")
	}
	midIdx := size / 2
	endIdx := size - 1
	startPoint := benchmarks.Point{X: 0, Y: 0}
	midPoint := benchmarks.Point{X: midIdx, Y: midIdx}
	endPoint := benchmarks.Point{X: endIdx, Y: endIdx}

	if reflect.ValueOf(slice).Index(0).Interface() != startPoint {
		panic("Slice check failed")
	}

	if reflect.ValueOf(slice).Index(midIdx).Interface() != midPoint {
		panic("Slice check failed")
	}

	if reflect.ValueOf(slice).Index(endIdx).Interface() != endPoint {
		panic("Slice check failed")
	}
}

func benchmarkNativeSliceAppend(sliceSize int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		var slice []benchmarks.Point
		for x := 0; x < sliceSize; x++ {
			slice = append(slice, benchmarks.Point{X: x, Y: x})
		}
		verifySlice(slice)
	}
}

func benchmarkInterfaceSliceAppend(sliceSize int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		var slice []interface{}
		for x := 0; x < sliceSize; x++ {
			slice = append(slice, benchmarks.Point{X: x, Y: x})
		}
		verifySlice(slice)
	}
}

func benchmarkReferenceSliceAppend(sliceSize int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		arrayList := benchmarks.NewArrayList(benchmarks.Point{})
		p := benchmarks.Point{X: 0, Y: 0}
		pointRef := reference.ToPointer(&p)
		for x := 0; x < sliceSize; x++ {
			p = benchmarks.Point{X: x, Y: x}
			arrayList.AppendRef(pointRef.Ref())
		}
		verifySlice(arrayList.GetSlice())
	}
}

func benchmarkReflectSliceAppend(sliceSize int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		arrayList := benchmarks.NewArrayList(benchmarks.Point{})
		p := benchmarks.Point{X: 0, Y: 0}
		pPtr := &p
		pointVal := reflect.ValueOf(pPtr).Elem()
		for x := 0; x < sliceSize; x++ {
			p = benchmarks.Point{X: x, Y: x}
			arrayList.AppendValue(&pointVal)
		}
		verifySlice(arrayList.GetSlice())
	}
}

func BenchmarkNativeSlice100(b *testing.B)       { benchmarkNativeSliceAppend(100, b) }
func BenchmarkNativeSlice1000(b *testing.B)      { benchmarkNativeSliceAppend(1000, b) }
func BenchmarkNativeSlice10000(b *testing.B)     { benchmarkNativeSliceAppend(10000, b) }
func BenchmarkNativeSlice100000(b *testing.B)    { benchmarkNativeSliceAppend(100000, b) }
func BenchmarkNativeSlice1000000(b *testing.B)   { benchmarkNativeSliceAppend(1000000, b) }
func BenchmarkNativeSlice10000000(b *testing.B)  { benchmarkNativeSliceAppend(10000000, b) }
func BenchmarkNativeSlice100000000(b *testing.B) { benchmarkNativeSliceAppend(100000000, b) }

func BenchmarkReferenceSlice100(b *testing.B)       { benchmarkReferenceSliceAppend(100, b) }
func BenchmarkReferenceSlice1000(b *testing.B)      { benchmarkReferenceSliceAppend(1000, b) }
func BenchmarkReferenceSlice10000(b *testing.B)     { benchmarkReferenceSliceAppend(10000, b) }
func BenchmarkReferenceSlice100000(b *testing.B)    { benchmarkReferenceSliceAppend(100000, b) }
func BenchmarkReferenceSlice1000000(b *testing.B)   { benchmarkReferenceSliceAppend(1000000, b) }
func BenchmarkReferenceSlice10000000(b *testing.B)  { benchmarkReferenceSliceAppend(10000000, b) }
func BenchmarkReferenceSlice100000000(b *testing.B) { benchmarkReferenceSliceAppend(100000000, b) }

func BenchmarkReflectSlice100(b *testing.B)       { benchmarkReflectSliceAppend(100, b) }
func BenchmarkReflectSlice1000(b *testing.B)      { benchmarkReflectSliceAppend(1000, b) }
func BenchmarkReflectSlice10000(b *testing.B)     { benchmarkReflectSliceAppend(10000, b) }
func BenchmarkReflectSlice100000(b *testing.B)    { benchmarkReflectSliceAppend(100000, b) }
func BenchmarkReflectSlice1000000(b *testing.B)   { benchmarkReflectSliceAppend(1000000, b) }
func BenchmarkReflectSlice10000000(b *testing.B)  { benchmarkReflectSliceAppend(10000000, b) }
func BenchmarkReflectSlice100000000(b *testing.B) { benchmarkReflectSliceAppend(100000000, b) }

func BenchmarkInterfaceSlice100(b *testing.B)       { benchmarkInterfaceSliceAppend(100, b) }
func BenchmarkInterfaceSlice1000(b *testing.B)      { benchmarkInterfaceSliceAppend(1000, b) }
func BenchmarkInterfaceSlice10000(b *testing.B)     { benchmarkInterfaceSliceAppend(10000, b) }
func BenchmarkInterfaceSlice100000(b *testing.B)    { benchmarkInterfaceSliceAppend(100000, b) }
func BenchmarkInterfaceSlice1000000(b *testing.B)   { benchmarkInterfaceSliceAppend(1000000, b) }
func BenchmarkInterfaceSlice10000000(b *testing.B)  { benchmarkInterfaceSliceAppend(10000000, b) }
func BenchmarkInterfaceSlice100000000(b *testing.B) { benchmarkInterfaceSliceAppend(100000000, b) }
