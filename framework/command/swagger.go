package command

import (
	"fmt"
	"path/filepath"

	"github.com/swaggo/swag/gen"
	"github.com/y19941115mx/ygo/framework/cobra"
	"github.com/y19941115mx/ygo/framework/contract"
	"github.com/y19941115mx/ygo/framework/util"
)

func initSwaggerCommand() *cobra.Command {
	swaggerCommand.AddCommand(swaggerGenCommand)
	return swaggerCommand
}

var swaggerCommand = &cobra.Command{
	Use:   "swagger",
	Short: "swagger对应命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

// swaggerGenCommand 生成具体的swagger文档
var swaggerGenCommand = &cobra.Command{
	Use:   "gen",
	Short: "生成对应的swagger文件, contain swagger.yaml, doc.go",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		httpFolder := appService.HttpFolder()
		outputDir := filepath.Join(httpFolder, "swagger")

		conf := &gen.Config{
			// 遍历需要查询注释的目录
			SearchDir: httpFolder,
			// 不包含哪些文件
			Excludes: "",
			// 输出目录
			OutputDir: outputDir,
			// 输出类型
			OutputTypes: []string{"go", "json", "yaml"},
			// 整个swagger接口的说明文档注释
			MainAPIFile: "swagger.go",
			// 名字的显示策略，比如首字母大写等
			PropNamingStrategy: "",
			// 是否要解析vendor目录
			ParseVendor: false,
			// 是否要解析外部依赖库的包
			ParseDependency: false,
			// 是否要解析标准库的包
			ParseInternal: false,
			// 是否要查找markdown文件，这个markdown文件能用来为tag增加说明格式
			MarkdownFilesDir: "",
			// 是否应该在docs.go中生成时间戳
			GeneratedTime: false,
		}
		err := gen.New().Build(conf)
		if err != nil {
			fmt.Println(err)
		}
		return util.RebuildApp()

	},
}
