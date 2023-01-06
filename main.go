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
	"gitee.com/y19941115mx/ygo/framework/provider/distributed"
	"gitee.com/y19941115mx/ygo/framework/provider/env"
	"gitee.com/y19941115mx/ygo/framework/provider/kernel"
)

func main() {
	// 初始化框架级别的服务容器
	container := framework.NewYgoContainer()
	// 绑定App服务
	container.Bind(&app.YgoAppProvider{})
	// 绑定本地分布式定时任务服务
	container.Bind(&distributed.LocalDistributedProvider{})
	// 绑定env服务
	container.Bind(&env.YgoEnvProvider{})
	// 绑定config服务
	container.Bind(&config.YgoConfigProvider{})

	// 将HTTP引擎初始化,并且作为服务提供者绑定到服务容器中
	if engine, err := http.NewHttpEngine(); err == nil {
		container.Bind(&kernel.YgoKernelProvider{HttpEngine: engine})
	}

	// 运行根 cmd
	console.RunCommand(container)
}
