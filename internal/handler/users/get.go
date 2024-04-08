package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ikarizxc/http-server/internal/entities/users"
	response "github.com/ikarizxc/http-server/internal/handler/responce"
)

type UserGetter interface {
	GetById(id int) (*users.User, error)
}

func Get(userGetter UserGetter) func(c *gin.Context) {
	return func(c *gin.Context) {
		op := "handlers.users.Get : "

		userId := c.Param("user_id")

		id, err := strconv.Atoi(userId)

		if err != nil {
			response.NewErrorResponce(c, http.StatusBadRequest, "id must be a number", op+"wrong id")
			return
		}

		user, err := userGetter.GetById(id)
		if err != nil {
			response.NewErrorResponce(c, http.StatusInternalServerError, "", op+err.Error())
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
