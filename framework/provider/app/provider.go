package app

import (
	"github.com/y19941115mx/ygo/framework"
	"github.com/y19941115mx/ygo/framework/contract"
)

type YgoAppProvider struct {
	BaseFolder string
}

// Register 注册HadeApp方法
func (h *YgoAppProvider) Register(container framework.Container) framework.NewInstance {
	return NewYgoApp
}

// Boot 启动调用
func (h *YgoAppProvider) Boot(container framework.Container) error {
	return nil
}

// IsDefer 是否延迟初始化
func (h *YgoAppProvider) IsDefer() bool {
	return false
}

// Params 获取初始化参数
func (h *YgoAppProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container, h.BaseFolder}
}

// Name 获取字符串凭证
func (h *YgoAppProvider) Name() string {
	return contract.AppKey
}
