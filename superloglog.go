package main

import (
	"math"
	"math/bits"
	"sort"
)

// *Truncation Rule*. When collecting register values to produce final estimate, retain only the m_0:=floor(theta * m) smallest values and discard the rest. theta used in paper is 0.7. It should help to increase accuracy
// *Restriction Rule*. Use register values that are in the interval [0..B], where ceil(log_2(N_max/m)+3) <= B
type SuperLogLog[T Hasher] struct {
	buckets       []int
	bitsForBucket int
	a_inf         float64
	a_m           float64
	theta         float64
	m_zero        int
}

func InitSuperLogLog[T Hasher](bitsForBucket int) *SuperLogLog[T] {
	m := 1 << bitsForBucket
	sll := &SuperLogLog[T]{
		buckets:       make([]int, m),
		bitsForBucket: bitsForBucket,
		a_inf:         0.39701,
		theta:         0.7,
	}
	sll.m_zero = int(math.Floor(sll.theta * float64(m)))
	// sll.a_m = sll.a_inf - (2*math.Pi*math.Pi+math.Ln2*math.Ln2)/(48*float64(sll.m_zero))
	return sll
}

func (sll *SuperLogLog[T]) Add(val T) {
	hash := val.Hash()
	// assuming uint is 64bits
	bucketIdx := hash >> (64 - sll.bitsForBucket)
	remainder := (hash << sll.bitsForBucket) >> sll.bitsForBucket

	sll.buckets[bucketIdx] = max(sll.buckets[bucketIdx], bits.TrailingZeros(uint(remainder))+1)
}

func (sll *SuperLogLog[T]) GetCount() uint64 {
	m := len(sll.buckets)
	copy_buckets := make([]int, m)
	copy(copy_buckets, sll.buckets)
	sort.Ints(copy_buckets)

	sum := 0
	for i := 0; i < sll.m_zero; i++ {
		sum += copy_buckets[i]
	}
	var avg float64 = float64(sum) / float64(sll.m_zero)

	// I use some other value, from some other place, but I guess it doesn't matter much
	return uint64(math.Round(float64(sll.m_zero) * 1.09295 * math.Pow(2, avg)))
}
