package main

import (
	"gin.example.com/entity/config"
	"gin.example.com/logic"
	"gin.example.com/middleware"
	"gin.example.com/middleware/log"
	"github.com/gin-gonic/gin"
)

func main() {
	configPath := "config.yaml"
	r := gin.New()
	config.InitConfig(configPath)
	log.InitLogger(config.GetLogConfig())
	r.Use(log.GinLogger(), log.GinRecovery(true))

	// v1 := r.Group("/v1").Handlers()
	v1 := r.Group("/v1")
	v1.GET("/ping", middleware.HandleBindings(logic.Ping))
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
