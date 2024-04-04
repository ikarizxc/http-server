package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ikarizxc/http-server/internal/entities/user"
	response "github.com/ikarizxc/http-server/internal/handler/responce"
)

type UserGetter interface {
	GetById(id int) (user.User, error)
}

func RequireAuth(userGetter UserGetter) func(c *gin.Context) {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("accessToken")
		if err != nil {
			response.NewErrorResponce(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			response.NewErrorResponce(c, http.StatusUnauthorized, err.Error())
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if float64(time.Now().Unix()) >= claims["exp"].(float64) {
				response.NewErrorResponce(c, http.StatusUnauthorized, "unauthorized")
				return
			}

			user, err := userGetter.GetById(int(claims["sub"].(float64)))
			if err != nil {
				if time.Now().Unix() >= claims["exp"].(int64) {
					response.NewErrorResponce(c, http.StatusInternalServerError, "unauthorized")
					return
				}
			}

			c.Set("userId", user.Id)
			c.Set("userIsAdmin", user.IsAdmin)

			c.Next()
		} else {
			response.NewErrorResponce(c, http.StatusUnauthorized, "unauthorized")
			return
		}
	}
}
