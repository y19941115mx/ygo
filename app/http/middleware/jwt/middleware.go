package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/y19941115mx/ygo/app/http/httputil"
	"github.com/y19941115mx/ygo/framework/gin"
)

type MyClaims struct {
	UserId uint `json:"userid"`
	jwt.StandardClaims
}

func GenerateToken(c *gin.Context, userId uint) (string, error) {
	configer := c.MustMakeConfig()

	jwtConfigMap := configer.GetStringMapString("app.jwt")
	tokenExpireDuration, err := time.ParseDuration(jwtConfigMap["token_expire_duration"])
	if err != nil {
		return "", err
	}

	// 创建一个我们自己的声明
	claims := MyClaims{
		userId, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpireDuration).Unix(), // 过期时间
			Issuer:    jwtConfigMap["issuer"],                     // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(jwtConfigMap["secret_signature"])
}

func parseToken(secret string, tokenString string) (*jwt.Token, *MyClaims, error) {
	Claims := &MyClaims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return secret, nil
	})
	return token, Claims, err
}

// JwtMiddleware 代表中间件函数
func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		configer := c.MustMakeConfig()

		jwtConfigMap := configer.GetStringMapString("app.jwt")
		authHeader := c.Request.Header.Get(jwtConfigMap["auth_header_key"])

		var err error // 异常信息
		if authHeader == "" {
			err = httputil.BusinessError{
				Code: httputil.ERROR_TOKEN_NOT_EXIST,
			}
			httputil.FailWithError(c, err)
			c.Abort()
			return
		}

		getToken, claims, err := parseToken(jwtConfigMap["secret_signature"], authHeader)
		if getToken.Valid {
			primarykey := claims.UserId
			c.Set(jwtConfigMap["session_key"], primarykey)
			c.Next()
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// token已过期
				err = httputil.BusinessError{
					Code: httputil.ERROR_TOKEN_EXPIRE,
				}
			} else {
				// 错误的token
				err = httputil.BusinessError{
					Code: httputil.ERROR_TOKEN_WRONG,
				}
			}
			httputil.FailWithError(c, err)
			c.Abort()
		}
	}
}
