package middleware

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	response "github.com/ikarizxc/http-server/internal/handler/responce"
	"github.com/ikarizxc/http-server/pkg/tokens"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type TokenManager interface {
	WriteRefreshToken(id int, refreshToken string) error
	UpdateRefreshToken(id int, refreshToken string) error
	ReadRefreshToken(id int) (string, error)
}

func RefreshTokens(tokenManager TokenManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		isAuthenticated := c.GetBool("Authenticated")

		if isAuthenticated {
			// всё четко
			c.Next()
		} else {
			// рефреш токенов

			cookieAccessToken, err := c.Cookie("access_token")
			if err != nil {
				response.NewErrorResponce(c, http.StatusUnauthorized, "unauthorized")
				return
			}

			cookieRefreshToken, err := c.Request.Cookie("refresh_token")
			if err != nil {
				response.NewErrorResponce(c, http.StatusUnauthorized, "unauthorized")
				return
			}

			refreshTokenEncrypted := cookieRefreshToken.Value

			// check if refreshtoken is expired
			if cookieRefreshToken.MaxAge < 0 {
				response.NewErrorResponce(c, http.StatusUnauthorized, "unauthorized")
				return
			}

			refreshTokenDecryptedBytes, err := base64.StdEncoding.DecodeString(refreshTokenEncrypted)
			if err != nil {
				response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
				return
			}

			refreshTokenDecrypted := string(refreshTokenDecryptedBytes)

			// match refreshtoken with accesstoken
			if refreshTokenDecrypted[len(refreshTokenDecrypted)-8:] != cookieAccessToken[len(cookieAccessToken)-8:] {
				response.NewErrorResponce(c, http.StatusUnauthorized, "unauthorized")
				return
			}

			id := c.GetInt("userId")
			isAdmin := c.GetBool("isAdmin")

			// match refreshtokendecrypted with refreshtokenhash from db
			refreshTokenHashFromStorage, err := tokenManager.ReadRefreshToken(id)
			if err != nil {
				response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
				return
			}

			if err := bcrypt.CompareHashAndPassword([]byte(refreshTokenHashFromStorage), refreshTokenDecryptedBytes); err != nil {
				response.NewErrorResponce(c, http.StatusUnauthorized, err.Error())
				return
			}

			claims := jwt.MapClaims{
				"sub":     id,
				"exp":     time.Now().Add(time.Duration(viper.GetInt("accessTokenTTL")) * time.Second).Unix(),
				"isAdmin": isAdmin,
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
			if rt, err := tokenManager.ReadRefreshToken(id); rt == "" && err == nil {
				// id does not exist
				err = tokenManager.WriteRefreshToken(id, string(refreshTokenBcrypt))
				if err != nil {
					response.NewErrorResponce(c, http.StatusInternalServerError, err.Error())
					return
				}
			} else if rt != "" && err == nil {
				// id exist
				err = tokenManager.UpdateRefreshToken(id, string(refreshTokenBcrypt))
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

			c.Next()
		}

	}
}
