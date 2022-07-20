package middleware

import (
	"strconv"
	"net/http"

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

func (h *Controller) GetPosts(c *gin.Context) {
	queryString := "SELECT uid, context, likes, dislikes FROM posts;"

	if queryLimitString, ok := c.GetQuery("last"); ok {
		if _, err := strconv.ParseUint(queryLimitString, 10, 64); err == nil {
			queryString += " LIMIT " + queryLimitString
		} else {
			c.String(http.StatusBadRequest, "Parameter `last` is invalid.\n`last`=" + queryLimitString)
			// TODO log error
			return
		}
	}

	rows, err := h.DB.Query(queryString)
	if err != nil {
		// TODO log error
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	total := 0
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
		"total":total,
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

	if _, err := h.DB.Exec(queryString, u); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Controller) PostDislike(c *gin.Context) {
	queryString := "UPDATE posts SET dislikes = dislikes + 1 WHERE uuid = $1;"

	u, ok := c.GetQuery("uuid")
	if !ok || !isValidUUID(u) {
		c.String(http.StatusBadRequest, "Provide valid `uuid` parameter")
		return
	}

	if _, err := h.DB.Exec(queryString, u); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Controller) PostNewPost(c *gin.Context) {
	queryString := "INSERT INTO posts(content) VALUES ($1);"

	var req newPostRequestBody

	if err := c.BindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if _, err := h.DB.Exec(queryString, req.Content); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
