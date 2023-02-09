package utils

type BusinessError struct {
	Code int
}

func (e BusinessError) Error() string {
	return getErrMsg(e.Code)
}

const (
	// code= 100x... 用户模块的错误
	ERROR_USERNAME_USED        = 1001
	ERROR_PASSWORD_WRONG       = 1002
	ERROR_ABNORMAL_PERMISSIONS = 1003

	// code= 200x... token 相关错误
	ERROR_TOKEN_NOT_EXIST = 2001
	ERROR_TOKEN_WRONG     = 2002
	ERROR_TOKEN_EXPIRE    = 2003

	// code= xxx... 其他业务模块的错误

)

var codeMsg = map[int]string{
	ERROR_USERNAME_USED:        "用户名已存在！",
	ERROR_PASSWORD_WRONG:       "用户名或密码错误",
	ERROR_ABNORMAL_PERMISSIONS: "用户权限异常",

	ERROR_TOKEN_NOT_EXIST: "TOKEN不存在",
	ERROR_TOKEN_WRONG:     "TOKEN不正确,请重新登陆",
	ERROR_TOKEN_EXPIRE:    "TOKEN已过期,请重新登陆",
}

func getErrMsg(code int) string {
	return codeMsg[code]
}
