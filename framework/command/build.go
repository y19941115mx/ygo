package command

import (
	"github.com/y19941115mx/ygo/framework/cobra"
	"github.com/y19941115mx/ygo/framework/util"
)

var buildCommand = &cobra.Command{
	Use:   "build",
	Short: "编译项目",
	RunE: func(c *cobra.Command, args []string) error {
		return util.RebuildApp()
	},
}
