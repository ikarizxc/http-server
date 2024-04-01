package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ikarizxc/http-server/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.POST("/user", h.createUser)
	router.GET("/user/:user_id", h.getUser)
	router.GET("/users", h.getUsers)
	router.DELETE("/user/:user_id", h.deleteUser)

	return router
}