// Copyright 2021 jianfengye.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package main

import (
	"gitee.com/y19941115mx/ygo/app/console"
	"gitee.com/y19941115mx/ygo/app/http"
	"gitee.com/y19941115mx/ygo/framework"
	"gitee.com/y19941115mx/ygo/framework/provider/app"
	"gitee.com/y19941115mx/ygo/framework/provider/config"
	"gitee.com/y19941115mx/ygo/framework/provider/env"
	"gitee.com/y19941115mx/ygo/framework/provider/kernel"
	"gitee.com/y19941115mx/ygo/framework/provider/log"
	"gitee.com/y19941115mx/ygo/framework/provider/trace"
)

func main() {
	// 初始化框架级别的服务容器
	container := framework.NewYgoContainer()
	// 绑定目录服务
	container.Bind(&app.YgoAppProvider{})
	// 绑定环境变量服务
	container.Bind(&env.YgoEnvProvider{})
	// 绑定配置服务
	container.Bind(&config.YgoConfigProvider{})
	// 绑定日志服务
	container.Bind(&log.YgoLogServiceProvider{})
	// 绑定trace全链条日志
	container.Bind(&trace.YgoTraceProvider{})

	// 将HTTP引擎初始化,并且作为服务提供者绑定到服务容器中
	if engine, err := http.NewHttpEngine(); err == nil {
		container.Bind(&kernel.YgoKernelProvider{HttpEngine: engine})
	}

	// 运行根 cmd
	console.RunCommand(container)
}
