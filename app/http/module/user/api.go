package user

import (
	"github.com/y19941115mx/ygo/app/http/middleware/jwt"
	"github.com/y19941115mx/ygo/app/provider/user"
	"github.com/y19941115mx/ygo/framework/gin"
)

type UserApi struct{}

func RegisterRoutes(r *gin.Engine) error {
	api := &UserApi{}
	if !r.IsBind(user.UserKey) {
		r.Bind(&user.UserProvider{})
	}

	// 注册
	r.POST("/user/register", api.Register)
	// 注册验证
	r.GET("/user/register-verify", api.Verify)
	// 登录
	r.POST("/user/login", api.Login)
	// 获取登录用户信息
	r.GET("/user/userinfo", jwt.JwtMiddleware(), api.UserInfo)
	// 创建测试用户
	r.GET("/user/mock-test-user", api.RegisterMockUser)

	return nil
}
