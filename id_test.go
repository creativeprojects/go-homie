package homie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidID(t *testing.T) {
	testData := []string{
		"valid", "valid-id", "valid123", "valid-123", "AlsoValid",
	}
	for _, id := range testData {
		t.Run(id, func(t *testing.T) {
			assert.True(t, IsValidID(id))
		})
	}
}

func TestInvalidID(t *testing.T) {
	testData := []string{
		"-invalid", "invalid-", "not_valid", "also not",
	}
	for _, id := range testData {
		t.Run(id, func(t *testing.T) {
			assert.False(t, IsValidID(id))
		})
	}
}
