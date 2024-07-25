package main

import (
	"gin.example.com/entity/config"
	"gin.example.com/middleware/log"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	logConfig := &config.LogConfig{
		Level:      "debug",
		Filename:   "gin_example.log",
		MaxSize:    30,
		MaxAge:     30,
		MaxBackups: 10,
	}
	log.InitLogger(logConfig)

	r.Use(log.GinLogger(), log.GinRecovery(true))
	r.GET("/ping", func(c *gin.Context) {
		log.Debug("ping request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
