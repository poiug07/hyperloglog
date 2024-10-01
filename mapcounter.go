package main

import (
	"math"
	"math/bits"
	"math/rand"
)

type MapCounter[T Hasher] struct {
	Map map[T]bool
}

func InitMapCounter[T Hasher](capacity int) *MapCounter[T] {
	return &MapCounter[T]{
		Map: make(map[T]bool, capacity),
	}
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

// Not thread safe
type MorrisCountingAlgo struct {
	Rand    rand.Rand
	counter uint64
}

func InitMorrisCountingAlgo(rng rand.Rand) *MorrisCountingAlgo {
	return &MorrisCountingAlgo{Rand: rng}
}

func (mca *MorrisCountingAlgo) Add() {
	// TODO: more efficient code
	var v uint = 1
	var i uint64
	for i = 0; i < mca.counter; i++ {
		v &= uint(rand.Intn(2))
	}
	if v == 1 {
		mca.counter++
	}
}

func (mca *MorrisCountingAlgo) GetCount() uint64 {
	return 1 << mca.counter
}

type LogLog[T Hasher] struct {
	buckets       []int
	bitsForBucket int
	a_inf         float64
	a_m           float64
}

func InitLogLog[T Hasher](bitsForBucket int) *LogLog[T] {
	m := 1 << bitsForBucket
	ll := &LogLog[T]{
		buckets:       make([]int, m),
		bitsForBucket: bitsForBucket,
		a_inf:         0.39701,
	}
	ll.a_m = ll.a_inf - (2*math.Pi*math.Pi+math.Log10(2)*math.Log10(2))/(48*float64(m))
	return ll
}

func (ll *LogLog[T]) Add(val T) {
	hash := val.Hash()
	// assuming uint is 64bits
	bucketIdx := hash >> (64 - ll.bitsForBucket)
	remainder := (hash << ll.bitsForBucket) >> ll.bitsForBucket

	ll.buckets[bucketIdx] = max(ll.buckets[bucketIdx], bits.TrailingZeros(uint(remainder))+1)
}

func (ll *LogLog[T]) GetCount() uint64 {
	m := len(ll.buckets)

	sum := 0
	for i := 0; i < m; i++ {
		sum += ll.buckets[i]
	}
	var am float64 = float64(sum) / float64(m)

	return uint64(math.Round(ll.a_m * float64(m) * math.Pow(2, am)))
}
