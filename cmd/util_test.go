package cmd

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseWatchFrequency(t *testing.T) {

	t.Run("with an invalid value", func(t *testing.T) {
		f, err := parseWatchFrequency("foobar")
		assert.Nil(t, f)
		assert.Equal(t, errors.New("invalid watch frequency \"foobar\""), err)
	})

	t.Run("with a value of seconds", func(t *testing.T) {
		f, err := parseWatchFrequency("27s")
		assert.Nil(t, err)
		assert.Equal(t, 27*time.Second, *f)
	})

	t.Run("with a value of minutes", func(t *testing.T) {
		f, err := parseWatchFrequency("32m")
		assert.Nil(t, err)
		assert.Equal(t, 32*time.Minute, *f)
	})

	t.Run("with a value of hours", func(t *testing.T) {
		f, err := parseWatchFrequency("164h")
		assert.Nil(t, err)
		assert.Equal(t, 164*time.Hour, *f)
	})

	t.Run("with a value of days", func(t *testing.T) {
		f, err := parseWatchFrequency("65d")
		assert.Nil(t, err)
		assert.Equal(t, 65*24*time.Hour, *f)
	})
}
