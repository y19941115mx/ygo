package user

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/y19941115mx/ygo/app/http/middleware/jwt"
	"github.com/y19941115mx/ygo/framework"
	"github.com/y19941115mx/ygo/framework/contract"
	"github.com/y19941115mx/ygo/framework/gin"
	"github.com/y19941115mx/ygo/framework/provider/orm"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type UserService struct {
	Service
	container framework.Container
	logger    contract.Log
	configer  contract.Config
	db        *gorm.DB
	cache     contract.CacheService
}

//生成随机验证码
func genCaptcha(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

//验证注册用户的合法性：邮箱，用户名唯一
func (u *UserService) isRegisterUserValid(user *User) error {
	userDB := &User{}
	if u.db.Where(&User{Email: user.Email}).First(userDB).Error != gorm.ErrRecordNotFound {
		return errors.New("邮箱已注册用户，不能重复注册")
	}
	if u.db.Where(&User{UserName: user.UserName}).First(userDB).Error != gorm.ErrRecordNotFound {
		return errors.New("用户名已经被注册，请换一个用户名")
	}
	return nil
}

func NewUserService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	logger := container.MustMake(contract.LogKey).(contract.Log)
	configer := container.MustMake(contract.ConfigKey).(contract.Config)
	cache := container.MustMake(contract.CacheKey).(contract.CacheService)
	db, err := container.MustMake(contract.ORMKey).(contract.ORMService).GetDB(orm.WithConfigPath("database.default"))

	if err != nil {
		return nil, err
	}

	if err = db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}

	return &UserService{container: container, logger: logger, configer: configer, db: db, cache: cache}, nil
}

func (u *UserService) Register(ctx context.Context, user *User) (*User, error) {
	if err := u.isRegisterUserValid(user); err != nil {
		return nil, err
	}

	user.Captcha = genCaptcha(10)

	// 将请求注册写入缓存，保存一小时
	key := fmt.Sprintf("user:register:%v", user.Captcha)
	if err := u.cache.SetObj(ctx, key, user, 1*time.Hour); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) SendRegisterMail(ctx context.Context, user *User) error {
	// 配置服务中获取发送邮件需要的参数
	host := u.configer.GetString("app.smtp.host")
	port := u.configer.GetInt("app.smtp.port")
	username := u.configer.GetString("app.smtp.username")
	password := u.configer.GetString("app.smtp.password")
	from := u.configer.GetString("app.smtp.from")
	domain := u.configer.GetString("app.domain")

	// 实例化gomail
	d := gomail.NewDialer(host, port, username, password)

	// 组装message
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetAddressHeader("To", user.Email, user.UserName)
	m.SetHeader("Subject", "感谢您注册我们的hadecast")
	link := fmt.Sprintf("%v/user/register/verify?captcha=%v", domain, user.Captcha)
	m.SetBody("text/html", fmt.Sprintf("请点击下面的链接完成注册：%s", link))

	// 发送电子邮件
	if err := d.DialAndSend(m); err != nil {
		u.logger.Error(ctx, "send email error", map[string]interface{}{
			"err":     err,
			"message": m,
		})
		return err
	}
	return nil
}

func (u *UserService) VerifyRegister(ctx context.Context, captcha string) error {
	//验证token
	key := fmt.Sprintf("user:register:%v", captcha)
	user := &User{}
	if err := u.cache.GetObj(ctx, key, user); err != nil {
		return err
	}

	if err := u.isRegisterUserValid(user); err != nil {
		return err
	}

	// 验证成功将密码存储数据库之前需要加密，不能原文存储进入数据库
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	// 具体在数据库创建用户
	if err := u.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserService) Login(ctx context.Context, user *User) (*User, error) {
	userDB := &User{}
	if err := u.db.Where("username=?", user.UserName).First(userDB).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password)); err != nil {
		return nil, err
	}

	userDB.Password = ""
	// 生成 token
	jwtConfigMap := u.configer.GetStringMapString("app.jwt")
	token, err := jwt.GenerateToken(jwtConfigMap, userDB.ID)
	if err != nil {
		return nil, err
	}

	userDB.Token = token
	return userDB, nil
}

func (u *UserService) GetUser(ctx context.Context, userID uint) (*User, error) {
	user := &User{}
	if err := u.db.WithContext(ctx).First(user, userID).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) GetLoginUser(ctx context.Context) (*User, error) {
	sessionKey := u.configer.GetString("app.jwt.session_key")
	userid, _ := ctx.(*gin.Context).Get(sessionKey)
	return u.GetUser(ctx, userid.(uint))
}
