package random

import (
	"errors"
	"math/rand"
)

var ErrEmptySlice = errors.New("empty slice")

func Choice[T comparable](elems []T) (T, error) {
	var zero T
	if len(elems) == 0 {
		return zero, ErrEmptySlice
	}
	idx := rand.Intn(len(elems))
	return elems[idx], nil
}
