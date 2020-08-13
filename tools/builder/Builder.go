package builder

import (
	"EasyGo/context"
	"EasyGo/kernel"
	"EasyGo/route"
	"EasyGo/tools/helper"
)

func init() {
	kernel.Inject(func() {
		route.Route("EasyBuilder/:fileType/:fileName", Create)
	})
}

//创建项目文件
func Create(ctx *context.Context) {
	types := ctx.Route.Input("fileType")
	fileName := ctx.Route.Input("fileName")
	if types == "" || fileName == "" {
		ctx.Output.WriteString("文件信息缺失")
		return
	}
	switch types {
	case "controller":
		//创建控制器文件
		handler := NewController(helper.UcFirst(fileName.Strings()))
		err := handler.build()
		if err != nil {
			ctx.Response.Error(err.Error())
			return
		}
		ctx.Output.WriteString("create successful")
	case "model":
		//创建model
		handler := NewModel(helper.UcFirst(fileName.Strings()))
		err := handler.build()
		if err != nil {
			ctx.Response.Error(err.Error())
			return
		}
		ctx.Output.WriteString("create successful")
	case "middleware":
		//创建中间件
		handler := NewMiddle(helper.UcFirst(fileName.Strings()))
		err := handler.build()
		if err != nil {
			ctx.Output.WriteString(err.Error())
			return
		}
		ctx.Output.WriteString("create successful")
	default:
		ctx.Output.WriteString("invalid file type")
	}
}
