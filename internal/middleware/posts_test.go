package middleware

import (
	"io"
	"bytes"
	"errors"
	"regexp"
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

type posts struct {
	Total	int				`json:"total"`
	Data	[]models.Post	`json:"data"`
}

type badPost struct {
	WrongContent string `json:"wrongcontent"`
}

type emptyPost struct {} // yes should be empty

func TestInit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	gin.EnableJsonDecoderDisallowUnknownFields()
}

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

// SELECT on empty table (sql.ErrNoRows)
func TestGetPostsEmpty(t *testing.T) {
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
		ExpectQuery("SELECT uuid, content, likes, dislikes FROM posts").
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
	assert.EqualValues(t, []models.Post{}, p.Data)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// query on failed db pool
func TestGetPostsNoDB(t *testing.T) {
	// Mock init
	connString := "postgres://postgres:postgres@123.123.123.123:123/data?connect_timeout=1"
	db, err := sql.Open("postgres", connString)
	assert.NoError(t, err)
	defer db.Close()

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

// SELECT Correct
func TestGetPostsOKsingle(t *testing.T) {
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

	mockPost := models.Post {
		UUID: "123",
		Content: "simple text",
		Likes: 123,
		Dislikes: 321,
	}

	rows := sqlmock.
		NewRows([]string{"uuid", "content", "likes", "dislikes"}).
		AddRow(mockPost.UUID, mockPost.Content, mockPost.Likes, mockPost.Dislikes)

	mock.
		ExpectQuery("SELECT uuid, content, likes, dislikes FROM posts").
		WillReturnRows(rows)

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
	assert.Equal(t, 1, p.Total)
	assert.EqualValues(t, []models.Post{mockPost}, p.Data)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// SELECT Correct
func TestGetPostsOKmultiple(t *testing.T) {
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

	mockPost := models.Post {
		UUID: "123",
		Content: "simple text",
		Likes: 123,
		Dislikes: 321,
	}

	rows := sqlmock.
		NewRows([]string{"uuid", "content", "likes", "dislikes"}).
		AddRow(mockPost.UUID, mockPost.Content, mockPost.Likes, mockPost.Dislikes).
		AddRow(mockPost.UUID, mockPost.Content, mockPost.Likes, mockPost.Dislikes).
		AddRow(mockPost.UUID, mockPost.Content, mockPost.Likes, mockPost.Dislikes)

	mock.
		ExpectQuery("SELECT uuid, content, likes, dislikes FROM posts").
		WillReturnRows(rows)

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
	assert.Equal(t, 3, p.Total)
	assert.EqualValues(t, []models.Post{mockPost, mockPost, mockPost}, p.Data)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// SELECT on wrong fields
func TestGetPostsScanErr(t *testing.T) {
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

	mockPost := models.Post {
		UUID: "123",
		Content: "simple text",
		Likes: 123,
		Dislikes: 321,
	}

	rows := sqlmock.
		NewRows([]string{"uuid", "content", "likes", "dislikes"}).
		AddRow(mockPost.UUID, mockPost.Likes, mockPost.Content, mockPost.Dislikes)

	mock.
		ExpectQuery("SELECT uuid, content, likes, dislikes FROM posts").
		WillReturnRows(rows)

	// mock request
	request, err := http.NewRequest(http.MethodGet, "/posts", nil)
	assert.NoError(t, err)

	// make request
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// SELECT Correct with parameter
func TestGetPostsOKParam(t *testing.T) {
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

	mockPost := models.Post {
		UUID: "123",
		Content: "simple text",
		Likes: 123,
		Dislikes: 321,
	}

	rows := sqlmock.
		NewRows([]string{"uuid", "content", "likes", "dislikes"}).
		AddRow(mockPost.UUID, mockPost.Content, mockPost.Likes, mockPost.Dislikes)

	mock.
		ExpectQuery("SELECT uuid, content, likes, dislikes FROM posts LIMIT 1").
		WillReturnRows(rows)

	// mock request
	request, err := http.NewRequest(http.MethodGet, "/posts?last=1", nil)
	assert.NoError(t, err)

	// make request
	router.ServeHTTP(rr, request)

	var p posts

	// convert body to `posts`
	err = json.NewDecoder(rr.Body).Decode(&p)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, 1, p.Total)
	assert.EqualValues(t, []models.Post{mockPost}, p.Data)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// SELECT Correct with parameter from multiple
func TestGetPostsOKParamMultiple(t *testing.T) {
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

	mockPost := models.Post {
		UUID: "123",
		Content: "simple text",
		Likes: 123,
		Dislikes: 321,
	}

	mockPost2 := models.Post {
		UUID: "456",
		Content: "hard text",
		Likes: 0,
		Dislikes: 999,
	}

	rows2 := sqlmock.
		NewRows([]string{"uuid", "content", "likes", "dislikes"}).
		AddRow(mockPost.UUID, mockPost.Content, mockPost.Likes, mockPost.Dislikes).
		AddRow(mockPost2.UUID, mockPost2.Content, mockPost2.Likes, mockPost2.Dislikes)

	_ = sqlmock.
		NewRows([]string{"uuid", "content", "likes", "dislikes"}).
		AddRow(mockPost.UUID, mockPost.Content, mockPost.Likes, mockPost.Dislikes).
		AddRow(mockPost2.UUID, mockPost2.Content, mockPost2.Likes, mockPost2.Dislikes).
		AddRow(mockPost.UUID, mockPost.Content, mockPost.Likes, mockPost.Dislikes).
		AddRow(mockPost.UUID, mockPost.Content, mockPost.Likes, mockPost.Dislikes).
		AddRow(mockPost.UUID, mockPost.Content, mockPost.Likes, mockPost.Dislikes).
		AddRow(mockPost.UUID, mockPost.Content, mockPost.Likes, mockPost.Dislikes).
		AddRow(mockPost.UUID, mockPost.Content, mockPost.Likes, mockPost.Dislikes).
		AddRow(mockPost.UUID, mockPost.Content, mockPost.Likes, mockPost.Dislikes)

	mock.
		ExpectQuery("SELECT uuid, content, likes, dislikes FROM posts LIMIT 2").
		WillReturnRows(rows2)

	// mock request
	request, err := http.NewRequest(http.MethodGet, "/posts?last=2", nil)
	assert.NoError(t, err)

	// make request
	router.ServeHTTP(rr, request)

	var p posts

	// convert body to `posts`
	err = json.NewDecoder(rr.Body).Decode(&p)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)
	// assert.Equal(t, 2, p.Total)
	// assert.EqualValues(t, []models.Post{mockPost, mockPost2}, p.Data)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// SELECT Correct with parameter from multiple
func TestGetPostsOKParamMultipleZero(t *testing.T) {
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

	mockPost := models.Post {
		UUID: "123",
		Content: "simple text",
		Likes: 123,
		Dislikes: 321,
	}

	mockPost2 := models.Post {
		UUID: "456",
		Content: "hard text",
		Likes: 0,
		Dislikes: 999,
	}

	_ = sqlmock.
		NewRows([]string{"uuid", "content", "likes", "dislikes"}).
		AddRow(mockPost.UUID, mockPost.Content, mockPost.Likes, mockPost.Dislikes).
		AddRow(mockPost2.UUID, mockPost2.Content, mockPost2.Likes, mockPost2.Dislikes).
		AddRow(mockPost.UUID, mockPost.Content, mockPost.Likes, mockPost.Dislikes)

	mock.
		ExpectQuery("SELECT uuid, content, likes, dislikes FROM posts LIMIT 0").
		WillReturnError(sql.ErrNoRows)

	// mock request
	request, err := http.NewRequest(http.MethodGet, "/posts?last=0", nil)
	assert.NoError(t, err)

	// make request
	router.ServeHTTP(rr, request)

	var p posts

	// convert body to `posts`
	err = json.NewDecoder(rr.Body).Decode(&p)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, 0, p.Total)
	assert.EqualValues(t, []models.Post{}, p.Data)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// SELECT Bad parameter
func TestGetPostsEmptyParam(t *testing.T) {
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

	mockPost := models.Post {
		UUID: "123",
		Content: "simple text",
		Likes: 123,
		Dislikes: 321,
	}

	_ = sqlmock.
		NewRows([]string{"uuid", "content", "likes", "dislikes"}).
		AddRow(mockPost.UUID, mockPost.Content, mockPost.Likes, mockPost.Dislikes)

	// mock request
	request, err := http.NewRequest(http.MethodGet, "/posts?last=", nil)
	assert.NoError(t, err)

	// make request
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// SELECT Bad parameter
func TestGetPostsBadParam(t *testing.T) {
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

	mockPost := models.Post {
		UUID: "123",
		Content: "simple text",
		Likes: 123,
		Dislikes: 321,
	}

	_ = sqlmock.
		NewRows([]string{"uuid", "content", "likes", "dislikes"}).
		AddRow(mockPost.UUID, mockPost.Content, mockPost.Likes, mockPost.Dislikes)

	// mock request
	request, err := http.NewRequest(http.MethodGet, "/posts?last=asd", nil)
	assert.NoError(t, err)

	// make request
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// SELECT Bad parameter
func TestGetPostsNegativeParam(t *testing.T) {
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

	mockPost := models.Post {
		UUID: "123",
		Content: "simple text",
		Likes: 123,
		Dislikes: 321,
	}

	_ = sqlmock.
		NewRows([]string{"uuid", "content", "likes", "dislikes"}).
		AddRow(mockPost.UUID, mockPost.Content, mockPost.Likes, mockPost.Dislikes)

	// mock request
	request, err := http.NewRequest(http.MethodGet, "/posts?last=-1", nil)
	assert.NoError(t, err)

	// make request
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSingleTransactionNoDB(t *testing.T) {
	// Mock init
	connString := "postgres://postgres:postgres@123.123.123.123:123/data?connect_timeout=1"
	db, err := sql.Open("postgres", connString)
	assert.NoError(t, err)
	defer db.Close()

	ctrl := Controller{
		DB: db,
	}

	// register request
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	assert.NoError(t, err)

	stmt := "SELECT 1;"

	// actual function call
	singleTransaction(&ctrl, c, stmt)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSingleTransactionBadStmt(t *testing.T) {
	// Mock init
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	ctrl := Controller{
		DB: db,
	}

	// register request
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	assert.NoError(t, err)

	stmt := "SELECT =;"

	mock.ExpectBegin()
	mock.
		ExpectPrepare(stmt).
		WillReturnError(errors.New(""))
	mock.ExpectRollback()

	// actual function call
	singleTransaction(&ctrl, c, stmt)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSingleTransactionBadArgs(t *testing.T) {
	// Mock init
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	ctrl := Controller{
		DB: db,
	}

	// register request
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	assert.NoError(t, err)

	stmt := "SELECT * FROM t WHERE id = $1;"

	mock.ExpectBegin()
	mock.
		ExpectPrepare(regexp.QuoteMeta(stmt)).
		ExpectExec().
		WillReturnError(errors.New(""))
	mock.ExpectRollback()

	// actual function call
	singleTransaction(&ctrl, c, stmt)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSingleTransactionOKCommit(t *testing.T) {
	// Mock init
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	ctrl := Controller{
		DB: db,
	}

	// register request
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	assert.NoError(t, err)

	stmt := "SELECT * FROM t WHERE id = $1;"

	mock.ExpectBegin()
	mock.
		ExpectPrepare(regexp.QuoteMeta(stmt)).
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(1, 1)) // firt result, 1 row affected
	mock.ExpectCommit()

	// actual function call
	singleTransaction(&ctrl, c, stmt, "123")

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostLikeOK(t *testing.T) {
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
	router.GET("/like", ctrl.PostLike)

	stmt := "UPDATE posts SET likes = likes + 1 WHERE uuid = $1;"

	mock.ExpectBegin()
	mock.
		ExpectPrepare(regexp.QuoteMeta(stmt)).
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(1, 1)) // firt result, 1 row affected
	mock.ExpectCommit()

	// mock request
	request, err := http.NewRequest(http.MethodGet, "/like?uuid=78204138-90c6-49f7-90d9-1461d5d640f8", nil)
	assert.NoError(t, err)

	// make request
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostLikeNoUUID(t *testing.T) {
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
	router.GET("/like", ctrl.PostLike)

	_ = "UPDATE posts SET likes = likes + 1 WHERE uuid = $1;"

	// mock request
	request, err := http.NewRequest(http.MethodGet, "/like", nil)
	assert.NoError(t, err)

	// make request
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostLikeBadUUID(t *testing.T) {
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
	router.GET("/like", ctrl.PostLike)

	_ = "UPDATE posts SET likes = likes + 1 WHERE uuid = $1;"

	// mock request
	request, err := http.NewRequest(http.MethodGet, "/like?uuid=78204138-90c6-49f7-90d9", nil)
	assert.NoError(t, err)

	// make request
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostDislikeOK(t *testing.T) {
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
	router.GET("/dislike", ctrl.PostDislike)

	stmt := "UPDATE posts SET dislikes = dislikes + 1 WHERE uuid = $1;"

	mock.ExpectBegin()
	mock.
		ExpectPrepare(regexp.QuoteMeta(stmt)).
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(1, 1)) // firt result, 1 row affected
	mock.ExpectCommit()

	// mock request
	request, err := http.NewRequest(http.MethodGet, "/dislike?uuid=78204138-90c6-49f7-90d9-1461d5d640f8", nil)
	assert.NoError(t, err)

	// make request
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostDisikeNoUUID(t *testing.T) {
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
	router.GET("/dislike", ctrl.PostDislike)

	_ = "UPDATE disposts SET likes = dislikes + 1 WHERE uuid = $1;"

	// mock request
	request, err := http.NewRequest(http.MethodGet, "/dislike", nil)
	assert.NoError(t, err)

	// make request
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostDislikeBadUUID(t *testing.T) {
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
	router.GET("/dislike", ctrl.PostDislike)

	_ = "UPDATE posts SET dislikes = dislikes + 1 WHERE uuid = $1;"

	// mock request
	request, err := http.NewRequest(http.MethodGet, "/dislike?uuid=78204138-90c6-49f7-90d9", nil)
	assert.NoError(t, err)

	// make request
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostNewPostOK(t *testing.T) {
	// Mock init
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	ctrl := Controller{
		DB: db,
	}

	// register request
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)

	np := newPostRequestBody {
		Content: "New message",
	}

	c.Request = &http.Request {
		Header: make(http.Header),
	}
	assert.Equal(t, true, c.Request != nil) //non-zero return

	c.Request.Method = http.MethodPost
	c.Request.Header.Set("Content-Type", "application/json")

	jbytes, err := json.Marshal(np)
	assert.NoError(t, err)

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jbytes))

	stmt := "INSERT INTO posts(content) VALUES ($1);"

	mock.ExpectBegin()
	mock.
		ExpectPrepare(regexp.QuoteMeta(stmt)).
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(1, 1)) // firt result, 1 row affected
	mock.ExpectCommit()

	// actual function call
	ctrl.PostNewPost(c)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostNewPostBadJson(t *testing.T) {
	// Mock init
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	ctrl := Controller{
		DB: db,
	}

	// register request
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)

	np := badPost {
		WrongContent: "New message",
	}

	c.Request = &http.Request {
		Header: make(http.Header),
	}
	assert.Equal(t, true, c.Request != nil) //non-zero return

	c.Request.Method = http.MethodPost
	c.Request.Header.Set("Content-Type", "application/json")

	jbytes, err := json.Marshal(np)
	assert.NoError(t, err)

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jbytes))

	_ = "INSERT INTO posts(content) VALUES ($1);"

	// actual function call
	ctrl.PostNewPost(c)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostNewPostEmptyJson(t *testing.T) {
	// Mock init
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	ctrl := Controller{
		DB: db,
	}

	// register request
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)

	np := emptyPost {}

	c.Request = &http.Request {
		Header: make(http.Header),
	}
	assert.Equal(t, true, c.Request != nil) //non-zero return

	c.Request.Method = http.MethodPost
	c.Request.Header.Set("Content-Type", "application/json")

	jbytes, err := json.Marshal(np)
	assert.NoError(t, err)

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jbytes))

	_ = "INSERT INTO posts(content) VALUES ($1);"

	// actual function call
	ctrl.PostNewPost(c)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
