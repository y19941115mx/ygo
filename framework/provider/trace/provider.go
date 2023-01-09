package trace

import (
	"github.com/y19941115mx/ygo/framework"
	"github.com/y19941115mx/ygo/framework/contract"
)

type YgoTraceProvider struct {
}

// Register registe a new function for make a service instance
func (provider *YgoTraceProvider) Register(c framework.Container) framework.NewInstance {
	return NewYgoTraceService
}

// Boot will called when the service instantiate
func (provider *YgoTraceProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *YgoTraceProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *YgoTraceProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}

/// Name define the name for this service
func (provider *YgoTraceProvider) Name() string {
	return contract.TraceKey
}
