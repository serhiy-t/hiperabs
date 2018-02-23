package benchmarks

import (
	"reflect"

	"github.com/stykhanskyy/hiperabs/reference"
)

// Hashable ...
type Hashable interface {
	Hash() int
}

type hashBucket struct {
	begin, end, size int
}

// HashSet ...
type HashSet struct {
	t                      reflect.Type
	data                   interface{}
	dataHashes             []int
	dataRef                reference.SliceRef
	dataSize, dataCapacity int

	buckets     []hashBucket
	bucketsSize int // should be divisible by 2

	size int
}

// NewHashSet ...
func NewHashSet(protoObject interface{}) HashSet {
	if _, ok := protoObject.(Hashable); !ok {
		panic("Non-hashable object")
	}
	t := reflect.TypeOf(protoObject)
	data := reflect.MakeSlice(reflect.SliceOf(t), 0, 0).Interface()
	dataRef := reference.ToSlice(data)
	buckets := make([]hashBucket, 8, 8)

	return HashSet{
		t:            t,
		data:         data,
		dataHashes:   make([]int, 0, 0),
		dataRef:      dataRef,
		dataSize:     0,
		dataCapacity: 0,
		buckets:      buckets,
		bucketsSize:  8,
		size:         0,
	}
}

func (set *HashSet) allocateBucket(capacity int) hashBucket {
	if set.dataSize+capacity > set.dataCapacity {
		newCapacity := set.dataCapacity * 2
		if set.dataSize+capacity > newCapacity {
			newCapacity = set.dataSize + capacity
		}
		newData := reflect.MakeSlice(reflect.SliceOf(set.t), newCapacity, newCapacity).Interface()
		newDataRef := reference.ToSlice(newData)
		newDataHashes := make([]int, newCapacity, newCapacity)
		// reflect.Copy(reflect.ValueOf(newData), reflect.ValueOf(set.data))
		// copy(newDataHashes, set.dataHashes)

		newDataPtr := 0
		for bucketIdx := 0; bucketIdx < set.bucketsSize; bucketIdx++ {
			bucket := &set.buckets[bucketIdx]
			for dataIdx := 0; dataIdx < bucket.size; dataIdx++ {
				reference.Copy(
					newDataRef.ElementRef(newDataPtr+dataIdx),
					set.dataRef.ElementRef(bucket.begin+dataIdx))
				newDataHashes[newDataPtr+dataIdx] = set.dataHashes[bucket.begin+dataIdx]
			}
			bucket.begin = newDataPtr
			bucket.end = newDataPtr + bucket.size + 1
			newDataPtr += bucket.size + 1
		}

		// compact copy
		set.data = newData
		set.dataHashes = newDataHashes
		set.dataCapacity = newCapacity
		set.dataRef = newDataRef
	}
	return hashBucket{set.dataSize, set.dataSize + capacity, 0}
}

func (set *HashSet) rebalance() {
	if set.size*4 < set.bucketsSize*3 {
		// if set.size*13 < set.bucketsSize*2 {
		return
	}

	newBucketsSize := set.bucketsSize * 2
	newBuckets := make([]hashBucket, newBucketsSize, newBucketsSize)

	for bucketIdx := 0; bucketIdx < set.bucketsSize; bucketIdx++ {
		bucket := set.buckets[bucketIdx]
		rehash1 := 0
		rehash1Count := 0
		rehash2 := 0
		rehash2Count := 0

		for dataIdx := 0; dataIdx < bucket.size; dataIdx++ {
			rehash := set.dataHashes[bucket.begin+dataIdx] % newBucketsSize
			if rehash < bucket.size {
				rehash1 = rehash
				rehash1Count++
			} else {
				rehash2 = rehash
				rehash2Count++
			}
		}

		if rehash1Count == 0 && rehash2Count == 0 {
			// no data, nothing to do
		} else if rehash1Count == 0 {
			// reuse bucket
			newBuckets[rehash2] = bucket
		} else if rehash2Count == 0 {
			// reuse bucket
			newBuckets[rehash1] = bucket
		} else {
			// split bucket
			newBucket1 := set.allocateBucket(rehash1Count)
			newBucket2 := set.allocateBucket(rehash2Count)

			for dataIdx := 0; dataIdx < bucket.size; dataIdx++ {
				rehash := set.dataHashes[bucket.begin+dataIdx] % newBucketsSize
				if rehash < bucket.size {
					reference.Copy(
						set.dataRef.ElementRef(newBucket1.begin+newBucket1.size),
						set.dataRef.ElementRef(bucket.begin+dataIdx))
					newBucket1.size++
				} else {
					reference.Copy(
						set.dataRef.ElementRef(newBucket2.begin+newBucket2.size),
						set.dataRef.ElementRef(bucket.begin+dataIdx))
					newBucket2.size++
				}
			}
			newBuckets[rehash1] = newBucket1
			newBuckets[rehash2] = newBucket2
		}
	}

	set.buckets = newBuckets
	set.bucketsSize = newBucketsSize
}

func (set *HashSet) Add(element reference.PointerRef) {
	hash := element.Object().(Hashable).Hash()
	bucketNumber := hash % set.bucketsSize
	bucket := &set.buckets[bucketNumber]
	if bucket.size+1 >= bucket.end-bucket.begin {
		newBucket := set.allocateBucket((bucket.end-bucket.begin+1)*2 - 1)
		for dataIdx := 0; dataIdx < bucket.size; dataIdx++ {
			reference.Copy(
				set.dataRef.ElementRef(newBucket.begin+dataIdx),
				set.dataRef.ElementRef(bucket.begin+dataIdx))
			set.dataHashes[newBucket.begin+dataIdx] = set.dataHashes[bucket.begin+dataIdx]
		}
		newBucket.size = bucket.size
		*bucket = newBucket
	}
	reference.Copy(
		set.dataRef.ElementRef(bucket.begin+bucket.size),
		element.Ref())
	set.dataHashes[bucket.begin+bucket.size] = hash
	bucket.size++
	set.size++

	set.rebalance()
}
