package user

import (
	provider "github.com/y19941115mx/ygo/app/provider/user"
	"github.com/y19941115mx/ygo/framework/gin"
	"github.com/y19941115mx/ygo/framework/util/httputil"
)

type loginParam struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,gte=6"`
}

// Login 代表登录
// @Summary 用户登录
// @Description 用户登录接口, 使用 data 字段返回用户token
// @Accept  json
// @Produce  json
// @Tags user
// @Param loginParam body loginParam  true "login with param"
// @Success 200 {object} httputil.Response{data=string}
// @Failure 500  {object}  httputil.HTTPError
// @Router /user/login [post]
func (api *UserApi) Login(c *gin.Context) {
	// 验证参数
	userService := c.MustMake(provider.UserKey).(provider.Service)

	param := &loginParam{}
	if valid := httputil.ValidateBind(c, param); !valid {
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
