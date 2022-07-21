package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"database/sql"
	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"

	"my-service/internal/models"
)

// google has underlying function covered
// need to validate true/false only
func TestIsValidUUID(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		uuid := "eaeaa9c9-85c0-4c53-9309-9d499c6c0026"
		assert.Equal(t, true, isValidUUID(uuid))
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		uuid := "asd"
		assert.Equal(t, false, isValidUUID(uuid))
	})
}

type posts struct {
	Total		int		`json:"total"`
	Data		[]models.Post	`json:"data"`
}

// SELECT on empty table (sql.ErrNoRows)
func TestGetPostsEmpty(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock init
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	ctrl := Controller{
		DB: db,
	}

	// register request
	rr := httptest.NewRecorder()

	// set up test router
	router := gin.Default()
	router.GET("/posts", ctrl.GetPosts)

	mock.
		ExpectQuery("SELECT uid, content, likes, dislikes FROM posts").
		WillReturnError(sql.ErrNoRows)

	// mock request
	request, err := http.NewRequest(http.MethodGet, "/posts", nil)
	assert.NoError(t, err)

	// make request
	router.ServeHTTP(rr, request)

	var p posts

	// convert body to `posts`
	err = json.NewDecoder(rr.Body).Decode(&p)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, 0, p.Total)
	assert.EqualValues(t, []models.Post(nil), p.Data)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// query on failed db pool
func TestGetPostsNoDB(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock init
	connString := "postgres://postgres:postgres@123.123.123.123:123/data?connect_timeout=1"
	db, err := sql.Open("postgres", connString)
	assert.NoError(t, err)

	ctrl := Controller {
		DB: db,
	}

	// register request
	rr := httptest.NewRecorder()

	// set up test router
	router := gin.Default()
	router.GET("/posts", ctrl.GetPosts)

	// mock request
	request, err := http.NewRequest(http.MethodGet, "/posts", nil)
	assert.NoError(t, err)

	// make request
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
