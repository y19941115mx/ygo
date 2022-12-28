package cobra

import "gitee.com/y19941115mx/ygo/framework"

// SetContainer 设置服务容器
func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

// GetContainer 获取容器
func (c *Command) GetContainer() framework.Container {
	return c.Root().container
}
