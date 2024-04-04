package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ikarizxc/http-server/internal/entities/user"
	response "github.com/ikarizxc/http-server/internal/handler/responce"
)

type UsersGetter interface {
	GetAll() ([]*user.User, error)
}

func GetAll(usersCreator UsersGetter) func(c *gin.Context) {
	return func(c *gin.Context) {
		users, err := usersCreator.GetAll()
		if err != nil {
			response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, users)
	}
}
