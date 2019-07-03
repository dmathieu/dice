package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupFlagValue(t *testing.T) {
	assert.Equal(t, flagValue, "roll")
	Setup(FlagValue("roll-test"))
	assert.Equal(t, "roll-test", flagValue)
	Setup(FlagValue("roll"))
}
