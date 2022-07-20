package models

import (
	"testing"
	"os"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	var e EnvVar
	key, val := "FOO", "BAR"

	// prepare playground
	os.Setenv(key, val)

	t.Run("FOO BAR", func (t *testing.T) {
		t.Parallel()

		// seems bad, but no way around
		e.GetEnv(key)
		assert.Equal(t, val, e.String())
	})

	t.Run("panic", func (t *testing.T) {
		t.Parallel()

		// expected not to appear in env
		key2 := "BUZZ"
		os.Unsetenv(key2)
		assert.Panics(t, func() { e.GetEnv(key2) }, "`BUZZ` is not supposed to be set")
	})
}
