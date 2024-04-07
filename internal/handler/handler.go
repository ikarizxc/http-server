package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ikarizxc/http-server/internal/handler/authentication"
	usersHandler "github.com/ikarizxc/http-server/internal/handler/users"
	"github.com/ikarizxc/http-server/internal/middleware"
	"github.com/ikarizxc/http-server/internal/repository/tokens"
	"github.com/ikarizxc/http-server/internal/repository/users"
)

type Handler struct {
	usersStorage  users.Storage
	tokensStorage tokens.Storage
}

func NewHandler(usersStorage users.Storage, tokensStorage tokens.Storage) *Handler {
	return &Handler{
		usersStorage:  usersStorage,
		tokensStorage: tokensStorage,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/signup", authentication.SignUp(h.usersStorage))
		auth.GET("/signin", authentication.SignIn(h.usersStorage, h.tokensStorage))
		auth.GET("/logout", authentication.Logout())
	}

	users := router.Group("/users", middleware.RequireAuth(h.usersStorage), middleware.RefreshTokens(h.tokensStorage))
	{
		users.GET("/:user_id", usersHandler.Get(h.usersStorage))
		users.GET("/", usersHandler.GetAll(h.usersStorage))

		users.PATCH("/:user_id", usersHandler.Update(h.usersStorage))
		users.DELETE("/:user_id", usersHandler.Delete(h.usersStorage))
	}

	return router
}
