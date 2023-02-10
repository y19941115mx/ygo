package user

import (
	"github.com/y19941115mx/ygo/app/http/httputil"
	provider "github.com/y19941115mx/ygo/app/provider/user"
	"github.com/y19941115mx/ygo/framework/gin"
)

type UserResponse struct {
	httputil.Meta
	Data UserDTO
}

// UserInfo 获取登录用户信息
// @Summary 获取登录用户信息
// @Description 获取登录用户信息接口
// @Accept  json
// @Produce  json
// @Tags user
// @Success 200 {object} UserResponse
// @Failure 200  {object}  httputil.HTTPError
// @Router /user/userinfo [get]
func (api *UserApi) UserInfo(c *gin.Context) {
	userService := c.MustMake(provider.UserKey).(provider.Service)

	user, err := userService.GetLoginUser(c)
	if err != nil {
		httputil.FailWithError(c, err)
		return
	}
	// 输出
	httputil.OkWithData(c, ConvertUserToDTO(user))
}
