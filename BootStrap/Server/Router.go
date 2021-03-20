package Server

import (
	"xwork/App/Https/frontend"
	"xwork/BootStrap/Artisan"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/pprof"
)

// 未找到页面
// @param ctx iris.Context
// @author wf
// @time 2020-11-16
func notFound(ctx iris.Context) {
	ctx.JSON(Artisan.JsonError(Artisan.NOTFOUND, "not found", nil))
	ctx.StopExecution()
}

// 网络服务错误
// @param ctx iris.Context
// @author wf
// @time 2020-11-16
func internalServerError(ctx iris.Context) {
	ctx.JSON(Artisan.JsonError(Artisan.SERVICEFAIL, "internal server error", nil))
	ctx.StopExecution()
}

// 初始化路由
// @param app *iris.Application
// @return *iris.Application
// @author wf
// @time 2020-11-16
func initRouter(app *iris.Application) *iris.Application {
	//404页面
	app.OnErrorCode(iris.StatusNotFound, notFound)
	//服务器错误
	app.OnErrorCode(iris.StatusInternalServerError, internalServerError)
	//程序监控路由
	//app = pprofRouter(app)
	//内部接口路由
	app = frontendRouter(app)
	return app
}

// 程序监控路由
// @param ctx iris.Context
// @author wf
// @time 2020-11-16
func pprofRouter(app *iris.Application) *iris.Application {
	p := pprof.New()
	app.Any("/debug/pprof", p)
	app.Any("/debug/pprof/{action:path}", p)
	return app
}

func frontendRouter(app *iris.Application) *iris.Application {
	var apiV1 = app.Party("/v1").AllowMethods(iris.MethodOptions)
	apiV1.Get("/", frontend.Index)
	return app
}
