package main

import (
	"math"
	"math/bits"
)

type MapCounter[T Hasher] struct {
	Map map[T]bool
}

func (mc *MapCounter[T]) Init(capacity int) {
	mc.Map = make(map[T]bool, capacity)
}

func (mc *MapCounter[T]) Add(val T) {
	if mc.Map == nil {
		mc.Map = make(map[T]bool)
	}

	mc.Map[val] = true
}

func (mc *MapCounter[T]) GetCount() uint64 {
	return uint64(len(mc.Map))
}

type LogLog[T Hasher] struct {
	buckets       []int
	bitsForBucket int
	constant      float64
}

func InitLogLog[T Hasher](bitsForBucket int) *LogLog[T] {
	return &LogLog[T]{
		buckets:       make([]int, 1<<(bitsForBucket)),
		bitsForBucket: bitsForBucket,
		constant:      0.79,
	}
}

func (ll *LogLog[T]) Add(val T) {
	hash := val.Hash()
	// assuming int is 64bits
	bucketIdx := hash >> (64 - ll.bitsForBucket)
	remainder := (hash << ll.bitsForBucket) >> ll.bitsForBucket

	ll.buckets[bucketIdx] = max(ll.buckets[bucketIdx], bits.TrailingZeros(uint(remainder)))
}

func (ll *LogLog[T]) GetCount() uint64 {
	m := len(ll.buckets)

	sum := 0
	for i := 0; i < m; i++ {
		sum += ll.buckets[i]
	}
	var am float64 = float64(sum) / float64(m)

	return uint64(math.Round(ll.constant * float64(m) * math.Pow(2, am)))
}
