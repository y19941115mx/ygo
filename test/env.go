package test

import (
	"fmt"

	"github.com/y19941115mx/ygo/framework"
	"github.com/y19941115mx/ygo/framework/provider/app"
	"github.com/y19941115mx/ygo/framework/provider/env"
	"github.com/y19941115mx/ygo/framework/util"
)

const (
	BasePath = "C:\\Users\\19608\\Desktop\\project\\ygo"
)

func InitBaseContainer() framework.Container {
	// 初始化服务容器
	container := framework.NewYgoContainer()
	fmt.Println(util.GetExecDirectory())
	container.Bind(&app.YgoAppProvider{BaseFolder: BasePath})
	// 后续初始化需要绑定的服务提供者...
	container.Bind(&env.YgoTestingEnvProvider{})
	return container
}
