package middleware

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"

	"my-service/internal/models"
)

func TestGetHealthz(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := models.Config {
		ServiceVersion: "Test",
	}

	t.Run("GetHealthz", func(t *testing.T){
		rr := httptest.NewRecorder()

		router := gin.Default()

		request, err := http.NewRequest(http.MethodGet, "/healthz", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "Service is ready.", rr.Body.String())
	})
}
