package middleware

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"

	"feed-service/internal/models"
)

func TestGetHealthz(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock init
	cfg := models.Config {
		ServiceVersion: "Test",
	}
	ctrl := Controller{
		Cfg: &cfg,
		DB: nil,
	}

	// record request
	rr := httptest.NewRecorder()

	// test router
	router := gin.Default()
	router.GET("/healthz", ctrl.GetHealthz)

	// mock request
	request, err := http.NewRequest(http.MethodGet, "/healthz", nil)
	assert.NoError(t, err)

	// make request
	router.ServeHTTP(rr, request)

	// expect OK with `Service is ready. Test`
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Service is ready. " + cfg.ServiceVersion.String(), rr.Body.String())
}
