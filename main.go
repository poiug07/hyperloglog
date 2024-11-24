package main

import (
	"fmt"
	"hash/maphash"
)

type HashInt int

var seed maphash.Seed = maphash.MakeSeed()

func (val HashInt) Hash() uint {
	return uint(maphash.Bytes(seed, []byte(fmt.Sprintf("%d", val))))
}

func main() {
	mc := InitMapCounter[HashInt](1024)
	ll := InitLogLog[HashInt](10)
	sll := InitSuperLogLog[HashInt](10)

	for i := 0; i < 10000000; i++ {
		mc.Add(HashInt(i))
		ll.Add(HashInt(i))
		sll.Add(HashInt(i))

		if i%1000000 == 0 {
			fmt.Printf("%d %d %d\n", mc.GetCount(), ll.GetCount(), sll.GetCount())
		}
	}

	fmt.Printf("MapCounter: %d\n", mc.GetCount())
	fmt.Printf("LogLog: %d\n", ll.GetCount())
	fmt.Printf("SuperLogLog %d\n", sll.GetCount())
}
