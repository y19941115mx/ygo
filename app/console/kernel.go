package console

import (
	"github.com/y19941115mx/ygo/app/console/command/demo"
	"github.com/y19941115mx/ygo/framework"
	"github.com/y19941115mx/ygo/framework/cobra"
	"github.com/y19941115mx/ygo/framework/command"
)

// RunCommand  初始化根 Command 并运行
func RunCommand(container framework.Container) error {
	// 根 Command
	var rootCmd = &cobra.Command{
		// 定义根命令的关键字
		Use: "ygo",
		// 简短介绍
		Short: "ygo 命令",
		// 根命令的详细介绍
		Long: "ygo 框架提供的命令行工具，使用这个命令行工具能很方便执行框架自带命令，也能很方便编写业务命令",
		// 根命令的执行函数
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.InitDefaultHelpFlag()
			return cmd.Help()
		},
		// 不需要出现 cobra 默认的 completion 子命令
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
	// 为根 Command 设置服务容器
	rootCmd.SetContainer(container)
	// 绑定框架的命令
	command.AddKernelCommands(rootCmd)
	// 绑定业务的命令
	addAppCommand(rootCmd)
	// 执行 RootCommand
	return rootCmd.Execute()
}

// 绑定业务的命令
func addAppCommand(rootCmd *cobra.Command) {
	// AddCronCommand 添加定时执行命令 这里代表每秒调用一次Foo命令
	rootCmd.AddCronCommand("* * * * * *", demo.FooCommand)
}
