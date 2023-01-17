package random

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestChoice(t *testing.T) {
	v := Choice([]int{1})
	assert.Equal(t, 1, v, "random choice")
}

func TestContains(t *testing.T) {
	has := Contains([]int{1, 2, 3}, 1)
	assert.Equal(t, true, has, "contains 1")

	has = Contains([]int{1, 2, 3}, 4)
	assert.Equal(t, false, has, "not contains 4")
}
