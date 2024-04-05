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

		var fieldToUpdate = make(map[string]string)

		if err := c.BindJSON(&fieldToUpdate); err != nil {
			response.NewErrorResponce(c, http.StatusBadRequest, err.Error())
			return
		}

		if _, ok := fieldToUpdate["password"]; ok {
			fieldToUpdate["password_hash"], err = password.GeneratePasswordHash(fieldToUpdate["password"])
			if err != nil {
				response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
				return
			}
			delete(fieldToUpdate, "password")
		}

		err = userUpdater.Update(id, fieldToUpdate)
		if err != nil {
			response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}
