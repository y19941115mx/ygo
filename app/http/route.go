package http

import (
	"github.com/y19941115mx/ygo/app/http/module/user"
	"github.com/y19941115mx/ygo/framework/contract"
	"github.com/y19941115mx/ygo/framework/gin"
	ginSwagger "github.com/y19941115mx/ygo/framework/middleware/gin-swagger"
	"github.com/y19941115mx/ygo/framework/middleware/gin-swagger/swaggerFiles"
)

// Routes 绑定业务层路由
func Routes(r *gin.Engine) {
	container := r.GetContainer()
	configService := container.MustMake(contract.ConfigKey).(contract.Config)

	// 如果配置了swagger，则显示swagger的中间件
	if configService.GetBool("app.swagger") {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	// 动态路由定义
	user.RegisterRoutes(r)
}
