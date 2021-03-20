package frontend

import (
	"xwork/BootStrap/Artisan"
	"github.com/kataras/iris/v12"
)

func Index(ctx iris.Context) {
	ctx.JSON(Artisan.JsonData(200,"sss","sss"))
}
