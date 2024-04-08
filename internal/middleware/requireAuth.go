package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ikarizxc/http-server/internal/entities/users"
	response "github.com/ikarizxc/http-server/internal/handler/responce"
)

type UserGetter interface {
	GetById(id int) (*users.User, error)
}

func RequireAuth(userGetter UserGetter) func(c *gin.Context) {
	return func(c *gin.Context) {
		op := "handler.middleware.RequireAuth : "

		c.Set("Authenticated", false)

		accessToken, err := c.Cookie("access_token")
		if err != nil {
			response.NewErrorResponce(c, http.StatusUnauthorized, "unauthorized", op+"no access token in cookie")
			return
		}

		token, _ := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		// распарсили клеймсы
		if claims, ok := token.Claims.(jwt.MapClaims); ok {

			c.Set("userId", int(claims["sub"].(float64)))
			c.Set("isAdmin", claims["isAdmin"])

			// проверяем на просрочку
			if float64(time.Now().Unix()) >= claims["exp"].(float64) {
				// response.NewErrorResponce(c, http.StatusUnauthorized, "unauthorized")
				return
			}

			// получаем юзера по айди из бд
			_, err := userGetter.GetById(int(claims["sub"].(float64)))
			if err != nil {
				// response.NewErrorResponce(c, http.StatusUnauthorized, "unauthorized")
				return
			}

			c.Set("Authenticated", true)
			c.Next()
		} else {
			// response.NewErrorResponce(c, http.StatusUnauthorized, "unauthorized")
			return
		}
	}
}
