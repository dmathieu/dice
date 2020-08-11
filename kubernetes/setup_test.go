package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupFlagValue(t *testing.T) {
	assert.Equal(t, flagValue, "roll")
	assert.NoError(t, Setup(FlagValue("roll-test")))
	assert.Equal(t, "roll-test", flagValue)
	assert.NoError(t, Setup(FlagValue("roll")))
}
