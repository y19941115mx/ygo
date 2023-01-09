package http

import (
	"github.com/y19941115mx/ygo/app/http/module/demo"
	"github.com/y19941115mx/ygo/framework/gin"
	"github.com/y19941115mx/ygo/framework/middleware"
)

// Routes 绑定业务层路由
func Routes(r *gin.Engine) {

	r.Static("/dist/", "./dist/")

	r.Use(middleware.Trace())

	demo.Register(r)
}
