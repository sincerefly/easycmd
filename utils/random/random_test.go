package random

import (
	"errors"
	"testing"
)

func TestChoice(t *testing.T) {
	v, err := Choice([]int{1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 1 {
		t.Errorf("choice {1} expected 1, but %d got", v)
	}

	v, err = Choice([]int{1, 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 1 && v != 2 {
		t.Errorf("choice {1,2} expected 1 or 2, but %d got", v)
	}

	f, err := Choice([]float64{1.1, 2.2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != 1.1 && f != 2.2 {
		t.Errorf("choice {1.1,2.2} expected 1.1 or 2.2, but %f got", f)
	}
}

func TestChoiceEmptySlice(t *testing.T) {
	_, err := Choice([]int{})
	if !errors.Is(err, ErrEmptySlice) {
		t.Fatalf("expected ErrEmptySlice, got %v", err)
	}
}

func TestContains(t *testing.T) {
	if !Contains([]int{1, 2, 3}, 1) {
		t.Error("contains {1,2,3} value 1 expected true")
	}
	if Contains([]int{1, 2, 3}, 4) {
		t.Error("contains {1,2,3} value 4 expected false")
	}
}
