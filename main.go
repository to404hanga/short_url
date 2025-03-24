package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	initViperWatch()

	app := Init()

	app.cron.Start()
	defer func() {
		<-app.cron.Stop().Done()
	}()

	server := app.server
	server.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})
	server.Run(viper.GetString("app.addr"))
}

func initViperWatch() {
	cfile := pflag.String("config",
		"config/config.yaml", "配置文件路径")
	pflag.Parse()
	// 直接指定文件路径
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
