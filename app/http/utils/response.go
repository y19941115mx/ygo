package utils

import (
	"net/http"

	"github.com/y19941115mx/ygo/framework/gin"
)

type Response struct {
	Data interface{} `json:"data"`
	Meta `json:"meta"`
}
type Meta struct {
	Msg  string `json:"msg"`
	Code int    `json:"status"`
}

const (
	ERROR   = 500
	SUCCESS = 200
)

func result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		data,
		Meta{Code: code, Msg: msg},
	})
}

func Ok(c *gin.Context) {
	result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithData(data interface{}, c *gin.Context) {
	result(SUCCESS, data, "操作成功", c)
}

func Fail(c *gin.Context) {
	result(ERROR, map[string]interface{}{}, "操作失败", c)
}

func FailWithError(err error, c *gin.Context) {
	// 统一异常处理
	businessErr, ok := err.(BusinessError)
	if ok {
		result(businessErr.Code, map[string]interface{}{}, businessErr.Error(), c)
	} else {
		result(ERROR, map[string]interface{}{}, err.Error(), c)
	}
}
