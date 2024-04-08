package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ikarizxc/http-server/internal/entities/users"
	response "github.com/ikarizxc/http-server/internal/handler/responce"
	"github.com/ikarizxc/http-server/pkg/hash/password"
)

type UserCreator interface {
	Create(user *users.User) (int, error)
}

func Create(userCreator UserCreator) func(c *gin.Context) {
	return func(c *gin.Context) {
		op := "handlers.users.Create : "

		var input *users.User

		if err := c.BindJSON(&input); err != nil {
			response.NewErrorResponce(c, http.StatusBadRequest, "incorrect input data", op+"error occured while binding json")
			return
		}

		pass, err := password.GeneratePasswordHash(input.Password)
		if err != nil {
			response.NewErrorResponce(c, http.StatusInternalServerError, "", op+err.Error())
			return
		}
		input.Password = pass

		id, err := userCreator.Create(input)
		if err != nil {
			response.NewErrorResponce(c, http.StatusInternalServerError, "", op+err.Error())
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"id": id,
		})
	}
}
