package http

import (
	"gitee.com/y19941115mx/ygo/app/http/module/demo"
	"gitee.com/y19941115mx/ygo/framework/gin"
)

// Routes 绑定业务层路由
func Routes(r *gin.Engine) {

	r.Static("/dist/", "./dist/")

	demo.Register(r)
}
