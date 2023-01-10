package command

import "github.com/y19941115mx/ygo/framework/cobra"

func AddKernelCommands(root *cobra.Command) {
	// app
	root.AddCommand(initAppCommand())

	// cobar
	root.AddCommand(initCronCommand())

	// env
	root.AddCommand(initEnvCommand())

	// config
	root.AddCommand(initConfigCommand())

	// provider
	root.AddCommand(initProviderCommand())

	// cmd
	root.AddCommand(initCmdCommand())

	// middeware
	root.AddCommand(initMiddlewareCommand())

	// new
	root.AddCommand(initNewCommand())

	// Swagger
	root.AddCommand(initSwaggerCommand())
}
