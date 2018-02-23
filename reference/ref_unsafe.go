package reference

import (
	"reflect"
	"sync"
	"unsafe"
)

type typeDescriptor struct {
	size uintptr
	t    reflect.Type
	tPtr unsafe.Pointer
}

// Ref ...
type Ref struct {
	pointer    unsafe.Pointer
	object     *interface{}
	descriptor *typeDescriptor
}

// PointerRef ...
type PointerRef struct {
	ref *Ref
}

// SliceRef ...
type SliceRef struct {
	basePointer   unsafe.Pointer
	ref1          *Ref
	ref2          *Ref
	sliceLen      int
	typeSizeCache uintptr
}

// Ref ...
func (pRef *PointerRef) Ref() *Ref {
	return pRef.ref
}

var zero int

// ElementRef ...
func (sRef *SliceRef) ElementRef(x int) *Ref {
	if x < 0 || x >= sRef.sliceLen {
		// panic with inline
		zero = zero / zero
	}
	sRef.ref1, sRef.ref2 = sRef.ref2, sRef.ref1
	sRef.ref1.pointer = unsafe.Pointer(uintptr(sRef.basePointer) + uintptr(x)*sRef.typeSizeCache)
	return sRef.ref1
}

type bytes1Ptr [1]byte
type bytes2Ptr [2]byte
type bytes4Ptr [4]byte
type bytes8Ptr [8]byte
type bytes16Ptr [16]byte
type bytes32Ptr [32]byte

// Copy ...
func Copy(dst, src *Ref) {
	if src.descriptor != dst.descriptor {
		panic("attempt to copy references of different types")
	}
	typedmemmove(src.descriptor.tPtr, dst.pointer, src.pointer)
}

// map[reflect.Type]*typeDescriptor
var typeDescriptorTable sync.Map

func makeTypeDescriptor(t reflect.Type) *typeDescriptor {
	return &typeDescriptor{t.Size(), t, unsafe.Pointer(reflect.ValueOf(t).Pointer())}
}

func getTypeDescriptor(t reflect.Type) *typeDescriptor {
	descriptor, ok := typeDescriptorTable.Load(t)

	if ok {
		return descriptor.(*typeDescriptor)
	}

	newDescriptor := makeTypeDescriptor(t)
	descriptor, ok = typeDescriptorTable.LoadOrStore(t, newDescriptor)
	return descriptor.(*typeDescriptor)
}

// ToPointer ...
func ToPointer(valuePtr interface{}) PointerRef {
	valuePtrType := reflect.TypeOf(valuePtr)
	if valuePtrType.Kind() != reflect.Ptr {
		panic("valuePtr is not a Ptr")
	}

	descriptor := getTypeDescriptor(valuePtrType.Elem())

	pointer := unsafe.Pointer(reflect.ValueOf(valuePtr).Pointer())
	ref := Ref{pointer, &valuePtr, descriptor}
	return PointerRef{&ref}
}

// ToSlice ...
func ToSlice(slice interface{}) SliceRef {
	sliceType := reflect.TypeOf(slice)
	if sliceType.Kind() != reflect.Slice {
		panic("slice is not a Slice")
	}

	descriptor := getTypeDescriptor(sliceType.Elem())
	basePointer := unsafe.Pointer(reflect.ValueOf(slice).Pointer())
	ref1 := Ref{basePointer, &slice, descriptor}
	ref2 := Ref{basePointer, &slice, descriptor}

	return SliceRef{basePointer, &ref1, &ref2, reflect.ValueOf(slice).Len(), descriptor.size}
}
