package response

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Error struct {
	Message string `json:"message"`
}

func NewErrorResponce(c *gin.Context, statusCode int, message, devMessage string) {
	logrus.Error(devMessage)
	if message == "" {
		c.AbortWithStatusJSON(statusCode, map[string]string{})
	}
	c.AbortWithStatusJSON(statusCode, Error{Message: message})
}
