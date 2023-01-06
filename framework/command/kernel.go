package command

import "gitee.com/y19941115mx/ygo/framework/cobra"

func AddKernelCommands(root *cobra.Command) {
	// app
	root.AddCommand(initAppCommand())

	// cobar
	root.AddCommand(initCronCommand())

	// env
	root.AddCommand(initEnvCommand())

}
