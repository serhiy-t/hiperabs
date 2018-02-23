package benchmarks

import (
	"reflect"

	"github.com/stykhanskyy/hiperabs/reference"
)

// ArrayList ...
type ArrayList struct {
	slice          interface{}
	t              reflect.Type
	size, capacity int
	sliceRef       reference.SliceRef
	sliceValue     reflect.Value
}

// NewArrayList ...
func NewArrayList(protoObject interface{}) ArrayList {
	t := reflect.TypeOf(protoObject)
	slice := reflect.MakeSlice(reflect.SliceOf(t), 0, 0).Interface()
	return ArrayList{
		slice:      slice,
		sliceRef:   reference.ToSlice(slice),
		sliceValue: reflect.ValueOf(slice),
		t:          t}
}

func (list *ArrayList) reserveCapacityForNextElement() {
	if list.size == list.capacity {
		newCapacity := list.capacity * 2
		if newCapacity == 0 {
			newCapacity++
		}
		newSlice := reflect.MakeSlice(reflect.SliceOf(list.t), newCapacity, newCapacity).Interface()
		reflect.Copy(reflect.ValueOf(newSlice), reflect.ValueOf(list.slice))
		list.slice = newSlice
		list.sliceRef = reference.ToSlice(newSlice)
		list.sliceValue = reflect.ValueOf(newSlice)
		list.capacity = newCapacity
	}
}

// AppendRef ...
func (list *ArrayList) AppendRef(element *reference.Ref) {
	list.reserveCapacityForNextElement()
	reference.Copy(list.sliceRef.ElementRef(list.size), element)
	list.size++
}

// AppendValue ...
func (list *ArrayList) AppendValue(element *reflect.Value) {
	list.reserveCapacityForNextElement()
	list.sliceValue.Index(list.size).Set(*element)
	list.size++
}

// Length ...
func (list *ArrayList) Length() int {
	return list.size
}

// GetSlice ...
func (list *ArrayList) GetSlice() interface{} {
	return reflect.ValueOf(list.slice).Slice(0, list.size).Interface()
}
