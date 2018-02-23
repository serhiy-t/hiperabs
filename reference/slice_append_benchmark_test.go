package reference_test

import (
	"testing"

	"github.com/stykhanskyy/hiperabs/reference"
)

type point struct {
	x, y int
}

func checkSlice(slice []point) {
	if len(slice) < 1 {
		panic("Empty slice")
	}
	midIdx := len(slice) / 2
	endIdx := len(slice) - 1
	startPoint := point{0, 0}
	midPoint := point{midIdx, midIdx}
	endPoint := point{endIdx, endIdx}

	if slice[0] != startPoint {
		panic("Slice check failed")
	}

	if slice[midIdx] != midPoint {
		panic("Slice check failed")
	}

	if slice[endIdx] != endPoint {
		panic("Slice check failed")
	}
}

func benchmarkNativeSliceAppend(sliceSize int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		var slice []point
		for x := 0; x < sliceSize; x++ {
			slice = append(slice, point{x, x})
		}
		checkSlice(slice)
	}
}

func benchmarkInterfaceSliceAppend(sliceSize int, b *testing.B) {
	for i := 0; i < b.N; i++ {

		var slice []interface{}
		for x := 0; x < sliceSize; x++ {
			slice = append(slice, point{x, x})
		}
		// checkSlice(slice)
	}
}

type appendOnlyArrayList struct {
	slice          []point
	size, capacity int
	sliceRef       reference.SliceRef
}

func newAppendOnlyArrayList() appendOnlyArrayList {
	var slice []point
	return appendOnlyArrayList{slice: slice, sliceRef: reference.ToSlice(slice)}
}

func (list *appendOnlyArrayList) append(element *reference.Ref) {
	if list.size == list.capacity {
		newCapacity := list.capacity * 2
		if newCapacity == 0 {
			newCapacity++
		}
		newSlice := make([]point, newCapacity, newCapacity)
		copy(newSlice, list.slice)
		list.slice = newSlice
		list.sliceRef = reference.ToSlice(newSlice)
		list.capacity = newCapacity
	}

	reference.Copy(list.sliceRef.ElementRef(list.size), element)
	list.size++
}

func benchmarkRefSliceAppend(sliceSize int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		arrayList := newAppendOnlyArrayList()
		p := point{0, 0}
		pointRef := reference.ToPointer(&p)
		for x := 0; x < sliceSize; x++ {
			p = point{x, x}
			arrayList.append(pointRef.Ref())
		}
		checkSlice(arrayList.slice[:sliceSize])
	}
}

func BenchmarkNativeSlice100(b *testing.B)       { benchmarkNativeSliceAppend(100, b) }
func BenchmarkNativeSlice1000(b *testing.B)      { benchmarkNativeSliceAppend(1000, b) }
func BenchmarkNativeSlice10000(b *testing.B)     { benchmarkNativeSliceAppend(10000, b) }
func BenchmarkNativeSlice100000(b *testing.B)    { benchmarkNativeSliceAppend(100000, b) }
func BenchmarkNativeSlice1000000(b *testing.B)   { benchmarkNativeSliceAppend(1000000, b) }
func BenchmarkNativeSlice10000000(b *testing.B)  { benchmarkNativeSliceAppend(10000000, b) }
func BenchmarkNativeSlice100000000(b *testing.B) { benchmarkNativeSliceAppend(100000000, b) }

func BenchmarkRefSlice100(b *testing.B)       { benchmarkRefSliceAppend(100, b) }
func BenchmarkRefSlice1000(b *testing.B)      { benchmarkRefSliceAppend(1000, b) }
func BenchmarkRefSlice10000(b *testing.B)     { benchmarkRefSliceAppend(10000, b) }
func BenchmarkRefSlice100000(b *testing.B)    { benchmarkRefSliceAppend(100000, b) }
func BenchmarkRefSlice1000000(b *testing.B)   { benchmarkRefSliceAppend(1000000, b) }
func BenchmarkRefSlice10000000(b *testing.B)  { benchmarkRefSliceAppend(10000000, b) }
func BenchmarkRefSlice100000000(b *testing.B) { benchmarkRefSliceAppend(100000000, b) }

func BenchmarkInterfaceSlice100(b *testing.B)       { benchmarkInterfaceSliceAppend(100, b) }
func BenchmarkInterfaceSlice1000(b *testing.B)      { benchmarkInterfaceSliceAppend(1000, b) }
func BenchmarkInterfaceSlice10000(b *testing.B)     { benchmarkInterfaceSliceAppend(10000, b) }
func BenchmarkInterfaceSlice100000(b *testing.B)    { benchmarkInterfaceSliceAppend(100000, b) }
func BenchmarkInterfaceSlice1000000(b *testing.B)   { benchmarkInterfaceSliceAppend(1000000, b) }
func BenchmarkInterfaceSlice10000000(b *testing.B)  { benchmarkInterfaceSliceAppend(10000000, b) }
func BenchmarkInterfaceSlice100000000(b *testing.B) { benchmarkInterfaceSliceAppend(100000000, b) }
