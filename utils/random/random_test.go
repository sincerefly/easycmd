package random

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestChoice(t *testing.T) {
	if v := Choice([]int{1}); v != 1 {
		t.Errorf("choice {1} expected 1, but %d got", v)
	}
	if v := Choice([]int{1, 2}); v != 1 && v != 2 {
		t.Errorf("choice {1,2} expected 1 or 2, but %d got", v)
	}
	if v := Choice([]float64{1.1, 2.2}); v != 1.1 && v != 2.2 {
		t.Errorf("choice {1,2} expected 1 or 2, but %f got", v)
	}
}

func TestContains(t *testing.T) {
	has := Contains([]int{1, 2, 3}, 1)
	assert.Equal(t, true, has, "contains {1,2,3} value 1 expected true")

	has = Contains([]int{1, 2, 3}, 4)
	assert.Equal(t, false, has, "contains {1,2,3} value 4 expected false")
}
