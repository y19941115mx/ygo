package user

import (
	provider "github.com/y19941115mx/ygo/app/provider/user"
	"github.com/y19941115mx/ygo/framework/gin"
)

// Verify 代表验证注册信息
// @Summary 验证注册信息
// @Description 使用token验证用户注册信息
// @Accept  json
// @Produce  json
// @Tags user
// @Param captcha query string true "注册的验证码"
// @Success 200 {string} Message "注册成功，请进入登录页面"
// @Router /user/register/verify [get]
func (api *UserApi) Verify(c *gin.Context) {
	// 验证参数
	userService := c.MustMake(provider.UserKey).(provider.Service)
	captcha := c.Query("captcha")
	if captcha == "" {
		c.ISetStatus(400).IText("参数错误")
		return
	}

	verified, err := userService.VerifyRegister(c, captcha)
	if err != nil {
		c.ISetStatus(500).IText(err.Error())
		return
	}

	if !verified {
		c.ISetStatus(500).IText("验证错误")
		return
	}

	// 输出
	c.IRedirect("/#/login").IText("注册成功，请进入登录页面")
}
