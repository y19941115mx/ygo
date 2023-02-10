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
	BasePath = "D:\\project\\ygo"
)

func InitBaseContainer() framework.Container {
	// 初始化服务容器
	container := framework.NewYgoContainer()
	fmt.Println(util.GetExecDirectory())
	container.Bind(&app.YgoAppProvider{BaseFolder: BasePath})

	container.Bind(&env.YgoTestingEnvProvider{})
	container.Bind(&config.YgoConfigProvider{})
	container.Bind(&log.YgoLogServiceProvider{})
	return container
}
