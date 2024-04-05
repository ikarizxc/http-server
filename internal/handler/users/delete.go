package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	response "github.com/ikarizxc/http-server/internal/handler/responce"
)

type UserDeleter interface {
	Delete(id int) error
}

func Delete(userDeleter UserDeleter) func(c *gin.Context) {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		id, err := strconv.Atoi(userId)

		if err != nil {
			response.NewErrorResponce(c, http.StatusBadRequest, "id must be a number")
			return
		}

		curUserId := c.GetInt("userId")
		userIsAdmin := c.GetBool("userIsAdmin")
		if !userIsAdmin && id != curUserId {
			response.NewErrorResponce(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		err = userDeleter.Delete(id)
		if err != nil {
			response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success",
		})
	}
}
