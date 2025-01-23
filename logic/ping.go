package logic

import (
	"gin.example.com/middleware/log"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context, input *string) (*string, error) {
	res := "pong"
	if *input != "" {
		res = *input + "pong"
	}
	log.Debugf(c, "Ping|res: %s", res)
	return &res, nil
}
