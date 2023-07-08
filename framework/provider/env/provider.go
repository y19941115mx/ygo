package env

import (
	"github.com/y19941115mx/ygo/framework"
	"github.com/y19941115mx/ygo/framework/contract"
)

type YgoEnvProvider struct {
	IsTest bool
}

// Register registe a new function for make a service instance
func (provider *YgoEnvProvider) Register(c framework.Container) framework.NewInstance {
	return NewYgoEnv
}

// Boot will called when the service instantiate
func (provider *YgoEnvProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *YgoEnvProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *YgoEnvProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c, provider.IsTest}
}

// / Name define the name for this service
func (provider *YgoEnvProvider) Name() string {
	return contract.EnvKey
}
