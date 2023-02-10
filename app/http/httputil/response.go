package httputil

import (
	"net/http"

	"github.com/y19941115mx/ygo/framework/gin"
)

type Response struct {
	Data interface{} `json:"data"`
	Meta
}

type Meta struct {
	Msg  string `json:"msg" example:"200"`
	Code int    `json:"code" example:"操作成功"`
}

// HTTPError example
type HTTPError struct {
	Code int    `json:"code" example:"500"`
	Msg  string `json:"msg" example:"操作失败"`
}

const (
	ERROR   = 500
	SUCCESS = 200
)

func result(c *gin.Context, code int, data interface{}, msg string) {
	c.JSON(http.StatusOK, Response{
		data,
		Meta{Code: code, Msg: msg},
	})
}

func Ok(c *gin.Context) {
	result(c, SUCCESS, map[string]interface{}{}, "操作成功")
}

func OkWithData(c *gin.Context, data interface{}) {
	result(c, SUCCESS, data, "操作成功")
}

func Fail(c *gin.Context) {
	result(c, ERROR, map[string]interface{}{}, "操作失败")
}

// 统一异常处理
func FailWithError(c *gin.Context, err error) {
	businessErr, ok := err.(BusinessError)
	if ok {
		result(c, businessErr.Code, map[string]interface{}{}, businessErr.Error())
	} else {
		result(c, ERROR, map[string]interface{}{}, err.Error())
	}
}

// 验证参数
func ValidateParameter(c *gin.Context, param interface{}) bool {
	if err := c.ShouldBind(param); err != nil {
		err = BusinessError{Code: ERROR_PARAMETER_VALIDATION}
		FailWithError(c, err)
		return false
	}
	return true
}
