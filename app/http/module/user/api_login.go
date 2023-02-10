package user

import (
	"github.com/y19941115mx/ygo/app/http/httputil"
	provider "github.com/y19941115mx/ygo/app/provider/user"
	"github.com/y19941115mx/ygo/framework/gin"
)

type loginParam struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,gte=6"`
}

type LoginResponse struct {
	httputil.Meta
	Data string `json:"data" example:"token"`
}

// Login 代表登录
// @Summary 登录
// @Description 用户登录接口
// @Accept  json
// @Produce  json
// @Tags user
// @Param loginParam body loginParam  true "login with param"
// @Success 200 {object} LoginResponse
// @Failure 200  {object}  httputil.HTTPError
// @Router /user/login [post]
func (api *UserApi) Login(c *gin.Context) {
	// 验证参数
	userService := c.MustMake(provider.UserKey).(provider.Service)

	param := &loginParam{}
	if valid := httputil.ValidateParameter(c, param); !valid {
		return
	}

	// 登录
	model := &provider.User{
		UserName: param.UserName,
		Password: param.Password,
	}
	userWithToken, err := userService.Login(c, model)
	if err != nil {
		httputil.FailWithError(c, err)
		return
	}
	// 输出
	httputil.OkWithData(c, userWithToken.Token)
}
