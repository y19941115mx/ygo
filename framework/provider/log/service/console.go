package service

import (
	"os"

	"gitee.com/y19941115mx/ygo/framework"
	"gitee.com/y19941115mx/ygo/framework/contract"
)

type YgoConsoleLog struct {
	YgoLog
}

// NewYgoConsoleLog 实例化YgoConsoleLog
func NewYgoConsoleLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	log := &YgoConsoleLog{}

	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	// 最重要的将内容输出到控制台
	log.SetOutput(os.Stdout)
	log.c = c
	return log, nil
}
