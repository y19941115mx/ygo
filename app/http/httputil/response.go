package httputil

import (
	"strconv"

	"github.com/y19941115mx/ygo/framework/gin"
)

type Response struct {
	Data interface{} `json:"data"`
	Meta
}

type Meta struct {
	Msg  string `json:"msg" example:"操作成功"`
	Code int    `json:"code" example:"200"`
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

func successResult(c *gin.Context, code int, data interface{}, msg string) {
	c.JSON(SUCCESS, Response{
		data,
		Meta{Code: code, Msg: msg},
	})
}

func failResult(c *gin.Context, code int, data interface{}, msg string) {
	c.JSON(ERROR, Response{
		data,
		Meta{Code: code, Msg: msg},
	})
}

func Ok(c *gin.Context) {
	successResult(c, SUCCESS, map[string]interface{}{}, "操作成功")
}

func OkWithData(c *gin.Context, data interface{}) {
	successResult(c, SUCCESS, data, "操作成功")
}

func Fail(c *gin.Context) {
	failResult(c, ERROR, map[string]interface{}{}, "操作失败")
}

// 统一异常处理
func FailWithError(c *gin.Context, err error) {
	businessErr, ok := err.(BusinessError)
	if ok {
		failResult(c, businessErr.Code, map[string]interface{}{}, businessErr.Error())
	} else {
		failResult(c, ERROR, map[string]interface{}{}, err.Error())
	}
}

// bind请求体并进行参数验证
func ValidateBind(c *gin.Context, param interface{}) bool {
	logger := c.MustMakeLog()
	if err := c.ShouldBind(param); err != nil {
		logger.Error(c, err.Error(), nil)
		err = BusinessError{Code: ERROR_PARAMETER_VALIDATION}
		FailWithError(c, err)
		return false
	}
	return true
}

// query 参数验证
func ValidateQuery(c *gin.Context, query string) bool {
	if id := c.Query(query); id == "" {
		err := BusinessError{Code: ERROR_QUERY_NOT_EXIST}
		FailWithError(c, err)
		return false
	}
	return true
}

// query 参数验证
func ValidateIntQuery(c *gin.Context, query string, result *int) bool {
	var id int
	var err error

	if !ValidateQuery(c, query) {
		return false
	}

	if id, err = strconv.Atoi(c.Query(query)); err != nil {
		err := BusinessError{Code: ERROR_PARAMETER_VALIDATION}
		FailWithError(c, err)
		return false
	}

	*result = id
	return true
}
