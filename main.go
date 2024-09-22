package main

import (
	"fmt"
	"hash/fnv"
)

type HashInt int

func (val HashInt) Hash() uint {
	// h := crc64.New(crc64.MakeTable(128))
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%d", val)))

	hashValue := h.Sum64()
	return uint(hashValue)
}

func main() {
	mc := new(MapCounter[HashInt])
	mc.Init(128)

	ll := InitLogLog[HashInt](4)
	for i := 0; i < 100000; i++ {
		mc.Add(HashInt(i))
		ll.Add(HashInt(i))
	}

	fmt.Printf("MapCounter: %d\n", mc.GetCount())
	fmt.Printf("LogLog: %d\n", ll.GetCount())
}
