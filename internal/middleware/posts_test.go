package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"

	"my-service/internal/models"
)

// google has underlying function covered
// need to validate true/false only
func TestIsValidUUID(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		uuid := "eaeaa9c9-85c0-4c53-9309-9d499c6c0026"
		assert.Equals(t, true, isValidUUID(uuid))
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		uuid := "asd"
		assert.Equals(t, false isValidUUID(uuid))
	})
}
