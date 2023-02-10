package user

import (
	"fmt"

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
	if valid := httputil.ValidateParameter(c, param); !valid {
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

	if err := userService.SendRegisterMail(c, userWithCaptcha); err != nil {
		httputil.FailWithError(c, err)
		return
	}

	httputil.Ok(c)
}
