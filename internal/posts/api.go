package posts

import (
	"strconv"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"my-service/internal/entity"
)

type newPostRequestBody struct {
	Context		string		`json:"content"`
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func (h *Controller) getPosts(c *gin.Context) {
	queryString := "SELECT uid, context, likes, dislikes FROM posts;"

	if queryLimitString, ok := c.GetQuery("last"); ok {
		if _, err := ParseUint(queryLimitString, 10, 64); err == nil {
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
	response := gin.H {
		"data": make([]Post, 0, 32)
	}
	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.UUID, &post.Context, &post.Likes, &post.Dislikes)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			// TODO log error
			return
		}
		response["data"] = append(response["data"], post)
	}
	response["total"] = total
	// TODO log success
	c.JSON(http.StatusOK, response)
}

func (h *Controller) postLike(c *gin.Context) {
	queryString := "UPDATE posts SET likes = likes + 1 WHERE uuid = $1;"

	u, ok := c.GetQuery("uuid")
	if !ok || !isValidUUID(u) {
		c.String(http.StatusBadRequest, "Provide valid `uuid` parameter")
		return
	}

	if _, err = h.DB.Exec(queryString, u); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Controller) postDislike() {
	queryString := "UPDATE posts SET dislikes = dislikes + 1 WHERE uuid = $1;"

	u, ok := c.GetQuery("uuid")
	if !ok || !isValidUUID(u) {
		c.String(http.StatusBadRequest, "Provide valid `uuid` parameter")
		return
	}

	if _, err = h.DB.Exec(queryString, u); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Controller) postNewPost() {
	queryString := "INSERT INTO posts(content) VALUES ($1);"

	var req newPostRequestBody

	if err := c.BindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if _, err = h.DB.Exec(queryString, req.content); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
