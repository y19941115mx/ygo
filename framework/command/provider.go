package command

import (
	"fmt"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jianfengye/collection"
	"github.com/pkg/errors"
	"github.com/y19941115mx/ygo/framework"
	"github.com/y19941115mx/ygo/framework/cobra"
	"github.com/y19941115mx/ygo/framework/contract"
	"github.com/y19941115mx/ygo/framework/util"
)

func initProviderCommand() *cobra.Command {
	providerCommand.AddCommand(providerListCommand)
	providerCommand.AddCommand(providerCreateCommand)

	return providerCommand
}

// providerCommand 一级命令
var providerCommand = &cobra.Command{
	Use:   "provider",
	Short: "服务提供相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

// providerListCommand 列出容器内的所有服务
var providerListCommand = &cobra.Command{
	Use:   "list",
	Short: "列出容器内的所有服务",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		ygoContainer := container.(*framework.YgoContainer)
		// 获取字符串凭证
		list := ygoContainer.NameList()
		// 打印
		for _, line := range list {
			println(line)
		}
		return nil
	},
}

// providerCreateCommand 创建一个新的服务，包括服务提供者，服务接口协议，服务实例
var providerCreateCommand = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "init"},
	Short:   "创建一个服务",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		fmt.Println("创建一个服务")
		var name string
		var folder string
		{
			prompt := &survey.Input{
				Message: "请输入服务名称(服务凭证)：",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}
		}

		{
			prompt := &survey.Input{
				Message: "请输入服务所在目录名称(默认: 同服务名称):",
			}
			err := survey.AskOne(prompt, &folder)
			if err != nil {
				return err
			}
		}

		// 检查服务是否存在
		providers := container.(*framework.YgoContainer).NameList()
		providerColl := collection.NewStrCollection(providers)
		if providerColl.Contains(name) {
			fmt.Println("服务名称已经存在")
			return nil
		}

		if folder == "" {
			folder = name
		}

		app := container.MustMake(contract.AppKey).(contract.App)
		pFolder := app.ProviderFolder()
		folderPath := filepath.Join(pFolder, folder)

		//  创建contract.go 需要检查目录是否存在
		if err := util.CreateFileTemlate(true, folderPath, "contract.go", contractTmp, name); err != nil {
			return errors.Cause(err)
		}

		//  创建provider.go
		if err := util.CreateFileTemlate(false, folderPath, "provider.go", providerTmp, name); err != nil {
			return errors.Cause(err)
		}

		//  创建service.go
		if err := util.CreateFileTemlate(false, folderPath, "service.go", serviceTmp, name); err != nil {
			return errors.Cause(err)
		}

		fmt.Println("创建服务成功, 文件夹地址:", folderPath)
		return nil
	},
}

var contractTmp string = `package {{.}}

const {{.|title}}Key = "{{.}}"

type Service interface {
	// 请在这里定义你的方法
    Foo() string
}
`

var providerTmp string = `package {{.}}

import (
	"github.com/y19941115mx/ygo/framework"
)

type {{.|title}}Provider struct {
	framework.ServiceProvider
}

func (sp *{{.|title}}Provider) Name() string {
	return {{.|title}}Key
}

func (sp *{{.|title}}Provider) Register(c framework.Container) framework.NewInstance {
	return New{{.|title}}Service
}

func (sp *{{.|title}}Provider) IsDefer() bool {
	return true 
}

func (sp *{{.|title}}Provider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}

func (sp *{{.|title}}Provider) Boot(c framework.Container) error {
	return nil
}

`

var serviceTmp string = `package {{.}}

import "github.com/y19941115mx/ygo/framework"

type {{.|title}}Service struct {
	Service
	container framework.Container
}

func New{{.|title}}Service(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	return &{{.|title}}Service{container: container}, nil
}

func (s *{{.|title}}Service) Foo() string {
    return ""
}
`
