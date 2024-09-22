package main

type Hasher interface {
	comparable
	Hash() uint
}

type Counter[T Hasher] interface {
	Add(T)
	GetCount() uint64
}
