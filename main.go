package main

import (
	"fmt"
	"hash/fnv"
)

type HashInt int

func (val HashInt) Hash() uint {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%d", val)))

	hashValue := h.Sum64()
	return uint(hashValue)
}

func main() {
	mc := InitMapCounter[HashInt](1024)
	ll := InitLogLog[HashInt](4)

	for i := 0; i < 100000; i++ {
		mc.Add(HashInt(i))
		ll.Add(HashInt(i))

		if i%1000 == 0 {
			fmt.Printf("%d %d\n", mc.GetCount(), ll.GetCount())
		}
	}

	fmt.Printf("MapCounter: %d\n", mc.GetCount())
	fmt.Printf("LogLog: %d\n", ll.GetCount())
}
