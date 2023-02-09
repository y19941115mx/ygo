package env

import (
	"github.com/y19941115mx/ygo/framework"
	"github.com/y19941115mx/ygo/framework/contract"
)

type YgoTestingEnvProvider struct {
	Folder string
}

// Register registe a new function for make a service instance
func (provider *YgoTestingEnvProvider) Register(c framework.Container) framework.NewInstance {
	return NewYgoTestingEnv
}

// Boot will called when the service instantiate
func (provider *YgoTestingEnvProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *YgoTestingEnvProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *YgoTestingEnvProvider) Params(c framework.Container) []interface{} {
	return []interface{}{}
}

/// Name define the name for this service
func (provider *YgoTestingEnvProvider) Name() string {
	return contract.EnvKey
}
