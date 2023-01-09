package demo

import (
	demoService "github.com/y19941115mx/ygo/app/provider/demo"
	"github.com/y19941115mx/ygo/framework/contract"
	"github.com/y19941115mx/ygo/framework/gin"
)

type DemoApi struct {
	service *Service
}

func Register(r *gin.Engine) error {
	api := NewDemoApi()
	r.Bind(&demoService.DemoProvider{})
	group := r.Group("/demo")

	group.GET("/demo", api.Demo)
	group.GET("/demo1", api.Demo1)
	group.GET("/demo2", api.Demo2)
	group.POST("/demo_post", api.DemoPost)
	return nil
}

func NewDemoApi() *DemoApi {
	service := NewService()
	return &DemoApi{service: service}
}

// Demo godoc
// @Summary 获取所有用户
// @Description 获取所有用户
// @Produce  json
// @Tags demo
// @Success 200 array []UserDTO
// @Router /demo/demo [get]
func (api *DemoApi) Demo(c *gin.Context) {
	// users := api.service.GetUsers()
	// usersDTO := UserModelsToUserDTOs(users)
	// c.JSON(200, usersDTO)
	configService := c.MustMakeConfig()
	password := configService.GetString("database.mysql.password")

	logger := c.MustMakeLog()
	logger.Info(c, "demo test info", map[string]interface{}{
		"api":  "demo/demo",
		"user": "jianfengye",
	})

	c.JSON(200, password)
}

func (api *DemoApi) Demo1(c *gin.Context) {
	// 获取password
	configService := c.MustMake(contract.ConfigKey).(contract.Config)
	password := configService.GetString("database.mysql.password")
	// 打印出来
	c.JSON(200, password)
}

// Demo godoc
// @Summary 获取所有学生
// @Description 获取所有学生
// @Produce  json
// @Tags demo
// @Success 200 array []UserDTO
// @Router /demo/demo2 [get]
func (api *DemoApi) Demo2(c *gin.Context) {
	demoProvider := c.MustMake(demoService.DemoKey).(demoService.IService)
	students := demoProvider.GetAllStudent()
	usersDTO := StudentsToUserDTOs(students)
	c.JSON(200, usersDTO)
}

func (api *DemoApi) DemoPost(c *gin.Context) {
	type Foo struct {
		Name string
	}
	foo := &Foo{}
	err := c.BindJSON(&foo)
	if err != nil {
		c.AbortWithError(500, err)
	}
	c.JSON(200, nil)
}
