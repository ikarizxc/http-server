package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ikarizxc/http-server/internal/entities/user"
	response "github.com/ikarizxc/http-server/internal/handler/responce"
)

func (h *Handler) createUser(c *gin.Context) {
	var input *user.User

	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponce(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.CreateUser(input)
	if err != nil {
		response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getUser(c *gin.Context) {
	userId := c.Param("user_id")

	id, err := strconv.Atoi(userId)

	if err != nil {
		response.NewErrorResponce(c, http.StatusBadRequest, "id must be a number")
		return
	}

	user, err := h.services.GetUser(id)
	if err != nil {
		response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) getUsers(c *gin.Context) {
	users, err := h.services.GetUsers()
	if err != nil {
		response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) deleteUser(c *gin.Context) {
	userId := c.Param("user_id")

	id, err := strconv.Atoi(userId)

	if err != nil {
		response.NewErrorResponce(c, http.StatusBadRequest, "id must be a number")
		return
	}

	err = h.services.DeleteUser(id)
	if err != nil {
		response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}
