package main

import (
	"gin.example.com/entity/config"
	"gin.example.com/middleware/log"
	"github.com/gin-gonic/gin"
)

func main() {
	configPath := "config.yaml"
	r := gin.New()
	config.InitConfig(configPath)
	log.InitLogger(config.GetLogConfig())
	r.Use(log.GinLogger(), log.GinRecovery(true))
	r.GET("/ping", func(c *gin.Context) {
		log.Debug("ping request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
