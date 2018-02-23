package benchmarks_test

import (
	"reflect"
	"testing"

	"github.com/stykhanskyy/hiperabs/reference/benchmarks"

	"github.com/stykhanskyy/hiperabs/reference"
)

func verifyReversedSlice(slice []benchmarks.Point) {
	size := reflect.ValueOf(slice).Len()

	beginPoint := benchmarks.Point{X: size - 1, Y: size - 1}
	endPoint := benchmarks.Point{X: 0, Y: 0}

	if slice[0] != beginPoint {
		panic("Slice check failed")
	}

	if slice[size-1] != endPoint {
		panic("Slice check failed")
	}
}

func createSliceToReverse(size int) []benchmarks.Point {
	slice := make([]benchmarks.Point, size, size)
	for x := 0; x < len(slice); x++ {
		slice[x] = benchmarks.Point{X: x, Y: x}
	}
	return slice
}

func benchmarkNativeReverse(size int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		slice := createSliceToReverse(size)
		middle := size / 2
		for x := 0; x < middle; x++ {
			slice[x], slice[size-1-x] = slice[size-1-x], slice[x]
		}
		verifyReversedSlice(slice)
	}
}

func benchmarkSwapperReverse(size int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		slice := createSliceToReverse(size)
		middle := size / 2
		swapper := reflect.Swapper(slice)
		for x := 0; x < middle; x++ {
			swapper(x, size-1-x)
		}
		verifyReversedSlice(slice)
	}
}

func benchmarkReferenceReverse(size int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		slice := createSliceToReverse(size)
		middle := size / 2
		sliceRef := reference.ToSlice(slice)
		tmpPoint := benchmarks.Point{}
		tmpRef := reference.ToPointer(&tmpPoint)
		for x := 0; x < middle; x++ {
			reference.Copy(tmpRef.Ref(), sliceRef.ElementRef(x))
			reference.Copy(sliceRef.ElementRef(x), sliceRef.ElementRef(size-x-1))
			reference.Copy(sliceRef.ElementRef(size-x-1), tmpRef.Ref())
		}
		verifyReversedSlice(slice)
	}
}

func benchmarkReflectReverse(size int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		slice := createSliceToReverse(size)
		middle := size / 2
		sliceValue := reflect.ValueOf(slice)
		tmpPoint := benchmarks.Point{}
		tmpPointPtr := &tmpPoint
		tmpValue := reflect.ValueOf(tmpPointPtr).Elem()
		for x := 0; x < middle; x++ {
			tmpValue.Set(sliceValue.Index(x))
			sliceValue.Index(x).Set(sliceValue.Index(size - x - 1))
			sliceValue.Index(size - x - 1).Set(tmpValue)
		}
		verifyReversedSlice(slice)
	}
}

func BenchmarkNativeReverse1000(b *testing.B)     { benchmarkNativeReverse(1000, b) }
func BenchmarkNativeReverse100000(b *testing.B)   { benchmarkNativeReverse(100000, b) }
func BenchmarkNativeReverse10000000(b *testing.B) { benchmarkNativeReverse(10000000, b) }

func BenchmarkSwapperReverse1000(b *testing.B)     { benchmarkSwapperReverse(1000, b) }
func BenchmarkSwapperReverse100000(b *testing.B)   { benchmarkSwapperReverse(100000, b) }
func BenchmarkSwapperReverse10000000(b *testing.B) { benchmarkSwapperReverse(10000000, b) }

func BenchmarkReferenceReverse1000(b *testing.B)     { benchmarkReferenceReverse(1000, b) }
func BenchmarkReferenceReverse100000(b *testing.B)   { benchmarkReferenceReverse(100000, b) }
func BenchmarkReferenceReverse10000000(b *testing.B) { benchmarkReferenceReverse(10000000, b) }

func BenchmarkReflectReverse1000(b *testing.B)     { benchmarkReflectReverse(1000, b) }
func BenchmarkReflectReverse100000(b *testing.B)   { benchmarkReflectReverse(100000, b) }
func BenchmarkReflectReverse10000000(b *testing.B) { benchmarkReflectReverse(10000000, b) }
