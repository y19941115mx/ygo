package test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	provider "github.com/y19941115mx/ygo/app/provider/user"
	"github.com/y19941115mx/ygo/framework/contract"
	"github.com/y19941115mx/ygo/framework/provider/cache"
	"github.com/y19941115mx/ygo/framework/provider/orm"
)

// 测试正常的注册登录流程
func Test_UserRegisterLogin(t *testing.T) {
	container := InitBaseContainer()
	container.Bind(&orm.GormProvider{})
	container.Bind(&cache.CacheProvider{})

	ormService := container.MustMake(contract.ORMKey).(contract.ORMService)
	db, err := ormService.GetDB()
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&provider.User{}); err != nil {
		t.Fatal(err)
	}

	tmp, err := provider.NewUserService(container)
	if err != nil {
		t.Fatal(err)
	}
	us := tmp.(*provider.UserService)
	ctx := context.Background()

	user1 := &provider.User{
		UserName: "jianfengye",
		Password: "123456",
		Email:    "1960892068@qq.com",
	}

	Convey("正常注册登录流程", t, func() {

		Convey("注册用户", func() {
			userWithCaptcha, err := us.Register(ctx, user1)
			So(err, ShouldBeNil)
			user1.Captcha = userWithCaptcha.Captcha
		})

		Convey("发送邮件", func() {
			err := us.SendRegisterMail(ctx, user1)
			So(err, ShouldBeNil)
		})

		Convey("验证注册信息", func() {
			err := us.VerifyRegister(ctx, user1.Captcha)
			So(err, ShouldBeNil)
			// 数据库有数据
			userDB := &provider.User{}
			err = db.Where("username=?", user1.UserName).First(userDB).Error
			So(err, ShouldBeNil)
			So(userDB.ID, ShouldNotBeZeroValue)
		})

		Convey("用户登录", func() {
			userWithToken, err := us.Login(ctx, user1)
			So(err, ShouldBeNil)
			So(userWithToken, ShouldNotBeNil)
			user1.Token = userWithToken.Token
		})
	})
}
