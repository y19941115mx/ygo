package config

import (
	"github.com/y19941115mx/ygo/framework"
	"github.com/y19941115mx/ygo/framework/contract"
)

type YgoConfigProvider struct {
}

// Register registe a new function for make a service instance
func (provider *YgoConfigProvider) Register(c framework.Container) framework.NewInstance {
	return NewYgoConfig
}

// Boot will called when the service instantiate
func (provider *YgoConfigProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *YgoConfigProvider) IsDefer() bool {
	return true
}

// Params define the necessary params for NewInstance
func (provider *YgoConfigProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}

// / Name define the name for this service
func (provider *YgoConfigProvider) Name() string {
	return contract.ConfigKey
}
