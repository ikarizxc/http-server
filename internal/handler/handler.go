package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ikarizxc/http-server/internal/handler/authentication"
	usersHandler "github.com/ikarizxc/http-server/internal/handler/users"
	"github.com/ikarizxc/http-server/internal/middleware"
	"github.com/ikarizxc/http-server/internal/repository/users"
)

type Handler struct {
	usersStorage users.Storage
}

func NewHandler(usersStorage users.Storage) *Handler {
	return &Handler{usersStorage: usersStorage}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/signup", authentication.SignUp(h.usersStorage))
		auth.GET("/signin", authentication.SignIn(h.usersStorage))
	}

	router.GET("/validate", middleware.RequireAuth(h.usersStorage), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "logged",
		})
	})

	users := router.Group("/users", middleware.RequireAuth(h.usersStorage))
	{
		users.GET("/:user_id", usersHandler.Get(h.usersStorage))
		users.GET("/", usersHandler.GetAll(h.usersStorage))

		users.PATCH("/:user_id", usersHandler.Update(h.usersStorage))
		users.DELETE("/:user_id", usersHandler.Delete(h.usersStorage))
	}

	return router
}
