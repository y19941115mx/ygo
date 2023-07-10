package model

import (
	"github.com/y19941115mx/ygo/framework/cobra"
)

// 代表输出路径
var output string

// 代表数据库连接
var database string

// 代表表格
var table string

// InitModelCommand 获取model相关的命令
func InitModelCommand() *cobra.Command {

	// model test
	modelTestCommand.Flags().StringVarP(&database, "database", "d", "database.default", "模型连接的数据库配置")
	modelCommand.AddCommand(modelTestCommand)

	// model gen
	modelGenCommand.Flags().StringVarP(&output, "output", "o", "app/model", "模型文件输出的文件夹位置")
	modelGenCommand.Flags().StringVarP(&database, "database", "d", "database.default", "模型连接的数据库配置")
	modelCommand.AddCommand(modelGenCommand)

	// model api
	modelApiCommand.Flags().StringVarP(&database, "database", "d", "database.default", "连接的数据库配置")
	modelApiCommand.Flags().StringVarP(&output, "module", "m", "test", "模块名称")
	modelApiCommand.Flags().StringVarP(&table, "table", "t", "default", "模块连接的数据表")
	modelCommand.AddCommand(modelApiCommand)
	return modelCommand
}

// modelCommand 模型相关的命令
var modelCommand = &cobra.Command{
	Use:   "model",
	Short: "数据库模型相关的命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}
