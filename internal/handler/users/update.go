package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	response "github.com/ikarizxc/http-server/internal/handler/responce"
	"github.com/ikarizxc/http-server/pkg/hash/password"
)

type UserUpdater interface {
	Update(int, map[string]string) error
}

func Update(userUpdater UserUpdater) func(c *gin.Context) {
	return func(c *gin.Context) {
		op := "handlers.users.Update : "

		userId := c.Param("user_id")

		id, err := strconv.Atoi(userId)

		if err != nil {
			response.NewErrorResponce(c, http.StatusBadRequest, "id must be a number", op+"wrong id")
			return
		}

		curUserId := c.GetInt("userId")
		userIsAdmin := c.GetBool("userIsAdmin")
		if !userIsAdmin && id != curUserId {
			response.NewErrorResponce(c, http.StatusUnauthorized, "no rights to update user", op+"not admin tries to update another user")
			return
		}

		var fieldToUpdate = make(map[string]string)

		if err := c.BindJSON(&fieldToUpdate); err != nil {
			response.NewErrorResponce(c, http.StatusBadRequest, "incorrect input data", op+"error occured while binding json")
			return
		}

		if _, ok := fieldToUpdate["password"]; ok {
			fieldToUpdate["password_hash"], err = password.GeneratePasswordHash(fieldToUpdate["password"])
			if err != nil {
				response.NewErrorResponce(c, http.StatusInternalServerError, "", op+err.Error())
				return
			}
			delete(fieldToUpdate, "password")
		}

		err = userUpdater.Update(id, fieldToUpdate)
		if err != nil {
			response.NewErrorResponce(c, http.StatusInternalServerError, "", op+err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}
