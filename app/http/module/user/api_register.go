package user

import (
	"fmt"

	"github.com/guonaihong/gout"
	"github.com/y19941115mx/ygo/app/http/httputil"
	provider "github.com/y19941115mx/ygo/app/provider/user"
	"github.com/y19941115mx/ygo/framework/gin"
)

type registerParam struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,gte=6"`
	Email    string `json:"email" binding:"required,gte=6"`
}

// Register 注册
// @Summary 用户注册
// @Description 用户注册接口
// @Accept  json
// @Produce  json
// @Tags user
// @Param registerParam body registerParam true "注册参数"
// @Success 200 {object} httputil.Response
// @Failure 500  {object}  httputil.HTTPError
// @Router /user/register [post]
func (api *UserApi) Register(c *gin.Context) {
	userService := c.MustMake(provider.UserKey).(provider.Service)
	logger := c.MustMakeLog()

	param := &registerParam{}
	if valid := httputil.ValidateBind(c, param); !valid {
		return
	}

	// 注册
	model := &provider.User{
		UserName: param.UserName,
		Password: param.Password,
		Email:    param.Email,
	}
	userWithCaptcha, err := userService.Register(c, model)
	if err != nil {
		logger.Error(c, err.Error(), map[string]interface{}{
			"stack": fmt.Sprintf("%+v", err),
		})
		httputil.FailWithError(c, err)
		return
	}

	if userWithCaptcha.UserName == "admin" {
		httputil.OkWithData(c, userWithCaptcha)
		return
	}

	if err := userService.SendRegisterMail(c, userWithCaptcha); err != nil {
		httputil.FailWithError(c, err)
		return
	}

	httputil.Ok(c)
}

// Register 添加测试用户
// @Summary 添加测试用户
// @Description 添加测试用户 用户名：admin 密码：admin123 邮箱：admin@123.com
// @Accept  json
// @Produce  json
// @Tags user
// @Success 200 {object} httputil.Response
// @Failure 500  {object}  httputil.HTTPError
// @Router /user/mock-test-user [get]
func (api *UserApi) RegisterMockUser(c *gin.Context) {
	config := c.MustMakeConfig()
	port := config.GetInt("app.address")
	param := &registerParam{
		UserName: "admin",
		Password: "admin123",
		Email:    "admin@123.com",
	}
	rsp := &httputil.Response{}

	url := fmt.Sprintf(":%d/user/register", port)
	// 调用注册接口
	err := gout.POST(url).SetJSON(param).BindJSON(rsp).Do()
	if err != nil {
		httputil.Fail(c)
		return
	}

	// 调用验证接口
	url = fmt.Sprintf(":%d/user/register-verify", port)
	u := rsp.Data.(map[string]interface{})
	// 调用注册接口
	err = gout.GET(url).SetQuery(gout.H{"captcha": u["Captcha"]}).Do()
	if err != nil {
		httputil.Fail(c)
		return
	}

	httputil.Ok(c)
}
