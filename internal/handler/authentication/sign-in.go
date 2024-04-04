package authentication

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ikarizxc/http-server/internal/entities/user"
	response "github.com/ikarizxc/http-server/internal/handler/responce"
	"github.com/ikarizxc/http-server/pkg/hash/password"
)

type UserGetter interface {
	GetByEmail(email string) (user.User, error)
}

func SignIn(userGetter UserGetter) func(c *gin.Context) {
	return func(c *gin.Context) {

		type UserLogin struct {
			Email    string
			Password string
		}

		var input UserLogin

		if err := c.BindJSON(&input); err != nil {
			response.NewErrorResponce(c, http.StatusBadRequest, err.Error())
			return
		}

		user, err := userGetter.GetByEmail(input.Email)
		if err != nil {
			response.NewErrorResponce(c, http.StatusUnauthorized, "no user with this email")
			return
		}

		if password.Compare(input.Password, user.Password) {
			response.NewErrorResponce(c, http.StatusUnauthorized, "wrong password")
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":     user.Id,
			"exp":     time.Now().Add(time.Hour).Unix(),
			"isAdmin": user.IsAdmin,
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("accessToken", tokenString, 3600, "", "", false, true)

		c.JSON(http.StatusOK, gin.H{})
	}
}
