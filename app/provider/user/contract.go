package user

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

const UserKey = "user"

type Service interface {

	// Register 注册用户,注意这里只是将用户注册, 并没有激活, 需要调用
	// 参数：user必填，username，password, email
	// 返回值： user 带上token
	Register(ctx context.Context, user *User) (*User, error)
	// SendRegisterMail 发送注册的邮件
	// 参数：user必填： username, password, email, token
	SendRegisterMail(ctx context.Context, user *User) error
	// VerifyRegister 注册用户，验证注册信息, 返回验证是否成功
	VerifyRegister(ctx context.Context, captcha string) error

	// Login 登录相关，使用用户名密码登录，获取完成User信息
	Login(ctx context.Context, user *User) (*User, error)

	// GetUser 获取用户信息
	GetUser(ctx context.Context, userID uint) (*User, error)

	// GetLoginUser 获取登录的用户信息
	GetLoginUser(ctx context.Context) (*User, error)
}

// User 代表一个用户，注意这里的用户信息字段在不同接口和参数可能为空
type User struct {
	gorm.Model
	UserName  string    `gorm:"column:username"`
	Password  string    `gorm:"column:password"`
	Email     string    `gorm:"column:email"`
	CreatedAt time.Time `gorm:"column:created_at"`

	Captcha string `gorm:"-"` // 注册验证码
	Token   string `gorm:"-"` // 用户token
}

// 使用到缓存功能 需要实现BinaryMarshaler接口
func (b *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(b)
}

// 使用到缓存功能 需要实现 BinaryUnMarshaler接口
func (b *User) UnmarshalBinary(bt []byte) error {
	return json.Unmarshal(bt, b)
}
