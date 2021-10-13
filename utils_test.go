package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_minPositiveInt(t *testing.T) {
	assert.Equal(t, minPositiveInt(5, 6), 5)
	assert.Equal(t, minPositiveInt(15, 6), 6)
	assert.Equal(t, minPositiveInt(5, -6), 5)
	assert.Equal(t, minPositiveInt(1, 20), 1)
	assert.Equal(t, minPositiveInt(-1, 0), 0)
	assert.Equal(t, minPositiveInt(22, 3), 3)
}

func Test_minInt(t *testing.T) {
	assert.Equal(t, minInt(5, 6), 5)
	assert.Equal(t, minInt(15, 6), 6)
	assert.Equal(t, minInt(5, -6), -6)
	assert.Equal(t, minInt(1, 20), 1)
	assert.Equal(t, minInt(-1, 0), -1)
	assert.Equal(t, minInt(22, 3), 3)
}
