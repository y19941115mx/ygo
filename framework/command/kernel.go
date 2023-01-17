package command

import "github.com/y19941115mx/ygo/framework/cobra"

func AddKernelCommands(root *cobra.Command) {
	// app
	root.AddCommand(initAppCommand())

	// cron
	root.AddCommand(initCronCommand())

	// env
	root.AddCommand(initEnvCommand())

	// build
	root.AddCommand(initBuildCommand())

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

	// deploy
	root.AddCommand(initDeployCommand())
}
