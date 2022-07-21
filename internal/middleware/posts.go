package middleware

import (
	"strconv"
	"net/http"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"my-service/internal/models"
)

type newPostRequestBody struct {
	Content		string		`json:"content"`
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func singleTransaction(h *Controller, c *gin.Context, queryString string, params ...interface{}) {
	tx, err := h.DB.Begin()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	stmt, err := tx.Prepare(queryString)
	if err != nil {
		// `_ =` to silence lint, no way to react to this
		_ = tx.Rollback()
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(params...)
	if err != nil {
		// `_ =` to silence lint, no way to react to this
		_ = tx.Rollback()
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		// `_ =` to silence lint, no way to react to this
		_ = tx.Rollback()
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Controller) GetPosts(c *gin.Context) {
	queryString := "SELECT uuid, content, likes, dislikes FROM posts"

	if queryLimitString, ok := c.GetQuery("last"); ok {
		if _, err := strconv.ParseUint(queryLimitString, 10, 64); err == nil {
			queryString += " LIMIT " + queryLimitString
		} else {
			c.String(http.StatusBadRequest, "Parameter `last` is invalid.\n`last`=" + queryLimitString)
			return
		}
	}

	rows, err := h.DB.Query(queryString)
	switch {
	// should return empty response with 200 StatusOK
	// REASON: could be SELECT on empty table
	case err == sql.ErrNoRows:
		c.JSON(http.StatusOK, gin.H { "total": 0, "data": []models.Post{} })
		return
	case err != nil:
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts := make([]models.Post, 0, 32)
	for rows.Next() {
		post := models.Post{}
		err = rows.Scan(&post.UUID, &post.Content, &post.Likes, &post.Dislikes)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	response := gin.H {
		"data": posts,
		"total":len(posts),
	}

	c.JSON(http.StatusOK, response)
}

func (h *Controller) PostLike(c *gin.Context) {
	queryString := "UPDATE posts SET likes = likes + 1 WHERE uuid = $1;"

	u, ok := c.GetQuery("uuid")
	if !ok || !isValidUUID(u) {
		c.String(http.StatusBadRequest, "Provide valid `uuid` parameter")
		return
	}

	singleTransaction(h, c, queryString, u)
}

func (h *Controller) PostDislike(c *gin.Context) {
	queryString := "UPDATE posts SET dislikes = dislikes + 1 WHERE uuid = $1;"

	u, ok := c.GetQuery("uuid")
	if !ok || !isValidUUID(u) {
		c.String(http.StatusBadRequest, "Provide valid `uuid` parameter")
		return
	}

	singleTransaction(h, c, queryString, u)
}

func (h *Controller) PostNewPost(c *gin.Context) {
	queryString := "INSERT INTO posts(content) VALUES ($1);"

	var req newPostRequestBody

	if err := c.BindJSON(&req); err != nil {
		// `_ =` to silence lint, no way to react to this
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	singleTransaction(h, c, queryString, req.Content)
}
