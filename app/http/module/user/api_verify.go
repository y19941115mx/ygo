package user

import (
	provider "github.com/y19941115mx/ygo/app/provider/user"
	"github.com/y19941115mx/ygo/framework/gin"
	"github.com/y19941115mx/ygo/framework/util/httputil"
)

// Verify 代表验证注册信息
// @Summary 验证注册信息
// @Description 使用token验证用户注册信息
// @Accept  json
// @Produce  json
// @Tags user
// @Param captcha query string true "注册的验证码"
// @Success 200 {object} httputil.Response
// @Failure 500  {object}  httputil.HTTPError
// @Router /user/register-verify [get]
func (api *UserApi) Verify(c *gin.Context) {
	// 验证参数
	userService := c.MustMake(provider.UserKey).(provider.Service)
	captcha := c.Query("captcha")
	if captcha == "" {
		err := httputil.BusinessError{Code: httputil.ERROR_PARAMETER_VALIDATION}
		httputil.FailWithError(c, err)
		return
	}

	err := userService.VerifyRegister(c, captcha)
	if err != nil {
		httputil.FailWithError(c, err)
		return
	}

	// 输出
	httputil.Ok(c)
}
