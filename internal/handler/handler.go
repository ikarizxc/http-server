package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ikarizxc/http-server/internal/handler/authentication"
	"github.com/ikarizxc/http-server/internal/handler/user"
	"github.com/ikarizxc/http-server/internal/middleware"
	userRepository "github.com/ikarizxc/http-server/internal/repository/user"
)

type Handler struct {
	userRepository *userRepository.Postgres
}

func NewHandler(userRepository *userRepository.Postgres) *Handler {
	return &Handler{userRepository: userRepository}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/signup", authentication.SignUp(h.userRepository))
		auth.GET("/signin", authentication.SignIn(h.userRepository))
	}

	router.GET("/validate", middleware.RequireAuth(h.userRepository), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "logged",
		})
	})

	users := router.Group("/users", middleware.RequireAuth(h.userRepository))
	{
		users.GET("/:user_id", user.Get(h.userRepository))
		users.GET("/", user.GetAll(h.userRepository))

		users.PATCH("/:user_id", user.Update(h.userRepository))
		users.DELETE("/:user_id", user.Delete(h.userRepository))
	}

	return router
}
