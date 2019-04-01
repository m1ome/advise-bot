package fga

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Random(t *testing.T) {
	c := New()
	advise, err := c.Random()

	assert.NoError(t, err)
	assert.True(t, len(advise) > 0)
}
