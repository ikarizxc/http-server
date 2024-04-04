package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ikarizxc/http-server/internal/entities/user"
	response "github.com/ikarizxc/http-server/internal/handler/responce"
	"github.com/ikarizxc/http-server/pkg/hash/password"
)

type UserCreator interface {
	Create(user *user.User) (int, error)
}

func Create(userCreator UserCreator) func(c *gin.Context) {
	return func(c *gin.Context) {
		var input *user.User

		if err := c.BindJSON(&input); err != nil {
			response.NewErrorResponce(c, http.StatusBadRequest, err.Error())
			return
		}

		pass, err := password.GeneratePasswordHash(input.Password)
		if err != nil {
			response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
			return
		}
		input.Password = pass

		id, err := userCreator.Create(input)
		if err != nil {
			response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"id": id,
		})
	}
}
