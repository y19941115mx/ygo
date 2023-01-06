package command

import (
	"fmt"

	"gitee.com/y19941115mx/ygo/framework/cobra"
	"gitee.com/y19941115mx/ygo/framework/contract"
	"github.com/kr/pretty"
)

func initConfigCommand() *cobra.Command {
	configCommand.AddCommand(configGetCommand)
	return configCommand
}

// configCommand 是命令行参数第一级为app的命令，它没有实际功能，只是打印帮助文档
var configCommand = &cobra.Command{
	Use:   "config",
	Short: "配置服务相关功能",
	RunE: func(c *cobra.Command, args []string) error {
		// 打印帮助文档
		c.Help()
		return nil
	},
}

// appStartCommand 启动一个Web服务
var configGetCommand = &cobra.Command{
	Use:     "get",
	Short:   "获取相关配置信息",
	Example: "./ygo config get database.mysql ",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		configService := container.MustMake(contract.ConfigKey).(contract.Config)

		if len(args) != 1 {
			fmt.Println("参数错误")
			return nil
		}

		configPath := args[0]
		val := configService.Get(configPath)

		if val == nil {
			fmt.Println("配置路径 ", configPath, " 不存在")
			return nil
		}

		fmt.Printf("%#v", pretty.Formatter(val))

		return nil
	},
}
