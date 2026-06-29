package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputePercentages(t *testing.T) {
	result := computePercentages(map[string]int{
		"Go":         300,
		"TypeScript": 100,
	})

	assert.Equal(t, 75.0, result["Go"])
	assert.Equal(t, 25.0, result["TypeScript"])

	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	assert.Len(t, keys, 2)
}

func TestComputePercentagesEmpty(t *testing.T) {
	result := computePercentages(map[string]int{})
	assert.Empty(t, result)
}

func TestVisibilityKey(t *testing.T) {
	assert.Equal(t, "private", visibilityKey(true))
	assert.Equal(t, "public", visibilityKey(false))
}
