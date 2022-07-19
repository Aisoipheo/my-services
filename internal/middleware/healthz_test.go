package middleware

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
)

func TestGetHealthz(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("GetHealthz", func(t *testing.T){
		rr := httptest.NewRecorder()

		router := gin.Default()

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "Service is ready. 0.0.1-alpha", rr.Body.String())
	})
}
