package posts

import (
	"strconv"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getPosts(c *gin.Context) {
	// TODO add DB query
	if queryLimitString, ok := c.GetQuery("last"); ok {
		if queryLimitUint64, err := ParseUint(queryLimitString, 10, 64); err == nil {
			// TODO add `limit queryLimitUint64`
		} else {
			c.String(http.StatusBadRequest, "Parameter \`last\` is invalid.\n\`last\`=" + queryLimitString)
			// TODO log error
		}
	}
}

func postLike() {

}

func postDislike() {

}

func postNewPost() {

}
