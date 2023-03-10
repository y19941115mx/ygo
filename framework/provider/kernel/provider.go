package kernel

import (
	"github.com/y19941115mx/ygo/framework"
	"github.com/y19941115mx/ygo/framework/contract"
	"github.com/y19941115mx/ygo/framework/gin"
)

// YgoKernelProvider 提供web引擎
type YgoKernelProvider struct {
	HttpEngine *gin.Engine
}

// Register 注册服务提供者
func (provider *YgoKernelProvider) Register(c framework.Container) framework.NewInstance {
	return NewYgoKernelService
}

// Boot 启动的时候判断是否由外界注入了Engine，如果注入的化，用注入的，如果没有，重新实例化
func (provider *YgoKernelProvider) Boot(c framework.Container) error {
	if provider.HttpEngine == nil {
		provider.HttpEngine = gin.Default()
	}
	return nil
}

// IsDefer 引擎的初始化我们希望开始就进行初始化
func (provider *YgoKernelProvider) IsDefer() bool {
	return false
}

// Params 参数就是一个HttpEngine
func (provider *YgoKernelProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.HttpEngine}
}

// Name 提供凭证
func (provider *YgoKernelProvider) Name() string {
	return contract.KernelKey
}
