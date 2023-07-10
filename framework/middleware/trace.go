package middleware

import (
	"github.com/y19941115mx/ygo/framework/contract"
	"github.com/y19941115mx/ygo/framework/gin"
)

// trace 机制，实现全链路日志
func Trace() gin.HandlerFunc {
	// 使用函数回调
	return func(c *gin.Context) {
		// 记录开始时间
		tracer := c.MustMake(contract.TraceKey).(contract.Trace)
		traceCtx := tracer.ExtractHTTP(c.Request)

		tracer.WithTrace(c, traceCtx)

		// 使用next执行具体的业务逻辑
		c.Next()

		// 访问其他服务的接口时，需要将全链路字段加入request
		// tc := tracer.StartSpan(c)
		// tracer.InjectHTTP(c.Request, tc)

	}
}
