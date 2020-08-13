package route

import (
	"EasyGo/context"
	"EasyGo/validate"
	"EasyGo/view"
	"net/http"
)

//基础路由控制器
type Controller struct {
	Context  *context.Context
	verifier *validate.Validate
}

func (this *Controller) Init(ctx *context.Context) {
	this.Context = ctx
	this.Context.Output.SuccessCode = CODE_OK
	this.Context.Output.ErrorCode = CODE_ERROR
	this.Context.Response.Status = http.StatusOK
}

func (c *Controller) Success(vars ...interface{}) context.Render {
	return c.Context.Output.Success(vars...)
}

func (c *Controller) Error(vars ...interface{}) context.Render {
	return c.Context.Output.Error(vars...)
}

func (c *Controller) Json(code int, msg string, data interface{}, status ...int) {
	c.Context.Output.Json(code, msg, data, status...)
}

func (c *Controller) Write(str string) {
	c.Context.Output.WriteString(str)
}

func (c *Controller) View(path string, data map[string]interface{}) {
	_ = view.View(c.Context, path, data)
}

/**
 * 验证表单参数
 * param 验证参数
 * rule 验证规则
 */
func (c *Controller) Validator(params map[string]string, ruleSet map[string][]string) (bool, string) {
	if c.verifier == nil {
		c.verifier = validate.NewValidate()
	}
	return c.verifier.Verify(params, ruleSet)
}
