package authentication

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
		c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	}
}
