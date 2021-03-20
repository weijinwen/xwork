package Server

import (
	"xwork/App/Middleware/recover"
	"xwork/BootStrap/DbBase"
	"xwork/BootStrap/LogInit"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/spf13/viper"
)

func InitIris() {
	app := iris.New()
	app = LogInit.LogInit(app)
	//错误定义
	app.Use(recover.New())
	//日志定义
	app.Use(logger.New())
	//定义允许头部
	app.Use(cors.New(cors.Options{
		AllowedOrigins:   viper.GetStringSlice("AllowedOrigins"),
		MaxAge:           viper.GetInt("MaxAge"),
		AllowCredentials: viper.GetBool("AllowCredentials"),
		AllowedMethods:   []string{iris.MethodGet, iris.MethodPost},
		AllowedHeaders:   viper.GetStringSlice("AllowedHeaders"),
	}))
	//初始化路由
	app = initRouter(app)
	//关闭服务器后处理
	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		DbBase.CloseDB()
		app.Shutdown(ctx)
	})
	//配置服务器
	app.Configure(iris.WithConfiguration(iris.YAML("./Config/iris.yaml")))
	log.Fatal(app.Run(iris.Server(&http.Server{Addr: ":" + viper.GetString("Port")}), iris.WithoutInterruptHandler))
}
