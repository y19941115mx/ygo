package demo

import (
	demoService "github.com/y19941115mx/ygo/app/provider/demo"
	"github.com/y19941115mx/ygo/framework/gin"
)

type api struct{}

func Register(r *gin.Engine) error {
	r.Bind(&demoService.DemoProvider{})
	group := r.Group("/demo")

	group.GET("/demo", Demo)
	group.POST("/demo_post", DemoPost)
	return nil
}

// Demo godoc
// @Summary 获取所有学生
// @Description 获取所有学生
// @Produce  json
// @Tags demo
// @Success 200 array []UserDTO
// @Router /demo/demo [get]
func Demo(c *gin.Context) {
	demoProvider := c.MustMake(demoService.DemoKey).(demoService.IService)
	students := demoProvider.GetAllStudent()
	usersDTO := StudentsToUserDTOs(students)
	c.JSON(200, usersDTO)
}

func DemoPost(c *gin.Context) {
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
