package test

import (
	"fmt"

	"github.com/y19941115mx/ygo/framework"
	"github.com/y19941115mx/ygo/framework/provider/app"
	"github.com/y19941115mx/ygo/framework/provider/config"
	"github.com/y19941115mx/ygo/framework/provider/env"
	"github.com/y19941115mx/ygo/framework/provider/log"
	"github.com/y19941115mx/ygo/framework/util"
)

const (
	BasePath = "C:\\Users\\19608\\Desktop\\project\\ygo" // 自定义
)

func InitBaseContainer() framework.Container {
	// 初始化服务容器
	container := framework.NewYgoContainer()
	fmt.Println(util.GetExecDirectory())
	container.Bind(&app.YgoAppProvider{BaseFolder: BasePath})

	container.Bind(&env.YgoEnvProvider{IsTest: true})
	container.Bind(&config.YgoConfigProvider{})
	container.Bind(&log.YgoLogServiceProvider{})
	return container
}
