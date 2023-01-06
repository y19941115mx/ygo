package kernel

import (
	"errors"
	"net/http"

	"gitee.com/y19941115mx/ygo/framework/contract"
	"gitee.com/y19941115mx/ygo/framework/gin"
)

// 引擎服务
type YgoKernelService struct {
	contract.Kernel
	engine *gin.Engine
}

// 初始化 web 引擎服务实例
func NewYgoKernelService(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("kernel param error")
	}
	httpEngine := params[0].(*gin.Engine)
	return &YgoKernelService{engine: httpEngine}, nil
}

// 返回 web 引擎
func (s *YgoKernelService) HttpEngine() http.Handler {
	return s.engine
}
