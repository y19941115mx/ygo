package orm

import (
	"context"
	"time"

	"github.com/y19941115mx/ygo/framework/contract"
	"gorm.io/gorm/logger"
)

// OrmLogger orm的日志实现类, 实现了gorm.Logger.Interface
type OrmLogger struct {
	logger contract.Log // 有一个logger对象存放ygo的log服务
}

// NewOrmLogger 初始化一个ormLogger,
func NewOrmLogger(logger contract.Log) *OrmLogger {
	return &OrmLogger{logger: logger}
}

// LogMode 什么都不实现，日志级别完全依赖ygo的日志定义
func (o *OrmLogger) LogMode(level logger.LogLevel) logger.Interface {
	return o
}

// Info 对接ygo的info输出
func (o *OrmLogger) Info(ctx context.Context, s string, i ...interface{}) {
	fields := map[string]interface{}{
		"fields": i,
	}
	o.logger.Info(ctx, s, fields)
}

// Warn 对接ygo的Warn输出
func (o *OrmLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	fields := map[string]interface{}{
		"fields": i,
	}
	o.logger.Warn(ctx, s, fields)
}

// Error 对接ygo的Error输出
func (o *OrmLogger) Error(ctx context.Context, s string, i ...interface{}) {
	fields := map[string]interface{}{
		"fields": i,
	}
	o.logger.Error(ctx, s, fields)
}

// Trace 对接ygo的Trace输出
func (o *OrmLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	elapsed := time.Since(begin)
	fields := map[string]interface{}{
		"begin": begin,
		"error": err,
		"sql":   sql,
		"rows":  rows,
		"time":  elapsed,
	}

	s := "orm trace sql"
	o.logger.Trace(ctx, s, fields)
}
