package random

import (
	"math/rand"
)

func Choice[T comparable](elems []T) T {
	idx := rand.Intn(len(elems))
	return elems[idx]
}

func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
