package authentication

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ikarizxc/http-server/internal/entities/users"
	response "github.com/ikarizxc/http-server/internal/handler/responce"
	"github.com/ikarizxc/http-server/pkg/hash/password"
	"github.com/ikarizxc/http-server/pkg/tokens"
	"github.com/spf13/viper"
)

type UserGetter interface {
	GetByEmail(email string) (*users.User, error)
}

type TokenManager interface {
	WriteRefreshToken(id int, refreshToken string) error
	UpdateRefreshToken(id int, refreshToken string) error
	ReadRefreshToken(id int) (string, error)
}

func SignIn(userGetter UserGetter, tokenManager TokenManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		type UserSignIn struct {
			Email    string
			Password string
		}

		var input UserSignIn

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

		claims := jwt.MapClaims{
			"sub":     user.Id,
			"exp":     time.Now().Add(time.Duration(viper.GetInt("accessTokenTTL")) * time.Second).Unix(),
			"isAdmin": user.IsAdmin,
		}

		// генерация рефреш + акцес токена
		accessToken, err := tokens.GenerateAccessToken(claims)
		if err != nil {
			response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
			return
		}

		refreshToken, err := tokens.GenerateRefreshToken(accessToken)
		if err != nil {
			response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
			return
		}

		// содаём хэш токена для хранения в монге
		refreshTokenBcrypt, err := tokens.GenerateHashToken(refreshToken)
		if err != nil {
			response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
			return
		}

		// запись рефреш токена в монгу
		if rt, err := tokenManager.ReadRefreshToken(user.Id); rt == "" && err == nil {
			// id does not exist
			err = tokenManager.WriteRefreshToken(user.Id, string(refreshTokenBcrypt))
			if err != nil {
				response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
				return
			}
		} else if rt != "" && err == nil {
			// id exist
			err = tokenManager.UpdateRefreshToken(user.Id, string(refreshTokenBcrypt))
			if err != nil {
				response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
				return
			}
		} else if err != nil {
			// error occurred
			response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
			return
		}

		// в кукисах храним в базе64
		refreshTokenBase64 := base64.StdEncoding.EncodeToString([]byte(refreshToken))

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("access_token", accessToken, viper.GetInt("refreshTokenTTL"), "", "", false, true)
		c.SetCookie("refresh_token", refreshTokenBase64, viper.GetInt("refreshTokenTTL"), "", "", false, true)

		c.JSON(http.StatusOK, gin.H{
			"access_token": accessToken,
		})
	}
}
