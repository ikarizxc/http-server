package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	response "github.com/ikarizxc/http-server/internal/handler/responce"
)

func RequireAdmin(userGetter UserGetter) func(c *gin.Context) {
	return func(c *gin.Context) {
		userIsAdmin := c.GetBool("userIsAdmin")
		if userIsAdmin {
			c.Next()
		} else {
			response.NewErrorResponce(c, http.StatusUnauthorized, "unauthorized")
			return
		}
	}
}
